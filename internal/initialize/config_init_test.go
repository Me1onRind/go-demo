package initialize

import (
	"context"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/stretchr/testify/assert"
)

func Test_Load_Local_File_Cfg(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Dir(filename)
	if assert.Empty(t, InitFileConfig(dir+"/config_test.yaml")(context.Background())) {
		t.Log(gconfig.LocalFileCfg.Etcd)
		assert.Equal(t, []string{"127.0.0.1:1234"}, gconfig.LocalFileCfg.Etcd.Endpoints)
		assert.Equal(t, time.Second*5, gconfig.LocalFileCfg.Etcd.DialTimeout)
	}
}
