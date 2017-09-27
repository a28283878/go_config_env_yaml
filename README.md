# configy
It is a go library for mapping data from yaml or enviroment variable.

## Support type
* string
* int
* bool
* float32
* float64
* slice(with all type above)

## Example
```go
package main

import (
	"fmt"

	"gitlab.paradise-soft.com.tw/rd/configy"
)

type Configuration struct {
	String string `yaml:"string" envv:"STRING"`
	Struct struct {
		Int  int  `yaml:"int" envv:"INT"`
		Bool bool `yaml:"bool" envv:"BOOL"`
	}
	Float32    float32   `yaml:"float32" envv:"FLOAT32"`
	Float64    float64   `yaml:"float64" envv:"FLOAT64"`
	SliceStr   []string  `yaml:"slice_str" envv:"SLICE_STR"`
	SliceInt   []int     `yaml:"slice_int" envv:"SLICE_INT"`
	SliceBool  []bool    `yaml:"slice_bool" envv:"SLICE_BOOL"`
	SliceFloat []float64 `yaml:"slice_float" envv:"SLICE_FLOAT"`
}

func main() {
	_config := Configuration{}
	configy.Load(&_config, "")
	fmt.Println(_config)
}
```
