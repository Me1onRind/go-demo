package mysql

import (
	"context"
	"io"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Me1onRind/go-demo/internal/infrastructure/unittest"
	"github.com/Me1onRind/go-demo/internal/model/po"
	sqlmockGormHelper "github.com/Me1onRind/sqlmock-gorm-helper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type TestModel2 struct {
	po.BaseModel
	Field1 string `gorm:"column:field1"`
}

func (t *TestModel2) TableName() string {
	return "test_tab"
}

var testModel2 *TestModel2 = &TestModel2{
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
		data *TestModel2
		err  error
	}{
		{
			name: "success",
			data: &TestModel2{Field1: "mck"},
		},
		{
			name: "fail",
			data: &TestModel2{Field1: "mck2"},
			err:  io.EOF,
		},
	}
	mock := NewMysqlMock("mock")
	defer assert.Empty(t, mock.ExpectationsWereMet())
	for _, test := range tests {
		mockProcess := mock.ExpectExec(sqlmockGormHelper.InsertSql(&TestModel2{}, `test_tab`)).
			WithArgs(unittest.Now{}, unittest.Now{}, test.data.Field1)
		t.Run(test.name, func(t *testing.T) {
			if test.err != nil {
				mockProcess.WillReturnError(test.err)
			} else {
				mockProcess.WillReturnResult(sqlmock.NewResult(100, 1))
			}

			value, err := Create(context.Background(), "mock", test.data)
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
		data *TestModel2
		err  error
	}{
		{
			name: "success",
			data: testModel2,
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
	mock := NewMysqlMock("mock")
	defer assert.Empty(t, mock.ExpectationsWereMet())
	for _, test := range tests {
		mockProcess := mock.ExpectQuery("SELECT \\* FROM `test_tab` WHERE id=\\? LIMIT 1").
			WithArgs(123456)
		t.Run(test.name, func(t *testing.T) {
			if test.err != nil {
				mockProcess.WillReturnError(test.err)
			} else {
				mockProcess.WillReturnRows(sqlmockGormHelper.ModelToRows(test.data))
			}

			value, err := Take[TestModel2](context.Background(), "mock", func(db *gorm.DB) *gorm.DB {
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
		data []*TestModel2
		err  error
	}{
		{
			name: "success",
			data: []*TestModel2{testModel2},
		},
		{
			name: "notFound",
			err:  gorm.ErrRecordNotFound,
		},
		{
			name: "empty",
			data: []*TestModel2{},
		},
		{
			name: "ioError",
			err:  io.EOF,
		},
	}
	mock := NewMysqlMock("mock")
	defer assert.Empty(t, mock.ExpectationsWereMet())
	for _, test := range tests {
		mockProcess := mock.ExpectQuery("SELECT \\* FROM `test_tab` WHERE field1=\\?").
			WithArgs("mock")
		t.Run(test.name, func(t *testing.T) {
			if test.err != nil {
				mockProcess.WillReturnError(test.err)
			} else {
				mockProcess.WillReturnRows(sqlmockGormHelper.ModelToRows(test.data))
			}

			value, err := Find[TestModel2](context.Background(), "mock", func(db *gorm.DB) *gorm.DB {
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
