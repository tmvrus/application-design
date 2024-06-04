package notify

import "log/slog"

type Service struct{}

func New() Service {
	return Service{}
}

func (Service) Send(address string, body string) {
	slog.Info("got a email", "to", address, "content", body)
}
