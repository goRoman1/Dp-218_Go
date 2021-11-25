package dto

import "net/http"

type Scooter struct {
	Id 				 int 			`json:"scooterid"`
	LocationId		 int			`json:"locationid"`
	Lattitude 		 float64		`json:"lattitude"`
	Longtitude 		 float64		`json:"longtitude"`
	ScooterModel 	 string			`json:"scootermodel"`
	MaxWeight 		 float64 	    `json:"maxweight"`
	BatteryRemain    float64   		`json:"batteryremain"`
	CanBeRent		 bool			`json:"canberent"`
}

type ScooterList struct {
	Scooters []Scooter `json:"scooters"`
}

func (s *Scooter) Bind(r *http.Request) error {
	return nil
}

func (*ScooterList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*Scooter) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}



