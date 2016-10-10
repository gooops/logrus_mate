package syslog

import (
	"log/syslog"

	"github.com/Sirupsen/logrus"
	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"

	"github.com/gooops/logrus_mate"
)

type SyslogHookConfig struct {
	Network  string `json:"network" yaml:"network"`
	Address  string `json:"address" yaml:"address"`
	Priority string `json:"priority" yaml:"priority"`
	Tag      string `json:"tag" yaml:"tag"`
}

func init() {
	logrus_mate.RegisterHook("syslog", NewSyslogHook)
}

func NewSyslogHook(options logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := SyslogHookConfig{}
	if err = options.ToObject(&conf); err != nil {
		return
	}

	return logrus_syslog.NewSyslogHook(
		conf.Network,
		conf.Address,
		toPriority(conf.Priority),
		conf.Tag)
}

func toPriority(priority string) syslog.Priority {
	switch priority {
	case "LOG_EMERG":
		return syslog.LOG_EMERG
	case "LOG_ALERT":
		return syslog.LOG_ALERT
	case "LOG_CRIT":
		return syslog.LOG_CRIT
	case "LOG_ERR":
		return syslog.LOG_ERR
	case "LOG_WARNING":
		return syslog.LOG_WARNING
	case "LOG_NOTICE":
		return syslog.LOG_NOTICE
	case "LOG_INFO":
		return syslog.LOG_INFO
	case "LOG_DEBUG":
		return syslog.LOG_DEBUG
	}
	return syslog.LOG_DEBUG
}
