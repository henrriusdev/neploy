package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"neploy.dev/config"
	"neploy.dev/pkg/email"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type User interface {
	Create(ctx context.Context, user model.CreateUserRequest) error
	Get(ctx context.Context, id string) (model.User, error)
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset uint) ([]model.TeamMemberResponse, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error)
	GetProvider(ctx context.Context, userID string) (string, error)
	InviteUser(ctx context.Context, req model.InviteUserRequest) error
	AcceptInvitation(ctx context.Context, token string) (model.Invitation, error)
	GetInvitationByToken(ctx context.Context, token string) (model.Invitation, error)
	AddUserRole(ctx context.Context, email, roleID string) error
	UpdateProfile(ctx context.Context, profileReq model.ProfileRequest, userID string) error
	UpdatePassword(ctx context.Context, req model.PasswordRequest, userID string) error
	UpdateTechStacks(ctx context.Context, req model.SelectUserTechStacksRequest) error
	GetAll(ctx context.Context) ([]model.User, error)
	NewPasswordLink(ctx context.Context, email, language string) (string, error)
}

type user struct {
	repos repository.Repositories
	email *email.Email
}

func NewUser(repos repository.Repositories) User {
	return &user{repos: repos, email: email.NewEmail()}
}

func (u *user) Create(ctx context.Context, req model.CreateUserRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		DOB:       req.DOB,
		Phone:     req.Phone,
		Address:   req.Address,
		Provider:  model.Provider(req.Provider), // Convert string to Provider type
	}

	newUser, err := u.repos.User.Create(ctx, user)
	if err != nil {
		return err
	}

	for _, roleName := range req.Roles {
		role, err := u.repos.Role.GetByName(ctx, roleName)
		if err != nil {
			logger.Error("error getting role: %v role=%s", err, roleName)
			return err
		}

		roleID := role.ID
		userRole := model.UserRoles{
			UserID: newUser.ID,
			RoleID: roleID,
		}
		if _, err := u.repos.UserRole.Insert(ctx, userRole); err != nil {
			return err
		}
	}

	return nil
}

func (u *user) Get(ctx context.Context, id string) (model.User, error) {
	user, err := u.repos.User.Get(ctx, id)
	if err != nil {
		logger.Error("failed to get user: user_id=%s, error=%v", id, err)
		return model.User{}, err
	}
	user.Password = ""
	return user, nil
}

func (u *user) Update(ctx context.Context, user model.User) error {
	return u.repos.User.Update(ctx, user)
}

func (u *user) Delete(ctx context.Context, id string) error {
	return u.repos.User.Delete(ctx, id)
}

func (u *user) List(ctx context.Context, limit, offset uint) ([]model.TeamMemberResponse, error) {
	users, err := u.repos.User.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var teamMembers []model.TeamMemberResponse

	// Convert each user to TeamMemberResponse
	for _, user := range users {
		// Get provider
		provider, err := u.GetProvider(ctx, user.ID)
		if err != nil {
			logger.Error("failed to get provider for user: user_id=%s, error=%v", user.ID, err)
			continue
		}

		// Get roles
		userRoles, err := u.repos.UserRole.GetByUserID(ctx, user.ID)
		if err != nil {
			logger.Error("failed to get roles for user: user_id=%s, error=%v", user.ID, err)
			continue
		}

		// Get tech stacks
		userTechStacks, err := u.repos.UserTechStack.GetByUserID(ctx, user.ID)
		if err != nil {
			logger.Error("failed to get tech stacks for user: user_id=%s, error=%v", user.ID, err)
			continue
		}
		var roles []model.Role

		for _, userRole := range userRoles {
			role, err := u.repos.Role.GetByID(ctx, userRole.RoleID)
			if err != nil {
				logger.Error("failed to get role: role_id=%s, error=%v", userRole.RoleID, err)
				continue
			}

			roles = append(roles, role)
		}

		var techStacks []model.TechStack

		for _, userTechStack := range userTechStacks {
			techStack, err := u.repos.TechStack.GetByID(ctx, userTechStack.TechStackID)
			if err != nil {
				logger.Error("failed to get tech stack: tech_stack_id=%s, error=%v", userTechStack.TechStackID, err)
				continue
			}

			techStacks = append(techStacks, techStack)
		}

		member := model.TeamMemberResponse{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Provider:   provider,
			Roles:      roles,
			TechStacks: techStacks,
		}

		teamMembers = append(teamMembers, member)
	}

	return teamMembers, nil
}

func (u *user) GetByEmail(ctx context.Context, email string) (model.User, error) {
	return u.repos.User.GetByEmail(ctx, email)
}

func (u *user) Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error) {
	user, err := u.GetByEmail(ctx, req.Email)
	if err != nil {
		return model.LoginResponse{}, err
	}

	// Critical security fix: Verify password before generating token
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return model.LoginResponse{}, errors.New("invalid credentials")
	}

	t, err := u.generateToken(ctx, user)
	if err != nil {
		return model.LoginResponse{}, err
	}

	// Don't return the password hash in the response
	user.Password = ""

	return model.LoginResponse{
		Token: t,
		User:  user,
	}, nil
}

func (u *user) GetProvider(ctx context.Context, userID string) (string, error) {
	user, err := u.repos.User.GetOneById(ctx, userID)
	if err != nil {
		logger.Error("failed to get user: user_id=%s, error=%v", userID, err)
		return "", err
	}

	return string(user.Provider), nil
}

func (u *user) InviteUser(ctx context.Context, req model.InviteUserRequest) error {
	// Check if user already exists
	_, err := u.GetByEmail(ctx, req.Email)
	if err == nil {
		logger.Error("user already exists: email=%s", req.Email)
		return err
	}

	// Generate invitation token
	token := generateInviteToken()

	role, err := u.repos.Role.GetByName(ctx, req.Role)
	if err != nil {
		logger.Error("failed to get role: role=%s, error=%v", req.Role, err)
		return err
	}

	// Create invitation record
	invitation := model.Invitation{
		Email:     req.Email,
		Role:      role.ID,
		Token:     token,
		ExpiresAt: model.Date{Time: time.Now().Add(7 * 24 * time.Hour)}, // 7 days
	}

	// Save invitation
	if err := u.repos.User.CreateInvitation(ctx, invitation); err != nil {
		logger.Error("failed to create invitation: email=%s, role=%s, error=%v", req.Email, req.Role, err)
		return err
	}

	teamName, err := u.repos.Metadata.GetTeamName(ctx)
	if err != nil {
		logger.Error("failed to get team name: error=%v", err)
		return err
	}

	// Send invitation email
	inviteLink := fmt.Sprintf("%s/users/invite/%s", config.Env.BaseURL, token)

	language, err := u.repos.Metadata.GetLanguage(ctx)
	if err != nil {
		logger.Error("failed to get language: error=%v", err)
		return err
	}
	if err := u.email.SendInvitation(ctx, req.Email, teamName, req.Role, inviteLink, language); err != nil {
		// Log the error but don't fail the invitation creation
		logger.Error("failed to send invitation email: email=%s, role=%s, language=%s, error=%v", req.Email, req.Role, language, err)
		return err
	}

	return nil
}

func (u *user) AcceptInvitation(ctx context.Context, token string) (model.Invitation, error) {
	// Obtener la invitación por token
	invitation, err := u.repos.User.GetInvitationByToken(ctx, token)
	if err != nil {
		logger.Error("failed to get invitation by token: token=%s, error=%v", token, err)
		return model.Invitation{}, err
	}

	// Verificar si la invitación ha expirado
	if time.Now().After(invitation.ExpiresAt.Time) {
		return model.Invitation{}, errors.New("invitation has expired")
	}

	// Marcar la invitación como aceptada
	now := model.Date{Time: time.Now()}
	invitation.AcceptedAt = &now
	if err := u.repos.User.UpdateInvitation(ctx, invitation); err != nil {
		return model.Invitation{}, err
	}

	return invitation, nil
}

func generateInviteToken() string {
	// Generar 32 bytes de datos aleatorios
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// En caso de error, usar un fallback menos seguro
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	// Convertir a base64URL (seguro para URLs y sin caracteres especiales)
	return base64.URLEncoding.EncodeToString(b)
}

func (u *user) GetInvitationByToken(ctx context.Context, token string) (model.Invitation, error) {
	return u.repos.User.GetInvitationByToken(ctx, token)
}

func (u *user) AddUserRole(ctx context.Context, email, roleID string) error {
	user, err := u.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	userRole := model.UserRoles{
		UserID: user.ID,
		RoleID: roleID,
	}

	_, err = u.repos.UserRole.Insert(ctx, userRole)
	return err
}

func ValidateJWT(token string) (model.JWTClaims, bool, error) {
	claims := model.JWTClaims{}
	t, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Env.JWTSecret), nil
	})
	if err != nil {
		return model.JWTClaims{}, false, err
	}
	return claims, t.Valid, nil
}

func (s *user) UpdateProfile(ctx context.Context, profileReq model.ProfileRequest, userID string) error {
	// Validate the ProfileRequest struct fields (You can enhance this as per your needs)
	if profileReq.Email == "" || profileReq.FirstName == "" || profileReq.LastName == "" {
		return errors.New("missing required fields: email, first name, last name")
	}

	// Create User model from ProfileRequest
	user := model.User{
		Email:     profileReq.Email,
		FirstName: profileReq.FirstName,
		LastName:  profileReq.LastName,
		DOB:       profileReq.Dob,
		Address:   profileReq.Address,
		Phone:     profileReq.Phone,
	}

	// Update the user profile in the repository
	_, err := s.repos.User.UpdateOneById(ctx, userID, user)
	if err != nil {
		return err
	}

	// Success
	return nil
}

// UpdatePassword updates the user's password
func (s *user) UpdatePassword(ctx context.Context, req model.PasswordRequest, userID string) error {
	// Fetch the user from the repository
	user, err := s.repos.User.Get(ctx, userID)
	if err != nil {
		return err
	}

	if !req.Reset {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
		if err != nil {
			return errors.New("current password is incorrect")
		}
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the password in the repository
	user.Password = string(hashedPassword)
	err = s.repos.User.Update(ctx, user)
	if err != nil {
		return err
	}

	// Success
	return nil
}

func (u *user) UpdateTechStacks(ctx context.Context, req model.SelectUserTechStacksRequest) error {
	userTStacks, err := u.repos.UserTechStack.GetByUserID(ctx, req.UserId)
	if err != nil {
		logger.Error("error getting user tech stacks %v", err)
		return err
	}

	for _, stack := range userTStacks {
		if err := u.repos.UserTechStack.Delete(ctx, stack.UserID, stack.TechStackID); err != nil {
			logger.Error("error deleting existing user tech stacks")
			return err
		}
	}

	newStacks := []model.UserTechStack{}
	for _, stack := range req.TechStackIDs {
		newStacks = append(newStacks, model.UserTechStack{
			TechStackID: stack,
			UserID:      req.UserId,
		})
	}

	if _, err := u.repos.UserTechStack.InsertMany(ctx, newStacks); err != nil {
		logger.Error("error creating user tech stacks %v", err)
		return err
	}

	return nil
}

func (u *user) GetAll(ctx context.Context) ([]model.User, error) {
	users, err := u.repos.User.GetAll(ctx)
	if err != nil {
		logger.Error("error getting all users %v", err)
		return nil, err
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (u *user) NewPasswordLink(ctx context.Context, userEmail, language string) (string, error) {
	user, err := u.repos.User.GetByEmail(ctx, userEmail)
	if err != nil {
		logger.Error("error getting user by email %v", err)
		return "", err
	}

	metadata, err := u.repos.Metadata.Get(ctx)
	if err != nil {
		logger.Error("error getting metadata %v", err)
		return "", err
	}

	// 2. Crear token
	token, err := u.generateToken(ctx, user)
	if err != nil {
		logger.Error("error generating password reset token %v", err)
		return "", err
	}

	// 3. Preparar los datos del email
	// Parse the base URL to handle it properly
	baseURL, err := url.Parse(config.Env.BaseURL)
	if err != nil {
		logger.Error("error parsing base URL %v", err)
		return "", err
	}

	// If the base URL doesn't have a port, add the configured port
	if baseURL.Port() == "" {
		baseURL.Host = baseURL.Host + ":" + config.Env.Port
	}

	// Set the path and query parameters
	baseURL.Path = "/password/change"
	query := baseURL.Query()
	query.Set("token", token)
	baseURL.RawQuery = query.Encode()

	resetURL := baseURL.String()

	emailData := email.PasswordResetData{
		UserName:     user.FirstName,
		CompanyName:  metadata.TeamName,
		LogoURL:      metadata.LogoURL,
		ResetURL:     resetURL,
		ResetToken:   token,
		BaseURL:      config.Env.BaseURL,
		CurrentYear:  time.Now().Year(),
		Translations: email.GetTranslations(language),
		Language:     language,
	}

	// 4. Enviar email
	if err := u.email.SendPasswordReset(user.Email, emailData); err != nil {
		logger.Error("error sending password reset email %v", err)
		return "", err
	}

	// retornar un mensaje de éxito
	return "Password reset email sent successfully", nil
}

func (u *user) generateToken(ctx context.Context, user model.User) (string, error) {
	roles, err := u.repos.UserRole.GetByUserID(ctx, user.ID)
	if err != nil {
		return "", err
	}

	roleNames := make([]string, len(roles))
	roleIDs := make([]string, len(roles))
	roleNamesLower := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Role.Name
		roleIDs[i] = role.Role.ID
		roleNamesLower[i] = strings.ToLower(role.Role.Name)
	}

	// create the JWT access token here and return it
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.JWTClaims{
		ID:         user.ID,
		Email:      user.Email,
		Roles:      roleNames,
		RoleIDs:    roleIDs,
		RolesLower: roleNamesLower,
		Name:       user.FirstName + " " + user.LastName,
		Username:   user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 5)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return token.SignedString([]byte(config.Env.JWTSecret))
}
