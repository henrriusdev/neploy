package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
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
	Create(ctx context.Context, user model.CreateUserRequest, oauthID string) error
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
}

type user struct {
	repos repository.Repositories
	email *email.Email
}

func NewUser(repos repository.Repositories) User {
	return &user{repos: repos, email: email.NewEmail()}
}

func (u *user) Create(ctx context.Context, req model.CreateUserRequest, oauthID string) error {
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

	// Create OAuth connection if we have an oauth_id
	if oauthID != "" {
		oauth := model.UserOAuth{
			UserID:   newUser.ID,
			OAuthID:  oauthID,
			Provider: model.Provider(req.Provider),
		}

		if err := u.repos.UserOauth.Insert(ctx, oauth); err != nil {
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return model.LoginResponse{}, err
	}

	roles, err := u.repos.UserRole.GetByUserID(ctx, user.ID)
	if err != nil {
		return model.LoginResponse{}, err
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

	t, err := token.SignedString([]byte(config.Env.JWTSecret))
	if err != nil {
		return model.LoginResponse{}, err
	}

	return model.LoginResponse{
		Token:     t,
		User:      user,
		RoleIDs:   roleIDs,
		RoleNames: roleNamesLower,
	}, nil
}

func (u *user) GetProvider(ctx context.Context, userID string) (string, error) {
	oauth, err := u.repos.UserOauth.GetByUserID(ctx, userID)
	if err != nil {
		logger.Error("failed to get user oauth: user_id=%s, error=%v", userID, err)
		return "", err
	}

	return string(oauth.Provider), nil
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
	inviteLink := fmt.Sprintf("%s:%s/users/invite/%s", config.Env.BaseURL, config.Env.Port, token)
	if err := u.email.SendInvitation(ctx, req.Email, teamName, req.Role, inviteLink); err != nil {
		// Log the error but don't fail the invitation creation
		logger.Error("failed to send invitation email: email=%s, role=%s, error=%v", req.Email, req.Role, err)
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

	// Compare the current password with the stored password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		return errors.New("current password is incorrect")
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
