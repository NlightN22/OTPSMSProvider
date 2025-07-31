// main.go
package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	api "github.com/NlightN22/OTPSMSProvider/api"
	config "github.com/NlightN22/OTPSMSProvider/config"
	_ "github.com/NlightN22/OTPSMSProvider/docs"
	"github.com/NlightN22/OTPSMSProvider/middleware"
	logger "github.com/NlightN22/OTPSMSProvider/pkg"
	service "github.com/NlightN22/OTPSMSProvider/service"
	storage "github.com/NlightN22/OTPSMSProvider/storage"
	"github.com/NlightN22/OTPSMSProvider/validator"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
)

// @title TOTP SMS Auth API
// @version 1.0
// @description API for generating and verifying OTP codes sent via SMS
// @host localhost:8080
// @BasePath /
func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	if err := logger.Init(cfg.LogLevel); err != nil {
		panic(fmt.Errorf("logger init: %w", err))
	}

	mainLog := logger.New("main")

	defer mainLog.Sync()

	mainLog.Infow("Loaded configuration", "config", cfg)

	store := storage.NewMemoryStorage()

	algo := otp.AlgorithmSHA1
	switch strings.ToUpper(cfg.Algorithm) {
	case "SHA256":
		algo = otp.AlgorithmSHA256
	case "SHA512":
		algo = otp.AlgorithmSHA512
	}

	var notifier service.Notifier
	if cfg.Debug {
		notifier = service.NewNoopNotifier()
	} else {
		notifier = service.NewSMSCService(cfg.SMSC.Login, cfg.SMSC.Password, cfg.PrefixText)
	}

	svc := service.NewTotpService(
		store,
		"TOTP Service",
		uint(cfg.Period),
		otp.Digits(cfg.Digits),
		algo,
		uint(cfg.Skew),
		time.Duration(cfg.Interval),
		notifier,
	)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	whitelistMw := middleware.NewWhitelistMiddleware(cfg.WhiteList)
	r.Use(whitelistMw.Handler())

	validator.RegisterCustomValidations()

	api := api.NewAPI(svc)
	api.RegisterRoutes(r)

	mainLog.Fatal(r.Run(cfg.BindAddr))
}
