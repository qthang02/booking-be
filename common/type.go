package common

type UserType int

const (
	UserTypeAdmin UserType = iota
	UserTypeStaff
)

type Status int

const (
	StatusActive Status = iota
	StatusInactive
)
