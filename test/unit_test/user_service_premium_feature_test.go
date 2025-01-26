package unit_test

import (
	"errors"
	"testing"

	"github.com/GradiyantoS/go-dealls-test-app/models"
	"github.com/GradiyantoS/go-dealls-test-app/services"
	userMock "github.com/GradiyantoS/go-dealls-test-app/test/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEnablePremiumFeature(t *testing.T) {
	mockRepo := new(userMock.MockUserRepository)
	service := services.NewUserService(mockRepo)

	testCases := []struct {
		name          string
		setupMocks    func()
		userID        int
		feature       string
		expectedError string
	}{
		{
			name: "Success - Unlimited Swipes",
			setupMocks: func() {
				mockRepo.On("GetUserByID", 1).Return(&models.User{
					ID:         1,
					IsInactive: false,
					PremiumFeatures: struct {
						UnlimitedSwipes bool `json:"unlimited_swipes"`
						IsVerified      bool `json:"profile_boost"`
					}{
						UnlimitedSwipes: false,
					},
				}, nil)

				// Correctly mock the UpdateUser method
				mockRepo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(nil)
			},
			userID:        1,
			feature:       "remove_swipe_limit",
			expectedError: "",
		},
		{
			name: "Error - User Not Found",
			setupMocks: func() {
				mockRepo.On("GetUserByID", 1).Return(nil, errors.New("user not found"))
			},
			userID:        1,
			feature:       "remove_swipe_limit",
			expectedError: "user not found",
		},
		{
			name: "Error - User Is Inactive",
			setupMocks: func() {
				mockRepo.On("GetUserByID", 1).Return(&models.User{
					ID:         1,
					IsInactive: true,
				}, nil)
			},
			userID:        1,
			feature:       "remove_swipe_limit",
			expectedError: "user account is inactive",
		},
		{
			name: "Error - Feature Already Active",
			setupMocks: func() {
				mockRepo.On("GetUserByID", 1).Return(&models.User{
					ID:         1,
					IsInactive: false,
					PremiumFeatures: struct {
						UnlimitedSwipes bool `json:"unlimited_swipes"`
						IsVerified      bool `json:"profile_boost"`
					}{
						UnlimitedSwipes: true,
					},
				}, nil)
			},
			userID:        1,
			feature:       "remove_swipe_limit",
			expectedError: "unlimited swipes is already active",
		},
		{
			name: "Error - Invalid Feature",
			setupMocks: func() {
				mockRepo.On("GetUserByID", 1).Return(&models.User{
					ID:         1,
					IsInactive: false,
				}, nil)
			},
			userID:        1,
			feature:       "invalid_feature",
			expectedError: "invalid premium feature",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mock calls
			mockRepo.ExpectedCalls = nil
			tc.setupMocks()

			// Call EnablePremiumFeature
			err := service.EnablePremiumFeature(tc.userID, tc.feature)

			// Assertions
			if tc.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
