package models

type OrderStatus int64

const (
	Created       OrderStatus = 0
	InPreparation OrderStatus = 1
	Prepared      OrderStatus = 2
	Delivered     OrderStatus = 3
)
