package logrus_mate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gogap/env_json"
	"gopkg.in/yaml.v2"
)

type Environments struct {
	RunEnv string `json:"run_env" yaml:"run_env"`
}

type FormatterConfig struct {
	Name    string  `json:"name" yaml:"name"`
	Options Options `json:"options" yaml:"options"`
}

type HookConfig struct {
	Name    string  `json:"name" yaml:"name"`
	Options Options `json:"options" yaml:"options"`
}

type WriterConfig struct {
	Name    string  `json:"name" yaml:"name"`
	Options Options `json:"options" yaml:"options"`
}

type LoggerItem struct {
	Name   string                  `json:"name" yaml:"name"`
	Config map[string]LoggerConfig `json:"config" yaml:"config"`
}

type LoggerConfig struct {
	Out       WriterConfig    `json:"out" yaml:"out"`
	Level     string          `json:"level" yaml:"level"`
	Hooks     []HookConfig    `json:"hooks" yaml:"hooks"`
	Formatter FormatterConfig `json:"formatter" yaml:"formatter"`
}

type LogrusMateConfig struct {
	EnvironmentKeys Environments `json:"env_keys" yaml:"env_keys"`
	Loggers         []LoggerItem `json:"loggers" yaml:"loggers"`
	configType      string
}

func (p *LogrusMateConfig) Serialize() (data []byte, err error) {
	if p.configType == "json" {
		data, err = json.Marshal(p)
	} else if p.configType == "yaml" {
		data, err = yaml.Marshal(p)
	}
	return
}

func LoadLogrusMateConfig(filename string) (conf LogrusMateConfig, err error) {
	var data []byte

	if data, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	if strings.HasSuffix(filename, ".json") || strings.HasSuffix(filename, ".conf") {
		conf.configType = "json"
		envJSON := env_json.NewEnvJson(env_json.ENV_JSON_KEY, env_json.ENV_JSON_EXT)

		if err = envJSON.Unmarshal(data, &conf); err != nil {
			return
		}
	} else if strings.HasSuffix(filename, ".yml") || strings.HasSuffix(filename, ".yaml") {
		conf.configType = "yaml"
		if err = yaml.Unmarshal(data, &conf); err != nil {
			return
		}
	}

	return
}

func (p *LogrusMateConfig) Validate() (err error) {
	for _, logger := range p.Loggers {
		for envName, conf := range logger.Config {
			if _, err = logrus.ParseLevel(conf.Level); err != nil {
				return
			}

			if conf.Hooks != nil {
				for id, hook := range conf.Hooks {
					if newFunc, exist := newHookFuncs[hook.Name]; !exist {
						err = fmt.Errorf("logurs mate: hook not registered, env: %s, id: %d, name: %s", envName, id, hook)
						return
					} else if newFunc == nil {
						err = fmt.Errorf("logurs mate: hook's func is damaged, env: %s, id: %d name: %s", envName, id, hook)
						return
					}
				}
			}

			if conf.Formatter.Name != "" {
				if newFunc, exist := newFormatterFuncs[conf.Formatter.Name]; !exist {
					err = fmt.Errorf("logurs mate: formatter not registered, env: %s, name: %s", envName, conf.Formatter.Name)
					return
				} else if newFunc == nil {
					err = fmt.Errorf("logurs mate: formatter's func is damaged, env: %s, name: %s", envName, conf.Formatter.Name)
					return
				}
			}

			if conf.Out.Name != "" {
				if newFunc, exist := newWriterFuncs[conf.Out.Name]; !exist {
					err = fmt.Errorf("logurs mate: writter not registered, env: %s, name: %s", envName, conf.Out.Name)
					return
				} else if newFunc == nil {
					err = fmt.Errorf("logurs mate: writter's func is damaged, env: %s, name: %s", envName, conf.Out.Name)
					return
				}
			}
		}
	}
	return
}

func (p *LogrusMateConfig) RunEnv() string {
	env := os.Getenv(p.EnvironmentKeys.RunEnv)
	if env == "" {
		env = "development"
	}
	return env
}
