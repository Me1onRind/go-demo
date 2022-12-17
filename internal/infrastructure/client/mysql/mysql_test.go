package mysql

import (
	"context"
	"io"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	sqlmock_gorm_helper "github.com/Me1onRind/sqlmock-gorm-helper"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TestModel struct {
	Id   uint64 `gorm:"id"`
	Name string `gorm:"name"`
}

func Test_Read_Write_Split(t *testing.T) {
	master, masterMock, _ := sqlmock.New()
	slave, slaveMock, _ := sqlmock.New()

	masterMock.ExpectQuery("SELECT VERSION").WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).AddRow("5.7.32"))
	slaveMock.ExpectQuery("SELECT VERSION").WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).AddRow("5.7.32"))

	masterDB, _ := newGormDBByDialector(mysql.New(mysql.Config{Conn: master}), &dbDesc{role: "master", dbLabel: "mock_db"})
	slaveDB, _ := newGormDBByDialector(mysql.New(mysql.Config{Conn: slave}), &dbDesc{role: "slave", dbLabel: "mock_db"})
	dbs.register("mock_db", newGormDBInfo(masterDB, []*gorm.DB{slaveDB}))

	defer func() {
		assert.Empty(t, masterMock.ExpectationsWereMet())
		assert.Empty(t, slaveMock.ExpectationsWereMet())
	}()

	ctx := context.Background()
	t.Run("test_read_replica", func(tt *testing.T) {
		slaveMock.ExpectQuery("SELECT \\* FROM `test_models` LIMIT 1").WillReturnError(gorm.ErrRecordNotFound)
		assert.Equal(tt, gorm.ErrRecordNotFound, GetReadDB(ctx, "mock_db").Take(&TestModel{}).Error)
	})

	t.Run("test_write_master", func(tt *testing.T) {
		masterMock.ExpectExec("UPDATE `test_models` SET `name`=\\? WHERE id=\\?").WithArgs("test", 1).
			WillReturnResult(sqlmock.NewResult(0, 0))
		assert.Empty(tt, GetWriteDB(ctx, "mock_db").Model(&TestModel{}).Where("id=?", 1).
			Updates(map[string]any{"name": "test"}).Error)
	})
}

func Test_Transaction(t *testing.T) {
	mock := NewMysqlMock("t_label")
	mock2 := NewMysqlMock("t_label_1")

	defer func() {
		assert.Empty(t, mock.ExpectationsWereMet())
		assert.Empty(t, mock2.ExpectationsWereMet())
	}()

	t.Run("no_t=>t=>no_t", func(t *testing.T) {
		ctx := context.Background()
		// select no use transaction
		mock.ExpectQuery("SELECT \\* FROM `test_models` LIMIT 1").WillReturnError(gorm.ErrRecordNotFound)
		assert.Equal(t, gorm.ErrRecordNotFound, GetWriteDB(ctx, "t_label").Take(&TestModel{}).Error)

		// transaction, begin -> select -> update -> commit
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT \\* FROM `test_models` LIMIT 1").WillReturnRows(sqlmock_gorm_helper.ModelToRows(TestModel{Id: 1, Name: "test"}))
		mock.ExpectExec("UPDATE `test_models` SET `name`=\\? WHERE id>\\?").WithArgs("new name", 50).
			WillReturnResult(sqlmock.NewResult(0, 100))
		mock.ExpectCommit()
		fc := func(ctx context.Context) error {
			assert.Empty(t, GetWriteDB(ctx, "t_label").Take(&TestModel{}).Error)
			assert.Empty(t, GetWriteDB(ctx, "t_label").Model(&TestModel{}).Where("id>?", 50).
				Updates(map[string]any{"name": "new name"}).Error)
			return nil
		}
		assert.Empty(t, Transaction(ctx, "t_label", fc))
		// select no use transaction
		mock.ExpectQuery("SELECT \\* FROM `test_models` LIMIT 1").WillReturnError(gorm.ErrRecordNotFound)
		assert.Equal(t, gorm.ErrRecordNotFound, GetWriteDB(ctx, "t_label").Take(&TestModel{}).Error)
	})

	t.Run("transaction rollback", func(t *testing.T) {
		// transaction, begin -> select -> rollback
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT \\* FROM `test_models` LIMIT 1").WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()
		fc2 := func(ctx context.Context) error {
			return GetWriteDB(ctx, "t_label").Take(&TestModel{}).Error
		}
		assert.Equal(t, gorm.ErrRecordNotFound, Transaction(context.Background(), "t_label", fc2))
	})

	t.Run("transaction rollback error", func(t *testing.T) {
		// transaction, begin -> select -> rollback
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT \\* FROM `test_models` LIMIT 1").WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback().WillReturnError(io.EOF)
		fc2 := func(ctx context.Context) error {
			return GetWriteDB(ctx, "t_label").Take(&TestModel{}).Error
		}
		assert.Equal(t, gorm.ErrRecordNotFound, Transaction(context.Background(), "t_label", fc2))
	})

	t.Run("nest_transaction", func(t *testing.T) {
		// transaction nest
		mock.ExpectBegin()
		mock.ExpectCommit()
		fc3 := func(ctx context.Context) error {
			return nil
		}
		fc4 := func(ctx context.Context) error {
			return Transaction(ctx, "t_label", fc3)
		}
		assert.Empty(t, Transaction(context.Background(), "t_label", fc4))
	})

	t.Run("nest transaction diff db", func(t *testing.T) {
		// transaction nest
		mock.ExpectBegin()
		mock2.ExpectBegin()
		mock2.ExpectCommit()
		mock.ExpectCommit()
		fc3 := func(ctx context.Context) error {
			return nil
		}
		fc4 := func(ctx context.Context) error {
			return Transaction(ctx, "t_label_1", fc3)
		}
		assert.Empty(t, Transaction(context.Background(), "t_label", fc4))
	})

	t.Run("transaction mixed not transaction db", func(t *testing.T) {
		// transaction nest
		mock.ExpectQuery("SELECT \\* FROM `test_models` LIMIT 1").WillReturnError(gorm.ErrRecordNotFound)
		mock2.ExpectBegin()
		mock2.ExpectCommit()
		fc5 := func(ctx context.Context) error {
			assert.Equal(t, gorm.ErrRecordNotFound, GetWriteDB(ctx, "t_label").Take(&TestModel{}).Error)
			return nil
		}
		assert.Empty(t, Transaction(context.Background(), "t_label_1", fc5))
	})

	t.Run("commit return err", func(t *testing.T) {
		// transaction nest
		mock.ExpectQuery("SELECT \\* FROM `test_models` LIMIT 1").WillReturnError(gorm.ErrRecordNotFound)
		mock2.ExpectBegin()
		mock2.ExpectCommit().WillReturnError(gorm.ErrInvalidTransaction)
		fc5 := func(ctx context.Context) error {
			assert.Equal(t, gorm.ErrRecordNotFound, GetWriteDB(ctx, "t_label").Take(&TestModel{}).Error)
			return nil
		}
		assert.Equal(t, gorm.ErrInvalidTransaction, Transaction(context.Background(), "t_label_1", fc5))
	})

}

func Test_GetDB_Panic(t *testing.T) {
	assert.Panics(t, func() {
		GetWriteDB(context.Background(), "db")
	},
	)
}

//func Test_SetDBLabel(t *testing.T) {
//assert.Empty(t, setDBLabel(&gorm.DB{}, "d", true))
//tests := []struct {
//name  string
//label string
//}{
//{
//name:  "label zero",
//label: "",
//},
//{
//name:  "duplicate",
//label: "d",
//},
//}

//for _, test := range tests {
//t.Run(test.name, func(t *testing.T) {
//assert.NotEmpty(t, setDBLabel(&gorm.DB{}, test.label, true))
//})
//}
//}
