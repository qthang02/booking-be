package types

type CategoryType int

const (
	ORIGINAL CategoryType = iota
	VIP1
	VIP2
	POPULAR
	SEPARATELY
	DELUXE
)

type RoomStatus int

const (
	Ready RoomStatus = iota
	Occupied
	DueOut
)
