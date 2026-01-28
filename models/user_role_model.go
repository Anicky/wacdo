package models

type UserRole string

const (
	Admin       UserRole = "admin"
	Greeter     UserRole = "greeter"
	OrderPicker UserRole = "order_picker"
	Manager     UserRole = "manager"
)
