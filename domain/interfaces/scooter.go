package interfaces

import (
	model "Dp218Go/domain/dto"
	"Dp218Go/repositories"
)

type ScooterRepo interface {
	GetAllScooters() (*model.ScooterList, error)
	GetScooterById(scooterId int) (model.Scooter, error)
	SendPosition(scooter model.Scooter)
	SendAtStart(uid int, client *repositories.Client) (error, int, int)
	SendAtEnd(tripId, locId int, client *repositories.Client) error
}
