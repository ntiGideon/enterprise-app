package middleware

import (
	"Enterprise/data"
	"Enterprise/helpers"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
)

func RoleBasedAuthMiddleware(allowedRoles []string, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		bearerToken := r.Header.Get("Authorization")

		if bearerToken == "" {
			helpers.WriteResponseBody(w, &data.WebResponse{
				Code:    http.StatusUnauthorized,
				Message: "Authorization token is missing!",
				Data:    nil,
			}, http.StatusUnauthorized)
			return
		}

		claims, err := helpers.ValidateToken(bearerToken)

		if err != nil {
			helpers.WriteResponseBody(w, &data.WebResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
				Data:    nil,
			}, http.StatusUnauthorized)
			return
		}
		userRole := claims.Role
		if !isAllowedRole(userRole, allowedRoles) {
			helpers.WriteResponseBody(w, &data.WebResponse{
				Code:    http.StatusForbidden,
				Message: "Access Denied: Insufficient permissions!",
				Data:    nil,
			}, http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims.Id)
		next(w, r.WithContext(ctx), params)
	}
}

func isAllowedRole(userRole string, allowedRoles []string) bool {
	for _, role := range allowedRoles {
		if role == userRole {
			return true
		}
	}
	return false
}
