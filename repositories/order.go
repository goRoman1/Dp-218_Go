package repositories

import "Dp218Go/models"

type OrderRepo interface {
	CreateOrder(user models.User, scooter models.Scooter) (models.Order, error)
	SetOrderStart(order *models.Order, status models.ScooterStatusInRent) error
	SetOrderEnd(order *models.Order, status models.ScooterStatusInRent) error
	UpdateOrder(order *models.Order) error
}
