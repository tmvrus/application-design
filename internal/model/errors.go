package model

import "fmt"

var (
	ErrNoRoomsAvailable   = fmt.Errorf("there is no rooms by order")
	ErrInventoryCorrupted = fmt.Errorf("invertory corrupted somehow")
)
