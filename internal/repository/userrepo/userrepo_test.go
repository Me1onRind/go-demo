package userrepo

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/unittest"
	"github.com/Me1onRind/go-demo/internal/model/po/userpo"
	sqlmockGormHelper "github.com/Me1onRind/sqlmock-gorm-helper"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	gconfig.DynamicCfg.DefaultDB.Label = "default"
	os.Exit(m.Run())
}

func Test_GetUserByUserId(t *testing.T) {
	tests := []struct {
		name string
		data *userpo.User
		err  error
	}{
		{
			name: "success",
			data: UserJ,
		},
		{
			name: "not found",
			err:  gorm.ErrRecordNotFound,
		},
		{
			name: "io error",
			err:  io.EOF,
		},
	}

	userRepo := NewUserRepo()
	mock := mysql.NewMysqlMock(userRepo.DBLabel())
	defer assert.Empty(t, mock.ExpectationsWereMet())
	for _, test := range tests {
		mockProcess := mock.ExpectQuery("SELECT \\* FROM `user_tab` WHERE user_id=\\? LIMIT 1").
			WithArgs(166)

		t.Run(test.name, func(t *testing.T) {
			if test.err != nil {
				mockProcess.WillReturnError(test.err)
			} else {
				mockProcess.WillReturnRows(sqlmockGormHelper.ModelToRows(test.data))
			}

			user, err := userRepo.GetUser(context.Background(), WithUserId(166))
			if test.err != nil {
				assert.Empty(t, user)
				assert.ErrorIs(t, err, test.err)
			} else {
				assert.Empty(t, err)
				assert.Equal(t, true, assert.ObjectsAreEqual(user, test.data))
			}
		})
	}
}

func Test_CreateUser(t *testing.T) {
	tests := []struct {
		name string
		data *userpo.User
		err  error
	}{
		{
			name: "success",
			data: &userpo.User{Name: "test_name", UserId: 167, Email: "test@google.com"},
		},
		{
			name: "fail",
			data: &userpo.User{Name: "test_name", UserId: 167, Email: "test2@google.com"},
			err:  &mysqlDriver.MySQLError{Number: 1062, Message: "Duplicate entry"},
		},
	}

	userRepo := NewUserRepo()
	mock := mysql.NewMysqlMock(userRepo.DBLabel())
	defer assert.Empty(t, mock.ExpectationsWereMet())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockProcess := mock.ExpectExec(sqlmockGormHelper.InsertSql(userpo.User{}, `user_tab`)).
				WithArgs(unittest.Now{}, unittest.Now{}, test.data.UserId, test.data.Email, test.data.Name)
			if test.err != nil {
				mockProcess.WillReturnError(test.err)
			} else {
				mockProcess.WillReturnResult(sqlmock.NewResult(100, 1))
			}

			user, err := userRepo.CreateUser(context.Background(), test.data)
			if test.err != nil {
				assert.Empty(t, user)
				assert.ErrorIs(t, err, test.err)
			} else {
				if assert.Empty(t, err) {
					assert.EqualValues(t, 100, user.Id)
					assert.Equal(t, true, user.CreateTime > 0)
					assert.Equal(t, true, user.UpdateTime > 0)
				}
			}

		})
	}
}

func Test_GetUserByEmail(t *testing.T) {
	tests := []struct {
		name string
		data *userpo.User
		err  error
	}{
		{
			name: "success",
			data: UserJ,
		},
		{
			name: "not found",
			err:  gorm.ErrRecordNotFound,
		},
		{
			name: "io error",
			err:  io.EOF,
		},
	}

	userRepo := NewUserRepo()
	mock := mysql.NewMysqlMock(userRepo.DBLabel())
	defer assert.Empty(t, mock.ExpectationsWereMet())
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockProcess := mock.ExpectQuery("SELECT \\* FROM `user_tab` WHERE email=\\? LIMIT 1").
				WithArgs("test@google.com")
			if test.err != nil {
				mockProcess.WillReturnError(test.err)
			} else {
				mockProcess.WillReturnRows(sqlmockGormHelper.ModelToRows(test.data))
			}

			user, err := userRepo.GetUser(context.Background(), WithEmail("test@google.com"))
			if test.err != nil {
				assert.Empty(t, user)
				assert.ErrorIs(t, err, test.err)
			} else {
				assert.Empty(t, err)
				assert.Equal(t, true, assert.ObjectsAreEqual(user, test.data))
			}
		})
	}
}
