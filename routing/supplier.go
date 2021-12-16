package routing

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
)

var supplierService *services.SupplierService

var scooterModelKeyRoutes = []Route{
	{
		Uri:     `/models`,
		Method:  http.MethodGet,
		Handler: getModels,
	},
	{
		Uri:     `/createModel`,
		Method:  http.MethodPost,
		Handler: createModel,
	},
	{
		Uri:     `/editPrice`,
		Method:  http.MethodPut,
		Handler: editPrice,
	},
	{
		Uri:     `/upload`,
		Method:  http.MethodPost,
		Handler: uploadFile,
	},
	{
		Uri:     `/getSuppliersScootersByModelId/{id}`,
		Method:  http.MethodGet,
		Handler: getSuppliersScootersByModelId,
	},
	{
		Uri:     `/addScooter/{id}`,
		Method:  http.MethodPost,
		Handler: addSuppliersScooter,
	},
	{
		Uri:     `/deleteScooter/{id}`,
		Method:  http.MethodDelete,
		Handler: deleteSuppliersScooter,
	},
}

type FileHeader struct {
	Filename string
	Header   textproto.MIMEHeader
}

func AddSupplierHandler(router *mux.Router, service *services.SupplierService){
	supplierService = service
	for _, rt := range scooterModelKeyRoutes {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getModels(w http.ResponseWriter, r *http.Request) {
	var modelList = &models.ScooterModelDTOList{}
	var err error
	format := GetFormatFromRequest(r)

	r.ParseForm()

	modelList, err = supplierService.GetModels()
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, modelList, HTMLPath+"supplier.html")
}

func createModel(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)
	model := &models.ScooterModelDTO{}

	DecodeRequest(FormatJSON, w, r, model, decodePriceRequest)

	if err := supplierService.AddModel(model); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, model)
}

func editPrice(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	modelId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	modelData, err := supplierService.SelectModel(modelId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	DecodeRequest(FormatJSON, w, r, modelData, decodePriceRequest)
	if err := supplierService.ChangePrice(modelData); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, modelData)
}

func decodePriceRequest(r *http.Request, data interface{}) error {
	r.ParseForm()
	modelData := data.(*models.ScooterModelDTO)

	if _, ok := r.Form["price"]; ok {
		modelPrice, err := strconv.Atoi(r.FormValue("price"))
		if err != nil {
			return err
		}
		modelData.Price = modelPrice
	}
	if _, ok := r.Form["modelName"]; ok {
		modelData.ModelName= r.FormValue("modelName")
	}
	if _, ok := r.Form["maxWeight"]; ok {
		modelMaxWeight, err := strconv.Atoi(r.FormValue("maxWeight"))
		if err != nil {
			return err
		}
		modelData.MaxWeight = modelMaxWeight
	}
	if _, ok := r.Form["speed"]; ok {
		modelSpeed, err := strconv.Atoi(r.FormValue("speed"))
		if err != nil {
			return err
		}
		modelData.MaxWeight = modelSpeed
	}

	return nil
}

/////////////////////

func getSuppliersScootersByModelId(w http.ResponseWriter, r *http.Request) {
	var scooters = &models.SuppliersScooterList{}
	var err error
	format := GetFormatFromRequest(r)

	modelId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	scooters, err = supplierService.GetSuppliersScootersByModelId(modelId)

	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, scooters, HTMLPath+"scooters-list.html")
}

func addSuppliersScooter(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)
	modelId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	scooter := &models.SuppliersScooter{}
	DecodeRequest(FormatJSON, w, r, scooter, nil)

	if err := supplierService.AddSuppliersScooter(modelId, scooter); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, scooter)
}

func deleteSuppliersScooter(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	userId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	err = supplierService.DeleteSuppliersScooter(userId)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}
	EncodeAnswer(format, w, ErrorRenderer(fmt.Errorf(""), "success", http.StatusOK))
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	filepath := "./internal/"+handler.Filename
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	supplierService.InsertScootersToDb(filepath)
}
