package service

import (
	"Enterprise/data"
	"Enterprise/helpers"
	"Enterprise/mail"
	"Enterprise/model"
	"Enterprise/prisma/db"
	"Enterprise/repository"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserService struct {
	Db          *db.PrismaClient
	RedisClient *redis.Client
}

func NewUserService(db *db.PrismaClient, redisClient *redis.Client) *UserService {
	return &UserService{
		Db:          db,
		RedisClient: redisClient,
	}
}

func (p *UserService) RoleCreation(ctx context.Context, roleModel *model.RoleCreationModel) *data.WebResponse {
	validator := helpers.RequestValidators(roleModel)
	if validator != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validator.Error(),
		}
	}
	roleName := strings.ToUpper(roleModel.Name)
	permissionsJson, err := json.Marshal(roleModel.Permissions)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: "Invalid permissions format",
			Data:    nil,
		}
	}

	existingRoleByName, _ := p.Db.Role.FindFirst(db.Role.Name.Equals(roleName)).Exec(ctx)
	if existingRoleByName != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Role already exists",
			Data:    nil,
		}
	}

	_, err = p.Db.Role.CreateOne(
		db.Role.Name.Set(roleName),
		db.Role.Permissions.Set(permissionsJson),
	).Exec(ctx)

	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	err = repository.AuditLogs(ctx, p.Db, roleModel.AuditId, "Role created", "This action was performed by")
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	return &data.WebResponse{
		Code:    http.StatusCreated,
		Message: "Role created",
		Data:    nil,
	}
}

func (p *UserService) CreateUserByAdmin(ctx context.Context, userModel *model.UserCreationModel) *data.WebResponse {
	validator := helpers.RequestValidators(userModel)
	if validator != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validator.Error(),
		}
	}

	existingUser, _ := p.Db.User.FindUnique(db.User.Email.Equals(userModel.Email)).Exec(ctx)
	if existingUser != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "User already exists",
			Data:    nil,
		}
	}

	user, err := p.Db.User.CreateOne(
		db.User.Email.Set(userModel.Email),
		db.User.FirstName.Set(userModel.Name),
		db.User.Role.Link(db.Role.ID.Equals(userModel.RoleId)),
	).Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	err = repository.AuditLogs(ctx, p.Db, userModel.UserId, "User created", "This action was performed by")
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	var emailToken = uuid.New().String()
	redisContext := context.Background()
	p.RedisClient.Set(redisContext, emailToken, user.ID, time.Hour*3)

	mailInputs := &data.MailInputs{
		Email:    user.Email,
		Code:     emailToken,
		Username: user.FirstName,
	}

	err = mail.ResetPassword(mailInputs)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	return &data.WebResponse{
		Code:    http.StatusCreated,
		Message: "User created",
		Data:    nil,
	}
}

func (p *UserService) CreateUserPassword(ctx context.Context, userDto *model.UserPasswordCreationModel) *data.WebResponse {
	validator := helpers.RequestValidators(userDto)
	if validator != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validator.Error(),
		}
	}

	if userDto.Password != userDto.ConfirmPassword {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Passwords do not match",
			Data:    nil,
		}
	}

	userId := p.RedisClient.Get(ctx, userDto.Code).Val()
	if userId == "" {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Code invalid!",
			Data:    nil,
		}
	}

	id, _ := strconv.ParseInt(userId, 10, 0)
	user, _ := p.Db.User.FindUnique(db.User.ID.Equals(int(id))).Exec(ctx)
	if user == nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Data:    nil,
		}
	}

	if user.State == db.StateEnumVerified {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "User password already set!",
			Data:    nil,
		}
	}

	hashPassword := helpers.HashPassword(userDto.Password)

	_, _ = p.Db.User.FindUnique(db.User.ID.Equals(int(id))).Update(
		db.User.Password.Set(hashPassword),
		db.User.State.Set(db.StateEnumVerified),
	).Exec(ctx)

	p.RedisClient.Del(ctx, userDto.Code)

	return &data.WebResponse{
		Code:    http.StatusCreated,
		Message: "Password updated!",
		Data:    nil,
	}
}

func (p *UserService) UpdateUserInfo(ctx context.Context, userDto *model.UserCreationModel) *data.WebResponse {
	validator := helpers.RequestValidators(userDto)
	if validator != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validator.Error(),
		}
	}

	existingUser, _ := p.Db.User.FindUnique(db.User.Email.Equals(userDto.Email)).Exec(ctx)
	if existingUser != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "User email already in use",
			Data:    nil,
		}
	}

	_, err := p.Db.User.UpsertOne(db.User.ID.Equals(userDto.UserId)).Update(
		db.User.Email.Set(userDto.Email),
		db.User.Role.Link(db.Role.ID.Equals(userDto.RoleId)),
		db.User.FirstName.Set(userDto.Name),
	).Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}

	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "User updated",
		Data:    nil,
	}
}

func (p *UserService) DeactivateUser(ctx context.Context, userId int) *data.WebResponse {
	_, err := p.Db.User.FindUnique(db.User.ID.Equals(userId)).Update(db.User.State.Set(db.StateEnumDisabled)).Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "User deactivated",
		Data:    nil,
	}
}

func (p *UserService) DeleteUser(ctx context.Context, userId int) *data.WebResponse {
	_, err := p.Db.User.FindUnique(db.User.ID.Equals(userId)).Update(db.User.State.Set(db.StateEnumDeleted)).Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "User deleted!",
		Data:    nil,
	}
}

func (p *UserService) Login(ctx context.Context, userDto *model.LoginUserModel) *data.WebResponse {
	validator := helpers.RequestValidators(userDto)
	if validator != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validator.Error(),
		}
	}

	existingUser, err := p.Db.User.FindFirst(db.User.Email.Equals(userDto.Email)).With(db.User.Role.Fetch()).Exec(ctx)

	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}
	if existingUser == nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Data:    nil,
		}
	}
	if existingUser.State == db.StateEnumFresh {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Account not activated!",
			Data:    nil,
		}
	}
	if existingUser.State == db.StateEnumDisabled {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "User account deleted!",
			Data:    nil,
		}
	}
	if existingUser.State == db.StateEnumDeleted {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "User account deleted!",
			Data:    nil,
		}
	}
	userPassword, _ := existingUser.Password()

	correctPassword := helpers.CheckPasswordHash(userDto.Password, userPassword)
	if !correctPassword {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Password do not match",
			Data:    nil,
		}
	}

	roleName := existingUser.Role().Name

	jwtPayload := &model.JWTPayload{
		Email: existingUser.Email,
		Id:    existingUser.ID,
		Role:  roleName,
	}

	accessToken, refreshToken := helpers.GenerateAuthToken(jwtPayload, userDto.RememberMe)

	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "User login!",
		Data: struct {
			Id           int    `json:"id"`
			Email        string `json:"email"`
			Role         string `json:"role"`
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}{
			Id:           existingUser.ID,
			Email:        existingUser.Email,
			Role:         roleName,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
}

func (p *UserService) ChangeUserInfo(ctx context.Context, userDto *model.UpdateUserInfoModel) *data.WebResponse {
	validator := helpers.RequestValidators(userDto)
	if validator != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    validator.Error(),
		}
	}
	existingUser, err := p.Db.User.FindUnique(db.User.ID.Equals(userDto.UserId)).With(db.User.Role.Fetch()).Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}
	if existingUser == nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Data:    nil,
		}
	}
	_, _ = p.Db.User.FindUnique(db.User.ID.Equals(userDto.UserId)).Update(
		db.User.FirstName.Set(userDto.FirstName),
		db.User.LastName.Set(userDto.LastName),
	).Exec(ctx)
	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "Data updated",
		Data:    nil,
	}
}

func (p *UserService) GetAllUsers(ctx context.Context) *data.WebResponse {
	users, err := p.Db.User.FindMany().Select(
		db.User.Email.Field(),
		db.User.ID.Field(),
		db.User.FirstName.Field(),
		db.User.LastName.Field(),
	).With(
		db.User.Role.Fetch(),
	).Exec(ctx)
	if err != nil {
		return &data.WebResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}

	var UserResponses []model.UserResponse
	for _, user := range users {
		lastName, _ := user.LastName()
		UserResponses = append(UserResponses, model.UserResponse{
			Id:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  lastName,
			CreatedAt: user.CreatedAt,
			Role: struct {
				Id          int    `json:"id"`
				Name        string `json:"name"`
				Permissions string `json:"permissions"`
			}{
				Id:          user.Role().ID,
				Name:        user.Role().Name,
				Permissions: string(user.Role().Permissions),
			},
		})
	}

	return &data.WebResponse{
		Code:    http.StatusOK,
		Message: "Users",
		Data:    UserResponses,
	}
}
