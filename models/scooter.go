package models

import "time"

type Scooter struct {
	ID 				 int 			`json:"scooter_id"`
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

//type Scooter struct {
//	ID 				int 		`json:"id"`
//	ModelID 		int 		`json:"model_id"`
//	OwnerID 		int 		`json:"owner_id"`
//	SerialNumber 	int 		`json:"serial_number"`
//}
//
//
//type ScooterModel struct {
//	ID 				int 		`json:"id"`
//	PaymentTypeID	int 		`json:"payment_type_id"`
//	ModelName 		string 		`json:"model_name"`
//	MaxWeight 		float64 	`json:"max_weight"`
//	Speed 			int 		`json:"speed"`
//}
//

type ScooterStatus struct {
	Scooter 		Scooter 	`json:"scooter"`
	Location 		Coordinate 	`json:"location"`
	BatteryRemain	int 		`json:"battery_remain"`
	StationID 		int 		`json:"station_id"`
}

type ScooterStatusInRent struct {
	ID 				int 		`json:"id"`
	StationID 		int 		`json:"station_id"`
	DateTime 		time.Time 	`json:"date_time"`
	Location 		Coordinate	`json:"location"`
}

//
//type Location struct {
//	ID 				int 		`json:"id"`
//	Latitude 		float64 	`json:"latitude"`
//	Longitude 		float64 	`json:"longitude"`
//	Label 			string 		`json:"label"`
//}