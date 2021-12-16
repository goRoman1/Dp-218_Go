package postgres

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jszwec/csvutil"
	"io"
	"os"
	"strings"
)

var userId = 1

type SupplierRepoDB struct {
	db repositories.AnyDatabase
}

func NewSupplierRepoDB(db repositories.AnyDatabase) *SupplierRepoDB {
	return &SupplierRepoDB{db}
}

// Соберает список моделей скуторв для отображения на странице суплаера, будут отображаться модели с ценой тарифа
func (s *SupplierRepoDB)GetModels()(*models.ScooterModelDTOList, error){
	modelsOdtList := &models.ScooterModelDTOList{}
	pricesList := &models.SupplierPricesODTList{}

	pricesList, err := s.getPrices()
	if err != nil {
		return modelsOdtList, err
	}

	querySQL := `SELECT * FROM scooter_models ORDER BY id DESC;`
	rows, err := s.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return modelsOdtList, err
	}

	for rows.Next() {
		var paymentTypeID int
		var model models.ScooterModelDTO
		err := rows.Scan(&model.ID, &paymentTypeID, &model.ModelName, &model.MaxWeight, &model.Speed)
		if err != nil {
			return modelsOdtList, err
		}

		model.Price, err = s.findSupplierPricesList(pricesList, paymentTypeID, userId)
		if err != nil {
			return modelsOdtList, err
		}
		modelsOdtList.ScooterModelsDTO = append(modelsOdtList.ScooterModelsDTO, model)
	}
	return modelsOdtList, nil
}

//возвращает модель выбранную на странице супплаеров, в дальнейшем она будет передедана для редактирования цены
func (s *SupplierRepoDB)SelectModel(id int)(*models.ScooterModelDTO, error){
	modelODT := &models.ScooterModelDTO{}

	querySQL := `SELECT id, payment_type_id, model_name, max_weight, speed  FROM scooter_models WHERE id = $1;`
	row := s.db.QueryResultRow(context.Background(), querySQL, id)

	var paymentTypeId int
	err := row.Scan(&modelODT.ID, &paymentTypeId, &modelODT.ModelName, &modelODT.MaxWeight, &modelODT.Speed)
	if err != nil {
		return modelODT, err
	}

	modelODT.Price, err = s.getPrice(paymentTypeId, userId)

	return modelODT, err
}

// добавляет модель и цену,  введенную вручную, в базу данных
func (s *SupplierRepoDB)AddModel(modelData *models.ScooterModelDTO)error{

	paymentTypeId, err := s.addPaymentTypeId(modelData.ModelName)
	if err != nil {
		return err
	}
	var modelId int
	querySQL := `INSERT INTO scooter_models(payment_type_id, model_name, max_weight, speed)
	   		VALUES($1, $2, $3, $4)
	   		RETURNING id;`
	err = s.db.QueryResultRow(context.Background(), querySQL, &paymentTypeId, modelData.ModelName, modelData.MaxWeight, modelData.Speed).Scan(&modelId)
	if err != nil {
		return  err
	}

	var priceId int
	querySQL = `INSERT INTO supplier_prices(price, payment_type_id, user_id)
	   		VALUES($1, $2, $3)
	   		RETURNING id;`
	err = s.db.QueryResultRow(context.Background(), querySQL, modelData.Price, paymentTypeId, userId).Scan(&priceId)
	if err != nil {
		return  err
	}
	return  nil
}

//изменяет цену аренды, которая привязана к модели
func (s *SupplierRepoDB) EditPrice(modelData *models.ScooterModelDTO) error{
	price := &models.ScooterModelDTO{}
	paymentTypeId, err := s.getPaymentTypeByModelName(modelData.ModelName)
	if err != nil {
		return err
	}

	querySQL := `UPDATE supplier_prices SET price=$1 WHERE payment_type_id = $2 AND user_id = $3 RETURNING price;`
	err = s.db.QueryResultRow(context.Background(), querySQL, modelData.Price, paymentTypeId, userId).Scan(&price.Price)
	if err != nil {
		return err
	}

	return nil
}

func (s *SupplierRepoDB) getPrices()(*models.SupplierPricesODTList, error){
	list := &models.SupplierPricesODTList{}

	querySQL := `SELECT * FROM supplier_prices ORDER BY id DESC;`
	rows, err := s.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var supplierPriceODT models.SupplierPricesDTO
		err := rows.Scan(&supplierPriceODT.ID, &supplierPriceODT.Price, &supplierPriceODT.PaymentTypeID, &supplierPriceODT.UserId)

		if err != nil {
			return list, err
		}

		list.SupplierPricesDTO = append(list.SupplierPricesDTO, supplierPriceODT)
	}
	return list, nil
}

func (s *SupplierRepoDB)findSupplierPricesList(supplierPrice *models.SupplierPricesODTList, paymentTypeId int, userId int )(int, error){
	for _, v := range supplierPrice.SupplierPricesDTO {
		if v.PaymentTypeID == paymentTypeId && v.UserId == userId{
			return v.Price, nil
		}
	}
	return 0, fmt.Errorf("not found paymentType id=%d", paymentTypeId)
}

func (s *SupplierRepoDB) getPrice(paymentTypeId, userId int) (int, error) {
	price := models.ScooterModelDTO{}
	querySQL := `SELECT price FROM supplier_prices WHERE payment_type_id = $1 AND user_id = $2;`
	row := s.db.QueryResultRow(context.Background(), querySQL, paymentTypeId, userId)
	err := row.Scan(&price.Price)

	return price.Price, err
}

func (s *SupplierRepoDB) addPaymentTypeId(modelName string)(int,error){
	var paymentTypeId int
	querySQL := `INSERT INTO payment_types (name) VALUES ($1) RETURNING id;`
	err := s.db.QueryResultRow(context.Background(), querySQL, modelName).Scan(&paymentTypeId)
	if err != nil {
		return 0,err
	}
	return  paymentTypeId, nil
}

func (s *SupplierRepoDB) getPaymentTypeByModelName(modelName string) (int, error) {
	paymentType := models.PaymentType{}
	querySQL := `SELECT * FROM payment_types WHERE name = $1;`
	row := s.db.QueryResultRow(context.Background(), querySQL, modelName)
	err := row.Scan(&paymentType.ID, &paymentType.Name)
	return paymentType.ID, err
}

//получить скутеры относящийся к модели скутеров
func (s *SupplierRepoDB) GetSuppliersScootersByModelId(modelId int) (*models.SuppliersScooterList, error) {
	list := &models.SuppliersScooterList{}

	querySQL := `SELECT id, serial_number FROM scooters WHERE model_id = $1 ORDER BY id DESC;`
	rows, err := s.db.QueryResult(context.Background(), querySQL, modelId)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var scooter models.SuppliersScooter
		err := rows.Scan(&scooter.ID, &scooter.SerialNumber)
		if err != nil {
			return list, err
		}

		list.Scooters = append(list.Scooters, scooter)
	}
	return list, nil
}

//добавляет скутер в таблицу скутуров, его серийный номер будет отображаться в списке скутеров
func (s *SupplierRepoDB) AddSuppliersScooter(modelId int, scooterData *models.SuppliersScooter) error {
	var id int
	querySQL := `INSERT INTO scooters(model_id, owner_id, serial_number)
	   		VALUES($1, $2, $3)
	   		RETURNING id;`
	err := s.db.QueryResultRow(context.Background(), querySQL, modelId, userId, scooterData.SerialNumber).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}
//удаляет скутер из списка скутеров и из базы
func (s *SupplierRepoDB) DeleteSuppliersScooter(id int) error {
	querySQL := `DELETE FROM scooters WHERE id = $1;`
	_, err := s.db.QueryExec(context.Background(), querySQL, id)
	return err
}

//конвертирование полученных данных из .csv в структуру для дальнейшей работы
func (s *SupplierRepoDB) ConvertToStruct(path string) []models.SuppliersScooter {

	csvFile, _ := os.Open(path)
	reader := csv.NewReader(csvFile)
	reader.Comma = ';'

	scooterHeader, _ := csvutil.Header(models.SuppliersScooter{}, "csv")
	dec, _ := csvutil.NewDecoder(reader, scooterHeader...)

	var fileData []models.SuppliersScooter
	for {
		var s models.SuppliersScooter
		if err := dec.Decode(&s.SerialNumber); err == io.EOF {
			break
		}
		fileData = append(fileData, s)
		fmt.Println(fileData)
	}
	return fileData
}

//вносить полученые из файла данные в базу
func (s *SupplierRepoDB) InsertToDb(modelId int, scooters []models.SuppliersScooter) error{

	valueStrings := make([]string, 0, len(scooters))
	valueArgs := make([]interface{}, 0, len(scooters) * 1)
	for i, scooter := range scooters {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d)", i*1+1))
		valueArgs = append(valueArgs, scooter.SerialNumber)
	}

	stmt := fmt.Sprintf("INSERT INTO scooters(scooter_model, user_id, serial_number) VALUES %s", strings.Join(valueStrings, ","))
	if _, err := s.db.QueryExec(context.Background(),stmt, valueArgs...)
		err != nil {
		fmt.Println("Unable to insert due to: ", err)
		return err
	}
	return nil
}
