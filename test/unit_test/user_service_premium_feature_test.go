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
		duration      int
		features      []string
		expectedError string
	}{
		{
			name: "Success - Enable Unlimited Swipes and Verified Badge",
			setupMocks: func() {
				mockRepo.On("GetUserByID", 1).Return(&models.User{
					ID:            1,
					IsInactive:    false,
					PremiumExpiry: nil,
					PremiumFeatures: models.PremiumFeatures{
						UnlimitedSwipes: false,
						IsVerified:      false,
					},
				}, nil)

				mockRepo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(nil)
			},
			userID:        1,
			duration:      30,
			features:      []string{"UnlimitedSwipes", "IsVerified"},
			expectedError: "",
		},
		{
			name: "Error - User Not Found",
			setupMocks: func() {
				mockRepo.On("GetUserByID", 1).Return(nil, errors.New("user not found"))
			},
			userID:        1,
			duration:      30,
			features:      []string{"UnlimitedSwipes"},
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
			duration:      30,
			features:      []string{"UnlimitedSwipes"},
			expectedError: "user account is inactive",
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
			duration:      30,
			features:      []string{"InvalidFeature"},
			expectedError: "invalid premium feature: InvalidFeature",
		},
		{
			name: "Error - Unlimited Swipes Already Active",
			setupMocks: func() {
				mockRepo.On("GetUserByID", 1).Return(&models.User{
					ID:         1,
					IsInactive: false,
					PremiumFeatures: models.PremiumFeatures{
						UnlimitedSwipes: true,
						IsVerified:      false,
					},
				}, nil)
			},
			userID:        1,
			duration:      30,
			features:      []string{"UnlimitedSwipes"},
			expectedError: "unlimited swipes is already active",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mock calls
			mockRepo.ExpectedCalls = nil
			tc.setupMocks()

			// Call EnablePremiumFeature
			err := service.EnablePremiumFeature(tc.userID, tc.duration, tc.features)

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
