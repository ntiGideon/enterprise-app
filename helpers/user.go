package helpers

import (
	"Enterprise/model"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateAuthToken(jwtPayload *model.JWTPayload, rememberMe bool) (string, string) {
	refreshToken, jwtId := GenerateRefreshToken(jwtPayload, rememberMe)
	accessToken := GenerateAccessToken(jwtPayload, jwtId)
	return accessToken, refreshToken
}

func GenerateAccessToken(jwtPayload *model.JWTPayload, jwtId string) string {
	expirationDate := time.Now().Add(1 * time.Hour)
	subject := uuid.New()
	claims := &model.JWTClaim{
		Role:       jwtPayload.Role,
		Email:      jwtPayload.Email,
		Id:         jwtPayload.Id,
		JwtId:      jwtId,
		CustomData: map[string]interface{}{"subject": subject.String()},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationDate.Unix(),
			Audience:  os.Getenv("FRONTEND_URL"),
			Issuer:    os.Getenv("BACKEND_URL"),
			Subject:   subject.String(),
			Id:        strconv.Itoa(jwtPayload.Id),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}

func GenerateRefreshToken(jwtPayload *model.JWTPayload, rememberMe bool) (string, string) {
	expirationDate := time.Now().Add(1 * time.Hour)
	subject := uuid.New()
	jwtId := uuid.New()
	if rememberMe {
		expirationDate = time.Now().Add(24 * time.Hour)
	}
	claims := &model.JWTClaim{
		Role:       jwtPayload.Role,
		Email:      jwtPayload.Email,
		Id:         jwtPayload.Id,
		JwtId:      jwtId.String(),
		CustomData: map[string]interface{}{"subject": subject.String()},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationDate.Unix(),
			Audience:  os.Getenv("FRONTEND_URL"),
			Issuer:    os.Getenv("BACKEND_URL"),
			Subject:   subject.String(),
			Id:        strconv.Itoa(jwtPayload.Id),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString, jwtId.String()
}

func ValidateToken(signedToken string) (*model.JWTClaim, error) {
	tokenString, err := jwt.ParseWithClaims(
		signedToken,
		&model.JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	if err != nil {
		return nil, err
	}
	claims, ok := tokenString.Claims.(*model.JWTClaim)
	if !ok || tokenString.Valid == false {
		return nil, err
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}
	return claims, nil
}
