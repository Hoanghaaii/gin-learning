package entities

import "time"

type EventStatus string

const (
	StatusUpcoming EventStatus = "event_upcoming"
	StatusOngoing  EventStatus = "event_ongoing"
	StatusClosed   EventStatus = "event_closed"
)

type Event struct {
	ID        uint
	ProductID uint
	SellerID  uint
	StartTime time.Time
	EndTime   time.Time
	MaxSlots  uint
	Status    EventStatus
}
