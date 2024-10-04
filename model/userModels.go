package model

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserCreationModel struct {
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	RoleId int    `json:"roleId" validate:"required"`
	UserId int    `json:"userId"`
}

type RoleCreationModel struct {
	Name        string                 `json:"name" validate:"required"`
	Permissions map[string]interface{} `json:"permissions" validate:"required"`
	AuditId     int                    `json:"auditId"`
}

type UserPasswordCreationModel struct {
	Code            string `json:"code"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type UpdateUserInfoModel struct {
	FirstName string `json:"firstName" validate:"required,min=5,max=32"`
	LastName  string `json:"lastName" validate:"required,min=5,max=32"`
	UserId    int    `json:"userId"`
}

type LoginUserModel struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"rememberMe" validate:"required"`
}

type JWTPayload struct {
	Email string `json:"email"`
	Id    int    `json:"id"`
	Role  string `json:"role"`
}

type JWTClaim struct {
	Id         int         `json:"id"`
	Role       string      `json:"role"`
	Email      string      `json:"email"`
	JwtId      string      `json:"jwtid"`
	CustomData interface{} `json:"customData"`
	jwt.StandardClaims
}

type AuditLogResponse struct {
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
	User      struct {
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
	} `json:"user"`
}

type UserResponse struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	Role      struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Permissions string `json:"permissions"`
	}
}
