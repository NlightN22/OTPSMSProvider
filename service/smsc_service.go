package service

import (
	"fmt"

	logger "github.com/NlightN22/OTPSMSProvider/pkg"
	"github.com/koorgoo/smsc"
	"go.uber.org/zap"
)

type SMSCService struct {
	client *smsc.Client
	log    *zap.SugaredLogger
}

func NewSMSCService(login string,
	password string,
	URL string,
) *SMSCService {
	svcLog, err := logger.New("SMSCService")
	if err != nil {
		panic(err)
	}
	smscClient, err := smsc.New(smsc.Config{
		Login:    login,
		Password: password,
		URL:      URL,
	})
	if err != nil {
		svcLog.Fatalw("smsc: Can not initialize smscClient", err)
	}
	return &SMSCService{client: smscClient, log: svcLog}
}

func (s *SMSCService) Send(to, text string) error {
	s.log.Debugw("Start sending SMS")
	result, err := s.client.Send(text, []string{to})
	if err != nil {
		s.log.Errorw("smsc: send error:", err)
		return fmt.Errorf("smsc: send error: %+v", err)
	}
	if result == nil || result.ID == 0 {
		s.log.Errorw("smsc: invalid result:", result)
		return fmt.Errorf("smsc: invalid result: %+v", result)
	}

	s.log.Infow("SMS sucessfully sended")
	return nil
}
