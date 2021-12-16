package repositories

import (
	"Dp218Go/models"
)

type SupplierRepoI interface {
	GetModels()(*models.ScooterModelDTOList, error)
	SelectModel(id int)(*models.ScooterModelDTO, error)
	AddModel(modelData *models.ScooterModelDTO)error
	EditPrice(modelData *models.ScooterModelDTO) error

	GetSuppliersScootersByModelId(modelId int) (*models.SuppliersScooterList, error)
	AddSuppliersScooter(modelId int, scooter *models.SuppliersScooter) error
	DeleteSuppliersScooter(id int) error
	ConvertToStruct(path string)[]models.SuppliersScooter
}

