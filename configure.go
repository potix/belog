package belog

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"strings"
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

func setupInstance(instance interface{}, configStruct *configStruct) (err error) {
	for _, structSetter := range configStruct.StructSetters {
		instanceValue := reflect.ValueOf(instance)
		methodValue := instanceValue.MethodByName(strings.TrimSpace(structSetter.SetterName))
		if !methodValue.IsValid() {
			return errors.Errorf("unexpected Method (%v)", loggerConfig.Formatter.FormatterName)
		}
		methodType := methodValue.Type()
		argsNum := methodType.NumIn()
		if len(structSetter.SetterParams) != argsNum {
			return errors.Errorf("parameter count mismatch of setter method")
		}
		outNum := typeOfFunc.NumOut()
		if outNum > 1 {
			return errors.Errorf("return value is too many of setter method")
		}
		methodArgs := make([]Value, 0, len(instanceArgs))
		for i, setterParam := range structSetter.SetterParams {
			argType := methodType.In(i)
			var reflectValue reflect.Value
			switch argType.Kind() {
			case reflect.Bool:
				val, err := strconv.ParseBool(setterParam)
				if err != nil {
					return err
				}
				reflectValue = reflect.ValueOf(val)
			case reflect.Int:
				val, err := strconv.ParseInt(setterParam, 10, 0)
				if err != nil {
					return err
				}
				reflectValue = reflect.ValueOf(int(val))
			case reflect.Int8:
				val, err := strconv.ParseInt(setterParam, 10, 8)
				if err != nil {
					return err
				}
				reflectValue = reflect.ValueOf(int8(val))
			case reflect.Int16:
				val, err := strconv.ParseInt(setterParam, 10, 16)
				if err != nil {
					return err
				}
				reflectValue = reflect.ValueOf(int16(val))
			case reflect.Int32:
				val, err := strconv.ParseInt(setterParam, 10, 32)
				if err != nil {
					return err
				}
				reflectValue = reflect.ValueOf(int32(val))
			case reflect.Int64:
				val, err := strconv.ParseInt(setterParam, 10, 64)
				if err != nil {
					return err
				}
				reflectValue = reflect.ValueOf(val)
			case reflect.Uint:
				val, err := strconv.ParseUint(setterParam, 10, 0)
				if err != nil {
					return err
				}
				reflectValue = reflect.ValueOf(uint(val))
			case reflect.Uint8:
				val, err := strconv.ParseUint(setterParam, 10, 8)
				if err != nil {
					return err
				}
				reflectValue = reflect.ValueOf(uint8(val))
			case reflect.Uint16:
				val, err := strconv.ParseUint(setterParam, 10, 16)
				if err != nil {
					return err
				}
				reflectValue = reflect.ValueOf(uint16(val))
			case reflect.Uint32:
				val, err := strconv.ParseUint(setterParam, 10, 32)
				if err != nil {
					return err
				}
				reflectValue = reflect.ValueOf(uint32(val))
			case reflect.Uint64:
				val, err := strconv.ParseUint(setterParam, 10, 64)
				if err != nil {
					return err
				}
				reflectValue = reflect.ValueOf(val)
			case reflect.Float32:
				val, err := strconv.ParseFloat(setterParam, 32)
				if err != nil {
					return err
				}
				reflectValue = reflect.ValueOf(float32(val))
			case reflect.Float64:
				val, err := strconv.ParseFloat(setterParam, 64)
				if err != nil {
					return err
				}
				reflectValue = reflect.ValueOf(val)
			case reflect.String:
				reflectValue = reflect.ValueOf(setterParam)
			default:
				return errors.Errorf("unsupported kind of setter paramter", argType.Kind())
			}
			methodArgs = append(methodArgs, reflectValue)
		}
		outs := methodValue.Call(methodArgs)
		if len(outs) == 1 {
			out := outs[0]
			if out.Kind() != reflect.Interface {
				return errors.Errorf("return value of setter method is not interface of error")
			}
			outType := out.Type()
			errorInterface := reflect.TypeOf((*error)(nil)).Elem()
			if outType.Implements(errorInterface) {
				return out.Interface().(error)
			} else {
				return errors.Errorf("return value of setter method is not interface of error")
			}
		}
	}
}
