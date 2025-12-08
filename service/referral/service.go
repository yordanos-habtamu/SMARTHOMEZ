package referral

import (
	"database/sql"
	"fmt"

	"github.com/yordanos-habtamu/realstate/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateReferralCode(userID uint, code, shortCode, qrCodeURL string) error {
	_, err := s.db.Exec(
		"INSERT INTO referral_codes (user_id, code, short_code, qr_code_url) VALUES ($1, $2, $3, $4)",
		userID, code, shortCode, qrCodeURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create referral code: %v", err)
	}
	return nil
}

func (s *Store) GetReferralByCode(code string) (*types.ReferralCode, error) {
	var ref types.ReferralCode
	err := s.db.QueryRow(
		"SELECT id, user_id, code, short_code, qr_code_url, visit_count, signup_count, created_at FROM referral_codes WHERE code = $1",
		code,
	).Scan(&ref.ID, &ref.UserID, &ref.Code, &ref.ShortCode, &ref.QRCodeURL, &ref.VisitCount, &ref.SignupCount, &ref.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("referral code not found")
		}
		return nil, err
	}
	return &ref, nil
}

func (s *Store) GetReferralByShortCode(shortCode string) (*types.ReferralCode, error) {
	var ref types.ReferralCode
	err := s.db.QueryRow(
		"SELECT id, user_id, code, short_code, qr_code_url, visit_count, signup_count, created_at FROM referral_codes WHERE short_code = $1",
		shortCode,
	).Scan(&ref.ID, &ref.UserID, &ref.Code, &ref.ShortCode, &ref.QRCodeURL, &ref.VisitCount, &ref.SignupCount, &ref.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("referral code not found")
		}
		return nil, err
	}
	return &ref, nil
}

func (s *Store) GetReferralByUserID(userID uint) (*types.ReferralCode, error) {
	var ref types.ReferralCode
	err := s.db.QueryRow(
		"SELECT id, user_id, code, short_code, qr_code_url, visit_count, signup_count, created_at FROM referral_codes WHERE user_id = $1",
		userID,
	).Scan(&ref.ID, &ref.UserID, &ref.Code, &ref.ShortCode, &ref.QRCodeURL, &ref.VisitCount, &ref.SignupCount, &ref.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("referral code not found for user")
		}
		return nil, err
	}
	return &ref, nil
}

func (s *Store) TrackVisit(referralCodeID uint, ipAddress, userAgent string) error {
	_, err := s.db.Exec(
		"INSERT INTO referral_analytics (referral_code_id, event_type, ip_address, user_agent) VALUES ($1, 'visit', $2, $3)",
		referralCodeID, ipAddress, userAgent,
	)
	if err != nil {
		return fmt.Errorf("failed to track visit: %v", err)
	}
	return nil
}

func (s *Store) TrackSignup(referralCodeID, convertedUserID uint, ipAddress, userAgent string) error {
	_, err := s.db.Exec(
		"INSERT INTO referral_analytics (referral_code_id, event_type, ip_address, user_agent, converted_user_id) VALUES ($1, 'signup', $2, $3, $4)",
		referralCodeID, ipAddress, userAgent, convertedUserID,
	)
	if err != nil {
		return fmt.Errorf("failed to track signup: %v", err)
	}
	return nil
}

func (s *Store) UpdateVisitCount(referralCodeID uint) error {
	_, err := s.db.Exec(
		"UPDATE referral_codes SET visit_count = visit_count + 1 WHERE id = $1",
		referralCodeID,
	)
	if err != nil {
		return fmt.Errorf("failed to update visit count: %v", err)
	}
	return nil
}

func (s *Store) UpdateSignupCount(referralCodeID uint) error {
	_, err := s.db.Exec(
		"UPDATE referral_codes SET signup_count = signup_count + 1 WHERE id = $1",
		referralCodeID,
	)
	if err != nil {
		return fmt.Errorf("failed to update signup count: %v", err)
	}
	return nil
}

func (s *Store) GetReferralStats(userID uint) (*types.ReferralStats, error) {
	var stats types.ReferralStats
	err := s.db.QueryRow(
		"SELECT code, short_code, qr_code_url, visit_count, signup_count FROM referral_codes WHERE user_id = $1",
		userID,
	).Scan(&stats.Code, &stats.ShortCode, &stats.QRCodeURL, &stats.VisitCount, &stats.SignupCount)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no referral stats found for user")
		}
		return nil, err
	}
	return &stats, nil
}
