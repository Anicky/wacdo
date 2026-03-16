package models

type UserRole string

const (
	Admin       UserRole = "admin"
	Greeter     UserRole = "greeter"
	OrderPicker UserRole = "order_picker"
	Manager     UserRole = "manager"
)

func (role UserRole) IsValid() bool {
	switch role {
	case Admin, Greeter, OrderPicker, Manager:
		return true
	}

	return false
}
