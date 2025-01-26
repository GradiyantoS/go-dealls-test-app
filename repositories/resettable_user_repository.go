// repositories/resettable_user_repository.go
package repositories

import (
	"time"

	"github.com/GradiyantoS/go-dealls-test-app/models"
	"golang.org/x/crypto/bcrypt"
)

type ResettableUserRepository struct {
	*userRepository
}

// NewResettableUserRepository creates a resettable version of UserRepository.
func NewResettableUserRepository() *ResettableUserRepository {
	return &ResettableUserRepository{
		userRepository: NewUserRepository().(*userRepository),
	}
}

// ClearData resets all users in the repository.
func (r *ResettableUserRepository) ClearData() {
	r.users = make(map[int]*models.User)
	r.nextUserID = 1
}

// SeedTestData adds initial users to the repository for testing.
func (r *ResettableUserRepository) SeedTestData() {
	r.ClearData()

	hashedPassword1, _ := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
	hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("password2"), bcrypt.DefaultCost)

	r.SaveUser(&models.User{
		ID:        1,
		Email:     "test1@example.com",
		Phone:     "1234567890",
		Password:  string(hashedPassword1),
		Name:      "User 1",
		Gender:    "male",
		CreatedAt: time.Now().Add(-48 * time.Hour),
		UpdatedAt: time.Now(),
	})

	r.SaveUser(&models.User{
		ID:        2,
		Email:     "test2@example.com",
		Phone:     "0987654321",
		Password:  string(hashedPassword2),
		Name:      "User 2",
		Gender:    "female",
		CreatedAt: time.Now().Add(-72 * time.Hour),
		UpdatedAt: time.Now(),
	})
}
