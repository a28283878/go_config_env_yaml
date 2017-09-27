package configy

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

var (
	useYaml bool
)

//Load : load date from yaml or os
func Load(out interface{}, file_path string) {
	var err error
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		log.Fatalf("Input struct is null or not pointer")
	}

	useYaml = true
	file, err := readFile(file_path)
	if err != nil {
		useYaml = false
	}

	if useYaml {
		loadYaml(file, out)
	} else {
		loadEnv(reflect.ValueOf(out).Elem())
	}
}

func loadYaml(file []byte, out interface{}) {
	err := yaml.Unmarshal(file, out)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func loadEnv(out reflect.Value) {
	for i := 0; i < out.NumField(); i++ {
		switch out.Field(i).Kind() {
		case reflect.Struct:
			loadEnv(out.Field(i))
		default:
			err := setValue(out.Field(i), os.Getenv(out.Type().Field(i).Tag.Get("env")))
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
	}
}

func setValue(out reflect.Value, env string) error {
	if len(env) == 0 {
		return nil
	}
	switch out.Kind() {
	case reflect.String:
		out.SetString(env)
	case reflect.Int:
		num, err := strconv.Atoi(env)
		if err != nil {
			return err
		}
		i := int64(num)
		if !out.OverflowInt(i) {
			out.SetInt(i)
		}
	case reflect.Bool:
		b, err := strconv.ParseBool(env)
		if err != nil {
			return err
		}
		out.SetBool(b)
	case reflect.Float32:
		f, err := strconv.ParseFloat(env, 32)
		if err != nil {
			return err
		}
		if !out.OverflowFloat(f) {
			out.SetFloat(f)
		}
	case reflect.Float64:
		f, err := strconv.ParseFloat(env, 64)
		if err != nil {
			return err
		}
		if !out.OverflowFloat(f) {
			out.SetFloat(f)
		}
	case reflect.Array:
		return (errors.New(out.Type().String() + " needs to be slice"))
	case reflect.Slice:
		switch out.Type().Elem().Kind() {
		case reflect.Int:
			var ints []int
			err := json.Unmarshal([]byte(env), &ints)
			if err != nil {
				return err
			}
			out.Set(reflect.ValueOf(ints))
		case reflect.String:
			var strings []string
			err := json.Unmarshal([]byte(env), &strings)
			if err != nil {
				return err
			}
			out.Set(reflect.ValueOf(strings))
		case reflect.Bool:
			var bools []bool
			err := json.Unmarshal([]byte(env), &bools)
			if err != nil {
				return err
			}
			out.Set(reflect.ValueOf(bools))
		case reflect.Float32:
			var floats []float32
			err := json.Unmarshal([]byte(env), &floats)
			if err != nil {
				return err
			}
			out.Set(reflect.ValueOf(floats))
		case reflect.Float64:
			var floats []float64
			err := json.Unmarshal([]byte(env), &floats)
			if err != nil {
				return err
			}
			out.Set(reflect.ValueOf(floats))
		default:
			return (errors.New("undefined type"))
		}

	default:
		return (errors.New("undefined type"))
	}

	return nil
}

func readFile(file_path string) ([]byte, error) {
	rootDirPath, err := filepath.Abs(filepath.Dir(file_path))
	if err != nil {
		log.Fatalf("file error: %v", err)
	}

	configPath := filepath.Join(rootDirPath, file_path)
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return file, err
	}
	return file, nil
}
