package services

import "Dp218Go/repositories"

type OrderService struct {
	repoOrder   repositories.OrderRepo
	scooterRepo repositories.ScooterRepo
}

func NewOrderService(orderRepo repositories.OrderRepo, scooterRepo repositories.ScooterRepo) *OrderService {
	return &OrderService{repoOrder: orderRepo, scooterRepo: scooterRepo}
}

//TODO continue implementation
