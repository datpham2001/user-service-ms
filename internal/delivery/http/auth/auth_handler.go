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

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url := h.authService.GetGoogleAuthUrl()
	response.Redirect(c, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	var req dto.GoogleCallbackRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	loginResponse, err := h.authService.ProcessGoogleCallback(c.Request.Context(), &req)
	if err != nil {
		response.ErrorService(c, err)
		return
	}

	response.Success(c, loginResponse)
}
