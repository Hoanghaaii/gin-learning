package entities

import "time"

type EventParticipant struct {
	ID           uint
	UserID       uint
	EventID      uint
	Quantity     uint
	RegisteredAt time.Time
}
