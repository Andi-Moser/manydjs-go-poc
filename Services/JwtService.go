package Services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"time"
)

type authCustomClaims struct {
	Token *oauth2.Token `json:"token"`
	jwt.StandardClaims
}

type JwtService struct {
}

// NewJwtService returns a new instance of the JwtService
func NewJwtService() JwtService {
	return JwtService{}
}

// GetPayload returns the payload of the JWT token, an oauth2.Token.
func (s JwtService) GetPayload(c *gin.Context) *oauth2.Token {
	token := c.MustGet("token").(*authCustomClaims)

	return token.Token
}

// AuthorizeJWT checks whether the user has a valid JWT in his Authorization header
func (s JwtService) AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")

		if authHeader != "" {
			tokenString := authHeader[len(BEARER_SCHEMA):]

			token, err := s.validateToken(tokenString)
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			c.Set("token", token)
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

// GenerateToken generates a new JWT token with the payload (an oauth2.Token)
func (s JwtService) GenerateToken(tokenData *oauth2.Token) string {
	claims := &authCustomClaims{
		tokenData,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    "manydjs-poc",
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		panic(err)
	}
	return t
}

// validateToken checks if the given JWT token is valid
func (s JwtService) validateToken(encodedToken string) (*authCustomClaims, error) {
	token, err := jwt.ParseWithClaims(encodedToken, &authCustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*authCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
