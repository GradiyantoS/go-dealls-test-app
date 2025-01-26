package unit_test

import (
	"testing"
	"time"

	"github.com/GradiyantoS/go-dealls-test-app/models"
	"github.com/GradiyantoS/go-dealls-test-app/services"
	userMock "github.com/GradiyantoS/go-dealls-test-app/test/mock"
	"github.com/GradiyantoS/go-dealls-test-app/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRecordSwipe(t *testing.T) {
	mockRepo := new(userMock.MockUserRepository)
	service := services.NewSwipeService(mockRepo)

	today := time.Now().Truncate(24 * time.Hour)

	testCases := []struct {
		name          string
		setupMocks    func()
		swipe         *models.Swipe
		expectedError string
	}{
		{
			name: "Success - Regular Swipe",
			setupMocks: func() {
				mockRepo.On("GetSwipesForUser", 1).Return([]models.Swipe{})
				mockRepo.On("GetUserByID", 1).Return(&models.User{ID: 1, PremiumExpiry: nil}, nil)
				mockRepo.On("SaveSwipe", mock.AnythingOfType("*models.Swipe")).Return(nil)
			},
			swipe: &models.Swipe{UserID: 1, TargetUserID: 2},
		},
		{
			name: "Error - Already Swiped on Target User",
			setupMocks: func() {
				mockRepo.On("GetSwipesForUser", 1).Return([]models.Swipe{
					{UserID: 1, TargetUserID: 2, CreatedAt: today.Add(1 * time.Second)},
				})
			},
			swipe:         &models.Swipe{UserID: 1, TargetUserID: 2},
			expectedError: "you have already swiped on this profile today",
		},
		{
			name: "Error - Swipe Limit Reached",
			setupMocks: func() {
				mockRepo.On("GetSwipesForUser", 1).Return([]models.Swipe{
					{UserID: 1, TargetUserID: 2, CreatedAt: today.Add(1 * time.Hour)},
					{UserID: 1, TargetUserID: 3, CreatedAt: today.Add(2 * time.Hour)},
					{UserID: 1, TargetUserID: 4, CreatedAt: today.Add(3 * time.Hour)},
					{UserID: 1, TargetUserID: 5, CreatedAt: today.Add(4 * time.Hour)},
					{UserID: 1, TargetUserID: 6, CreatedAt: today.Add(5 * time.Hour)},
					{UserID: 1, TargetUserID: 7, CreatedAt: today.Add(6 * time.Hour)},
					{UserID: 1, TargetUserID: 8, CreatedAt: today.Add(7 * time.Hour)},
					{UserID: 1, TargetUserID: 9, CreatedAt: today.Add(8 * time.Hour)},
					{UserID: 1, TargetUserID: 10, CreatedAt: today.Add(9 * time.Hour)},
					{UserID: 1, TargetUserID: 11, CreatedAt: today.Add(10 * time.Hour)},
				})
				mockRepo.On("GetUserByID", 1).Return(&models.User{ID: 1, PremiumExpiry: nil}, nil)
			},
			swipe:         &models.Swipe{UserID: 1, TargetUserID: 12},
			expectedError: "daily swipe limit reached",
		},
		{
			name: "Success - Unlimited Swipes Premium User",
			setupMocks: func() {
				mockRepo.On("GetSwipesForUser", 1).Return([]models.Swipe{
					{UserID: 1, TargetUserID: 2, CreatedAt: today.Add(1 * time.Hour)},
					{UserID: 1, TargetUserID: 3, CreatedAt: today.Add(2 * time.Hour)},
					{UserID: 1, TargetUserID: 4, CreatedAt: today.Add(3 * time.Hour)},
					{UserID: 1, TargetUserID: 5, CreatedAt: today.Add(4 * time.Hour)},
					{UserID: 1, TargetUserID: 6, CreatedAt: today.Add(5 * time.Hour)},
					{UserID: 1, TargetUserID: 7, CreatedAt: today.Add(6 * time.Hour)},
					{UserID: 1, TargetUserID: 8, CreatedAt: today.Add(7 * time.Hour)},
					{UserID: 1, TargetUserID: 9, CreatedAt: today.Add(8 * time.Hour)},
					{UserID: 1, TargetUserID: 10, CreatedAt: today.Add(9 * time.Hour)},
					{UserID: 1, TargetUserID: 11, CreatedAt: today.Add(10 * time.Hour)},
				})
				mockRepo.On("GetUserByID", 1).Return(&models.User{
					ID:              1,
					PremiumExpiry:   utils.TimePtr(time.Now().UTC().Add(24 * time.Hour)),
					PremiumFeatures: models.PremiumFeatures{UnlimitedSwipes: true},
				}, nil)
				mockRepo.On("SaveSwipe", mock.AnythingOfType("*models.Swipe")).Return(nil)
			},
			swipe:         &models.Swipe{UserID: 1, TargetUserID: 12},
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mock expectations for each test case
			mockRepo.ExpectedCalls = nil
			tc.setupMocks()

			err := service.RecordSwipe(tc.swipe)

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
