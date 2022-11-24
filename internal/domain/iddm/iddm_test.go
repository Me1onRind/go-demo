package iddm

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/unittest"
	"github.com/Me1onRind/go-demo/internal/model/configmd"
	"github.com/Me1onRind/go-demo/internal/model/po"
	"github.com/Me1onRind/go-demo/internal/model/po/idpo"
	sqlmock_gorm_helper "github.com/Me1onRind/sqlmock-gorm-helper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	StepOne = &idpo.IdCreator{
		BaseModel: po.BaseModel{
			Id:         2,
			CreateTime: 1680000000,
			UpdateTime: 1680000000,
		},
		IdType: idpo.UserIdType,
		Offset: 10,
		Step:   1,
	}

	StepTen = &idpo.IdCreator{
		BaseModel: po.BaseModel{
			Id:         3,
			CreateTime: 1680000000,
			UpdateTime: 1680000000,
		},
		IdType: idpo.UserIdType,
		Offset: 10,
		Step:   10,
	}
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func Test_SingleInstance(t *testing.T) {
	assert.Equal(t, NewIdDomain(), NewIdDomain())
}

func Test_GetId(t *testing.T) {
	tests := []struct {
		name   string
		record *idpo.IdCreator
		step   int
	}{
		{
			name:   "step one",
			record: StepOne,
			step:   1,
		},
		{
			name:   "step ten",
			record: StepTen,
			step:   10,
		},
	}
	mock := mysql.NewMysqlMock(configmd.DefaultDBLabel)
	defer assert.Empty(t, mock.ExpectationsWereMet())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectQuery("SELECT \\* FROM `id_creator_tab` WHERE id_type=\\? LIMIT 1").WithArgs(1).
				WillReturnRows(sqlmock_gorm_helper.ModelToRows(test.record))

			mock.ExpectExec("UPDATE `id_creator_tab` SET `offset`=\\?,`update_time`=\\? WHERE id_type=\\? AND offset=\\?").
				WithArgs(10+test.step, unittest.Now{}, 1, 10).WillReturnResult(sqlmock.NewResult(0, 1))

			for i := 0; i < test.step; i++ {
				id, err := NewIdDomain().GetId(context.Background(), idpo.UserIdType, -1)
				//t.Log(id)
				assert.Empty(t, err)
				assert.EqualValues(t, 10+i, id)
			}
		})
	}

	t.Run("record not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM `id_creator_tab` WHERE id_type=\\? LIMIT 1").WithArgs(1).
			WillReturnError(gorm.ErrRecordNotFound)
		id, err := NewIdDomain().GetId(context.Background(), idpo.UserIdType, -1)
		assert.Empty(t, id)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("record step is zero", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM `id_creator_tab` WHERE id_type=\\? LIMIT 1").WithArgs(1).
			WillReturnRows(sqlmock_gorm_helper.ModelToRows(idpo.IdCreator{}))
		id, err := NewIdDomain().GetId(context.Background(), idpo.UserIdType, -1)
		assert.Empty(t, id)
		assert.ErrorIs(t, err, ErrPoolStepIsZero)
	})

	t.Run("Update zero row", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM `id_creator_tab` WHERE id_type=\\? LIMIT 1").WithArgs(1).
			WillReturnRows(sqlmock_gorm_helper.ModelToRows(StepOne))

		mock.ExpectExec("UPDATE `id_creator_tab` SET `offset`=\\?,`update_time`=\\? WHERE id_type=\\? AND offset=\\?").
			WithArgs(11, unittest.Now{}, 1, 10).WillReturnResult(sqlmock.NewResult(0, 0))
		id, err := NewIdDomain().GetId(context.Background(), idpo.UserIdType, -1)
		assert.Empty(t, id)
		assert.ErrorIs(t, err, ErrGetPoolFromDBFail)
	})

	t.Run("Update fail", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM `id_creator_tab` WHERE id_type=\\? LIMIT 1").WithArgs(1).
			WillReturnRows(sqlmock_gorm_helper.ModelToRows(StepOne))

		mock.ExpectExec("UPDATE `id_creator_tab` SET `offset`=\\?,`update_time`=\\? WHERE id_type=\\? AND offset=\\?").
			WithArgs(11, unittest.Now{}, 1, 10).WillReturnError(io.EOF)
		id, err := NewIdDomain().GetId(context.Background(), idpo.UserIdType, -1)
		assert.Empty(t, id)
		assert.ErrorIs(t, err, io.EOF)
	})
}
