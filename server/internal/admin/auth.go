package admin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	AdminSessionCookieName = "tn24_admin_session"
	AdminSessionDuration   = 24 * time.Hour
)

// GenerateAdminSession creates an HMAC-SHA256 signed session token
func GenerateAdminSession(username, secret string) string {
	exp := time.Now().Add(AdminSessionDuration).Unix()
	payload := fmt.Sprintf("%s:%d", username, exp)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))
	tokenRaw := fmt.Sprintf("%s:%s", payload, sig)
	return base64.URLEncoding.EncodeToString([]byte(tokenRaw))
}

// ValidateAdminSession verifies token signature and expiration
func ValidateAdminSession(tokenStr, secret string) (string, bool) {
	if tokenStr == "" {
		return "", false
	}
	rawBytes, err := base64.URLEncoding.DecodeString(tokenStr)
	if err != nil {
		return "", false
	}
	parts := strings.Split(string(rawBytes), ":")
	if len(parts) != 3 {
		return "", false
	}
	username := parts[0]
	expStr := parts[1]
	sig := parts[2]

	exp, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return "", false
	}

	payload := fmt.Sprintf("%s:%d", username, exp)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return "", false
	}
	return username, true
}

// CheckAdminAuth extracts and validates session from cookie or header
func (h *Handler) CheckAdminAuth(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(AdminSessionCookieName)
	if err == nil && cookie.Value != "" {
		if username, ok := ValidateAdminSession(cookie.Value, h.cfg.JWTSecret); ok {
			return username, true
		}
	}
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if username, ok := ValidateAdminSession(token, h.cfg.JWTSecret); ok {
			return username, true
		}
	}
	return "", false
}

// HandleLogin renders login page on GET, handles verification on POST
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if _, ok := h.CheckAdminAuth(r); ok {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(RenderLoginPage("")))
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"error":   "Invalid JSON request",
				})
				return
			}
		} else {
			_ = r.ParseForm()
			req.Username = r.FormValue("username")
			req.Password = r.FormValue("password")
		}

		expectedUser := strings.TrimSpace(h.cfg.AdminUsername)
		if expectedUser == "" {
			expectedUser = "admin"
		}
		expectedPass := strings.TrimSpace(h.cfg.AdminPassword)
		if expectedPass == "" {
			expectedPass = "tn24admin"
		}

		reqUser := strings.TrimSpace(req.Username)
		reqPass := strings.TrimSpace(req.Password)

		if reqUser != expectedUser || reqPass != expectedPass {
			if strings.Contains(contentType, "application/json") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"error":   "Invalid username or password",
				})
			} else {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(RenderLoginPage("தவறான பயனர் பெயர் அல்லது கடவுச்சொல் (Invalid username or password)")))
			}
			return
		}

		// Success: generate session
		token := GenerateAdminSession(reqUser, h.cfg.JWTSecret)
		http.SetCookie(w, &http.Cookie{
			Name:     AdminSessionCookieName,
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(AdminSessionDuration),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   r.TLS != nil,
		})

		if strings.Contains(contentType, "application/json") {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success":  true,
				"username": reqUser,
				"redirect": "/admin",
			})
		} else {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		}
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// HandleLogout clears the admin session cookie and redirects to login
func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     AdminSessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"redirect": "/admin/login",
		})
		return
	}

	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

// RequireAuth wraps admin HTTP handlers to ensure valid admin authentication
func (h *Handler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow public assets or login endpoints
		path := r.URL.Path
		if path == "/admin/login" || path == "/admin/logout" || path == "/admin/api/maps/svg" {
			next(w, r)
			return
		}

		username, ok := h.CheckAdminAuth(r)
		if !ok {
			if strings.HasPrefix(path, "/admin/api/") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"error":   "Unauthorized. Please login at /admin/login",
				})
				return
			}
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		// Store username in request context or header if needed
		r.Header.Set("X-Admin-User", username)
		next(w, r)
	}
}
