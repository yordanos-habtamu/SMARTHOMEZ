package types

import (
	"time"
)

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(User) error
	GetUserById(id int) (*User, error)
	GetAllUsers() ([]User, error)
}

type HouseStore interface {
	CreateHouse(House) error
	GetHouseByAddress(address string) (*House, error)
	GetHouseById(id uint) (*House, error)
	GetAllHouses() ([]House, error)
	UpdateHouse(id uint, payload RegisterHousePayload) (*House, error)
	DeleteHouse(id uint) error
	GetHousesByCategory(category string) ([]House, error)
	GetHousesByPriceRange(minPrice, maxPrice float64) ([]House, error)
	GetLatestHouses(limit int) ([]House, error)
}

type ReferralStore interface {
	CreateReferralCode(userID uint, code, shortCode, qrCodeURL string) error
	GetReferralByCode(code string) (*ReferralCode, error)
	GetReferralByShortCode(shortCode string) (*ReferralCode, error)
	GetReferralByUserID(userID uint) (*ReferralCode, error)
	TrackVisit(referralCodeID uint, ipAddress, userAgent string) error
	TrackSignup(referralCodeID, convertedUserID uint, ipAddress, userAgent string) error
	GetReferralStats(userID uint) (*ReferralStats, error)
	UpdateVisitCount(referralCodeID uint) error
	UpdateSignupCount(referralCodeID uint) error
}

type User struct {
	ID  uint  `json:"Id"`
	FirstName string `json:"firstName"`
	LastName string   `json:"lastName"`
	Sex string `json:"sex"`
	Email string     `json:"email"`
	DoB time.Time   `json:"DoB"`
	Contact string `json:"contact"`
	Password string  `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	Role      string    `json:"role"`
}
type House struct {
	ID           uint      `json:"id"`
	Address      string    `json:"address"`
	Price        float64   `json:"price"`
	NumBathrooms int       `json:"numBathrooms"`
	NumBedrooms  int       `json:"numRooms"`
	AreaSqFt     float64      `json:"area_sq_ft"`
	Category     string    `json:"category"`
	ImageURL     string    `json:"imageUrl"`
	IsSold       bool      `json:"isSold"`
	AgentID      uint      `json:"agentId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Description  string    `json:"description"`
	ImgURL       string    `json:"imgUrl"`
}

type RegisterHousePayload struct {
	Address      string  `json:"address" validate:"required"`
	NumBathrooms int     `json:"numBathrooms" validate:"required"`
	NumBedrooms  int     `json:"numBedRooms" validate:"required"`
	Description  string  `json:"description" validate:"required"`
	Price        float64 `json:"price" validate:"required"`
	Category     string  `json:"catagory" validate:"required"`
	ImgUrl       string  `json:"imgUrl" validate:"required"`
	IsSold       bool    `json:"isSold"`
	AgentID      uint    `json:"agentId"`
	AreaSqFt     float64     `json:"areaSqFt" validate:"required"`
}

type RegisterUserPayload struct {
	FirstName    string `json:"firstName" validate:"required"`
	LastName     string `json:"lastName" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Contact      string `json:"contact" validate:"required"`
	DoB          string `json:"DoB" validate:"required"`
	Sex          string `json:"sex" validate:"required"`
	Password     string `json:"password" validate:"required,min=6,max=12"`
	Role         string `json:"role"`
	ReferralCode string `json:"referralCode,omitempty"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ReferralCode struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"userId"`
	Code        string    `json:"code"`
	ShortCode   string    `json:"shortCode"`
	QRCodeURL   string    `json:"qrCodeUrl"`
	VisitCount  int       `json:"visitCount"`
	SignupCount int       `json:"signupCount"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ReferralAnalytics struct {
	ID              uint      `json:"id"`
	ReferralCodeID  uint      `json:"referralCodeId"`
	EventType       string    `json:"eventType"` // "visit" or "signup"
	IPAddress       string    `json:"ipAddress"`
	UserAgent       string    `json:"userAgent"`
	ConvertedUserID *uint     `json:"convertedUserId,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
}

type ReferralStats struct {
	Code        string `json:"code"`
	ShortCode   string `json:"shortCode"`
	QRCodeURL   string `json:"qrCodeUrl"`
	VisitCount  int    `json:"visitCount"`
	SignupCount int    `json:"signupCount"`
}
