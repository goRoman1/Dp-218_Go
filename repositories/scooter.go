package repositories

import "Dp218Go/models"

type ScooterRepo interface {
	GetAllScooters() (*models.ScooterListDTO, error)
	GetScooterById(scooterId int) (models.ScooterDTO, error)
	GetScooterStatus(scooterID int) (models.ScooterStatus, error)
	SendCurrentPosition(id int, lat, lon float64) error
}

