package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func IsUserLoggedIn(c *gin.Context) {
	// Check if the user is logged in
	token, err := c.Cookie("id_token")
	if err != nil || token == "" {
		// User is not logged in, return unauthorized
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Extract user ID from token
	userID, walletID := extractUserIDFromToken(token)
	// Add user ID to request context
	c.Set("userID", userID)
	c.Set("walletID", walletID)

	// User is logged in, proceed to the next handler
	c.Next()
}

// Function to extract user ID from token
func extractUserIDFromToken(token string) (int64, int64) {
	// Split the token string to extract the claims part
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, 0
	}
	claims := parts[1]

	// Base64 decode the claims part of the token
	decodedClaims, err := base64.RawStdEncoding.DecodeString(claims)
	if err != nil {
		return 0, 0
	}

	// Parse the decoded claims as JSON to extract the user ID
	var claimsData map[string]interface{}
	err = json.Unmarshal(decodedClaims, &claimsData)
	if err != nil {
		return 0, 0
	}

	// Extract user ID from the claims data
	userIDFloat, ok := claimsData["userID"].(float64)
	if !ok {
		return 0, 0
	}

	// Extract user ID from the claims data
	walletIDFloat, ok := claimsData["walletID"].(float64)
	if !ok {
		return 0, 0
	}

	// Convert the float64 user ID to int64
	userID := int64(userIDFloat)
	walletID := int64(walletIDFloat)
	return userID, walletID
}
