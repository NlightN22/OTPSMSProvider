package service

import (
	"testing"
	"time"

	"github.com/pquerna/otp"
)

// stubStorage implements storage.Storage
type stubStorage struct {
	secret    string
	hasSecret bool
	lastSend  time.Time
	hasLast   bool
}

func (s *stubStorage) GetSecret(key string) (string, bool) {
	return s.secret, s.hasSecret
}
func (s *stubStorage) SaveSecret(key, secret string) {
	s.secret = secret
	s.hasSecret = true
}
func (s *stubStorage) GetLastSend(key string) (time.Time, bool) {
	return s.lastSend, s.hasLast
}
func (s *stubStorage) SaveLastSend(key string, t time.Time) {
	s.lastSend = t
	s.hasLast = true
}

// stubNotifier implements Notifier
type stubNotifier struct {
	sentTo   string
	sentCode string
	err      error
}

func (n *stubNotifier) Send(to, code string) error {
	n.sentTo = to
	n.sentCode = code
	return n.err
}

func TestGenerateAndValidate(t *testing.T) {
	store := &stubStorage{}
	notifier := &stubNotifier{}
	svc := NewTotpService(store, "test", 30, otp.DigitsSix, otp.AlgorithmSHA1, 1, time.Second, notifier)

	code, err := svc.GenerateCode("123")
	if err != nil {
		t.Fatalf("GenerateCode error: %v", err)
	}
	if notifier.sentTo != "123" {
		t.Errorf("Notifier sent to = %s; want %s", notifier.sentTo, "123")
	}
	if notifier.sentCode != code {
		t.Errorf("Notifier code = %s; want %s", notifier.sentCode, code)
	}

	// Validate correct code
	if !svc.ValidateCode("123", code) {
		t.Errorf("ValidateCode returned false for correct code")
	}
	// Validate wrong code
	if svc.ValidateCode("123", "000000") {
		t.Errorf("ValidateCode returned true for wrong code")
	}
}

func TestCanSend(t *testing.T) {
	now := time.Now()
	store := &stubStorage{lastSend: now.Add(-500 * time.Millisecond), hasLast: true}
	notifier := &stubNotifier{}
	svc := NewTotpService(store, "test", 30, otp.DigitsSix, otp.AlgorithmSHA1, 1, 1*time.Second, notifier)

	ok, wait := svc.CanSend("any")
	if ok {
		t.Errorf("CanSend = true; want false due to rate limit")
	}
	if wait <= 0 || wait > 1*time.Second {
		t.Errorf("Wait duration = %v; want >0 and <=1s", wait)
	}

	// No last send
	store2 := &stubStorage{hasLast: false}
	svc2 := NewTotpService(store2, "test", 30, otp.DigitsSix, otp.AlgorithmSHA1, 1, 1*time.Second, notifier)
	ok2, wait2 := svc2.CanSend("any")
	if !ok2 || wait2 != 0 {
		t.Errorf("CanSend = %v, wait = %v; want true,0", ok2, wait2)
	}
}
