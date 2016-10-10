package bugsnag

import (
	"github.com/Shopify/logrus-bugsnag"
	"github.com/Sirupsen/logrus"
	"github.com/bugsnag/bugsnag-go"

	"github.com/gooops/logrus_mate"
)

type BugsnagHookConfig struct {
	Endpoint     string `json:"endpoint" yaml:"endpoint"`
	ReleaseStage string `json:"release_stage" yaml:"release_stage"`
	APIKey       string `json:"api_key" yaml:"api_key"`
	Synchronous  bool   `json:"synchronous" yaml:"synchronous"`
}

func init() {
	logrus_mate.RegisterHook("bugsnag", NewBugsnagHook)
}

func NewBugsnagHook(options logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := BugsnagHookConfig{}

	if err = options.ToObject(&conf); err != nil {
		return
	}

	bugsnag.Configure(bugsnag.Configuration{
		Endpoint:     conf.Endpoint,
		ReleaseStage: conf.ReleaseStage,
		APIKey:       conf.APIKey,
		Synchronous:  conf.Synchronous,
	})

	hook, err = logrus_bugsnag.NewBugsnagHook()
	return
}
