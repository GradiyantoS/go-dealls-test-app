package services

import (
	"errors"
	"time"

	"github.com/GradiyantoS/go-dealls-test-app/models"
	"github.com/GradiyantoS/go-dealls-test-app/repositories"
)

type SwipeService interface {
	RecordSwipe(swipe *models.Swipe) error
	GetSwipeCandidates(userID int) ([]models.User, error)
}

type swipeService struct {
	userRepo repositories.UserRepository
}

func NewSwipeService(userRepo repositories.UserRepository) SwipeService {
	return &swipeService{userRepo}
}

func (s *swipeService) RecordSwipe(swipe *models.Swipe) error {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	totalSwipes := 0

	for _, s := range s.userRepo.GetSwipesForUser(swipe.UserID) {
		if !s.CreatedAt.Before(today) { // Include swipes on or after "today"
			totalSwipes++
			if s.TargetUserID == swipe.TargetUserID {
				return errors.New("you have already swiped on this profile today")
			}
		}
	}

	user, err := s.userRepo.GetUserByID(swipe.UserID)
	if err != nil {
		return err
	}

	if user.PremiumExpiry != nil && user.PremiumExpiry.After(time.Now().UTC()) {
		if user.PremiumFeatures.UnlimitedSwipes {
			totalSwipes = 0
		}
	}

	if totalSwipes >= 10 {
		return errors.New("daily swipe limit reached")
	}
	swipe.CreatedAt = time.Now().UTC()
	s.userRepo.SaveSwipe(swipe)
	return nil
}

// GetSwipeCandidates retrieves profiles that the user has not swiped on today
func (s *swipeService) GetSwipeCandidates(userID int) ([]models.User, error) {
	today := time.Now().Truncate(24 * time.Hour)
	swipedUserIDs := map[int]bool{}

	// Validate user existence first
	currentUser, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Collect user IDs that have already been swiped on today
	for _, swipe := range s.userRepo.GetSwipesForUser(userID) {
		if swipe.CreatedAt.After(today) {
			swipedUserIDs[swipe.TargetUserID] = true
		}
	}

	// Determine the opposite gender
	oppositeGender := "male"
	if currentUser.Gender == "male" {
		oppositeGender = "female"
	}

	// Retrieve all users and filter by opposite gender and unswiped profiles
	candidates := []models.User{}
	for _, user := range s.userRepo.GetAllUsers() {
		if user.ID != userID &&
			user.Gender == oppositeGender &&
			!swipedUserIDs[user.ID] &&
			!user.IsInactive {
			candidates = append(candidates, *user)
		}
	}

	return candidates, nil
}
