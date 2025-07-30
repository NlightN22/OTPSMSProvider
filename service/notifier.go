package service

type Notifier interface {
	Send(to, text string) error
}
