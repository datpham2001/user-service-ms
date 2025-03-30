package auth

import (
	"encoding/json"
	"net/http"

	"github.com/datpham/user-service-ms/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	if err := h.authService.Signup(c.Request.Context(), &req); err != nil {
		response.ErrorService(c, err)
		return
	}

	response.Created(c, response.CREATED)
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.authService.GetAuthURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	userInfo, err := h.authService.GetUserInfo(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Here you would typically:
	// 1. Create or update user in your database
	// 2. Generate a session token or JWT
	// 3. Set cookies or return tokens

	// For now, we'll just return the user info as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

// Login handles standard user login requests.
// Note: This uses gin.Context because we'll call it from Gin routes.
func (h *AuthHandler) Login(c *gin.Context) {
	// TODO: Bind JSON request body (e.g., email, password)
	// TODO: Call authService.Login(credentials)
	// TODO: Handle errors
	// TODO: Generate token and return appropriate response
	c.JSON(http.StatusOK, gin.H{"message": "handler login placeholder"})
}

// Signup handles user registration requests.
// Note: This uses gin.Context because we'll call it from Gin routes.
func (h *AuthHandler) Signup(c *gin.Context) {
	// TODO: Bind JSON request body (e.g., username, email, password)
	// TODO: Call authService.Signup(userInfo)
	// TODO: Handle errors
	// TODO: Return appropriate response (e.g., created user ID)
	c.JSON(http.StatusCreated, gin.H{"message": "handler signup placeholder"})
}
