package entities

type DiscountTier struct {
	ID              uint
	EventID         uint
	MinParticipants uint
	DiscountPercent float64
}
