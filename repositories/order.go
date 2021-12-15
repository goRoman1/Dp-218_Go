package repositories

import "Dp218Go/models"

type OrderRepo interface {
	CreateOrder(user models.User, scooterID, startID, endID int, distance float64) (models.Order, error)
	UpdateOrder(orderID int, orderData models.Order) (models.Order, error)
	DeleteOrder(orderID int) error
	GetAllOrders() (*models.OrderList, error)
	GetOrderByID(orderID int) (models.Order, error)
	GetOrdersByUserID(userID int) (models.OrderList, error)
	GetOrdersByScooterID(scooterID int) (models.OrderList, error)
	GetScooterMileageByID(scooterID int) (float64, error)
	GetUserMileageByID(userID int) (float64, error)
}
