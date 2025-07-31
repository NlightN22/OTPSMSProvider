package service

import (
	"time"

	logger "github.com/NlightN22/OTPSMSProvider/pkg"
	storage "github.com/NlightN22/OTPSMSProvider/storage"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
)

// TotpService implements OTPService.
type TotpService struct {
	store    storage.Storage
	issuer   string
	period   uint
	digits   otp.Digits
	algo     otp.Algorithm
	skew     uint
	interval time.Duration
	log      *zap.SugaredLogger
	notifier Notifier
}

// NewTotpService constructs TotpService with its own logger.
func NewTotpService(
	store storage.Storage,
	issuer string,
	period uint,
	digits otp.Digits,
	algo otp.Algorithm,
	skew uint,
	interval time.Duration,
	notifier Notifier,
) *TotpService {

	svcLog := logger.New("TotpService")

	return &TotpService{
		store:    store,
		issuer:   issuer,
		period:   period,
		digits:   digits,
		algo:     algo,
		skew:     skew,
		interval: interval,
		log:      svcLog,
		notifier: notifier,
	}
}

func (s *TotpService) CanSend(key string) (bool, time.Duration) {
	s.log.Infow("CanSend called", "phone", key)
	if last, ok := s.store.GetLastSend(key); ok {
		since := time.Since(last)
		if since < s.interval {
			s.log.Infow("Rate limit", "since", since, "interval", s.interval)
			return false, s.interval - since
		}
	}
	return true, 0
}

func (s *TotpService) GenerateCode(key string) (string, error) {
	s.log.Infow("GenerateCode called", "phone", key)
	secret, ok := s.store.GetSecret(key)
	if !ok {
		opt := totp.GenerateOpts{
			Issuer:      s.issuer,
			AccountName: key,
			Period:      s.period,
			Digits:      s.digits,
			Algorithm:   s.algo,
		}
		token, err := totp.Generate(opt)
		if err != nil {
			s.log.Errorw("TOTP.Generate error", "err", err)
			return "", err
		}
		secret = token.Secret()
		s.store.SaveSecret(key, secret)
	}
	s.store.SaveLastSend(key, time.Now())

	code, err := totp.GenerateCodeCustom(secret, time.Now(),
		totp.ValidateOpts{Period: s.period, Skew: s.skew, Digits: s.digits, Algorithm: s.algo})
	if err != nil {
		s.log.Errorw("GenerateCodeCustom error", "err", err)
		return "", err
	}
	s.log.Debugw("Code generated", "code", code)

	s.log.Debugw("Starting send code", "phone", key)

	err = s.notifier.Send(key, code)
	if err != nil {
		s.log.Errorw("Send code error", "err", err)
		return "", err
	}

	s.log.Infow("Code generated and sended", "code", code)
	return code, nil
}

func (s *TotpService) ValidateCode(key, code string) bool {
	s.log.Infow("ValidateCode called", "phone", key, "code", code)
	secret, ok := s.store.GetSecret(key)
	if !ok {
		s.log.Warnw("No secret for phone", "phone", key)
		return false
	}
	valid, _ := totp.ValidateCustom(code, secret, time.Now(),
		totp.ValidateOpts{Period: s.period, Skew: s.skew, Digits: s.digits, Algorithm: s.algo})
	s.log.Infow("Validation result", "valid", valid)
	return valid
}
