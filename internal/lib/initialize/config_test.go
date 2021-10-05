package initialize

import (
	"os"
	"testing"

	"github.com/Me1onRind/go-demo/internal/core/config"
	"github.com/stretchr/testify/assert"
)

func Test_InitConfig(t *testing.T) {
	_ = InitLogger()
	os.Setenv("env", "local")
	err := InitLocalConfig("../../../conf")()
	if assert.Empty(t, err) {
		assert.Equal(t, []string{"localhost:2379"}, config.LocalConfig.Etcd.Endpoints)
	}
}
