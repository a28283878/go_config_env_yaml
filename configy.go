package configy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

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
		log.Printf("can't load yaml, load env instead. %s", err.Error())
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

		case reflect.Slice:
			if out.Field(i).Type().Elem().Kind() == reflect.Struct {
				err := readSlice(out.Field(i), out.Type().Field(i))
				if err != nil {
					log.Fatalf(err.Error())
				}
			} else {
				err := setValue(out.Field(i), getEnvValue(out.Type().Field(i)))
				if err != nil {
					log.Fatalf(err.Error())
				}
			}
		case reflect.Struct:
			loadEnv(out.Field(i))
		default:
			err := setValue(out.Field(i), getEnvValue(out.Type().Field(i)))
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
	}
}

func readSlice(out reflect.Value, outField reflect.StructField) error {
	env := os.Getenv(outField.Tag.Get("env"))
	tokens := getTokens(env)
	elType := out.Type().Elem()
	slice := reflect.MakeSlice(out.Type(), out.Len(), out.Cap())

	if len(tokens) >= 1 {
		for _, ele := range tokens {
			el := reflect.New(elType).Elem()
			ele = strings.TrimLeft(ele, "{")
			ele = strings.TrimRight(ele, "}")
			vals := strings.Split(ele, ",")
			for index, val := range vals {
				if index >= el.NumField() {
					return fmt.Errorf("too many arguments. in {%s}", ele)
				}
				if err := setValue(el.Field(index), val); err != nil {
					return err
				}
			}
			slice = reflect.Append(slice, el)
		}
	}

	out.Set(slice)
	return nil
}

func getEnvValue(field reflect.StructField) string {
	var tag string
	var env string

	tag = field.Tag.Get("env")
	if len(tag) == 0 {
		tag = field.Name
	}

	method := strings.Split(tag, ":")

	if len(method) > 1 {
		if strings.TrimSpace(method[0]) == "default" {
			return method[1]
		}
	}

	env = os.Getenv(tag)
	return env
}

func getTokens(env string) []string {
	var result []string
	inBrace := false
	inString := false
	var token string

	for _, char := range env {
		if inBrace {
			if string(char) == "\"" {
				inString = !inString
			}
			token = token + string(char)
			if string(char) == "}" && !inString {
				inBrace = false
			}
		} else {
			if string(char) == "{" && !inString {
				inBrace = true
				if len(token) > 0 {
					result = append(result, token)
				}
				token = "{"
			}
		}
	}

	if len(token) > 0 {
		result = append(result, token)
	}

	return result
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
			return fmt.Errorf("undefined type")
		}
	case reflect.Struct:

	default:
		return fmt.Errorf("undefined type")
	}

	return nil
}

func readFile(file_path string) ([]byte, error) {
	file, err := ioutil.ReadFile(file_path)
	if err != nil {
		return file, err
	}
	return file, nil
}
