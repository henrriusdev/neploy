package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"neploy.dev/config"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type User interface {
	Create(ctx context.Context, user model.CreateUserRequest, oauthID int) error
	Get(ctx context.Context, id string) (model.User, error)
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset uint) ([]model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error)
	GetProvider(ctx context.Context, userID string) (string, error)
	InviteUser(ctx context.Context, req model.InviteUserRequest) error
}

type user struct {
	repos repository.Repositories
	email Email
}

func NewUser(repos repository.Repositories, email Email) User {
	return &user{repos: repos, email: email}
}

func (u *user) Create(ctx context.Context, req model.CreateUserRequest, oauthID int) error {
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
		Address:   req.Address,
		Phone:     req.Phone,
	}

	user, err = u.repos.User.Create(ctx, user)
	if err != nil {
		return err
	}

	for _, roleName := range req.Roles {
		role, err := u.repos.Role.GetByName(ctx, roleName)
		if err != nil {
			log.Err(err).Msg("error getting role")
			return err
		}

		roleID := role.ID
		userRole := model.UserRoles{
			UserID: user.ID,
			RoleID: roleID,
		}
		if _, err := u.repos.UserRole.Insert(ctx, userRole); err != nil {
			return err
		}
	}

	oauth := model.UserOAuth{
		UserID:   user.ID,
		Provider: model.Provider(req.Provider),
		OAuthID:  strconv.Itoa(oauthID),
	}

	if err := u.repos.UserOauth.Insert(ctx, oauth); err != nil {
		return err
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
		return err
	}

	// Generate invitation token
	token := generateInviteToken()

	// Create invitation record
	invitation := model.Invitation{
		Email:     req.Email,
		TeamID:    req.TeamID,
		Role:      req.Role,
		Token:     token,
		ExpiresAt: model.Date{Time: time.Now().Add(7 * 24 * time.Hour)}, // 7 days
	}

	// Save invitation
	if err := u.repos.User.CreateInvitation(ctx, invitation); err != nil {
		return err
	}

	teamName, err := u.repos.Metadata.GetTeamName(ctx)
	if err != nil {
		return err
	}

	// Send invitation email
	inviteLink := config.Env.BaseURL + "/invite/" + token
	if err := u.email.SendInvitation(ctx, req.Email, teamName, req.Role, inviteLink); err != nil {
		// Log the error but don't fail the invitation creation
		log.Error().Err(err).
			Str("email", req.Email).
			Str("role", req.Role).
			Msg("Failed to send invitation email")
	}

	return nil
}

func generateInviteToken() string {
	// Generate a random token
	// In production, use a proper UUID or secure token generator
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
