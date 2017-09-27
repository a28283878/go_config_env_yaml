# configy
It is a go library for mapping data from yaml or enviroment variable.

## Support type
* string
* int
* bool
* float32
* float64
* slice(with all types above)

## Example
```go
package main

import (
	"fmt"

	"gitlab.paradise-soft.com.tw/rd/configy"
)

type Configuration struct {
	String string `yaml:"string" env:"STRING"`
	Struct struct {
		Int  int  `yaml:"int" env:"INT"`
		Bool bool `yaml:"bool" env:"BOOL"`
	}
	Float32    float32   `yaml:"float32" env:"FLOAT32"`
	Float64    float64   `yaml:"float64" env:"FLOAT64"`
	SliceStr   []string  `yaml:"slice_str" env:"SLICE_STR"`
	SliceInt   []int     `yaml:"slice_int" env:"SLICE_INT"`
	SliceBool  []bool    `yaml:"slice_bool" env:"SLICE_BOOL"`
	SliceFloat []float64 `yaml:"slice_float" env:"SLICE_FLOAT"`
}

func main() {
	_config := Configuration{}
	configy.Load(&_config, "./config.yml")
	fmt.Println(_config)

	_config = Configuration{}
	//use os enviroment variables
	configy.Load(&_config, "")
	fmt.Println(_config)
}
```
