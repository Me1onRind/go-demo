package idrepo

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/unittest"
	"github.com/Me1onRind/go-demo/internal/model/po"
	"github.com/Me1onRind/go-demo/internal/model/po/idpo"
	sqlmock_gorm_helper "github.com/Me1onRind/sqlmock-gorm-helper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	UserIdIdCreator = &idpo.IdCreator{
		BaseModel: po.BaseModel{
			Id:         1,
			CreateTime: 1680000000,
			UpdateTime: 1680000000,
		},
		IdType: idpo.UserIdType,
		Offset: 1000,
		Step:   10,
	}
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func Test_GetRecord(t *testing.T) {
	tests := []struct {
		name string
		data *idpo.IdCreator
		err  error
	}{
		{
			name: "success",
			data: UserIdIdCreator,
		},
		{
			name: "record not found",
			err:  gorm.ErrRecordNotFound,
		},
	}

	idRepo := NewIdRepo()
	mock := mysql.NewMysqlMock(idRepo.DBLabel())
	defer assert.Empty(t, mock.ExpectationsWereMet())
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.err == nil {
				mock.ExpectQuery("SELECT \\* FROM `id_creator_tab` WHERE id_type=\\? LIMIT 1").WithArgs(1).
					WillReturnRows(sqlmock_gorm_helper.ModelToRows(UserIdIdCreator))
			} else {
				mock.ExpectQuery("SELECT \\* FROM `id_creator_tab` WHERE id_type=\\? LIMIT 1").WithArgs(1).
					WillReturnError(test.err)
			}

			record, err := idRepo.GetRecord(context.Background(), WithIdType(idpo.UserIdType))
			if test.err == nil {
				assert.Empty(t, err)
				assert.Equal(t, true, assert.ObjectsAreEqual(record, test.data))
			} else {
				assert.Empty(t, record)
				assert.ErrorIs(t, err, test.err)
			}
		})
	}
}

func Test_UpdateRecord(t *testing.T) {
	tests := []struct {
		name string
		rows int64
		err  error
	}{
		{
			name: "success",
			rows: 1,
		},
		{
			name: "io error",
			err:  io.EOF,
		},
		{
			name: "update zero row",
			rows: 0,
		},
	}

	idRepo := NewIdRepo()
	mock := mysql.NewMysqlMock(idRepo.DBLabel())
	defer assert.Empty(t, mock.ExpectationsWereMet())
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.err == nil {
				mock.ExpectExec("UPDATE `id_creator_tab` SET `offset`=\\?,`update_time`=\\? WHERE id_type=\\? AND offset=\\?").
					WithArgs(1500, unittest.Now{}, 1, 1000).WillReturnResult(sqlmock.NewResult(0, test.rows))
			} else {
				mock.ExpectExec("UPDATE `id_creator_tab` SET `offset`=\\?,`update_time`=\\? WHERE id_type=\\? AND offset=\\?").
					WithArgs(1500, unittest.Now{}, 1, 1000).WillReturnError(test.err)
			}

			rows, err := idRepo.UpdateOffset(context.Background(), idpo.UserIdType, 1000, 500)
			if test.err == nil {
				assert.Empty(t, err)
				assert.Equal(t, test.rows, rows)
			} else {
				assert.Empty(t, rows)
				assert.ErrorIs(t, err, test.err)
			}
		})
	}
}
