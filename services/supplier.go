package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
)

type SupplierService struct {
	SupplierRepo repositories.SupplierRepoI
}

func NewSupplierService(SupplierRepo repositories.SupplierRepoI) *SupplierService {
	return &SupplierService{
		SupplierRepo : SupplierRepo,
	}
}

func (s *SupplierService)GetSuppliersScootersByModelId(modelId int)(*models.SuppliersScooterList, error) {
	return s.SupplierRepo.GetSuppliersScootersByModelId(modelId)
}

func (s *SupplierService)AddSuppliersScooter(modelId int, scooter *models.SuppliersScooter)error{
	return s.SupplierRepo.AddSuppliersScooter(modelId, scooter)
}

func (s *SupplierService)DeleteSuppliersScooter(id int) error{
	return s.SupplierRepo.DeleteSuppliersScooter(id)
}

func (s *SupplierService)InsertScootersToDb(path string){
	s.SupplierRepo.ConvertToStruct(path)
}

func (s *SupplierService)GetModels()(*models.ScooterModelDTOList, error) {
	return s.SupplierRepo.GetModels()
}

func (s *SupplierService)SelectModel(id int)(*models.ScooterModelDTO, error) {
	return s.SupplierRepo.SelectModel(id)
}

func (s *SupplierService)AddModel(modelData *models.ScooterModelDTO)error {
	return s.SupplierRepo.AddModel(modelData)
}

func (s *SupplierService)ChangePrice(modelData *models.ScooterModelDTO)error {
	return s.SupplierRepo.EditPrice(modelData)
}


