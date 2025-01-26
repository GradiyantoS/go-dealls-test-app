package routes

import (
	"github.com/GradiyantoS/go-dealls-test-app/controllers"
	"github.com/GradiyantoS/go-dealls-test-app/middlewares"
	"github.com/GradiyantoS/go-dealls-test-app/repositories"
	"github.com/GradiyantoS/go-dealls-test-app/services"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	userRepo := repositories.NewUserRepository()

	userService := services.NewUserService(userRepo)
	swipeService := services.NewSwipeService(userRepo)

	authController := controllers.NewAuthController(userService)
	userController := controllers.NewUserController(userService, swipeService)

	// Create a new router
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/signup", authController.SignUp).Methods("POST")
	router.HandleFunc("/login", authController.Login).Methods("POST")

	// Protected routes (requires JWT authentication)
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middlewares.AuthMiddleware)

	protected.HandleFunc("/purchase-premium", userController.PurchasePremium).Methods("POST")
	protected.HandleFunc("/swipe", userController.SwipeHandler).Methods("POST")
	protected.HandleFunc("/candidates", userController.SwipeCandidates).Methods("GET")

	return router
}

func SetupRouterWithRepo(userRepo repositories.UserRepository) *mux.Router {

	userService := services.NewUserService(userRepo)
	swipeService := services.NewSwipeService(userRepo)
	authController := controllers.NewAuthController(userService)
	userController := controllers.NewUserController(userService, swipeService)

	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/signup", authController.SignUp).Methods("POST")
	router.HandleFunc("/login", authController.Login).Methods("POST")

	// Protected routes
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middlewares.AuthMiddleware)

	protected.HandleFunc("/purchase-premium", userController.PurchasePremium).Methods("POST")
	protected.HandleFunc("/swipe", userController.SwipeHandler).Methods("POST")
	protected.HandleFunc("/candidates", userController.SwipeCandidates).Methods("GET")

	return router
}
