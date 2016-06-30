package belog

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	Loggers map[string]configLogger
}

type configLogger struct {
	Filter    *configFilter
	Formatter *configFormatter
	Handlers  []*configHandler
}

type configFilter struct {
	FilterName   string
	FilterParams []string
}

type configHandler struct {
	FormatterName   string
	FormatterParams []string
}

type configHandler struct {
	HandlerName   string
	HandlerParams []string
}

func LoadConfig(configFilePath string) (err error) {
	config = new(config)
	switch filePath.Ext(configFilePath) {
	case "tml":
		fallthrough
	case "toml":
		_, err := toml.DecodeFile(configFilePath, config)
		if err != nil {
			return err
		}
	case "yml":
		fallthrough
	case "yaml":
		buf, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(buf, config)
		if err != nil {
			return err
		}
	case "json":
		buf, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(buf, config)
		if err != nil {
			return err
		}
	default:
		return errors.Errorf("unexpected file extension")
	}
	setupLoggers(config)
}

func setupLoggers(config *config) (err error) {
	for name, loggerConfig := range config.Loggers {
		filter, err := filter.GetFilter(loggerConfig.Filter.FilterName)
		if err != nil {
			return errors.Errorf("not found filter (%v)", loggerConfig.Filter.FilterName)
		}
		for _, filterParam := range loggerConfig.Filter.FilterParams {
			filterMethod, filterArgs := parseParam(filterParam)
			filter.filterMethod(filterArgs)
		}

		formatter, err := formatter.GetFormatter(loggerConfig.Formatter.FormatterName)
		if err != nil {
			return errors.Errorf("not found formatter (%v)", loggerConfig.Formatter.FormatterName)
		}
		for _, formatterParam := range loggerConfig.Formatter.FormatterParams {
			formatterMethod, formatterArgs := parseParam(formatterParam)
			formatter.formtterMethod(formatterArgs)
		}

		for _, handlerConfig := range loggerConfig.handlers {
			handler, err := formatter.GetHAndler(loggerConfig.Handler.HandlerName)
			if err != nil {
				return errors.Errorf("not found formatter (%v)", loggerConfig.Handler.HandlerName)
			}
			for _, handlerParam := range loggerConfig.Handler.HandlerParams {
				handlerMethod, handlerArgs := parseParam(handlerParam)
				handler.handlerMethod(handlerArgs)
			}
		}
	}
	return nil
}
