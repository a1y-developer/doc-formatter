package auth

import (
	"net/http"

	"github.com/a1y/doc-formatter/internal/gateway/domain/request"
	"github.com/gin-gonic/gin"
)

// @Summary Signup
// @Description Signup user account
// @Tags auth
// @Accept json
// @Produce json
// @Param body body request.SignupRequest true "Signup payload"
// @Success 201 {object} response.SignUpResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req request.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authManager.Signup(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": resp.UserID,
	})
}

// @Summary Login
// @Description Login user account, return JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param body body request.LoginRequest true "Login payload"
// @Success 200 {object} response.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authManager.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": resp.AccessToken,
		"expiry_unix":  resp.ExpiryUnix,
	})
}
