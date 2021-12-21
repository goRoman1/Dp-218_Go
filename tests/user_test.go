package tests

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"Dp218Go/tests/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type usecasesUser struct {
	repoUser *mocks.MockUserRepo
	repoRole *mocks.MockRoleRepo
	userUC   *services.UserService
}

type testCase struct {
	name string
	test func(t *testing.T, mock *usecasesUser)
}

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

			mock := NewUseCasesMock(ctrl)

			tc.test(tt, mock)
		})
	}
}

func NewUseCasesMock(ctrl *gomock.Controller) *usecasesUser {
	repoUser := mocks.NewMockUserRepo(ctrl)
	repoRole := mocks.NewMockRoleRepo(ctrl)
	userUC := services.NewUserService(repoUser, repoRole)

	return &usecasesUser{
		repoRole: repoRole,
		repoUser: repoUser,
		userUC:   userUC,
	}
}

func TestUseCases_User_ChangeUsersBlockStatus(t *testing.T) {
	runTestCases(t, []testCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *usecasesUser) {

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
			test: func(t *testing.T, mock *usecasesUser) {

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
			test: func(t *testing.T, mock *usecasesUser) {

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
