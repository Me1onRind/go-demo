package userrepo

import (
	"context"
	"io"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Me1onRind/go-demo/internal/constant/dblabel"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/unittest"
	"github.com/Me1onRind/go-demo/internal/model/po"
	"github.com/Me1onRind/go-demo/internal/model/po/userpo"
	sqlmock_rows_helper "github.com/Me1onRind/sqlmock-rows-helper"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	UserJ = &userpo.User{
		BaseModel: po.BaseModel{
			Id:         1,
			CreateTime: 1680000000,
			UpdateTime: 1680000000,
		},
		UserId: 166,
		Name:   "J",
	}
)

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

	mock := mysql.NewMysqlMock(dblabel.Default)
	defer assert.Empty(t, mock.ExpectationsWereMet())
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userRepo := NewUserRepo()

			if test.err != nil {
				mock.ExpectQuery("SELECT \\* FROM `user_tab` WHERE user_id=\\? LIMIT 1").
					WithArgs(166).WillReturnError(test.err)
			} else {
				mock.ExpectQuery("SELECT \\* FROM `user_tab` WHERE user_id=\\? LIMIT 1").
					WithArgs(166).WillReturnRows(sqlmock_rows_helper.ModelToRows(test.data))
			}

			user, err := userRepo.GetUserByUserId(context.Background(), 166)
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
			data: &userpo.User{Name: "test_name", UserId: 167},
		},
		{
			name: "fail",
			data: &userpo.User{Name: "test_name", UserId: 167},
			err:  &mysqlDriver.MySQLError{Number: 1062, Message: "Duplicate entry"},
		},
	}

	mock := mysql.NewMysqlMock(dblabel.Default)
	defer assert.Empty(t, mock.ExpectationsWereMet())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userRepo := NewUserRepo()

			if test.err != nil {
				mock.ExpectExec("INSERT INTO `user_tab` \\(`create_time`,`update_time`,`user_id`,`name`\\) VALUES \\(\\?,\\?,\\?,\\?\\)").
					WithArgs(unittest.NewGreater(int64(0)), unittest.NewGreater(int64(0)), 167, "test_name").WillReturnError(test.err)
			} else {
				mock.ExpectExec("INSERT INTO `user_tab` \\(`create_time`,`update_time`,`user_id`,`name`\\) VALUES \\(\\?,\\?,\\?,\\?\\)").
					WithArgs(unittest.NewGreater(int64(0)), unittest.NewGreater(int64(0)), 167, "test_name").WillReturnResult(sqlmock.NewResult(100, 1))
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
