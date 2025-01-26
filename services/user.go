package services

import (
	"errors"
	"time"

	"github.com/GradiyantoS/go-dealls-test-app/models"
	"github.com/GradiyantoS/go-dealls-test-app/repositories"
	"github.com/GradiyantoS/go-dealls-test-app/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUp(user *models.User) error
	Login(creds models.Credentials) (string, error)
	EnablePremiumFeature(userID int, duration int, features []string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) SignUp(user *models.User) error {
	if _, err := s.userRepo.GetUserByEmail(user.Email); err == nil {
		return errors.New("email already exists")
	}

	if _, err := s.userRepo.GetUserByPhone(user.Phone); err == nil {
		return errors.New("phone number already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to save password")
	}

	user.ID = s.userRepo.GenerateUserID()
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := s.userRepo.SaveUser(user); err != nil {
		return err
	}
	return nil
}

func (s *userService) Login(creds models.Credentials) (string, error) {
	var user *models.User
	var err error

	if creds.Identifier != "" {
		if utils.IsEmail(creds.Identifier) {
			user, err = s.userRepo.GetUserByEmail(creds.Identifier)
		} else {
			user, err = s.userRepo.GetUserByPhone(creds.Identifier)
		}
		if err != nil {
			return "", errors.New("invalid email/phone or password")
		}
	} else {
		return "", errors.New("identifier is required")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		return "", errors.New("invalid email/phone or password")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

// PurchasePremium activates a premium feature for a user
func (s *userService) EnablePremiumFeature(userID int, duration int, features []string) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.IsInactive {
		return errors.New("user account is inactive")
	}

	// Check if PremiumExpiry is in the future; extend or set it
	newExpiry := time.Now().Add(time.Duration(duration) * 24 * time.Hour)
	if user.PremiumExpiry != nil && user.PremiumExpiry.After(time.Now()) {
		newExpiry = user.PremiumExpiry.Add(time.Duration(duration) * 24 * time.Hour)
	}
	user.PremiumExpiry = &newExpiry

	// Enable the specified features
	for _, feature := range features {
		switch feature {
		case "UnlimitedSwipes":
			if user.PremiumFeatures.UnlimitedSwipes {
				return errors.New("unlimited swipes is already active")
			}
			user.PremiumFeatures.UnlimitedSwipes = true
		case "IsVerified":
			if user.PremiumFeatures.IsVerified {
				return errors.New("user is already verified")
			}
			user.PremiumFeatures.IsVerified = true
		default:
			return errors.New("invalid premium feature: " + feature)
		}
	}

	// Update user and persist changes
	user.UpdatedAt = time.Now()
	s.userRepo.UpdateUser(user)
	return nil
}
