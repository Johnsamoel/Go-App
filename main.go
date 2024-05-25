package main

import (
	"example.com/fintech-app/db"
	"example.com/fintech-app/handlers"
	"example.com/fintech-app/middlewares"
	"example.com/fintech-app/pkg/users"
	"example.com/fintech-app/pkg/wallets"
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

	// Initialize user & wallet repositories
	userRepo := users.NewUserRepo(db.DB)
	walletRepo := wallets.NewWalletRepo(db.DB)

	// Initialize the user & wallet services
	userService := users.NewUserService(userRepo)
	walletService := wallets.NewWalletService(walletRepo)

	// Grouping authentication routes
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/login", handlers.LogInHandler(userService))
		authRoutes.POST("/logout", handlers.LogoutHandler(userService))
		authRoutes.POST("/register", handlers.CreateUserHandler(userService))
	}

	// Check user authentications
	router.Use(middlewares.IsUserLoggedIn)

	// Grouping user routes
	userRoutes := router.Group("/api/user")
	{
		userRoutes.POST("/create")
		userRoutes.DELETE("/delete/:id", handlers.DeleteUserHandler(userService))
		userRoutes.PATCH("/update/:id", handlers.UpdateUserHandler(userService))
		userRoutes.GET("/:id", handlers.GetUserByIdHandler(userService))
	}

	// Grouping wallet routes
	walletRoutes := router.Group("/api/wallet")
	{
		walletRoutes.POST("/activate", handlers.ActivateWalletHandler(walletService))
		walletRoutes.POST("/refund", handlers.RefundWalletHandler(walletService))
		walletRoutes.POST("/withdraw", handlers.WithdrawWalletHandler(walletService))
	}

	router.Run(":8080")
}
