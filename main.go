package main

import (
	"example.com/fintech-app/db"
	"example.com/fintech-app/handlers"
	"example.com/fintech-app/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {

	var err error

	db.DB, err = db.OpenDbConnection()

	if err != nil {
		panic("server is down")
	}

	defer db.DB.Close()

	router := gin.Default()

	// Grouping authentication routes
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/login", handlers.LogInHandler)
		authRoutes.POST("/logout", handlers.LogoutHandler)
	}

	// Check user authentications
	router.Use(middlewares.IsUserLoggedIn)

	// Grouping user routes
	userRoutes := router.Group("/api/user")
	{
		userRoutes.POST("/create", handlers.CreateUserHandler)
		userRoutes.DELETE("/delete/:id", handlers.DeleteUserHandler)
		userRoutes.PATCH("/update/:id", handlers.UpdateUserHandler)
		userRoutes.GET("/:id", handlers.GetUserById)
	}

	// Grouping wallet routes
	walletRoutes := router.Group("/api/wallet")
	{
		walletRoutes.POST("/activate", handlers.ActivateWalletHandler)
		walletRoutes.POST("/refund", handlers.RefundWalletHandler)
		walletRoutes.POST("/withdraw", handlers.WithdrawWalletHandler)
	}

	router.Run(":8080")
}
