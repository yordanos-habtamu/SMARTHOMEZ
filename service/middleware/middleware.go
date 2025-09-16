package middleware


import (
	"fmt"
	"net/http"
	"strings"
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/yordanos-habtamu/realstate/config"
	"github.com/yordanos-habtamu/realstate/utils"
)

var (
	// Secret key for sessions, should be stored in an environment variable for production.
	secretKey = []byte(config.Envs.JWT_SECRET)
	store     = sessions.NewCookieStore(secretKey)
)

type CustomClaims struct {
	Id   int    `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// GetJWTFromSession retrieves the JWT token from the session
func GetJWTFromSession(r *http.Request) (string, error) {
	session, err := store.Get(r, "realstate")
	if err != nil {
		return "", err
	}

	// Retrieve the token from the session
	tokenString, ok := session.Values["token"].(string)
	if !ok {
		return "", nil
	}
	return tokenString, nil
}

// SetJWTInSession stores the JWT token in the session
func SetJWTInSession(w http.ResponseWriter, r *http.Request, tokenString string) error {
	session, err := store.Get(r, "realstate")
	if err != nil {
		return err
	}

	// Store the JWT token in the session
	session.Values["token"] = tokenString
	session.Save(r, w)
	return nil
}

// JwtMiddleware is a middleware to check for valid JWT tokens and enforce RBAC
func JwtMiddleware(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteError(w,  http.StatusUnauthorized,fmt.Errorf("Authorization header is missing"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 {
				utils.WriteError(w,  http.StatusUnauthorized,fmt.Errorf("Invalid authorization header format"))
				return
			}

			tokenString := parts[1]
			token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})

			if err != nil || !token.Valid {
				utils.WriteError(w, http.StatusUnauthorized,fmt.Errorf("Invalid token"))
				return
			}

			claims, ok := token.Claims.(*CustomClaims)
			if !ok {
				utils.WriteError(w,  http.StatusUnauthorized,fmt.Errorf("Invalid token claims"))
				return
			}

			// Check if the user's role is in the allowed roles
			roleAllowed := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				utils.WriteError(w, http.StatusForbidden, fmt.Errorf("Forbidden: Insufficient permissions"))
				return
			}

			// Add user info to the request context
			ctx := context.WithValue(r.Context(), "user", claims)
			next(w, r.WithContext(ctx))
		}
	}
}
