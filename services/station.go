package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
)


type StationService struct {
	repoStation repositories.StationRepo
}

func NewStationService(repoStation repositories.StationRepo) *StationService {
	return &StationService{repoStation: repoStation}
}

func (db *StationService) GetAllStations() (*models.StationList, error) {
	return db.repoStation.GetAllStations()
}

func (db *StationService) AddStation(station *models.Station) error {
	return db.repoStation.AddStation(station)
}

func (db *StationService) GetStationById(stationId int) (models.Station, error) {
	return db.repoStation.GetStationById(stationId)
}

func (db *StationService) DeleteStation(stationId int) error {
	return db.repoStation.DeleteStation(stationId)
}