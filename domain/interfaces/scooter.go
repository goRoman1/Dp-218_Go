package interfaces

import (
	model "Dp218Go/domain/dto"
)

type ScooterRepo interface {
	GetAllScooters() (*model.ScooterList, error)
	GetScooterById(scooterId int) (model.Scooter, error)
	SendPosition(scooter model.Scooter)
}
