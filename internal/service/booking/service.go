package booking

import (
	"context"
	"fmt"

	"applicationDesignTest/internal/model"
)

type roomStorage interface {
	GetAvailability(context.Context, model.Order) ([]model.RoomAvailability, error)
	SaveOrder(context.Context, model.Order, []model.RoomAvailability) error
}

type asyncMailerQueue interface {
	Send(string, string)
}

type Service struct {
	storage roomStorage
	mailer  asyncMailerQueue
}

func New(s roomStorage, q asyncMailerQueue) Service {
	return Service{storage: s, mailer: q}
}

func (s Service) Book(ctx context.Context, o model.Order) error {
	// quire rooms for booking
	res, err := s.storage.GetAvailability(ctx, o)
	if err != nil {
		return fmt.Errorf("get availability: %w", err)
	}

	// confirm booking after check
	// - money
	// - user response
	// - ...

	// commit acquired rooms
	err = s.storage.SaveOrder(ctx, o, res)
	if err != nil {
		return fmt.Errorf("save order: %w", err)
	}

	// do not affect booking process
	s.mailer.Send(o.UserEmail, "Booked! Waiting for you!")
	return nil
}
