package conf

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"

	"gopkg.in/yaml.v2"
)

// ConfigParser  parser yaml config file into config struct
func ConfParser(file string, in interface{}) (Conf, error) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		msg := fmt.Sprintf("failed to read '%s' with err: %s", file, err.Error())
		return Conf{}, errors.New(msg)
	}
	err = yaml.UnmarshalStrict(yamlFile, in)
	if err != nil {
		msg := fmt.Sprintf("failed to unmarshal '%s' with err: %s", file, err.Error())
		return Conf{}, errors.New(msg)
	}

	configs := reflect.ValueOf(in)
	if configs.Kind() == reflect.Ptr {
		configs = configs.Elem()
	}

	var config Conf
	values := make([]interface{}, configs.NumField())
	for i := 0; i < configs.NumField(); i++ {
		values[i] = configs.Field(i).Interface()
		switch values[i].(type) {
		case Protocol:
			proto := values[i].(Protocol)
			if proto.Enable == true {
				config.Proto = proto
			}
		case Database:
			db := values[i].(Database)
			config.Db = db
		}
	}
	return config, nil
}
