package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	service "github.com/NlightN22/OTPSMSProvider/service"
)

// SendRequest represents request for /send endpoint.
type SendRequest struct {
	Phone string `json:"phone" binding:"required,e164"`
}

// VerifyRequest represents request for /verify endpoint.
type VerifyRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// API groups TOTP handlers. Comments in English.
type API struct {
	svc service.OTPService
}

// NewAPI creates a new API instance.
func NewAPI(svc service.OTPService) *API {
	return &API{svc: svc}
}

// RegisterRoutes attaches routes and Swagger UI to the router.
func (a *API) RegisterRoutes(r *gin.Engine) {
	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// send TOTP code
	r.POST("/send", a.send)

	// verify TOTP code
	r.POST("/verify", a.verify)
}

// send handles code generation and SMS dispatch.
// @Summary Generate and send TOTP code via SMS
// @Description Generates TOTP for given phone and records send time
// @Accept json
// @Produce json
// @Param data body SendRequest true "Phone"
// @Success 200 {string} string "Code sent"
// @Failure 400 {string} string "Invalid request"
// @Failure 429 {string} string "Too many requests"
// @Router /send [post]
func (a *API) send(c *gin.Context) {
	var req SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}
	ok, wait := a.svc.CanSend(req.Phone)
	if !ok {
		c.String(http.StatusTooManyRequests, "Please wait %s", wait)
		return
	}
	if _, err := a.svc.GenerateCode(req.Phone); err != nil {
		c.String(http.StatusInternalServerError, "Code generation error")
		return
	}
	c.String(http.StatusOK, "Code sent")
}

// verify handles code validation.
// @Summary Validate TOTP code
// @Description Checks provided TOTP code for validity
// @Accept json
// @Produce json
// @Param data body VerifyRequest true "Code"
// @Success 200 {string} string "Code valid"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Invalid code"
// @Router /verify [post]
func (a *API) verify(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	if a.svc.ValidateCode(req.Phone, req.Code) {
		c.String(http.StatusOK, "Code valid")
	} else {
		c.String(http.StatusUnauthorized, "Invalid code")
	}
}
