package models

import "time"

type PremiumFeatures struct {
	UnlimitedSwipes bool `json:"unlimited_swipes"`
	IsVerified      bool `json:"profile_boost"`
}

type User struct {
	ID              int             `json:"id"`
	Email           string          `json:"email"`
	Password        string          `json:"password"`
	Phone           string          `json:"phone"`
	Name            string          `json:"name"`
	Gender          string          `json:"gender"` // "male" or "female"
	IsInactive      bool            `json:"is_inactive"`
	PremiumExpiry   *time.Time      `json:"premium_expiry"`
	PremiumFeatures PremiumFeatures `json:"premium_features"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type Credentials struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type Swipe struct {
	UserID       int       `json:"user_id"`
	TargetUserID int       `json:"target_user_id"`
	Action       string    `json:"action"` // "like" or "pass"
	CreatedAt    time.Time `json:"created_at"`
}
