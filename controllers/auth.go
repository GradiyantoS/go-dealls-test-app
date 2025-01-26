package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/GradiyantoS/go-dealls-test-app/models"
	"github.com/GradiyantoS/go-dealls-test-app/services"
	"github.com/GradiyantoS/go-dealls-test-app/utils"
)

type AuthController interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

// authController is the concrete implementation of AuthController
type authController struct {
	userService services.UserService
}

// NewAuthController creates a new AuthController
func NewAuthController(userService services.UserService) AuthController {
	return &authController{userService}
}

// SignUpHandler handles user registration
func (c *authController) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	err := c.userService.SignUp(&user)
	if err != nil {
		utils.ErrorResponse(w, http.StatusConflict, err.Error())
		return
	}

	utils.DataSuccessResponse(w, http.StatusCreated, map[string]string{"message": "user has been added"})
}

// LoginHandler handles user login
func (c *authController) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	token, err := c.userService.Login(creds)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	response := map[string]string{
		"message": "login success",
		"token":   token,
	}

	utils.DataSuccessResponse(w, http.StatusCreated, response)
}
