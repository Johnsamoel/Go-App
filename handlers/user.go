package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"example.com/fintech-app/models"
	"example.com/fintech-app/pkg/users"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CreateUserHandler(Context *gin.Context) {
	// Read request body
	reqBody, err := ioutil.ReadAll(Context.Request.Body)
	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	// Parse request body into User struct
	var newUser models.User
	err = json.Unmarshal(reqBody, &newUser)
	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
		return
	}

	// Validate User struct
	validate := validator.New()
	err = validate.Struct(newUser)
	if err != nil {
		// Validation failed, return validation errors
		Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = newUser.HashPassword()

	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "something went wrong."})
		return
	}

	// Validation successful, proceed with creating the user
	result, err := users.CreateUserService(&newUser)

	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "failed to create new user"})
		return
	}

	Context.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": result})
}

func DeleteUserHandler(Context *gin.Context) {

	// Read query parameter
	userIDStr := Context.Param("id")

	// Read request parameter
	if userIDStr == "" {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
		return
	}
	// Parse userIDStr to int64
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
		return
	}

	err = users.DeleteUserService(userID)

	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "Failed deleting user"})
		return
	}

	Context.JSON(http.StatusOK, gin.H{"message": "User was deleted successfully"})
}

func UpdateUserHandler(Context *gin.Context) {
	// Read query parameter
	userIDStr := Context.Param("id")

	// Read request parameter
	if userIDStr == "" {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
		return
	}

	// Parse userIDStr to int64
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
		return
	}

	// Read request body
	reqBody, err := ioutil.ReadAll(Context.Request.Body)
	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	// Parse request body into map[string]interface{}
	var updateData map[string]interface{}
	err = json.Unmarshal(reqBody, &updateData)
	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
		return
	}

	// Call the function to update user data
	err = users.UpdateUserService(userID, updateData)
	if err != nil {
		Context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	Context.JSON(http.StatusOK, gin.H{"message": "user data updated successfully"})
}

func GetUserById(Context *gin.Context) {
	// Read query parameter
	userIDStr := Context.Param("id")

	// Read request parameter
	if userIDStr == "" {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
		return
	}

	// Parse userIDStr to int64
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
		return
	}

	user, err := users.GetUserService(userID)

	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
		return
	}

	Context.JSON(http.StatusOK, gin.H{"user": user})
}

func LogInHandler(Context *gin.Context) {
	// Read request body
	reqBody, err := ioutil.ReadAll(Context.Request.Body)
	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	// Parse request body into User struct
	var newUser models.AuthUser
	err = json.Unmarshal(reqBody, &newUser)
	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
		return
	}

	// Validate User struct
	validate := validator.New()
	err = validate.Struct(newUser)
	if err != nil {
		// Validation failed, return validation errors
		Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otp, err := users.LogInService(newUser.Email, newUser.Password, Context.Writer, Context.Request)

	if err != nil {
		fmt.Printf("error %v", err)
		Context.JSON(http.StatusBadRequest, gin.H{"error": "something went wrong."})
		return
	}

	if otp != "" {
		Context.JSON(http.StatusAccepted, gin.H{
			"message": "Logged In successfully",
			"OTP":     otp,
		})
		return
	}

	Context.JSON(http.StatusAccepted, gin.H{"message": "Logged In successfully"})
}

func LogoutHandler(c *gin.Context) {
	users.LogoutService(c.Writer)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
