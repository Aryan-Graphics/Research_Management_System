package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"research-management-system/internal/database"
	"research-management-system/internal/middleware"
	"research-management-system/internal/models"
)

type AuthHandler struct {
	store database.Store
}

func NewAuthHandler(store database.Store) *AuthHandler {
	return &AuthHandler{store: store}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
		return
	}

	user, err := h.store.GetUserByEmail(req.Email)
	if err != nil {
		sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
		return
	}

	token, err := middleware.GenerateToken(user)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
		return
	}

	// Set cookie as well
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(48 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	_ = h.store.AddActivity(user.Name, "User Login", "Logged into research portal")

	sendJSON(w, http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user.ToSafe(),
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Name, email and password are required"})
		return
	}

	if req.Role == "" {
		req.Role = "Researcher"
	}
	if req.Department == "" {
		req.Department = "Interdisciplinary Sciences"
	}
	if req.Institution == "" {
		req.Institution = "Global Institute of Advanced Sciences"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to process password"})
		return
	}

	newUser := &models.User{
		Name:        req.Name,
		Email:       req.Email,
		Password:    string(hash),
		Role:        req.Role,
		Department:  req.Department,
		Institution: req.Institution,
		Avatar:      "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?auto=format&fit=crop&w=250&q=80",
	}

	if err := h.store.CreateUser(newUser); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	token, err := middleware.GenerateToken(newUser)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(48 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	_ = h.store.AddActivity(newUser.Name, "Account Registered", "Created new "+newUser.Role+" account")

	sendJSON(w, http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  newUser.ToSafe(),
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})
	sendJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	if claims == nil {
		sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	user, err := h.store.GetUserByID(claims.UserID)
	if err != nil {
		sendJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	sendJSON(w, http.StatusOK, user.ToSafe())
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
