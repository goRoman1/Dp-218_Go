package repositories

import "Dp218Go/models"

type ScooterRepo interface {
	GetAllScooters() (*models.ScooterList, error)
	GetScooterById(scooterId int) (models.ScooterDTO, error)
	GetScooterStatus(scooterID int) (models.ScooterStatus, error)
	//SendAtStart(uID, sID int) (error, int)
	//SendAtEnd(tripId int, client *repositories.Client) error
}

