package model

import (
	"fmt"
	"net/mail"
	"time"
)

type Order struct {
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

func (o Order) Validate() error {
	if o.HotelID == "" {
		return fmt.Errorf("empty hotel_id")
	}
	if o.RoomID == "" {
		return fmt.Errorf("empty hotel_id")
	}

	_, err := mail.ParseAddress(o.UserEmail)
	if err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	n := time.Now()
	if o.From.Before(n) {
		return fmt.Errorf("from in past: %v", o.From)
	}
	if o.To.Before(n) {
		return fmt.Errorf("from in past: %v", o.From)
	}
	if o.From.After(o.To) {
		return fmt.Errorf("from %v is after to %v", o.From, o.To)
	}
	return nil
}
