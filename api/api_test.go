package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/NlightN22/OTPSMSProvider/validator"
	"github.com/gin-gonic/gin"
)

func init() {
	validator.RegisterCustomValidations()
}

type stubService struct {
	canSend bool
	wait    time.Duration
	genErr  error
	valid   bool
	code    string
}

func (s *stubService) CanSend(key string) (bool, time.Duration) { return s.canSend, s.wait }
func (s *stubService) GenerateCode(key string) (string, error)  { return s.code, s.genErr }
func (s *stubService) ValidateCode(key, code string) bool       { return s.valid && code == s.code }

func performRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestSendEndpoint_RateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &stubService{canSend: false, wait: 5 * time.Second}
	a := NewAPI(stub)
	router := gin.New()
	a.RegisterRoutes(router)

	w := performRequest(router, "POST", "/send", `{"phone":"+1234567890"}`)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("status = %d; want %d", w.Code, http.StatusTooManyRequests)
	}
}

func TestSendEndpoint_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &stubService{canSend: true, code: "123456"}
	a := NewAPI(stub)
	router := gin.New()
	a.RegisterRoutes(router)

	w := performRequest(router, "POST", "/send", `{"phone":"+1234567890"}`)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d; want %d", w.Code, http.StatusOK)
	}
}

func TestSendEndpoint_GenerationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &stubService{canSend: true, genErr: fmt.Errorf("fail")}
	a := NewAPI(stub)
	router := gin.New()
	a.RegisterRoutes(router)

	w := performRequest(router, "POST", "/send", `{"phone":"+1234567890"}`)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d; want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestVerifyEndpoint_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &stubService{valid: true, code: "123456"}
	a := NewAPI(stub)
	router := gin.New()
	a.RegisterRoutes(router)

	w := performRequest(router, "POST", "/verify", `{"phone":"+1234567890","code":"123456"}`)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d; want %d", w.Code, http.StatusOK)
	}
}

func TestVerifyEndpoint_Invalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &stubService{valid: false, code: "123456"}
	a := NewAPI(stub)
	router := gin.New()
	a.RegisterRoutes(router)

	w := performRequest(router, "POST", "/verify", `{"phone":"+1234567890","code":"wrong"}`)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d; want %d", w.Code, http.StatusUnauthorized)
	}
}
