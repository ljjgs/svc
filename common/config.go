package common

import (
	"github.com/asim/go-micro/plugins/config/source/consul/v3"
	"github.com/asim/go-micro/v3/config"
	"strconv"
)

func GetConsulConfig(host string, port int64, prefix string) (config.Config, error) {
	consulSource := consul.NewSource(
		consul.WithAddress(host+":"+strconv.FormatInt(port, 10)),
		consul.WithPrefix(prefix),
		consul.StripPrefix(true),
	)
	newConfig, err := config.NewConfig()

	if err != nil {
		return newConfig, err
	}
	err = newConfig.Load(consulSource)
	return newConfig, err
}
