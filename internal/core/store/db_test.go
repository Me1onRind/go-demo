package store

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

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

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	dns := "me1onrind:guapi123@tcp(172.31.1.100:3306)/go-frame?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = NewConnectPool(dns)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func Test_Create(t *testing.T) {
	testTab := &TestTab{
		Name: fmt.Sprintf("%s_%d", "test", time.Now().Unix()),
	}
	if err := db.WithContext(context.Background()).Create(testTab).Error; err != nil {
		t.Fatal(err)
	}
}

func Test_Update(t *testing.T) {
	testTab := &TestTab{}
	if err := db.WithContext(context.Background()).Model(testTab).Where("id = ?", 1).Update("name", "test").Error; err != nil {
		t.Fatal(err)
	}
}

func Test_Select(t *testing.T) {
	testTab := &TestTab{}
	if err := db.WithContext(context.Background()).Where("id = ?", 1).Find(testTab).Error; err != nil {
		t.Fatal(err)
	}
	t.Log(testTab)
}
