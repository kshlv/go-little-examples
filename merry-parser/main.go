package main

import (
	"fmt"
	"reflect"
	"strings"
)

// IndexConfig ...
type IndexConfig struct {
	StartPos int
	EndPos   int
	BitMask  byte
}

// mask
//
// 0 0 0 0 0 0 0 0 : byte
// 7 6 5 4 3 2 1 0
// n s i f o l v k : code
//
// 0 - key
// 1 - value
// 2 - list
// 3 - object
// 4 - float64
// 5 - intfmt.Println(s)
// 6 - string
// 7 - null
type mask byte

const (
	key mask = 1 << iota
	value
	list
	object
	float
	integer
	str
	null
)

// Index produces an additional index for data with which a client can
// decode JSON string easier
//
// Index consists of triplets separated by semicolon
// like so:
// 0,24,1001
func Index(data interface{}) ([]byte, string) {
	json := []byte{}
	configs := []IndexConfig{}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Struct:
		json = append(json, '{')

		dataValue := reflect.ValueOf(data)
		configs = make([]IndexConfig, 2*dataValue.NumField()+1)
		configs[0] = IndexConfig{
			StartPos: 0,
			EndPos:   -1,
			BitMask:  byte(object),
		}

		for i := 0; i < dataValue.NumField(); i++ {
			keyName := strings.ToLower(dataValue.Type().Field(i).Name)
			keyConfig := IndexConfig{
				StartPos: len(json) + 1,
				EndPos:   len(json) + len(keyName) + 1,
				BitMask:  byte(str | key),
			}
			configs[2*i+1] = keyConfig

			json = append(json, []byte(fmt.Sprintf("%q:", keyName))...)

			// debugString := fmt.Sprintf("%d: %s %s = %v\n", i, typeOfV.Field(i).Name, f.Type(), f.Interface())
			// fmt.Println(debugString)

			f := dataValue.Field(i)
			valueConfig := IndexConfig{}

			switch f.Type().Kind() {
			case reflect.String:
				q := f.String()
				valueConfig = IndexConfig{
					StartPos: len(json) + 1,
					EndPos:   len(json) + len(q) + 1,
					BitMask:  byte(str | value),
				}
				json = append(json, []byte(fmt.Sprintf("%q", q))...)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				d := fmt.Sprintf("%d", f.Int())
				valueConfig = IndexConfig{
					StartPos: len(json),
					EndPos:   len(json) + len(d),
					BitMask:  byte(integer | value),
				}
				json = append(json, []byte(d)...)
			default:
				fmt.Println("no such option")
			}

			configs[2*i+2] = valueConfig

			if i != dataValue.NumField()-1 {
				json = append(json, ',')
			}
		}

		json = append(json, '}')
		configs[0].EndPos = len(json)

		ss := []string{}
		for _, ic := range configs {
			ss = append(ss, fmt.Sprintf("%d,%d,%b", ic.StartPos, ic.EndPos, ic.BitMask))
		}

		return json, strings.Join(ss, ";")
	case reflect.Array:
	case reflect.Slice:
		fmt.Println("WIP")
	default:
		fmt.Println("error")
	}

	return nil, ""
}

// Parse ...
func Parse() {

}

func main() {
	fmt.Println("merry parser")
}
