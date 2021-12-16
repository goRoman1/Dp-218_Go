package repositories

import "Dp218Go/models"

//ScooterRepo the interface which implemented by functions which connect to the database.
type ScooterRepo interface {
	GetAllScooters() (*models.ScooterListDTO, error)
	GetScooterById(scooterId int) (models.ScooterDTO, error)
	GetScooterStatus(scooterID int) (models.ScooterStatus, error)
	SendCurrentStatus(id int, lat, lon,battery float64) error
	CreateScooterStatusInRent(scooterID int) (models.ScooterStatusInRent, error)
}