package mock

import (
	"github.com/GradiyantoS/go-dealls-test-app/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GenerateUserID() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockUserRepository) GetAllUsers() []*models.User {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*models.User)
	}
	return nil
}

func (m *MockUserRepository) GetUserByID(userID int) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByPhone(phone string) (*models.User, error) {
	args := m.Called(phone)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) SaveUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetSwipesForUser(userID int) []models.Swipe {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Swipe)
	}
	return nil
}

func (m *MockUserRepository) SaveSwipe(swipe *models.Swipe) error {
	args := m.Called(swipe)
	return args.Error(0)
}
