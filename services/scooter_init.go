package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
)

type ScooterInitService struct {
	scooterInitRepo repositories.ScooterInitRepoI
}

func NewScooterInitService(scooterInitRepo repositories.ScooterInitRepoI) *ScooterInitService {
	return &ScooterInitService{scooterInitRepo}
}

func (si *ScooterInitService) GetOwnersScooters() (*models.SuppliersScooterList, error) {
	return si.scooterInitRepo.GetOwnersScooters()
}

func (si *ScooterInitService) GetActiveStations()(*models.StationList, error) {
	return si.scooterInitRepo.GetActiveStations()
}

func (si *ScooterInitService) AddStatusesToScooters(scooterIds []int, station models.Station) error {
	return si.scooterInitRepo.AddStatusesToScooters(scooterIds, station)
}

func (si *ScooterInitService) ConvertForTemplateStruct()*models.ScootersStationsAllocation{
	list := &models.ScootersStationsAllocation{}

	scooters, err := si.scooterInitRepo.GetOwnersScooters()
	if err != nil {
		return nil
	}
	stations, err := si.scooterInitRepo.GetActiveStations()
	if err != nil {
		return nil
	}

	list.SuppliersScooterList = *scooters
	list.StationList = *stations

	return list
}