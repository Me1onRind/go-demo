package async

import "context"

type Sender interface {
	Send(ctx context.Context, data any) error
}

func NewSender() Sender {
	return &sender{}
}

type sender struct {
}

func (s *sender) Send(ctx context.Context, data any) error {
	return nil
}
