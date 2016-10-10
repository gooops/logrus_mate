package graylog

import (
	"github.com/Sirupsen/logrus"
	"github.com/gogap/logrus_mate/hooks/graylog"

	"github.com/gooops/logrus_mate"
)

type GraylogHookConfig struct {
	Address  string                 `json:"address" yaml:"address"`
	Facility string                 `json:"facility" yaml:"facility"`
	Extra    map[string]interface{} `json:"extra" yaml:"extra"`
}

func init() {
	logrus_mate.RegisterHook("graylog", NewGraylogHook)
}

func NewGraylogHook(options logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := GraylogHookConfig{}

	if err = options.ToObject(&conf); err != nil {
		return
	}

	hook = graylog.NewGraylogHook(
		conf.Address,
		conf.Facility,
		conf.Extra)

	return
}
