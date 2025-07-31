package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWhitelistMiddleware_AllowsAllowedIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mw := NewWhitelistMiddleware([]string{"127.0.0.1"})

	router := gin.New()
	router.Use(mw.Handler())
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}

func TestWhitelistMiddleware_BlocksOtherIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mw := NewWhitelistMiddleware([]string{"10.0.0.1"})

	router := gin.New()
	router.Use(mw.Handler())
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.100:5678"
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestWhitelistMiddleware_AllowsLoopbackEvenIfNotListed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mw := NewWhitelistMiddleware([]string{"10.0.0.1"})

	router := gin.New()
	router.Use(mw.Handler())
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:9999"
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}
