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

func TestSignUp(t *testing.T) {
	mockRepo := new(userMock.MockUserRepository)
	service := services.NewUserService(mockRepo)

	testCases := []struct {
		name          string
		setupMocks    func()
		input         models.User
		expectedError string
	}{
		{
			name: "Success",
			setupMocks: func() {
				mockRepo.On("GetUserByEmail", "newuser@example.com").Return((*models.User)(nil), errors.New("user not found"))
				mockRepo.On("GetUserByPhone", "1234567890").Return((*models.User)(nil), errors.New("user not found"))
				mockRepo.On("GenerateUserID").Return(1)
				mockRepo.On("SaveUser", mock.AnythingOfType("*models.User")).Return(nil)
			},
			input: models.User{
				Email:    "newuser@example.com",
				Password: "password123",
				Phone:    "1234567890",
				Name:     "New User",
				Gender:   "male",
			},
			expectedError: "",
		},
		{
			name: "Error - Email Already Exists",
			setupMocks: func() {
				mockRepo.On("GetUserByEmail", "existing@example.com").Return(&models.User{Email: "existing@example.com"}, nil)
			},
			input: models.User{
				Email: "existing@example.com",
			},
			expectedError: "email already exists",
		},
		{
			name: "Error - Phone Number Already Exists",
			setupMocks: func() {
				mockRepo.On("GetUserByEmail", "newuser@example.com").Return((*models.User)(nil), errors.New("user not found"))
				mockRepo.On("GetUserByPhone", "1234567890").Return(&models.User{Phone: "1234567890"}, nil)
			},
			input: models.User{
				Email:    "newuser@example.com",
				Password: "password123",
				Phone:    "1234567890",
				Name:     "New User",
				Gender:   "male",
			},
			expectedError: "phone number already exists",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil // Reset mocks
			tc.setupMocks()

			err := service.SignUp(&tc.input)

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
