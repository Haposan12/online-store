package pkg

import (
	"errors"
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/i18n"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GetLangVersion sets site language version.
func GetLangVersion(ctx *beegoContext.Context) string {
	// 1. Check URL arguments.
	lang := ctx.Input.Query("lang")

	// Check again in case someone modifies on purpose.
	if !i18n.IsExist(lang) {
		lang = ""
	}

	// 2. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := ctx.Request.Header.Get("Accept-Language")
		if i18n.IsExist(al) {
			lang = al
		}
	}

	// 3. Default language is Indonesia.
	if len(lang) == 0 {
		lang = "id"
	}

	// Set language properties.
	return lang
}

func GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["userID"].(float64))
		return userID, nil
	} else {
		return 0, err
	}
}
