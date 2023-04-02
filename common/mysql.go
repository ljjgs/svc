package common

import "github.com/asim/go-micro/v3/config"

type MysqlConfig struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	DataSource string `json:"datasource"`
}

func GetMysqlFromConsul(config config.Config, path ...string) *MysqlConfig {
	mysqlconfig := &MysqlConfig{}
	config.Get(path...).Scan(mysqlconfig)
	return mysqlconfig
}
