package mail

import (
	"github.com/Sirupsen/logrus"
	"github.com/zbindenren/logrus_mail"

	"github.com/gooops/logrus_mate"
)

type MailHookConfig struct {
	AppName  string `json:"app_name" yaml:"app_name"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	From     string `json:"from" yaml:"from"`
	To       string `json:"to" yaml:"to"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func init() {
	logrus_mate.RegisterHook("mail", NewMailHook)
}

func NewMailHook(options logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := MailHookConfig{}
	if err = options.ToObject(&conf); err != nil {
		return
	}

	hook, err = logrus_mail.NewMailAuthHook(
		conf.AppName,
		conf.Host,
		conf.Port,
		conf.From,
		conf.To,
		conf.Username,
		conf.Password)

	return
}
