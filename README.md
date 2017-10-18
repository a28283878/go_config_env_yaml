# configy
It is a go library for mapping data from yaml or enviroment variable.

## Support type
* string
* int
* bool
* float32
* float64
* struct
* slice(with all types above)

## Example
```go
package main

import (
	"fmt"

	"gitlab.paradise-soft.com.tw/rd/configy"
)

type Struct_Slice struct {
	Int  int  `yaml:"struct_int"`
	Bool bool `yaml:"struct_bool"`
}

type Configuration struct {
	StringDefault string `yaml:"string" env:"default:123"` //set default value
	Struct        struct {
		Sturct_Slice []Struct_Slice `env:"SLICE_STRUCT"`
		Int          int            `yaml:"int" env:"INT"`
		Bool         bool           `yaml:"bool" env:"BOOL"`
	}
	Float32    float32   `yaml:"float32" env:"FLOAT32"`
	Float64    float64   `yaml:"float64" env:"FLOAT64"`
	Slice_Str  []string  `yaml:"slice_str"` //use name to get env variable
	SliceInt   []int     `yaml:"slice_int" env:"SLICE_INT"`
	SliceBool  []bool    `yaml:"slice_bool" env:"SLICE_BOOL"`
	SliceFloat []float64 `yaml:"slice_float" env:"SLICE_FLOAT"`
}

func main() {
	_config := Configuration{}
	configy.Load(&_config, "./config.yml")
	fmt.Println(_config)

	_config = Configuration{}
	//if can't find target file or value is ""
	//will try to get value from env
	configy.Load(&_config, "")
	fmt.Println(_config)
}
```

## Example Yaml
* * *
```ymal
string: "Hi"
struct:
    int : 20
    bool : on
    sturct_slice :
      - struct_int: 50
        struct_bool: on
      - struct_int: 100
        struct_bool: false
float32: 3.1415926
float64: 1.23456789
slice_str: ["first","second"]
slice_int:  [1,2,4]
slice_bool: [true,false,true]
slice_float:    [1.2,2.3,45.2]
```

## Example ENV
* * *
```
STRING=Hi
SLICE_STRUCT={123,true},{255,false}
INT=50
BOOL=true
FLOAT32=21.5
FLOAT64=5.23215486
SLICE_STR=["Hi Again","AAAAA"]
SLICE_INT=[123,456]
SLICE_BOOL=[true,false,true]
SLICE_FLOAT=[1.2,2.3,45.2]
```
