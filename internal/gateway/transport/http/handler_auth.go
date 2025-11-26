package http

import (
    "net/http"

    "github.com/a1y/ai-doc-formatter/internal/gateway/clients"

    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authClient *clients.AuthClient
}

func NewAuthHandler(ac *clients.AuthClient) *AuthHandler {
    return &AuthHandler{authClient: ac}
}

type signupRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// @Summary Signup
// @Description Đăng ký user mới
// @Tags auth
// @Accept json
// @Produce json
// @Param body body signupRequest true "Signup payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
    var req signupRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.authClient.Signup(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "user_id": resp.UserId,
    })
}

// @Summary Login
// @Description Đăng nhập, trả về JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param body body loginRequest true "Login payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    var req loginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.authClient.Login(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        // ở đây đơn giản coi mọi lỗi là 401
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "access_token": resp.AccessToken,
        "expiry_unix":  resp.ExpiryUnix,
    })
}
