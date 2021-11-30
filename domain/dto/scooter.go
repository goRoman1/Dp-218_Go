package dto

import "net/http"

type Scooter struct {
	Id 				 int 			`json:"scooter_id"`
	Latitude 		 float64		`json:"latitude"`
	Longitude 		 float64		`json:"longitude"`
	ScooterModel 	 string			`json:"scooter_model"`
	MaxWeight 		 float64 	    `json:"max_weight"`
	BatteryRemain    float64   		`json:"battery_remain"`
	CanBeRent		 bool			`json:"can_be_rent"`
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



