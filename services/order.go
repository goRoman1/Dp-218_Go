package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
)

type OrderService struct {
	repoOrder   repositories.OrderRepo

}

func NewOrderService(orderRepo repositories.OrderRepo) *OrderService {
	return &OrderService{repoOrder: orderRepo}
}

func (ors *OrderService) CreateOrder(user models.User, scooterID, startID, endID int,
	distance float64) (models.Order, error) {
	return ors.repoOrder.CreateOrder(user, scooterID, startID, endID, distance)
}

func (ors *OrderService) GetAllOrders() (*models.OrderList, error) {
	return ors.repoOrder.GetAllOrders()
}

func (ors *OrderService) GetOrderByID(orderID int) (models.Order, error) {
	return ors.repoOrder.GetOrderByID(orderID)
}

func (ors *OrderService) GetOrdersByUserID(userID int) (models.OrderList, error) {
	return ors.repoOrder.GetOrdersByUserID(userID)
}

func (ors *OrderService)  GetOrdersByScooterID(scooterID int) (models.OrderList, error) {
	return ors.repoOrder.GetOrdersByScooterID(scooterID)
}

func (ors *OrderService) GetScooterMileageByID(scooterID int) (float64, error) {
	return ors.repoOrder.GetScooterMileageByID(scooterID)
}

func (ors *OrderService) GetUserMileageByID(userID int) (float64, error) {
	return ors.repoOrder.GetUserMileageByID(userID)
}

func (ors *OrderService) UpdateOrder(orderID int, orderData models.Order) (models.Order, error) {
	return ors.repoOrder.UpdateOrder(orderID, orderData)
}

func (ors *OrderService) DeleteOrder(orderID int) error {
	return ors.repoOrder.DeleteOrder(orderID)
}
