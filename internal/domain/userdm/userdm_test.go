package userdm

import (
	"context"
	"io"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Me1onRind/go-demo/internal/domain/iddm"
	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/unittest"
	"github.com/Me1onRind/go-demo/internal/model/configmd"
	"github.com/Me1onRind/go-demo/internal/model/errors"
	"github.com/Me1onRind/go-demo/internal/model/po/idpo"
	"github.com/Me1onRind/go-demo/internal/model/po/userpo"
	"github.com/Me1onRind/go-demo/internal/repository/userrepo"
	"github.com/Me1onRind/go-demo/protocol/code"
	sqlmockGormHelper "github.com/Me1onRind/sqlmock-gorm-helper"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mysql.NewMysqlMock(configmd.DefaultDBLabel)
	defer assert.Empty(t, mock.ExpectationsWereMet())

	userDomain := NewUserDomain()
	idDomain := iddm.NewMockIdDomain(ctrl)
	userDomain.IdDomain = idDomain

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM `user_tab` WHERE email=\\? LIMIT 1").WithArgs("test@google.com").WillReturnError(gorm.ErrRecordNotFound)
		idDomain.EXPECT().GetId(gomock.Any(), idpo.UserIdType, 5).Return(uint64(167), nil)
		mock.ExpectExec(sqlmockGormHelper.InsertSql(userpo.User{}, `user_tab`)).
			WithArgs(unittest.Now{}, unittest.Now{}, 167, "test@google.com", "test_j").
			WillReturnResult(sqlmock.NewResult(100, 1))

		user, err := userDomain.CreateUser(context.Background(), &userpo.User{
			Name:  "test_j",
			Email: "test@google.com",
		})
		assert.Empty(t, err)
		assert.EqualValues(t, 100, user.Id)
	})

	t.Run("write_db_error", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM `user_tab` WHERE email=\\? LIMIT 1").WithArgs("test@google.com").WillReturnError(gorm.ErrRecordNotFound)
		idDomain.EXPECT().GetId(gomock.Any(), idpo.UserIdType, 5).Return(uint64(167), nil)
		mock.ExpectExec(sqlmockGormHelper.InsertSql(userpo.User{}, `user_tab`)).
			WithArgs(unittest.Now{}, unittest.Now{}, 167, "test@google.com", "test_j").
			WillReturnError(io.EOF)

		user, err := userDomain.CreateUser(context.Background(), &userpo.User{
			Name:  "test_j",
			Email: "test@google.com",
		})
		assert.Empty(t, user)
		assert.ErrorIs(t, err, io.EOF)
	})

	t.Run("id_creator_error", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM `user_tab` WHERE email=\\? LIMIT 1").WithArgs("test@google.com").WillReturnError(gorm.ErrRecordNotFound)
		idDomain.EXPECT().GetId(gomock.Any(), idpo.UserIdType, 5).Return(uint64(0), iddm.ErrGetPoolFromDBFail)

		user, err := userDomain.CreateUser(context.Background(), &userpo.User{
			Name:  "test_j",
			Email: "test@google.com",
		})
		assert.Empty(t, user)
		assert.ErrorIs(t, err, iddm.ErrGetPoolFromDBFail)
	})

	t.Run("read_db_error", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM `user_tab` WHERE email=\\? LIMIT 1").WithArgs("test@google.com").WillReturnError(io.EOF)
		user, err := userDomain.CreateUser(context.Background(), &userpo.User{
			Name:  "test_j",
			Email: "test@google.com",
		})
		assert.Empty(t, user)
		assert.ErrorIs(t, err, io.EOF)
		assert.ErrorIs(t, err, gerror.ReadDBError)
	})

	t.Run("existed_error", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM `user_tab` WHERE email=\\? LIMIT 1").WithArgs("test@google.com").WillReturnRows(sqlmockGormHelper.ModelToRows(userrepo.UserJ))
		user, err := userDomain.CreateUser(context.Background(), &userpo.User{
			Name:  "test_j",
			Email: "test@google.com",
		})
		assert.Empty(t, user)
		assert.Equal(t, code.RecordExisted, errors.ExtractError(err).Code)
	})
}
