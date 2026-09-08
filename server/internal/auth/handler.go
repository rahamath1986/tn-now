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

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/register", h.HandleRegister)
	mux.HandleFunc("/auth/login", h.HandleLogin)
	mux.HandleFunc("/auth/refresh", h.HandleRefresh)
	mux.HandleFunc("/auth/logout", h.HandleLogout)
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

	resp, err := h.service.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	writeSuccess(w, http.StatusCreated, resp, "Registration successful")
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
		writeError(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	writeSuccess(w, http.StatusOK, resp, "Login successful")
}

func (h *Handler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
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

	writeSuccess(w, http.StatusOK, resp, "Token refreshed successfully")
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if req.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "RefreshToken is required", nil)
		return
	}

	err := h.service.Logout(r.Context(), req.RefreshToken)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	writeSuccess(w, http.StatusOK, nil, "Logged out successfully")
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
