package conf

import (
	"io/ioutil"

	"github.com/northfun/house/common/utils/db"
	"github.com/northfun/house/common/utils/logger"
	"github.com/northfun/house/common/utils/redis"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var c Config

type Config struct {
	DB    db.Config     `yaml:"db" json:"db"`
	Redis redis.Config  `yaml:"redis" json:"redis"`
	Log   logger.Config `yaml:"log" json:"log"`
	Store bool          `yaml:"store" json:"store"`

	Ac string `json:"ac" yaml:"ac"`

	StartModule ModuleSwitch `yaml:"start_module" json:"start_module"`
}

type ModuleSwitch struct {
	House bool `yaml:"house"`
}

func Init(addr string) error {
	fileData, err := ioutil.ReadFile(addr)
	if err != nil {
		logger.Error("[conf],readfile", zap.Error(err))
		return err
	}
	return yaml.Unmarshal(fileData, &c)
}

func C() *Config {
	return &c
}
