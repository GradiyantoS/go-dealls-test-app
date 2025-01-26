package repositories

import (
	"errors"

	"github.com/GradiyantoS/go-dealls-test-app/models"
)

type UserRepository interface {
	GenerateUserID() int
	GetAllUsers() []*models.User
	GetUserByID(userID int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByPhone(phone string) (*models.User, error)
	SaveUser(user *models.User) error
	UpdateUser(user *models.User) error
	GetSwipesForUser(userID int) []models.Swipe
	SaveSwipe(swipe *models.Swipe) error
}

type userRepository struct {
	users      map[int]*models.User
	swipes     []models.Swipe
	nextUserID int
}

// NewUserRepository creates a new instance of userRepository.
func NewUserRepository() UserRepository {
	return &userRepository{
		users:      make(map[int]*models.User),
		swipes:     []models.Swipe{},
		nextUserID: 1,
	}
}

// GenerateUserID generates the next unique user ID.
func (r *userRepository) GenerateUserID() int {
	id := r.nextUserID
	r.nextUserID++
	return id
}

// GetAllUsers retrieves all users.
func (r *userRepository) GetAllUsers() []*models.User {
	var result []*models.User
	for _, user := range r.users {
		result = append(result, user)
	}
	return result
}

// GetUserByID retrieves a user by their ID.
func (r *userRepository) GetUserByID(userID int) (*models.User, error) {
	user, exists := r.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetUserByEmail retrieves a user by their email.
func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// GetUserByPhone retrieves a user by their phone number.
func (r *userRepository) GetUserByPhone(phone string) (*models.User, error) {
	for _, user := range r.users {
		if user.Phone == phone {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// SaveUser saves a new user.
func (r *userRepository) SaveUser(user *models.User) error {
	if _, exists := r.users[user.ID]; exists {
		return errors.New("user ID already exists")
	}
	r.users[user.ID] = user
	r.nextUserID++
	return nil
}

// UpdateUser updates an existing user.
func (r *userRepository) UpdateUser(user *models.User) error {
	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	r.users[user.ID] = user
	return nil
}

// GetSwipesForUser retrieves all swipes for a specific user.
func (r *userRepository) GetSwipesForUser(userID int) []models.Swipe {
	var result []models.Swipe
	for _, swipe := range r.swipes {
		if swipe.UserID == userID {
			result = append(result, swipe)
		}
	}
	return result
}

// SaveSwipe saves a swipe action.
func (r *userRepository) SaveSwipe(swipe *models.Swipe) error {
	r.swipes = append(r.swipes, *swipe)
	return nil
}

func (r *userRepository) ClearData() {
	r.users = make(map[int]*models.User)
	r.nextUserID = 1
}
