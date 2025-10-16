package services

import (
	"context"
	"strings"
	"time"
	"user-service/config"
	"user-service/constants"
	errConstant "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/domain/models"
	"user-service/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.IRepositoryRegistry
}

type IUserService interface {
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, *dto.RegisterRequest) (*dto.RegisterRespose, error)
	Update(context.Context, *dto.UpdateUserRequest, string) (*dto.UserResponse, error)
	GetUserLogin(context.Context) (*dto.UserResponse, error)
	GetUserByUUID(context.Context, string) (*dto.UserResponse, error)
}

type Claims struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

func NewUserService(repository repositories.IRepositoryRegistry) IUserService {
	return &UserService{repository: repository}
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repository.GetUser().FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	expirationTime := time.Now().Add(time.Duration(config.Config.JwtExpireTime) * time.Minute).Unix()
	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        strings.ToLower(user.Role.Code),
	}
	claims := &Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	response := &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}
	return response, nil
}

func (s *UserService) isUserNameExist(ctx context.Context, username string) bool {
	user, err := s.repository.GetUser().FindByUsername(ctx, username)
	if err != nil {
		return false
	}
	if user != nil {
		return true
	}
	return false
}

func (s *UserService) isEmailExist(ctx context.Context, email string) bool {
	user, err := s.repository.GetUser().FindByUsername(ctx, email)
	if err != nil {
		return false
	}
	if user != nil {
		return true
	}
	return false
}
func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterRespose, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if s.isUserNameExist(ctx, req.Username) {
		return nil, errConstant.ErrUsernameExist
	}
	if s.isEmailExist(ctx, req.Email) {
		return nil, errConstant.ErrEmailExist
	}

	if req.Password != req.ConfirmPassword {
		return nil, errConstant.ErrPasswordDoesNotMatch
	}

	user, err := s.repository.GetUser().Register(ctx, &dto.RegisterRequest{
		Name:        req.Name,
		Username:    req.Username,
		Password:    string(hashedPassword),
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		RoleID:      constants.Customer,
	})
	if err != nil {
		return nil, err
	}

	response := &dto.RegisterRespose{
		User: dto.UserResponse{
			UUID:        user.UUID,
			Name:        user.Name,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		},
	}
	return response, nil
}

func (s *UserService) Update(ctx context.Context, req *dto.UpdateUserRequest, uuid string) (*dto.UserResponse, error) {
	var (
		password                  string
		checkUsername, checkEmail *models.User
		hashedPassword            []byte
		user, userResult          *models.User
		err                       error
		data                      dto.UserResponse
	)
	user, err = s.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	isUsernameExist := s.isUserNameExist(ctx, req.Username)
	if isUsernameExist && user.Username != req.Username {
		checkUsername, err = s.repository.GetUser().FindByUsername(ctx, req.Username)
		if err != nil {
			return nil, err
		}
		if checkUsername != nil {
			return nil, errConstant.ErrUsernameExist
		}
	}
	isEmailExist := s.isEmailExist(ctx, req.Email)
	if isEmailExist && user.Email != req.Email {
		checkEmail, err = s.repository.GetUser().FindByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if checkEmail != nil {
			return nil, errConstant.ErrEmailExist
		}
	}
	if req.Password != nil {
		if *req.Password != *req.ConfirmPassword {
			return nil, errConstant.ErrPasswordDoesNotMatch
		}
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		password = string(hashedPassword)
	}
	userResult, err = s.repository.GetUser().Update(ctx, &dto.UpdateUserRequest{
		Name:        req.Name,
		Username:    req.Username,
		Password:    &password,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
	}, uuid)
	if err != nil {
		return nil, err
	}

	data = dto.UserResponse{
		UUID:        userResult.UUID,
		Name:        userResult.Name,
		Username:    userResult.Username,
		Email:       userResult.Email,
		PhoneNumber: userResult.PhoneNumber,
	}
	return &data, nil
}

func (s *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	var (
		userLogin = ctx.Value(constants.UserLogin).(*dto.UserResponse)
		data      dto.UserResponse
	)
	data = dto.UserResponse{
		UUID:        userLogin.UUID,
		Name:        userLogin.Name,
		Username:    userLogin.Username,
		Email:       userLogin.Email,
		PhoneNumber: userLogin.PhoneNumber,
		Role:        userLogin.Role,
	}
	return &data, nil
}

func (s *UserService) GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	user, err := s.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}
	return data, nil
}
