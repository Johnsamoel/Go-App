package db

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// GenerateToken generates a JWT token with the user ID stored in the payload
func GenerateToken(userID , walletID int64, expirationDate time.Duration) (string, error) {
	// Define the token payload
	claims := jwt.MapClaims{
		"walletID": walletID,						// Store the user ID in the token payload
		"userID": userID,                           // Store the user ID in the token payload
		"exp":    time.Now().Add(expirationDate).Unix(), // Token expiration time (1 hour from now)
	}

	// Create the token with the payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
