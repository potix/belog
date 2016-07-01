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
	Filter    *configStruct
	Formatter *configStruct
	Handlers  []*configStruct
}

type configStruct struct {
	StructName    string
	StructSetters []*configStructSetter
}

type configStructSetter struct {
	SetterName   string
	SetterParams []string
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
	return setupLoggers(config)
}

func setupLoggers(config *config) (err error) {
	// XXXX
	tempLoggers := make(map[string]*logger)
	for name, loggerConfig := range config.Loggers {
		// create filter
		filter, err := filter.GetFilter(loggerConfig.Filter.StructName)
		if err != nil {
			return errors.Errorf("not found filter (%v)", loggerConfig.Filter.FilterName)
		}
		// setup filter
		if err := setupInstance(filter, loggerConfig.Filter); err != nil {
			return err
		}
		// create formatter
		formatter, err := formatter.GetFormatter(loggerConfig.Formatter.StructName)
		if err != nil {
			return errors.Errorf("not found formatter (%v)", loggerConfig.Formatter.FormatterName)
		}
		// setup formatter
		if err := setupInstance(formatter, loggerConfig.Formatter); err != nil {
			return err
		}
		handlers := make([]*handler.Handler, 0)
		for _, configStruct := range loggerConfig.Handlers {
			// create handler
			handler, err := formatter.GetHandler(configStruct.StructName)
			if err != nil {
				return errors.Errorf("not found handler (%v)", configStruct.StructName)
			}
			// setup formatter
			if err := setupInstance(handler, configStruct); err != nil {
				return err
			}
			handlers = append(handlers, handler)
		}
		newLogger := &logger{
			filter:    filter,
			formatter: formatter,
			handlers:  handlers,
		}
		temploggers[name] = newLogger
	}
	for name, newLogger := range tempLoggers {
		SetLogger(name, newLoger.filter, newLogger.formatter, newLogger.handlers)
	}
	return nil
}

func setupInstance(instance interface{}, configStruct *configStruct) {
	for _, structSetter := range configStruct.StructSetters {
		instanceValue := reflect.ValueOf(instance)
		methodValue := instanceValue.MethodByName(strings.TrimSpace(structSetter.SetterName))
		if !methodValue.IsValid() {
			return errors.Errorf("unexpected Method (%v)", loggerConfig.Formatter.FormatterName)
		}
		methodType := methodValue.Type()
		argsNum := methodType.NumIn()
		if len(structSetter.SetterParams) != argsNum {
			return errors.Errorf("parameter count mismatch")
		}
		methodArgs := make([]Value, 0, len(instanceArgs))
		for i, setterParam := range structSetter.SetterParams {
			argType := methodType.In(i)
			switch argType.Kind() {
			case int:
				// XXXXX:
				val, err := strconv(strings.Trimspace(setterParam), 10, 0)
				if err != nil {
					r
				}
				methodArgs = append(methodArgs, reflect.ValueOf(val))
			}
		}
		methodValue.Call(methodArgs)
	}
}
