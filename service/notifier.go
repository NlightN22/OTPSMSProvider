package service

type Notifier interface {
	Send(phone, code string) error
}
