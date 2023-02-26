package respository

import (
	"context"
	"io"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/unittest"
	"github.com/Me1onRind/go-demo/internal/model/po"
	sqlmockGormHelper "github.com/Me1onRind/sqlmock-gorm-helper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type TestModel struct {
	po.BaseModel
	Field1 string `gorm:"column:field1"`
}

func (t *TestModel) TableName() string {
	return "test_tab"
}

func (t *TestModel) DBLabel() string {
	return "mock"
}

var testModel *TestModel = &TestModel{
	BaseModel: po.BaseModel{
		Id:         123456,
		CreateTime: 1670000000,
		UpdateTime: 1680000000,
	},
	Field1: "mock",
}

func Test_Create(t *testing.T) {
	tests := []struct {
		name string
		data *TestModel
		err  error
	}{
		{
			name: "success",
			data: &TestModel{Field1: "mck"},
		},
		{
			name: "fail",
			data: &TestModel{Field1: "mck2"},
			err:  io.EOF,
		},
	}
	mock := mysql.NewMysqlMock((&TestModel{}).DBLabel())
	defer assert.Empty(t, mock.ExpectationsWereMet())
	repo := NewBaseRepo[*TestModel]()
	for _, test := range tests {
		mockProcess := mock.ExpectExec(sqlmockGormHelper.InsertSql(&TestModel{}, `test_tab`)).
			WithArgs(unittest.Now{}, unittest.Now{}, test.data.Field1)
		t.Run(test.name, func(t *testing.T) {
			if test.err != nil {
				mockProcess.WillReturnError(test.err)
			} else {
				mockProcess.WillReturnResult(sqlmock.NewResult(100, 1))
			}

			value, err := repo.Create(context.Background(), test.data)
			if test.err != nil {
				assert.Empty(t, value)
				assert.ErrorIs(t, err, test.err)
			} else {
				if assert.Empty(t, err) {
					assert.EqualValues(t, 100, value.Id)
					assert.Equal(t, true, value.CreateTime > 0)
					assert.Equal(t, true, value.UpdateTime > 0)
				}
			}
		})
	}
}

func Test_Take(t *testing.T) {
	tests := []struct {
		name string
		data *TestModel
		err  error
	}{
		{
			name: "success",
			data: testModel,
		},
		{
			name: "notFound",
			err:  gorm.ErrRecordNotFound,
		},
		{
			name: "ioError",
			err:  io.EOF,
		},
	}
	mock := mysql.NewMysqlMock((&TestModel{}).DBLabel())
	defer assert.Empty(t, mock.ExpectationsWereMet())
	repo := NewBaseRepo[*TestModel]()
	for _, test := range tests {
		mockProcess := mock.ExpectQuery("SELECT \\* FROM `test_tab` WHERE id=\\? LIMIT 1").
			WithArgs(123456)
		t.Run(test.name, func(t *testing.T) {
			if test.err != nil {
				mockProcess.WillReturnError(test.err)
			} else {
				mockProcess.WillReturnRows(sqlmockGormHelper.ModelToRows(test.data))
			}

			value, err := repo.Take(context.Background(), func(db *gorm.DB) *gorm.DB {
				return db.Where("id=?", 123456)
			})
			if test.err != nil {
				assert.Empty(t, value)
				assert.ErrorIs(t, err, test.err)
			} else {
				assert.Empty(t, err)
				assert.Equal(t, true, assert.ObjectsAreEqual(value, test.data))
			}
		})
	}
}

func Test_Find(t *testing.T) {
	tests := []struct {
		name string
		data []*TestModel
		err  error
	}{
		{
			name: "success",
			data: []*TestModel{testModel},
		},
		{
			name: "notFound",
			err:  gorm.ErrRecordNotFound,
		},
		{
			name: "empty",
			data: []*TestModel{},
		},
		{
			name: "ioError",
			err:  io.EOF,
		},
	}
	mock := mysql.NewMysqlMock((&TestModel{}).DBLabel())
	defer assert.Empty(t, mock.ExpectationsWereMet())
	repo := NewBaseRepo[*TestModel]()
	for _, test := range tests {
		mockProcess := mock.ExpectQuery("SELECT \\* FROM `test_tab` WHERE field1=\\?").
			WithArgs("mock")
		t.Run(test.name, func(t *testing.T) {
			if test.err != nil {
				mockProcess.WillReturnError(test.err)
			} else {
				mockProcess.WillReturnRows(sqlmockGormHelper.ModelToRows(test.data))
			}

			value, err := repo.Find(context.Background(), func(db *gorm.DB) *gorm.DB {
				return db.Where("field1=?", "mock")
			})
			if test.err != nil {
				assert.Empty(t, value)
				assert.ErrorIs(t, err, test.err)
			} else {
				assert.Empty(t, err)
				assert.Equal(t, true, assert.ObjectsAreEqual(value, test.data))
			}
		})
	}
}
