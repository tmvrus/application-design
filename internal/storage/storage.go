package storage

import (
	"context"
	"sync"
	"time"

	"applicationDesignTest/internal/model"
)

const someReasonableSize = 10

type Storage struct {
	orders       []model.Order
	availability []model.RoomAvailability
	lock         sync.Mutex
}

func New() *Storage {
	return &Storage{
		orders: make([]model.Order, 0, someReasonableSize),
		availability: []model.RoomAvailability{
			{"reddison", "lux", date(2024, 6, 1), 1},
			{"reddison", "lux", date(2024, 6, 2), 1},
			{"reddison", "lux", date(2024, 6, 3), 1},
			{"reddison", "lux", date(2024, 6, 4), 1},
			{"reddison", "lux", date(2024, 6, 5), 0},
		},
	}
}

func (s *Storage) GetAvailability(_ context.Context, o model.Order) ([]model.RoomAvailability, error) {
	toBook := datesToBook(o.From, o.To)
	available := make([]model.RoomAvailability, 0, len(toBook))

	for _, v := range s.availability {
		if v.Quota < 1 || v.HotelID != o.HotelID || v.RoomID != o.RoomID || !betweenOrSame(o.From, o.To, v.Date) {
			continue
		}
		r, ok := toBook[day(v.Date)]
		if !ok {
			continue
		}
		if r {
			return nil, model.ErrInventoryCorrupted
		}

		toBook[day(v.Date)] = true
		available = append(available, v)
	}

	for _, avl := range toBook {
		if !avl {
			return nil, model.ErrNoRoomsAvailable
		}
	}

	return available, nil
}

func (s *Storage) SaveOrder(_ context.Context, order model.Order, avl []model.RoomAvailability) error {
	for i := range avl {
		avl[i].Quota--
	}
	s.orders = append(s.orders, order)
	return nil
}

func betweenOrSame(from, to, d time.Time) bool {
	return (d.After(from) && d.Before(to)) || (from == d || to == d)
}

func datesToBook(from, to time.Time) map[time.Time]bool {
	res := map[time.Time]bool{}

	for d := day(from); !d.After(day(to)); d = d.AddDate(0, 0, 1) {
		res[d] = false
	}
	return res
}

func day(t time.Time) time.Time {
	return t.Truncate(time.Hour * 24)
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
