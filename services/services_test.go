package services

import (
	"Dp218Go/models"
	mock "Dp218Go/repositories/mock"
	mocks "Dp218Go/services/mock"
	"errors"
	"github.com/golang/mock/gomock"
	assert "github.com/stretchr/testify/require"
	"testing"
	"time"
)

//UseCasesMock is a struct which exists of repositories which are mocked and our service.
type UseCasesMock struct {
	AccountServiceUC       *AccountService
	RepoPaymentType        *mock.MockPaymentTypeRepo
	RepoAccountTransaction *mock.MockAccountTransactionRepo
	RepoAccount            *mock.MockAccountRepo
	Clock                  *mocks.MockClock

	repoUser *mock.MockUserRepo
	repoRole *mock.MockRoleRepo
	userUC   *UserService
}

type testCase struct {
	name string
	test func(t *testing.T, mock *UseCasesMock)
}

//We can keep this function without changes in our next test-cases. Except of 'mock' declaration.
func runTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					tt.Error(err)
				}
			}()

			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			//Here we should change if our struct name will be different to 'UseCasesMock'
			mock := NewUseCasesMock(ctrl)

			tc.test(tt, mock)
		})
	}
}

func NewUseCasesMock(ctrl *gomock.Controller) *UseCasesMock {
	repoAccount := mock.NewMockAccountRepo(ctrl)
	repoAccountTransaction := mock.NewMockAccountTransactionRepo(ctrl)
	repoPaymentType := mock.NewMockPaymentTypeRepo(ctrl)
	clock := mocks.NewMockClock(ctrl)

	//We created 'clock' for mocking 'time.Now()'
	//Transfer 'clock' here just because it doesn't work in any other way.
	accountServiceUC := NewAccountService(repoAccount, repoAccountTransaction, repoPaymentType, clock)


	repoUser := mock.NewMockUserRepo(ctrl)
	repoRole := mock.NewMockRoleRepo(ctrl)
	userUC := NewUserService(repoUser, repoRole)


	return &UseCasesMock{
		AccountServiceUC:       accountServiceUC,
		RepoPaymentType:        repoPaymentType,
		RepoAccountTransaction: repoAccountTransaction,
		RepoAccount:            repoAccount,
		Clock:                  clock,

		repoUser: repoUser,
		repoRole: repoRole,
		userUC: userUC,
	}
}

func TestUseCases_Account_AddMoneyToAccount(t *testing.T) {
	runTestCases(t, []testCase{
		{ //In this case we are going by happy path.
			name: "Correct",
			test: func(t *testing.T, mock *UseCasesMock) {

				//Create a variable with the exact time for mocking time.Now().
				var currentTime = time.Date(2021, 12, 19, 12, 21, 00, 00, time.UTC)

				//With help of mock we can call the functions of repositories without deployment.
				//'EXPECT' means that the function will be called.
				//The next we call the function we need ex:'GetPaymentTypeByID'
				//'Return' let us set the values which will be returned. We can also return an error.
				//With 'Times' we set how many times the function will be called.
				mock.RepoPaymentType.EXPECT().GetPaymentTypeById(2).
					Return(models.PaymentType{}, nil).Times(1)

				//Here we are mocking the time of our 'Clock' which is a wrapper of the system service 'Time'
				//With the value of 'currentTime'
				mock.Clock.EXPECT().Now().Return(currentTime).Times(1)

				//Into 'DateTime' we put the 'currentTime'.
				//So now for the test we have the same time into the struct and into the mocked 'Clock'.
				accTransaction := &models.AccountTransaction{
					DateTime:    currentTime,
					PaymentType: models.PaymentType{},
					AccountFrom: models.Account{},
					AccountTo:   models.Account{},
					Order:       models.Order{},
					AmountCents: 50}

				//We call this func like in the order of calls into the real 'AddMoneyToAccount'
				mock.RepoAccountTransaction.EXPECT().AddAccountTransaction(accTransaction).
					Return(nil).Times(1)

				//In this case we expect that function will be called without any errors.
				err := mock.AccountServiceUC.AddMoneyToAccount(accTransaction.AccountTo, 50)

				//Compare that expected value of error is nil.
				assert.Equal(t, nil, err)
			},
		}, { //In this case we are going by getting the error
			name: "Incorrect. Got error from GetPaymentTypeByID",
			test: func(t *testing.T, mock *UseCasesMock) {

				//Describe which error we'll get.
				expectedError := errors.New("expectedError")

				//Call 'GetPaymentTypeById' and return here our 'expectedError'
				mock.RepoPaymentType.EXPECT().GetPaymentTypeById(2).
					Return(models.PaymentType{}, expectedError).Times(1)

				accTransaction := &models.AccountTransaction{
					DateTime:    time.Now(),
					PaymentType: models.PaymentType{},
					AccountFrom: models.Account{},
					AccountTo:   models.Account{},
					Order:       models.Order{},
					AmountCents: 50}

				//Calling 'AddMoneyToAccount' will return us the error, because we had the error into the func before.
				err := mock.AccountServiceUC.AddMoneyToAccount(accTransaction.AccountTo, 50)

				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func TestUseCases_User_ChangeUsersBlockStatus(t *testing.T) {
	runTestCases(t, []testCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *UseCasesMock) {

				mock.repoUser.EXPECT().GetUserByID(1).
					Return(models.User{IsBlocked: false}, nil).Times(1)

				mock.repoUser.EXPECT().UpdateUser(1, models.User{IsBlocked: true}).
					Return(models.User{}, nil).Times(1)

				err := mock.userUC.ChangeUsersBlockStatus(1)
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "incorrect get by ID",
			test: func(t *testing.T, mock *UseCasesMock) {

				var someError = errors.New("error get by ID")

				mock.repoUser.EXPECT().GetUserByID(2).
					Return(models.User{IsBlocked: false}, someError).Times(1)

				err := mock.userUC.ChangeUsersBlockStatus(2)
				assert.Error(t, err)
				assert.Equal(t, someError, err)
			},
		},
		{
			name: "incorrect update user",
			test: func(t *testing.T, mock *UseCasesMock) {

				var someError = errors.New("error update user")

				mock.repoUser.EXPECT().GetUserByID(3).
					Return(models.User{IsBlocked: false}, nil).Times(1)

				mock.repoUser.EXPECT().UpdateUser(3, models.User{IsBlocked: true}).
					Return(models.User{}, someError).Times(1)

				err := mock.userUC.ChangeUsersBlockStatus(3)
				assert.Error(t, err)
				assert.Equal(t, someError, err)
			},
		},
	})
}

