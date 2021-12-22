package routing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
)

//type HttpMock struct {
//	ID int
//	StatusCode int
//	AccountServiceUC *services.AccountService
//	ScooterService *services.ScooterService
//	ScooterGrpcService *services.GrpcScooterService
//	OrderService *services.OrderService
//	Clock *mock.MockClock
//}
//
//func NewHttpMock(ctrl *gomock.Controller) *HttpMock {
//	repoAccount := mock.NewMockAccountRepo(ctrl)
//	repoAccountTransaction  := mock.NewMockAccountTransactionRepo(ctrl)
//	repoPaymentType := mock.NewMockPaymentTypeRepo(ctrl)
//	clock := mock.NewMockClock(ctrl)
//
//	//We created 'clock' for mocking 'time.Now()'
//	//Transfer 'clock' here just because it doesn't work in any other way.
//	accountServiceUC := services.NewAccountService(repoAccount, repoAccountTransaction, repoPaymentType,clock)
//
//	return &HttpMock{
//		AccountServiceUC: accountServiceUC,
//		RepoPaymentType: repoPaymentType,
//		RepoAccountTransaction: repoAccountTransaction,
//		RepoAccount: repoAccount,
//		Clock: clock,
//	}
//}

type testCaseForChoose struct {
	name string
	statusCode int
}

func TestChooseScooter(t *testing.T) {
	cases := []testCaseForChoose{
		{
			name: "Correct",
			statusCode: 200,
		},
	}

	for caseNum, item := range cases {
		bodyRequest := map[string]interface{}{
			"id": "3",
		}

		jsonBody, err := json.Marshal(bodyRequest)
		if err != nil {
			fmt.Println(err)
		}

		url :="/choose-scooter"
		req := httptest.NewRequest("POST", url, bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		ChooseScooter(w, req)

		if w.Code != item.statusCode {
			t.Errorf("[%d] wrong status code: got %d, expected %d", caseNum, w.Code, item.statusCode)
		}
	}
}
