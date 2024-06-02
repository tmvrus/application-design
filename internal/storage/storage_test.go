package storage

import (
	"testing"
	"time"
)

func Test_Night(t *testing.T) {
	t.Parallel()

	tt := []struct {
		from, to time.Time
		expected int
	}{
		{
			from:     time.Now(),
			to:       time.Now().Add(time.Hour * 24),
			expected: 1,
		},
		{
			from:     time.Now(),
			to:       time.Now().Add(time.Hour * 24 * 2),
			expected: 2,
		},
	}

	for i := range tt {
		if datesToBook(tt[i].from, tt[i].to) != tt[i].expected {
			t.Errorf("unexpected result on stio %d", i)
		}
	}
}

func Test_betweenOrSame(t *testing.T) {
	t.Parallel()

	tt := []struct {
		from, to, in time.Time
		expected     bool
	}{
		{
			from:     time.Now(),
			to:       time.Now().Add(time.Hour),
			in:       time.Now().Add(time.Minute),
			expected: true,
		},
		{
			to:       time.Now(),
			from:     time.Now().Add(time.Hour),
			in:       time.Now().Add(time.Minute),
			expected: false,
		},
		{
			to:       time.Now().Truncate(time.Hour),
			from:     time.Now().Add(time.Hour),
			in:       time.Now().Truncate(time.Hour),
			expected: true,
		},
	}

	for i := range tt {
		if tt[i].expected != betweenOrSame(tt[i].from, tt[i].to, tt[i].in) {
			t.Fatalf("invalid result on step %d", i)
		}
	}

}
