package mysql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Me1onRind/sqlmock-rows-helper"
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

	masterResource := mysql.New(mysql.Config{Conn: master})
	slaveResource := mysql.New(mysql.Config{Conn: slave})

	db, err := newMysqlClusterClientByDialector(masterResource, []gorm.Dialector{slaveResource})
	if !assert.Empty(t, err) {
		return
	}

	slaveMock.ExpectQuery("SELECT \\* FROM `test_models` LIMIT 1").WillReturnError(gorm.ErrRecordNotFound)
	assert.Equal(t, gorm.ErrRecordNotFound, db.Take(&TestModel{}).Error)

	masterMock.ExpectExec("UPDATE `test_models` SET `name`=\\? WHERE id=\\?").WithArgs("test", 1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	assert.Empty(t, db.Model(&TestModel{}).Where("id=?", 1).
		Updates(map[string]any{"name": "test"}).Error)

	assert.Empty(t, masterMock.ExpectationsWereMet())
	assert.Empty(t, slaveMock.ExpectationsWereMet())
}

func Test_Transaction(t *testing.T) {
	master, mock, _ := sqlmock.New()
	mock.ExpectQuery("SELECT VERSION").WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).AddRow("5.7.32"))
	masterResource := mysql.New(mysql.Config{Conn: master})
	master2, mock2, _ := sqlmock.New()
	mock2.ExpectQuery("SELECT VERSION").WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).AddRow("5.7.32"))
	masterResource2 := mysql.New(mysql.Config{Conn: master2})

	db, _ := newMysqlClusterClientByDialector(masterResource, []gorm.Dialector{})
	db2, _ := newMysqlClusterClientByDialector(masterResource2, []gorm.Dialector{})
	assert.Empty(t, setDBLabel(db, "t_label"))
	assert.Empty(t, setDBLabel(db2, "t_label_1"))

	t.Run("no_t=>t=>no_t", func(t *testing.T) {
		ctx := context.Background()
		// select no use transaction
		mock.ExpectQuery("SELECT \\* FROM `test_models` LIMIT 1").WillReturnError(gorm.ErrRecordNotFound)
		assert.Equal(t, gorm.ErrRecordNotFound, GetDB(ctx, "t_label").Take(&TestModel{}).Error)

		// transaction, begin -> select -> update -> commit
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT \\* FROM `test_models` LIMIT 1").WillReturnRows(sqlmock_rows_helper.ModelToRows(TestModel{Id: 1, Name: "test"}))
		mock.ExpectExec("UPDATE `test_models` SET `name`=\\? WHERE id>\\?").WithArgs("new name", 50).
			WillReturnResult(sqlmock.NewResult(0, 100))
		mock.ExpectCommit()
		fc := func(ctx context.Context) error {
			assert.Empty(t, GetDB(ctx, "t_label").Take(&TestModel{}).Error)
			assert.Empty(t, GetDB(ctx, "t_label").Model(&TestModel{}).Where("id>?", 50).
				Updates(map[string]any{"name": "new name"}).Error)
			return nil
		}
		assert.Empty(t, Transaction(ctx, "t_label", fc))
		// select no use transaction
		mock.ExpectQuery("SELECT \\* FROM `test_models` LIMIT 1").WillReturnError(gorm.ErrRecordNotFound)
		assert.Equal(t, gorm.ErrRecordNotFound, GetDB(ctx, "t_label").Take(&TestModel{}).Error)
	})

	t.Run("transaction rollback", func(t *testing.T) {
		// transaction, begin -> select -> rollback
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT \\* FROM `test_models` LIMIT 1").WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()
		fc2 := func(ctx context.Context) error {
			return GetDB(ctx, "t_label").Take(&TestModel{}).Error
		}
		assert.Equal(t, gorm.ErrRecordNotFound, Transaction(context.Background(), "t_label", fc2))
	})

	t.Run("nest transaction", func(t *testing.T) {
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
			assert.Equal(t, gorm.ErrRecordNotFound, GetDB(ctx, "t_label").Take(&TestModel{}).Error)
			return nil
		}
		assert.Empty(t, Transaction(context.Background(), "t_label_1", fc5))
	})

	assert.Empty(t, mock.ExpectationsWereMet())
}
