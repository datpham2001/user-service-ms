package auth

import (
	"errors"
	"net/http"
	"strconv"

	dto "github.com/datpham/user-service-ms/internal/dto/request"
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

func (h *AuthHandler) Logout(c *gin.Context) {
	response.Success(c, response.OK)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	loginResponse, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		response.ErrorService(c, err)
		return
	}

	response.Success(c, loginResponse)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	if err := h.authService.ForgotPassword(c.Request.Context(), &req); err != nil {
		response.ErrorService(c, err)
		return
	}

	response.Success(c, response.OK)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		response.Error(c, http.StatusBadRequest, errors.New("token is required"))
		return
	}

	token, err := strconv.Atoi(tokenStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errors.New("invalid token"))
		return
	}

	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), token, &req); err != nil {
		response.ErrorService(c, err)
		return
	}

	response.Success(c, response.OK)
}
