// service/smsc_service_alt.go
package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	logger "github.com/NlightN22/OTPSMSProvider/pkg"
	"go.uber.org/zap"
)

type SMSCService struct {
	login         string
	password      string
	prefixMessage string
	client        *http.Client
	log           *zap.SugaredLogger
}

func NewSMSCService(login, password, prefixMessage string) *SMSCService {
	svcLog := logger.New("SMSCService")

	return &SMSCService{
		login:         login,
		password:      password,
		prefixMessage: prefixMessage,
		client:        &http.Client{Timeout: 10 * time.Second},
		log:           svcLog,
	}
}

func (s *SMSCService) Send(phone, code string) error {
	message := s.prefixMessage + code
	s.log.Infow("smsc: message", message)
	params := url.Values{
		"login":  {s.login},
		"psw":    {s.password},
		"phones": {phone},
		"mes":    {message},
		"fmt":    {"3"},
	}
	endpoint := "https://smsc.ru/sys/send.php"
	request := endpoint + "?" + params.Encode()
	s.log.Debug("smsc: request: ", request)
	resp, err := s.client.Get(request)
	s.log.Debug("smsc: response: ", resp)
	if err != nil {
		return fmt.Errorf("smsc send request error: %w", err)
	}
	defer resp.Body.Close()

	var apiResp struct {
		ID        int    `json:"id"`
		Cnt       int    `json:"cnt"`
		ErrorCode int    `json:"error_code"`
		Error     string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("smsc response parse error: %w", err)
	}
	if apiResp.ErrorCode != 0 {
		return fmt.Errorf("smsc API error %d: %s", apiResp.ErrorCode, apiResp.Error)
	}
	s.log.Infof("smsc: sent id=%d, parts=%d", apiResp.ID, apiResp.Cnt)
	return nil
}
