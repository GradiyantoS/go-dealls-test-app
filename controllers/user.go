package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/GradiyantoS/go-dealls-test-app/middlewares"
	"github.com/GradiyantoS/go-dealls-test-app/models"
	"github.com/GradiyantoS/go-dealls-test-app/services"
	"github.com/GradiyantoS/go-dealls-test-app/utils"
)

type UserController interface {
	PurchasePremium(w http.ResponseWriter, r *http.Request)
	SwipeCandidates(w http.ResponseWriter, r *http.Request)
	SwipeHandler(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	userService  services.UserService
	swipeService services.SwipeService
}

func NewUserController(userService services.UserService, swipeService services.SwipeService) UserController {
	return &userController{userService, swipeService}
}

func (c *userController) PurchasePremium(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Failed to retrieve user ID")
		return
	}

	var input struct {
		Duration int      `json:"duration"` // Duration in days
		Features []string `json:"features"` // List of features to enable
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Duration <= 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid duration")
		return
	}

	err := c.userService.EnablePremiumFeature(userID, input.Duration, input.Features)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.DataSuccessResponse(w, http.StatusOK, map[string]string{"message": "Premium features enabled successfully"})
}

func (c *userController) SwipeCandidates(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Failed to retrieve user ID")
		return
	}

	candidates, err := c.swipeService.GetSwipeCandidates(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve swipe candidates")
		return
	}

	utils.DataSuccessResponse(w, http.StatusOK, candidates)
}

func (c *userController) SwipeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Failed to retrieve user ID")
		return
	}

	var swipe models.Swipe
	if err := json.NewDecoder(r.Body).Decode(&swipe); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	swipe.UserID = userID
	err := c.swipeService.RecordSwipe(&swipe)
	if err != nil {
		utils.ErrorResponse(w, http.StatusForbidden, err.Error())
		return
	}

	utils.DataSuccessResponse(w, http.StatusOK, map[string]string{"message": "Swipe recorded successfully"})
}
