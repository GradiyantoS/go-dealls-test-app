package integration_test

import (
	"time"

	"github.com/GradiyantoS/go-dealls-test-app/models"
	"github.com/GradiyantoS/go-dealls-test-app/repositories"
	"golang.org/x/crypto/bcrypt"
)

// ResettableTestRepository is a wrapper around UserRepository for testing purposes
type ResettableTestRepository struct {
	repo repositories.UserRepository
}

// NewResettableTestRepository creates a new instance of ResettableTestRepository
func NewResettableTestRepository(repo repositories.UserRepository) *ResettableTestRepository {
	return &ResettableTestRepository{repo: repo}
}

// ClearData resets the repository data
func (r *ResettableTestRepository) ClearData() {
	// Directly clear the internal map by reinitializing the repository
	r.repo = repositories.NewUserRepository()
}

// SeedTestData seeds predefined test data
func (r *ResettableTestRepository) SeedTestData() {
	r.ClearData()

	hashedPassword1, _ := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
	hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("password2"), bcrypt.DefaultCost)

	// Save users with seeded data
	r.repo.SaveUser(&models.User{
		Email:     "test1@example.com",
		Phone:     "1234567890",
		Password:  string(hashedPassword1),
		Name:      "User 1",
		Gender:    "male",
		CreatedAt: time.Now().Add(-48 * time.Hour),
		UpdatedAt: time.Now(),
	})

	r.repo.SaveUser(&models.User{
		Email:     "test2@example.com",
		Phone:     "0987654321",
		Password:  string(hashedPassword2),
		Name:      "User 2",
		Gender:    "female",
		CreatedAt: time.Now().Add(-72 * time.Hour),
		UpdatedAt: time.Now(),
	})
}

// GetRepository returns the wrapped repository
func (r *ResettableTestRepository) GetRepository() repositories.UserRepository {
	return r.repo
}
