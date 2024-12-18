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
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type User interface {
	Create(ctx context.Context, user model.CreateUserRequest, oauthID string) error
	Get(ctx context.Context, id string) (model.User, error)
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset uint) ([]model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error)
	GetProvider(ctx context.Context, userID string) (string, error)
	InviteUser(ctx context.Context, req model.InviteUserRequest) error
	AcceptInvitation(ctx context.Context, token string) (model.Invitation, error)
	GetInvitationByToken(ctx context.Context, token string) (model.Invitation, error)
	AddUserRole(ctx context.Context, email, roleID string) error
}

type user struct {
	repos repository.Repositories
	email Email
}

func NewUser(repos repository.Repositories, email Email) User {
	return &user{repos: repos, email: email}
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
	return u.repos.User.Get(ctx, id)
}

func (u *user) Update(ctx context.Context, user model.User) error {
	return u.repos.User.Update(ctx, user)
}

func (u *user) Delete(ctx context.Context, id string) error {
	return u.repos.User.Delete(ctx, id)
}

func (u *user) List(ctx context.Context, limit, offset uint) ([]model.User, error) {
	return u.repos.User.List(ctx, limit, offset)
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
		Token: t,
		User:  user,
	}, nil
}

func (u *user) GetProvider(ctx context.Context, userID string) (string, error) {
	oauth, err := u.repos.UserOauth.GetByUserID(ctx, userID)
	if err != nil {
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
