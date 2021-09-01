package store

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	sqlmock_rows_helper "github.com/Me1onRind/sqlmock-rows-helper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type TestTab struct {
	ID    uint64 `gorm:"column:id"`
	Name  string `gorm:"column:name"`
	CTime uint32 `gorm:"column:ctime"`
	MTime uint32 `gorm:"column:mtime"`
}

func (t *TestTab) TableName() string {
	return "test_tab"
}

func newTestDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	mock.ExpectQuery("SELECT VERSION").WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).AddRow("5.7.32"))
	gormDB, err := NewDBConnectPoolFRromDB(db)
	if err != nil {
		return nil, nil, err
	}
	return gormDB, mock, nil
}

func Test_Create(t *testing.T) {
	db, mock, err := newTestDB()
	if !assert.Empty(t, err) {
		return
	}
	mock.ExpectExec("CREATE").WillReturnResult(sqlmock.NewResult(1, 1))
	testTab := &TestTab{
		Name: fmt.Sprintf("%s_%d", "test", time.Now().Unix()),
	}
	if err := db.WithContext(context.Background()).Create(testTab).Error; err != nil {
		assert.Empty(t, err)
	}
}

func Test_Update(t *testing.T) {
	db, mock, err := newTestDB()
	if !assert.Empty(t, err) {
		return
	}
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(0, 1))
	if !assert.Empty(t, err) {
		return
	}
	testTab := &TestTab{}
	if err := db.WithContext(context.Background()).Model(testTab).Where("id = ?", 1).Update("name", "test").Error; err != nil {
		assert.Empty(t, err)
	}
}

func Test_Select(t *testing.T) {
	db, mock, err := newTestDB()
	if !assert.Empty(t, err) {
		return
	}
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock_rows_helper.ModelToRows(
		&TestTab{
			ID:    1,
			Name:  "test",
			CTime: 1630250445,
			MTime: 1630250445,
		},
	))
	if !assert.Empty(t, err) {
		return
	}
	testTab := &TestTab{}
	if err := db.WithContext(context.Background()).Where("id = ?", 1).Find(testTab).Error; err != nil {
		if !assert.Empty(t, err) {
			return
		}
	}
	assert.Equal(t, uint64(1), testTab.ID)
	assert.Equal(t, "test", testTab.Name)
	assert.Equal(t, uint32(1630250445), testTab.CTime)
}
