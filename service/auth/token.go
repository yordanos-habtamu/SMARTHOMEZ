package auth

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yordanos-habtamu/realstate/config"
	"github.com/yordanos-habtamu/realstate/types"
)

func CreateJWT(secret []byte,userID int , role string ) (string,error){
 expiration  := time.Second * time.Duration (config.Envs.JWTExpiration)

 token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
	"userID" : strconv.Itoa(userID),
    "role":role,
	"expiredAt":time.Now().Add(expiration).Unix(),
 })
 tokenString,err := token.SignedString(secret)
 if err != nil {
	return "" ,err
 }
 return tokenString,nil
}

func WithJwtAuth(handlerFunc http.HandlerFunc,store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		secret := []byte(config.Envs.JWT_SECRET)
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w,"Authorization header is required",http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
			return secret,nil
		})
		if err != nil {
			http.Error(w,"Invalid token",http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w,"Invalid token",http.StatusUnauthorized)
			return
		}
		userID, err := strconv.Atoi(claims["userID"].(string))
		if err != nil {
			http.Error(w,"Invalid token",http.StatusUnauthorized)
			return
		}
		user,err := store.GetUserById(userID)
		if err != nil {
			http.Error(w,"Invalid token",http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx,"user",user)
		r = r.WithContext(ctx)
		handlerFunc(w,r)
	}

}

func GetUserIdfromContext(ctx context.Context) uint {
	user := ctx.Value("user").(*types.User)
	return user.ID
}