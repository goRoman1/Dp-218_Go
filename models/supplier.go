package models

type ScooterModel struct {
	ID          		int         `json:"id"`
	PaymentType 		PaymentType `json:"payment_type"`
	ModelName         	string 		`json:"model_name"`
	MaxWeight         	int    		`json:"max_weight"`
	Speed 			  	int    		`json:"speed"`
}

type ScooterModelList struct {
	ScooterModels []ScooterModel `json:"scooter_models"`
}

type SuppliersScooter struct {
	ID      		int 			`json:"id"`
	ModelId         int 			`json:"model_id"`
	SerialNumber 	string 			`json:"serial_number"`
}

type SuppliersScooterList struct {
	Scooters []SuppliersScooter `json:"scooters"`
}

type SupplierPrices struct {
	ID  	int `json:"id"`
	Price 	int `json:"price"`
	PaymentType PaymentType `json:"payment_type"`
	User    User  `json:"user"`
}

type SupplierPricesList struct {
	SupplierPrices []SupplierPrices `json:"supplier_prices_list"`
}

type ScooterModelDTO struct {
	ID          		int         `json:"id"`
	Price 				int 		`json:"price"`
	ModelName         	string 		`json:"model_name"`
	MaxWeight         	int    		`json:"max_weight"`
	Speed 			  	int    		`json:"speed"`
}

type ScooterModelDTOList struct {
	ScooterModelsDTO []ScooterModelDTO `json:"models_dto"`
}

type SupplierPricesDTO struct {
	ID  			int `json:"id"`
	Price 			int `json:"price"`
	PaymentTypeID 	int `json:"payment_type_id"`
	UserId    		int `json:"user_id"`
}

type SupplierPricesODTList struct {
	SupplierPricesDTO []SupplierPricesDTO `json:"supplier_prices_odt_list"`
}
