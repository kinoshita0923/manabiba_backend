package jwt

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/golang-jwt/jwt/v4"
)

func GetTokenText(userId int64) string {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	_ = godotenv.Load(".env")
	SECRET_KEY := os.Getenv("SECRET_KEY")
	tokenText, _ := token.SignedString([]byte(SECRET_KEY))
	return tokenText
}