package lib

import (
	"app/domain/value/user"
	"fmt"
	"time"

	//"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// GenerateToken トークン生成
func GenerateToken(userId user.UserId, secretKey string, now time.Time, lifetime int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId.Value().String(),
		"exp":     now.Add(time.Hour * time.Duration(lifetime)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken トークンパース
func ParseToken(tokenString string, secretKey string) (user.UserId, int64, error) {
	userId := user.NewUserId(uuid.Nil)

	token, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return userId, 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return userId, 0, err
	}

	uid, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return userId, 0, err
	}

	userId = user.NewUserId(uid)

	return userId, int64(claims["exp"].(float64)), nil
}
