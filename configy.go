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
	useYaml  bool
	yamlName string
)

func init() {
	yamlName = "config.yml"
}

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
		log.Printf("cannot find the file \"%s\", use os environment variable instead", yamlName)
		useYaml = false
	}

	if useYaml {
		loadYaml(file, out)
	} else {
		loadEnv(reflect.ValueOf(out).Elem())
	}
}

//SetFileName : set custom yaml name
func SetFileName(newFileName string) {
	yamlName = newFileName
}

func loadYaml(file []byte, out interface{}) {
	err := yaml.Unmarshal(file, out)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
}

func loadEnv(out reflect.Value) {
	for i := 0; i < out.NumField(); i++ {
		switch out.Field(i).Kind() {
		case reflect.Struct:
			loadEnv(out.Field(i))
		default:
			err := setValue(out.Field(i), os.Getenv(out.Type().Field(i).Tag.Get("envv")))
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
	}
}

func setValue(out reflect.Value, envv string) error {
	if len(envv) == 0 {
		return nil
	}
	switch out.Kind() {
	case reflect.String:
		out.SetString(envv)
	case reflect.Int:
		num, err := strconv.Atoi(envv)
		if err != nil {
			return err
		}
		i := int64(num)
		if !out.OverflowInt(i) {
			out.SetInt(i)
		}
	case reflect.Bool:
		b, err := strconv.ParseBool(envv)
		if err != nil {
			return err
		}
		out.SetBool(b)
	case reflect.Float32:
		f, err := strconv.ParseFloat(envv, 32)
		if err != nil {
			return err
		}
		if !out.OverflowFloat(f) {
			out.SetFloat(f)
		}
	case reflect.Float64:
		f, err := strconv.ParseFloat(envv, 64)
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
			err := json.Unmarshal([]byte(envv), &ints)
			if err != nil {
				return err
			}
			out.Set(reflect.ValueOf(ints))
		case reflect.String:
			var strings []string
			err := json.Unmarshal([]byte(envv), &strings)
			if err != nil {
				return err
			}
			out.Set(reflect.ValueOf(strings))
		case reflect.Bool:
			var bools []bool
			err := json.Unmarshal([]byte(envv), &bools)
			if err != nil {
				return err
			}
			out.Set(reflect.ValueOf(bools))
		case reflect.Float32:
			var floats []float32
			err := json.Unmarshal([]byte(envv), &floats)
			if err != nil {
				return err
			}
			out.Set(reflect.ValueOf(floats))
		case reflect.Float64:
			var floats []float64
			err := json.Unmarshal([]byte(envv), &floats)
			if err != nil {
				return err
			}
			out.Set(reflect.ValueOf(floats))
		default:
			return (errors.New("undefine type"))
		}

	default:
		return (errors.New("undefine type"))
	}

	return nil
}

func readFile(file_path string) ([]byte, error) {
	rootDirPath, err := filepath.Abs(filepath.Dir(file_path))
	if err != nil {
		log.Fatalf("file error: %v", err)
	}

	configPath := filepath.Join(rootDirPath, yamlName)
	file, err := ioutil.ReadFile(configPath)

	return file, err
}
