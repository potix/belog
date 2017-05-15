package belog

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

//ConfigLoggers is config of Loggers
type ConfigLoggers struct {
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

//LoadConfig is load configration file
func LoadConfig(configFilePath string) (err error) {
	configLoggers := new(ConfigLoggers)
	ext := filepath.Ext(configFilePath)
	switch ext {
	case ".tml":
		fallthrough
	case ".toml":
		_, err := toml.DecodeFile(configFilePath, configLoggers)
		if err != nil {
			return err
		}
	case ".yml":
		fallthrough
	case ".yaml":
		buf, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(buf, configLoggers)
		if err != nil {
			return err
		}
	case ".jsn":
		fallthrough
	case ".json":
		buf, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(buf, configLoggers)
		if err != nil {
			return err
		}
	default:
		return errors.Errorf("unexpected file extension (%v)", ext)
	}
	return SetupLoggers(configLoggers)
}

// SetupLoggers is setup from configLoggets
func SetupLoggers(configLoggers *ConfigLoggers) (err error) {
	tmpLoggers := make(map[string]*logger)
	for name, loggerConfig := range configLoggers.Loggers {
		// create filter
		filter, err := getFilter(loggerConfig.Filter.StructName)
		if err != nil {
			return errors.Errorf("not found filter (%v)", loggerConfig.Filter.StructName)
		}
		// setup filter
		if err = setupInstance(filter, loggerConfig.Filter); err != nil {
			return err
		}
		// create formatter
		formatter, err := getFormatter(loggerConfig.Formatter.StructName)
		if err != nil {
			return errors.Errorf("not found formatter (%v)", loggerConfig.Formatter.StructName)
		}
		// setup formatter
		if err = setupInstance(formatter, loggerConfig.Formatter); err != nil {
			return err
		}
		handlers := make([]Handler, 0, 0)
		for _, configStruct := range loggerConfig.Handlers {
			// create handler
			handler, err := getHandler(configStruct.StructName)
			if err != nil {
				return errors.Errorf("not found handler (%v)", configStruct.StructName)
			}
			// setup formatter
			if err = setupInstance(handler, configStruct); err != nil {
				return err
			}
			handlers = append(handlers, handler)
		}
		newLogger := &logger{
			filter:    filter,
			formatter: formatter,
			handlers:  handlers,
		}
		tmpLoggers[name] = newLogger
	}
	for name, newLogger := range tmpLoggers {
		err := SetLogger(name, newLogger.filter, newLogger.formatter, newLogger.handlers)
		if err != nil {
			// XXX statistics
		}
	}
	return nil
}

func setupInstance(instance interface{}, configStruct *configStruct) (err error) {
	for _, structSetter := range configStruct.StructSetters {
		instanceValue := reflect.ValueOf(instance)
		methodValue := instanceValue.MethodByName(strings.TrimSpace(structSetter.SetterName))
		if !methodValue.IsValid() {
			return errors.Errorf("unexpected Method (%v)", structSetter.SetterName)
		}
		methodType := methodValue.Type()
		argsNum := methodType.NumIn()
		if len(structSetter.SetterParams) != argsNum {
			return errors.Errorf("parameter count mismatch of setter method (%v: exp %v != act %v)", structSetter.SetterName, argsNum, len(structSetter.SetterParams))
		}
		outNum := methodType.NumOut()
		if outNum > 1 {
			return errors.Errorf("return value is too many of setter method")
		}
		methodArgs := make([]reflect.Value, 0, argsNum)
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
				return errors.Errorf("unsupported kind of setter paramter (%v)", argType.Kind())
			}
			methodArgs = append(methodArgs, reflectValue.Convert(argType))
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
			}
			return errors.Errorf("return value of setter method is not interface of error")
		}
	}
	return nil
}
