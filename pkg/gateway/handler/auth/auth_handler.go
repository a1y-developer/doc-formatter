package auth

import (
	"net/http"

	"github.com/a1y/doc-formatter/pkg/gateway/domain/request"
	"github.com/a1y/doc-formatter/pkg/gateway/handler"
	"github.com/gin-gonic/gin"
)

// @Id				Signup
// @Summary		Signup
// @Description	Create a new user account
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			body	body		request.SignupRequest							true	"Signup payload"
// @Success		201		{object}	handler.Response{data=response.SignUpResponse}	"Success"
// @Failure		400		{object}	error											"Bad Request"
// @Failure		401		{object}	error											"Unauthorized"
// @Failure		429		{object}	error											"Too Many Requests"
// @Failure		404		{object}	error											"Not Found"
// @Failure		500		{object}	error											"Internal Server Error"
// @Router			/api/v1/auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var requestPayload request.SignupRequest
	if err := requestPayload.Decode(c); err != nil {
		handler.HandleResultWithStatus(c, err, nil, http.StatusBadRequest)
		return
	}

	// Validate request payload
	if err := requestPayload.Validate(); err != nil {
		handler.HandleResultWithStatus(c, err, nil, http.StatusBadRequest)
		return
	}

	response, err := h.authManager.Signup(c.Request.Context(), requestPayload)
	if err != nil {
		handler.HandleResult(c, err, response)
		return
	}
	handler.HandleResultWithStatus(c, nil, response, http.StatusCreated)
}

// @Id				Login
// @Summary		Login
// @Description	Login user and return JWT token
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			body	body		request.LoginRequest							true	"Login payload"
// @Success		200		{object}	handler.Response{data=response.LoginResponse}	"Success"
// @Failure		400		{object}	error											"Bad Request"
// @Failure		401		{object}	error											"Unauthorized"
// @Failure		429		{object}	error											"Too Many Requests"
// @Failure		404		{object}	error											"Not Found"
// @Failure		500		{object}	error											"Internal Server Error"
// @Router			/api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var requestPayload request.LoginRequest
	if err := requestPayload.Decode(c); err != nil {
		handler.HandleResultWithStatus(c, err, nil, http.StatusBadRequest)
		return
	}

	// Validate request payload
	if err := requestPayload.Validate(); err != nil {
		handler.HandleResultWithStatus(c, err, nil, http.StatusBadRequest)
		return
	}

	response, err := h.authManager.Login(c.Request.Context(), requestPayload)
	if err != nil {
		if err.Error() == "invalid credentials" {
			handler.HandleResultWithStatus(c, err, nil, http.StatusUnauthorized)
			return
		}
		handler.HandleResult(c, err, nil)
		return
	}

	handler.HandleResult(c, err, response)
}
