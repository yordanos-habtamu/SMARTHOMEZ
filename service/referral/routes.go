package referral

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yordanos-habtamu/realstate/config"
	"github.com/yordanos-habtamu/realstate/service/auth"
	"github.com/yordanos-habtamu/realstate/types"
	"github.com/yordanos-habtamu/realstate/utils"
)

type Handler struct {
	store     types.ReferralStore
	userStore types.UserStore
}

func NewHandler(store types.ReferralStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/referral/generate", auth.WithJwtAuth(h.handleGenerateReferral, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/referral/stats", auth.WithJwtAuth(h.handleGetStats, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/r/{shortCode}", h.handleTrackVisit).Methods(http.MethodGet)
}

func (h *Handler) handleGenerateReferral(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserfromContext(r.Context())

	// Check if user already has a referral code
	existing, err := h.store.GetReferralByUserID(user.ID)
	if err == nil && existing != nil {
		// User already has a referral code, return it
		utils.WriteJson(w, http.StatusOK, existing)
		return
	}

	// Generate new referral code and short code
	code, err := utils.GenerateReferralCode()
	if err != nil {
		log.Printf("Error generating referral code: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to generate referral code"))
		return
	}

	shortCode, err := utils.GenerateShortCode(8)
	if err != nil {
		log.Printf("Error generating short code: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to generate short code"))
		return
	}

	// Generate referral URL
	referralURL := fmt.Sprintf("%s/r/%s", config.Envs.PUBLIC_HOST, shortCode)

	// Generate QR code for the referral URL
	qrCodeURL, err := utils.GenerateQRCode(referralURL)
	if err != nil {
		log.Printf("Error generating QR code: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to generate QR code"))
		return
	}

	// Save to database
	err = h.store.CreateReferralCode(user.ID, code, shortCode, qrCodeURL)
	if err != nil {
		log.Printf("Error saving referral code: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to save referral code"))
		return
	}

	// Get the created referral code
	referralCode, err := h.store.GetReferralByUserID(user.ID)
	if err != nil {
		log.Printf("Error fetching referral code: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to fetch referral code"))
		return
	}

	utils.WriteJson(w, http.StatusCreated, referralCode)
}

func (h *Handler) handleGetStats(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserfromContext(r.Context())

	stats, err := h.store.GetReferralStats(user.ID)
	if err != nil {
		log.Printf("Error fetching referral stats: %v", err)
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("no referral code found"))
		return
	}

	utils.WriteJson(w, http.StatusOK, stats)
}

func (h *Handler) handleTrackVisit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	// Get referral code
	referralCode, err := h.store.GetReferralByShortCode(shortCode)
	if err != nil {
		log.Printf("Invalid referral code: %v", err)
		http.Redirect(w, r, config.Envs.PUBLIC_HOST, http.StatusFound)
		return
	}

	// Track the visit
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	err = h.store.TrackVisit(referralCode.ID, ipAddress, userAgent)
	if err != nil {
		log.Printf("Error tracking visit: %v", err)
	}

	// Update visit count
	err = h.store.UpdateVisitCount(referralCode.ID)
	if err != nil {
		log.Printf("Error updating visit count: %v", err)
	}

	// Redirect to registration page with referral code
	redirectURL := fmt.Sprintf("%s/register?ref=%s", config.Envs.PUBLIC_HOST, shortCode)

	// For now, return a simple HTML page with instructions
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Referral Link</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 600px;
            margin: 50px auto;
            padding: 20px;
            text-align: center;
        }
        .container {
            background: #f5f5f5;
            padding: 30px;
            border-radius: 10px;
        }
        h1 { color: #333; }
        p { color: #666; line-height: 1.6; }
        .code {
            background: white;
            padding: 15px;
            border-radius: 5px;
            font-size: 20px;
            font-weight: bold;
            color: #007bff;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Welcome to SMARTHOMEZ!</h1>
        <p>You've been referred by one of our agents.</p>
        <p>Use this referral code when registering:</p>
        <div class="code">%s</div>
        <p>Register at: <a href="%s">%s/register</a></p>
        <p>Include the referral code above in your registration to connect with your agent.</p>
    </div>
</body>
</html>
	`, shortCode, redirectURL, config.Envs.PUBLIC_HOST)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
