package config

import (
	"fmt"
	"time"
)

type MysqlResources struct {
	Master     MysqlConfig     `yaml:"master"`
	Slaves     []MysqlConfig   `yaml:"slaves"`
	MasterPool MysqlPoolConfig `yaml:"master_pool"`
	SlavePool  MysqlPoolConfig `yaml:"slave_pool"`
}

type MysqlPoolConfig struct {
	MaxIdleConns int           `yaml:"max_idle_conns"`
	MaxOpenConns int           `yaml:"max_open_conns"`
	MaxIdleTime  time.Duration `yaml:"max_idle_time"`
	MaxLifetime  time.Duration `yaml:"max_lifetime"`
}

type MysqlConfig struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	DBName       string `yaml:"dbname"`
	Timeout      string `yaml:"timeout"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
}

func (m *MysqlConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s",
		m.Username, m.Password, m.Host, m.Port, m.DBName, m.Timeout, m.ReadTimeout, m.WriteTimeout)
}
