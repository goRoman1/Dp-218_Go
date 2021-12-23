package repositories

import "Dp218Go/models"

type ScooterInitRepoI interface {
	GetOwnersScooters() (*models.SuppliersScooterList, error)
	GetActiveStations()(*models.StationList, error)
	AddStatusesToScooters(scooterIds []int, station models.Station) error
}