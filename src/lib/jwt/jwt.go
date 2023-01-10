package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(userId int64) string {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24 * 8).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	_ = godotenv.Load(".env")
	SECRET_KEY := os.Getenv("SECRET_KEY")
	tokenText, _ := token.SignedString([]byte(SECRET_KEY))
	return tokenText
}

func ParseToken(tokenText string) interface{} {
	_ = godotenv.Load(".env")
	
	token, _ := jwt.Parse(tokenText, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		SECRET_KEY := os.Getenv("SECRET_KEY")
		return []byte(SECRET_KEY), nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)
	userId := claims["user_id"]
	return userId
}