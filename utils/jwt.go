package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims personalizados (datos dentro del token)
type Claims struct {
	UserID int    `json:"user_id"`
	Rol    string `json:"rol"`
	jwt.RegisteredClaims
}

func GenerarToken(userID int, rol string) (string, error) {
	// Leer secret y expiración del .env
	secret := os.Getenv("JWT_SECRET")
	expirationStr := os.Getenv("JWT_EXPIRATION")

	expirationHours, err := strconv.Atoi(expirationStr)
	if err != nil {
		expirationHours = 24 // Default 24 horas
	}

	//Crear Claims
	claims := Claims{
		UserID: userID,
		Rol:    rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expirationHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	//crear token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Firmar token con el secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

// ValidarToken verifica si un token es válido
func ValidarToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")

	//parsear el token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el método de firma sea el correcto
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inválido")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extraer claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}
