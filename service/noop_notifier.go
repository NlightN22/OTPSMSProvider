package service

import (
	logger "github.com/NlightN22/OTPSMSProvider/pkg"
	"go.uber.org/zap"
)

type NoopNotifier struct {
	log *zap.SugaredLogger
}

func NewNoopNotifier() *NoopNotifier {
	l := logger.New("NoopNotifier")
	return &NoopNotifier{log: l}
}

func (n *NoopNotifier) Send(to, msg string) error {
	n.log.Infow("DEBUG mode — skipping SMS send", "to", to, "message", msg)
	return nil
}
