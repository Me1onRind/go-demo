package initialize

import (
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/Me1onRind/go-demo/internal/global/config"
	"github.com/stretchr/testify/assert"
)

func Test_Load_Local_File_Cfg(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Dir(filename)
	if assert.Empty(t, InitFileConfig(dir+"/config_test.yaml")()) {
		t.Log(config.LocalFileCfg.Etcd)
		assert.Equal(t, []string{"127.0.0.1:1234"}, config.LocalFileCfg.Etcd.Endpoints)
		assert.Equal(t, time.Second*5, config.LocalFileCfg.Etcd.DialTimeout)
	}
}
