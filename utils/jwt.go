// utils/jwt.go
package utils

import (
    "fmt"
    "time"

    "github.com/dgrijalva/jwt-go"
)

type Claims struct {
    UserID string `json:"user_id"`
    jwt.StandardClaims
}

// GenerateJWT generates a JWT token for authenticated users.
func GenerateJWT(userID string, secret string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
            IssuedAt:  time.Now().Unix(),
            Issuer:    "maliaki-backend",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

// ValidateJWT validates a JWT token and returns the user ID.
func ValidateJWT(tokenStr string, secret string) (string, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        // Ensure the token method is HMAC
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return []byte(secret), nil
    })

    if err != nil {
        return "", err
    }

    if !token.Valid {
        return "", fmt.Errorf("invalid token")
    }

    return claims.UserID, nil
}
