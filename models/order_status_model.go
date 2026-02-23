package models

type OrderStatus string

const (
	Created       OrderStatus = "created"
	InPreparation OrderStatus = "inPreparation"
	Prepared      OrderStatus = "prepared"
	Delivered     OrderStatus = "delivered"
)
