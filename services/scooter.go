package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
)

//ScooterService is the service which gives access to the ScooterRepo repository.
type ScooterService struct {
	repoScooter repositories.ScooterRepo
}

//NewScooterService creates the new ScooterService.
func NewScooterService(repoScooter repositories.ScooterRepo) *ScooterService {
	return &ScooterService{repoScooter: repoScooter}
}

//GetAllScooters gives the access to the ScooterRepo.GetAllScooters function.
func (ser *ScooterService) GetAllScooters() (*models.ScooterListDTO, error) {
	return ser.repoScooter.GetAllScooters()
}

//GetScooterById gives the access to the ScooterRepo.GetScooterById function.
func (ser *ScooterService) GetScooterById(uid int) (models.ScooterDTO, error) {
	return ser.repoScooter.GetScooterById(uid)
}

//GetScooterStatus gives the access to the ScooterRepo.GetScooterStatus function.
func (ser *ScooterService) GetScooterStatus(scooterID int) (models.ScooterStatus, error) {
	return ser.repoScooter.GetScooterStatus(scooterID)
}

//SendCurrentStatus gives the access to the ScooterRepo.SendCurrentStatus function.
func (ser *ScooterService) SendCurrentStatus(id int, lat, lon, battery float64) error {
	return ser.repoScooter.SendCurrentStatus(id, lat, lon, battery)
}

//CreateScooterStatusInRent gives the access to the ScooterRepo.CreateScooterStatusInRent function.
func (ser * ScooterService) CreateScooterStatusInRent(scooterID int) (models.ScooterStatusInRent, error) {
	return ser.repoScooter.CreateScooterStatusInRent(scooterID)
}

