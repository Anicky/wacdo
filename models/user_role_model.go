package models

type UserRole int64

const (
	Admin       UserRole = 0
	Greeter     UserRole = 1
	OrderPicker UserRole = 2
	Manager     UserRole = 3
)
