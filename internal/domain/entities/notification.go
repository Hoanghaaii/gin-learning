package entities

import "time"

type NotificationType string

const (
	NotificationEventRegistered  NotificationType = "EVENT_REGISTERED"
	NotificationThresholdReached NotificationType = "THRESHOLD_REACHED"
	NotificationEventEnded       NotificationType = "EVENT_ENDED"
)

type Notification struct {
	ID          uint
	UserID      uint
	Title       string
	Message     string
	Type        NotificationType
	ReferenceID uint
	IsRead      bool
	CreatedAt   time.Time
}
