package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
)

type StationService struct {
	repo repositories.AnyDatabase
}

func NewStationService(db repositories.AnyDatabase) *StationService {
	return &StationService{db}
}

func (db *StationService) GetAllStations() (*models.StationList, error) {
	return db.repo.GetAllStations()
}

func (db *StationService) AddStation(station *models.Station) error {
	return db.repo.AddStation(station)
}

func (db *StationService) GetStationById(stationId int) (models.Station, error) {
	return db.repo.GetStationById(stationId)
}

func (db *StationService) DeleteStation(stationId int) error {
	return db.repo.DeleteStation(stationId)
}