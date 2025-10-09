package utils

import (
	"os"
	"reflect"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func BindFromJson(dest any, filename, path string) error {
	v := viper.New()
	v.SetConfigType("json")
	v.AddConfigPath(path)
	v.SetConfigName(filename)

	err := v.ReadInConfig()
	if err != nil {
		return err
	}
	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("Failed to unmarshal config: %v", err)
		return err
	}
	return nil
}

func setEnvFromConsulKV(v *viper.Viper) error {
	env := make(map[string]any)
	err := v.Unmarshal(&env)
	if err != nil {
		logrus.Errorf("Failed to unmarshal config: %v", err)
		return err
	}
	for key, value := range env {
		var (
			vaLOf = reflect.ValueOf(value)
			val   string
		)
		switch vaLOf.Kind() {
		case reflect.String:
			val = vaLOf.String()
		case reflect.Int:
			val = strconv.Itoa(int(vaLOf.Int()))
		case reflect.Uint:
			val = strconv.Itoa(int(vaLOf.Uint()))
		case reflect.Float32:
			val = strconv.Itoa(int(vaLOf.Float()))
		case reflect.Bool:
			val = strconv.FormatBool(vaLOf.Bool())
		default:
			panic("unsupported type")
		}
		err := os.Setenv(key, val)
		if err != nil {
			logrus.Errorf("Failed to set env variable %s: %v", key, err)
			return err
		}
	}
	return nil
}

func BindFromConsul(dest any, endPoint, path string) error {
	v := viper.New()
	v.SetConfigType("json")
	err := v.AddRemoteProvider("consul", endPoint, path)
	if err != nil {
		logrus.Errorf("Failed to add remote provider: %v", err)
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		logrus.Errorf("Failed to read remote config: %v", err)
		return err
	}

	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("Failed to unmarshal config: %v", err)
		return err
	}

	err = setEnvFromConsulKV(v)
	if err != nil {
		logrus.Errorf("Failed to set env from Consul KV: %v", err)
		return err
	}
	return nil
}
