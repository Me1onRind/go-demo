package async

import "context"

type Sender interface {
	SendJob(ctx context.Context, data any) error
}

func NewSender() Sender {
	return &sender{}
}

type sender struct {
}

func (s *sender) SendJob(ctx context.Context, data any) error {
	return nil
}
