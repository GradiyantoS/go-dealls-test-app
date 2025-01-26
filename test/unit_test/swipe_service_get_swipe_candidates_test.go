package unit_test

import (
	"errors"
	"testing"
	"time"

	"github.com/GradiyantoS/go-dealls-test-app/models"
	"github.com/GradiyantoS/go-dealls-test-app/services"
	userMock "github.com/GradiyantoS/go-dealls-test-app/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetSwipeCandidates(t *testing.T) {
	mockRepo := new(userMock.MockUserRepository)
	service := services.NewSwipeService(mockRepo)

	today := time.Now().Truncate(24 * time.Hour)

	testCases := []struct {
		name          string
		setupMocks    func()
		userID        int
		expectedUsers []models.User
		expectedError string
	}{
		{
			name: "Success - Opposite Gender Candidates",
			setupMocks: func() {
				mockRepo.On("GetUserByID", 1).Return(&models.User{ID: 1, Gender: "male"}, nil)

				mockRepo.On("GetAllUsers").Return([]*models.User{
					{ID: 2, Gender: "female", IsInactive: false},
					{ID: 3, Gender: "female", IsInactive: false},
					{ID: 4, Gender: "male", IsInactive: false},
				})

				mockRepo.On("GetSwipesForUser", 1).Return([]models.Swipe{})
			},
			userID: 1,
			expectedUsers: []models.User{
				{ID: 2, Gender: "female", IsInactive: false},
				{ID: 3, Gender: "female", IsInactive: false},
			},
			expectedError: "",
		},
		{
			name: "Error - User Not Found",
			setupMocks: func() {
				mockRepo.On("GetUserByID", 1).Return((*models.User)(nil), errors.New("user not found"))
			},
			userID:        1,
			expectedUsers: nil,
			expectedError: "user not found",
		},
		{
			name: "Success - Filter Swiped and Inactive Users",
			setupMocks: func() {
				mockRepo.On("GetUserByID", 1).Return(&models.User{ID: 1, Gender: "female"}, nil)

				mockRepo.On("GetAllUsers").Return([]*models.User{
					{ID: 2, Gender: "male", IsInactive: false},
					{ID: 3, Gender: "male", IsInactive: true},
					{ID: 4, Gender: "male", IsInactive: false},
					{ID: 5, Gender: "female", IsInactive: false},
					{ID: 6, Gender: "male", IsInactive: false},
				})

				mockRepo.On("GetSwipesForUser", 1).Return([]models.Swipe{
					{UserID: 1, TargetUserID: 4, CreatedAt: today.Add(1 * time.Hour)},
				})
			},
			userID: 1,
			expectedUsers: []models.User{
				{ID: 2, Gender: "male", IsInactive: false},
				{ID: 6, Gender: "male", IsInactive: false},
			},
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tc.setupMocks()

			candidates, err := service.GetSwipeCandidates(tc.userID)

			if tc.expectedError == "" {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedUsers, candidates)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
