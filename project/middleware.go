package magic

import (
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Middleware structure
type Middleware struct {
	run func(*Context) error
}

// NewMiddleware function
// Create new Middleware
func NewMiddleware(handler func(context *Context) error) Middleware {
	return Middleware{
		run: handler,
	}
}

// NewJWTMiddleware function
// Create new JWT authorization middleware
// Claims will contains in context.Storage["claims"]
// contex.Storage["claims"] -> map[string]interface{}
func NewJWTMiddleware(secretKey, headerName string) Middleware {
	return NewMiddleware(func(context *Context) error {
		tokenWithBearer, err := context.Headers.ParseString(headerName)
		if err != nil {
			return context.SendErrorString("miss or invalid JWT")
		}

		mas := strings.Split(tokenWithBearer, " ")
		if len(mas) != 2 {
			return context.SendErrorString("miss or invalid JWT")
		}

		tokenStr := mas[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			tokenTime := (token.Claims.(jwt.MapClaims))["exp"].(float64)
			if time.Now().UnixNano()-int64(tokenTime) > 0 {
				return nil, context.SendErrorString("miss or invalid JWT")
			}
			return []byte(secretKey), nil
		})

		if err != nil || token.Valid == false {
			return context.SendErrorString("miss or invalid JWT")
		}

		claims := token.Claims.(jwt.MapClaims)
		context.Storage["claims"] = map[string]interface{}(claims)

		return nil
	})
}

// GenerateJWTToken function
// Use your claims, your method, secretKey and time
// Generate token and return it
func GenerateJWTToken(claims map[string]interface{}, method *jwt.SigningMethodHMAC, secretKey string, duration time.Duration) (string, error) {
	claims["exp"] = time.Now().UnixNano() + duration.Nanoseconds()
	token := jwt.NewWithClaims(method, jwt.MapClaims(claims))
	return token.SignedString([]byte(secretKey))
}
