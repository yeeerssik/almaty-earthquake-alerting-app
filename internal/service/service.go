package service

type Service interface {
	SendMessage(message string) error
}
