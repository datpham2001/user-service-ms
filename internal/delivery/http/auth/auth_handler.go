package auth

import (
	"net/http"

	dto "github.com/datpham/user-service-ms/internal/dto/auth"
	"github.com/datpham/user-service-ms/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService IAuthService
}

func New(authService IAuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req dto.UserSignupRequest
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

// func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
// 	url := h.authService.GetAuthURL()
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// }

// func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
// 	code := r.URL.Query().Get("code")
// 	if code == "" {
// 		http.Error(w, "Code not found", http.StatusBadRequest)
// 		return
// 	}

// 	userInfo, err := h.authService.GetUserInfo(code)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Here you would typically:
// 	// 1. Create or update user in your database
// 	// 2. Generate a session token or JWT
// 	// 3. Set cookies or return tokens

// 	// For now, we'll just return the user info as JSON
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(userInfo)
// }

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	loginResponse, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.ErrorService(c, err)
		return
	}

	response.Success(c, loginResponse)
}
