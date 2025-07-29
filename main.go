// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/NlightN22/OTPSMSProvider/docs" // Используй свой модульный путь после swag init
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var secrets = map[string]string{} // In-memory store (replace with Redis/db in production)
var lastSend = map[string]time.Time{}

// SendRequest represents request for /send
type SendRequest struct {
	Phone string `json:"phone"`
}

// VerifyRequest represents request for /verify
type VerifyRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

// @title TOTP SMS Auth API
// @version 1.0
// @description API для генерации и проверки TOTP-кодов, отправляемых по SMS
// @host localhost:8080
// @BasePath /
func main() {
	bindAddr := os.Getenv("TOTP_BIND")
	if bindAddr == "" {
		bindAddr = ":8080"
	}

	whiteList := strings.Split(os.Getenv("TOTP_WHITELIST"), ",")
	intervalStr := os.Getenv("TOTP_INTERVAL")
	interval := 30
	if i, err := time.ParseDuration(intervalStr + "s"); err == nil {
		interval = int(i.Seconds())
	}

	period := 60
	if val := os.Getenv("TOTP_PERIOD"); val != "" {
		if p, err := time.ParseDuration(val + "s"); err == nil {
			period = int(p.Seconds())
		}
	}

	digits := otp.DigitsSix
	if val := os.Getenv("TOTP_DIGITS"); val == "8" {
		digits = otp.DigitsEight
	}

	algo := otp.AlgorithmSHA1
	switch strings.ToUpper(os.Getenv("TOTP_ALGO")) {
	case "SHA256":
		algo = otp.AlgorithmSHA256
	case "SHA512":
		algo = otp.AlgorithmSHA512
	}

	skew := 0
	if val := os.Getenv("TOTP_SKEW"); val != "" {
		fmt.Sscanf(val, "%d", &skew)
	}

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		if len(whiteList) > 0 && whiteList[0] != "" {
			ip := c.ClientIP()
			allowed := false
			for _, w := range whiteList {
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

	r.POST("/send", func(c *gin.Context) { sendHandler(c, interval, period, digits, algo, skew) })
	r.POST("/verify", func(ctx *gin.Context) { verifyHandler(ctx, period, digits, algo, skew) })
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(r.Run(bindAddr))
}

func logger(service string, format string, a ...any) (n int, err error) {
	prefix := fmt.Sprintf("[GIN] %s %s ", time.Now().Format("2006/01/02 - 15:04:05"), service)
	return fmt.Fprintf(gin.DefaultWriter, prefix+format+"\n", a...)
}

// sendHandler godoc
// @Summary Сгенерировать и отправить TOTP-код по номеру телефона
// @Accept json
// @Produce json
// @Param data body SendRequest true "Phone"
// @Success 200 {string} string "Код отправлен"
// @Failure 400 {string} string "Ошибка"
// @Router /send [post]
func sendHandler(c *gin.Context, interval int, period int, digits otp.Digits, algo otp.Algorithm, skew int) {
	var req SendRequest
	if err := c.BindJSON(&req); err != nil || req.Phone == "" {
		c.String(http.StatusBadRequest, "Некорректный запрос")
		return
	}
	if last, ok := lastSend[req.Phone]; ok && time.Since(last) < time.Duration(interval)*time.Second {
		remainingTime := time.Duration(interval)*time.Second - time.Since(last)
		c.String(http.StatusTooManyRequests, "Подождите перед повторной отправкой %s", remainingTime)
		return
	}

	key := req.Phone
	secret := secrets[key]
	if secret == "" {
		opt := totp.GenerateOpts{
			Issuer:      "TOTPService",
			AccountName: key,
			Period:      uint(period),
			Digits:      digits,
			Algorithm:   algo,
		}
		token, err := totp.Generate(opt)
		if err != nil {
			c.String(http.StatusInternalServerError, "Ошибка генерации секрета")
			return
		}
		secret = token.Secret()
		secrets[key] = secret
	}

	opts := totp.ValidateOpts{
		Period:    uint(period),
		Skew:      uint(skew),
		Digits:    digits,
		Algorithm: algo,
	}
	code, err := totp.GenerateCodeCustom(secret, time.Now(), opts)
	if err != nil {
		logger("[SMS]", "Ошибка генерации кода: %v", err)
		c.String(http.StatusInternalServerError, "Ошибка генерации кода")
		return
	}
	lastSend[req.Phone] = time.Now()

	// fmt.Fprintf(gin.DefaultWriter, "[SMS] %s Телефон: %s, Код: %s\n", time.Now().Format("2006/01/02 - 15:04:05"), key, code) // Replace with real SMS API
	logger("[SMS]", "Создан Телефон: %s, Код: %s, Секрет: %s", key, code, secret)
	c.String(http.StatusOK, "Код отправлен")
}

// verifyHandler godoc
// @Summary Проверить введённый TOTP-код
// @Accept json
// @Produce json
// @Param data body VerifyRequest true "Проверка кода"
// @Success 200 {string} string "Код верен"
// @Failure 401 {string} string "Неверный код"
// @Router /verify [post]
func verifyHandler(c *gin.Context, period int, digits otp.Digits, algo otp.Algorithm, skew int) {
	var req VerifyRequest
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Ошибка запроса")
		return
	}
	secret := secrets[req.Phone]

	valid, _ := totp.ValidateCustom(req.Code, secret, time.Now(), totp.ValidateOpts{
		Period:    uint(period),
		Skew:      uint(skew),
		Digits:    digits,
		Algorithm: algo,
	})

	logger("[SMS]", "Проверен Телефон: %s, Код: %s, Cекрет: %s, Результат: %t", req.Phone, req.Code, secret, valid)

	if secret == "" || !valid {
		c.String(http.StatusUnauthorized, "Неверный код")
		return
	}

	c.String(http.StatusOK, "Код верен")
}
