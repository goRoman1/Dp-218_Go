package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"math"
)

type CustomerService struct {
	repoStation repositories.StationRepo
}

func NewCustomerService(repo repositories.StationRepo) *CustomerService {
	return &CustomerService{
		repoStation: repo,
	}
}

func (cs *CustomerService) ListStations() (*models.StationList, error) {
	return cs.repoStation.GetAllStations()
}

func (cs *CustomerService) ShowStation(id int) (*models.Station, error) {
	station, err := cs.repoStation.GetStationById(id)
	if err != nil {
		return nil, err
	}
	return &station, nil
}

func (cs *CustomerService) ShowNearestStation(x, y float64) (*models.Station, error) {

	stations, err := cs.repoStation.GetAllStations()
	if err != nil {
		return nil, err
	}

	nearest := calcNearest(x, y, stations.Station)
	return nearest, nil
}

func calcNearest(x, y float64, sts []models.Station) *models.Station {

	min := math.MaxFloat64
	var nearest models.Station

	for _, v := range sts {
		dis := calcDistance(x, y, v.Latitude, v.Longitude)
		if dis < min {
			min = dis
			nearest = v
		}
	}
	return &nearest
}

func calcDistance(x1, y1, x2, y2 float64) float64 {
	z := math.Pow(math.Abs(x1-x2), 2) + math.Pow(math.Abs(y1-y2), 2)
	return math.Sqrt(z)
}
