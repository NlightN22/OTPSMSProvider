// main.go
package main

import (
	"log"
	"net/http"
	"strings"

	api "github.com/NlightN22/OTPSMSProvider/api"
	config "github.com/NlightN22/OTPSMSProvider/config"
	_ "github.com/NlightN22/OTPSMSProvider/docs"
	logger "github.com/NlightN22/OTPSMSProvider/pkg"
	service "github.com/NlightN22/OTPSMSProvider/service"
	storage "github.com/NlightN22/OTPSMSProvider/storage"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
)

// @title TOTP SMS Auth API
// @version 1.0
// @description API для генерации и проверки TOTP-кодов, отправляемых по SMS
// @host localhost:8080
// @BasePath /
func main() {

	cfg, _ := config.LoadConfig()

	mainLog, err := logger.New("main")
	if err != nil {
		panic(err)
	}
	defer mainLog.Sync()

	// init service logger
	svcLog, err := logger.New("service")
	if err != nil {
		mainLog.Fatalw("failed to init service logger", "err", err)
	}

	store := storage.NewMemoryStorage()

	algo := otp.AlgorithmSHA1
	switch strings.ToUpper(cfg.Algorithm) {
	case "SHA256":
		algo = otp.AlgorithmSHA256
	case "SHA512":
		algo = otp.AlgorithmSHA512
	}

	periodSeconds := uint(cfg.Period.Seconds())
	svc := service.NewTotpService(
		store,
		"TOTP Service",
		periodSeconds,
		otp.Digits(cfg.Digits),
		algo, uint(cfg.Skew),
		cfg.Interval,
		svcLog,
	)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		if len(cfg.WhiteList) > 0 && cfg.WhiteList[0] != "" {
			ip := c.ClientIP()
			allowed := false
			for _, w := range cfg.WhiteList {
				if strings.TrimSpace(w) == ip {
					allowed = true
					break
				}
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "IP запрещён"})
				return
			}
		}
		c.Next()
	})

	api := api.NewAPI(svc)
	api.RegisterRoutes(r)

	log.Fatal(r.Run(cfg.BindAddr))
}
