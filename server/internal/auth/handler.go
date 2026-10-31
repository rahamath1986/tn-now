package auth

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type registerRequest struct {
	FullName string `json:"fullName"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type logoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type sendOTPRequest struct {
	Email string `json:"email"`
}

type verifyOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/register", h.HandleRegister)
	mux.HandleFunc("/auth/login", h.HandleLogin)
	mux.HandleFunc("/auth/refresh", h.HandleRefresh)
	mux.HandleFunc("/auth/logout", h.HandleLogout)
	mux.HandleFunc("/auth/me", h.HandleMe)
	mux.HandleFunc("/api/user/me", h.HandleMe)
	mux.HandleFunc("/auth/profile", h.HandleUpdateProfile)
	mux.HandleFunc("/api/user/profile", h.HandleUpdateProfile)
	mux.HandleFunc("/auth/otp/send", h.HandleSendOTP)
	mux.HandleFunc("/auth/otp/verify", h.HandleVerifyOTP)
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// Validate inputs
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.FullName = strings.TrimSpace(req.FullName)

	var validationErrors []map[string]string
	if req.Username == "" {
		validationErrors = append(validationErrors, map[string]string{"field": "username", "message": "Username is required"})
	}
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		validationErrors = append(validationErrors, map[string]string{"field": "email", "message": "Valid email is required"})
	}
	if len(req.Password) < 6 {
		validationErrors = append(validationErrors, map[string]string{"field": "password", "message": "Password must be at least 6 characters"})
	}

	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"data":    nil,
			"message": "Validation failed",
			"errors":  validationErrors,
		})
		return
	}

	resp, err := h.service.Register(r.Context(), req.Username, req.Email, req.Password, req.FullName)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	writeSuccess(w, http.StatusCreated, resp, "பதிவு வெற்றிகரமானது. உங்கள் மின்னஞ்சலுக்கு அனுப்பப்பட்ட OTP சரிபார்ப்புக் குறியீட்டை உள்ளிடவும்.")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	deviceInfo := r.UserAgent()
	resp, err := h.service.Login(r.Context(), req.Email, req.Password, deviceInfo)
	if err != nil {
		if err.Error() == "EMAIL_NOT_VERIFIED" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success":     false,
				"requiresOtp": true,
				"email":       req.Email,
				"message":     "உங்கள் மின்னஞ்சல் இன்னும் சரிபார்க்கப்படவில்லை. புதிய OTP அனுப்பப்பட்டுள்ளது.",
			})
			return
		}
		writeError(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	setSessionCookies(w, resp)
	writeSuccess(w, http.StatusOK, resp, "Login successful")
}

func (h *Handler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req refreshRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	if req.RefreshToken == "" {
		if cookie, err := r.Cookie("tn24_refresh_token"); err == nil && cookie.Value != "" {
			req.RefreshToken = cookie.Value
		}
	}

	if req.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "RefreshToken is required", nil)
		return
	}

	deviceInfo := r.UserAgent()
	resp, err := h.service.Refresh(r.Context(), req.RefreshToken, deviceInfo)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	setSessionCookies(w, resp)
	writeSuccess(w, http.StatusOK, resp, "Token refreshed successfully")
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req logoutRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	if req.RefreshToken == "" {
		if cookie, err := r.Cookie("tn24_refresh_token"); err == nil && cookie.Value != "" {
			req.RefreshToken = cookie.Value
		}
	}

	if req.RefreshToken != "" {
		_ = h.service.Logout(r.Context(), req.RefreshToken)
	}

	clearSessionCookies(w)
	writeSuccess(w, http.StatusOK, nil, "Logged out successfully")
}

func (h *Handler) HandleMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var token string
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	} else if cookie, err := r.Cookie("tn24_user_token"); err == nil && cookie.Value != "" {
		token = cookie.Value
	}

	if token == "" {
		writeError(w, http.StatusUnauthorized, "Not authenticated", nil)
		return
	}

	claims, err := h.service.VerifyToken(token)
	if err != nil {
		// Attempt auto-refresh if refresh cookie exists
		if refCookie, refErr := r.Cookie("tn24_refresh_token"); refErr == nil && refCookie.Value != "" {
			refResp, refErr2 := h.service.Refresh(r.Context(), refCookie.Value, r.UserAgent())
			if refErr2 == nil && refResp != nil {
				setSessionCookies(w, refResp)
				claims, err = h.service.VerifyToken(refResp.AccessToken)
			}
		}
	}

	if err != nil || claims == nil {
		writeError(w, http.StatusUnauthorized, "Invalid or expired token", nil)
		return
	}

	profile, err := h.service.GetMe(r.Context(), claims.UserID)
	if err != nil {
		writeError(w, http.StatusNotFound, "User not found", nil)
		return
	}

	writeSuccess(w, http.StatusOK, profile, "User profile retrieved")
}

func (h *Handler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var token string
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	} else if cookie, err := r.Cookie("tn24_user_token"); err == nil && cookie.Value != "" {
		token = cookie.Value
	}

	if token == "" {
		writeError(w, http.StatusUnauthorized, "உள்நுழைவு தேவை (Not authenticated)", nil)
		return
	}

	claims, err := h.service.VerifyToken(token)
	if err != nil {
		if refCookie, refErr := r.Cookie("tn24_refresh_token"); refErr == nil && refCookie.Value != "" {
			refResp, refErr2 := h.service.Refresh(r.Context(), refCookie.Value, r.UserAgent())
			if refErr2 == nil && refResp != nil {
				setSessionCookies(w, refResp)
				claims, err = h.service.VerifyToken(refResp.AccessToken)
			}
		}
	}

	if err != nil || claims == nil {
		writeError(w, http.StatusUnauthorized, "காலாவதியான அமர்வு. மீண்டும் உள்நுழையவும் (Session expired)", nil)
		return
	}

	var params UpdateProfileParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	updated, err := h.service.UpdateUserProfile(r.Context(), claims.UserID, params)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	writeSuccess(w, http.StatusOK, updated, "சுயவிவரம் வெற்றிகரமாக புதுப்பிக்கப்பட்டது (Profile updated successfully)")
}

func (h *Handler) HandleSendOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req sendOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "Email is required", nil)
		return
	}

	err := h.service.SendEmailOTP(r.Context(), req.Email)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	writeSuccess(w, http.StatusOK, map[string]string{
		"email": req.Email,
	}, "OTP சரிபார்ப்புக் குறியீடு உங்கள் மின்னஞ்சலுக்கு அனுப்பப்பட்டுள்ளது")
}

func (h *Handler) HandleVerifyOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req verifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.OTP = strings.TrimSpace(req.OTP)

	if req.Email == "" || req.OTP == "" {
		writeError(w, http.StatusBadRequest, "Email and OTP are required", nil)
		return
	}

	deviceInfo := r.UserAgent()
	resp, err := h.service.VerifyEmailOTP(r.Context(), req.Email, req.OTP, deviceInfo)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	setSessionCookies(w, resp)
	writeSuccess(w, http.StatusOK, resp, "OTP சரிபார்க்கப்பட்டு வெற்றிகரமாக உள்நுழைந்துள்ளீர்கள்")
}

func setSessionCookies(w http.ResponseWriter, resp *AuthResponse) {
	http.SetCookie(w, &http.Cookie{
		Name:     "tn24_user_token",
		Value:    resp.AccessToken,
		Path:     "/",
		HttpOnly: false,
		MaxAge:   7 * 24 * 3600,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "tn24_refresh_token",
		Value:    resp.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   7 * 24 * 3600,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearSessionCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "tn24_user_token",
		Value:    "",
		Path:     "/",
		HttpOnly: false,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "tn24_refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
}

func writeSuccess(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
		"message": message,
		"errors":  []interface{}{},
	})
}

func writeError(w http.ResponseWriter, statusCode int, message string, errors []interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if errors == nil {
		errors = []interface{}{}
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"data":    nil,
		"message": message,
		"errors":  errors,
	})
}
