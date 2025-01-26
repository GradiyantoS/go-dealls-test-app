package unit_test

import (
	"errors"
	"testing"

	"github.com/GradiyantoS/go-dealls-test-app/models"
	"github.com/GradiyantoS/go-dealls-test-app/services"
	userMock "github.com/GradiyantoS/go-dealls-test-app/test/mock"
	"github.com/GradiyantoS/go-dealls-test-app/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	mockRepo := new(userMock.MockUserRepository)
	service := services.NewUserService(mockRepo)

	// Mock GenerateJWT to return a static token
	originalGenerateJWT := utils.GenerateJWT
	utils.GenerateJWT = func(userID int) (string, error) {
		return "mocked-jwt-token", nil
	}
	defer func() { utils.GenerateJWT = originalGenerateJWT }() // Restore original function after the test

	testCases := []struct {
		name          string
		setupMocks    func()
		creds         models.Credentials
		expectedToken string
		expectedError string
	}{
		{
			name: "Success - Login with Email",
			setupMocks: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				mockRepo.On("GetUserByEmail", "test@example.com").Return(&models.User{
					ID:       1,
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}, nil)
			},
			creds: models.Credentials{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			expectedToken: "mocked-jwt-token",
			expectedError: "",
		},
		{
			name: "Success - Login with Phone",
			setupMocks: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				mockRepo.On("GetUserByPhone", "1234567890").Return(&models.User{
					ID:       1,
					Phone:    "1234567890",
					Password: string(hashedPassword),
				}, nil)
			},
			creds: models.Credentials{
				Identifier: "1234567890",
				Password:   "password123",
			},
			expectedToken: "mocked-jwt-token",
			expectedError: "",
		},
		{
			name:       "Error - Identifier Required",
			setupMocks: func() {},
			creds: models.Credentials{
				Identifier: "",
				Password:   "password123",
			},
			expectedToken: "",
			expectedError: "identifier is required",
		},
		{
			name: "Error - User Email Not Found",
			setupMocks: func() {
				mockRepo.On("GetUserByEmail", "unknown@example.com").Return(nil, errors.New("user not found"))
			},
			creds: models.Credentials{
				Identifier: "unknown@example.com",
				Password:   "password123",
			},
			expectedToken: "",
			expectedError: "invalid email/phone or password",
		},
		{
			name: "Error - User Phone Not Found",
			setupMocks: func() {
				mockRepo.On("GetUserByPhone", "1234567890").Return(nil, errors.New("user not found"))
			},
			creds: models.Credentials{
				Identifier: "1234567890",
				Password:   "password123",
			},
			expectedToken: "",
			expectedError: "invalid email/phone or password",
		},
		{
			name: "Error - Invalid Password",
			setupMocks: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
				mockRepo.On("GetUserByEmail", "test@example.com").Return(&models.User{
					ID:       1,
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}, nil)
			},
			creds: models.Credentials{
				Identifier: "test@example.com",
				Password:   "wrongpassword",
			},
			expectedToken: "",
			expectedError: "invalid email/phone or password",
		},
		{
			name: "Error - JWT Generation Failure",
			setupMocks: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				mockRepo.On("GetUserByEmail", "test@example.com").Return(&models.User{
					ID:       1,
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}, nil)

				utils.GenerateJWT = func(userID int) (string, error) {
					return "", errors.New("failed to generate token")
				}
			},
			creds: models.Credentials{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			expectedToken: "",
			expectedError: "failed to generate token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tc.setupMocks()

			token, err := service.Login(tc.creds)

			if tc.expectedError == "" {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedToken, token)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
