package usecases

import "Dp218Go/models"

type OrderUseCases interface {
	CountTripDistance(order models.Order) (int, error)
	CountTripAmountMoney(order models.Order) (int, error)
	CompleteOrder(order *models.Order) error
}
