package underscore

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// Underscore 对应的结构体
var Underscore = &underscore{}

type underscore struct{}

// IsEmpty 判断值是否为空
func (p *underscore) IsEmpty(val interface{}) bool {
	if val == nil {
		return true
	}
	switch reflect.TypeOf(val).Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return reflect.ValueOf(val).Len() == 0
	case reflect.Bool:
		return !reflect.ValueOf(val).Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(val).Int() == 0
	}
	return false
}

// IsNaN 判断一个数字是否为NaN
func (p *underscore) IsNaN(value interface{}) bool {
	if value == nil || p.IsString(value) {
		return false
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Float64:
		return math.IsNaN(value.(float64))
	}
	return false
}

// IsNumber 判断一个值是否为数字
func (p *underscore) IsNumber(value interface{}) bool {
	if value == nil {
		return false
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func (p *underscore) IsString(value interface{}) bool {
	if value == nil {
		return false
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return true
	default:
		return false
	}
}

func (p *underscore) IsObject(value interface{}) bool {
	if value == nil {
		return false
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Map, reflect.Struct, reflect.Slice:
		return true
	default:
		return false
	}
}

func (p *underscore) IsMap(value interface{}) bool {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Map:
		return true
	}
	return false
}

func (p *underscore) IsArray(value interface{}) bool {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Array:
		return true
	}
	return false
}

// Uniq 返回去重后的数组
func (p *underscore) Uniq(arr []interface{}) []interface{} {
	set := make(map[interface{}]struct{})
	results := []interface{}{}

	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	for _, item := range arr {
		if _, exists := set[item]; !exists {
			results = append(results, item)
			set[item] = struct{}{}
		}
	}

	return results
}

func (p *underscore) Len(value interface{}) int {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Array:
		return reflect.ValueOf(value).Len()
	case reflect.Map:
		return reflect.ValueOf(value).Len()
	case reflect.String:
		return len(value.(string))
	default:
		return 0
	}
}

func (p *underscore) Index(value interface{}, target interface{}) int {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(value)
		for i := 0; i < s.Len(); i++ {
			if s.Index(i).Interface() == target {
				return i
			}
		}
	case reflect.String:
		return strings.Index(value.(string), target.(string))
	}
	return -1
}

// Contains 数组中包含
func (p *underscore) Contains(arr interface{}, target interface{}) bool {
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(arr)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(s.Index(i).Interface(), target) {
				return true
			}
		}
	}
	return false
}

// SortBy 根据字段对数组进行排序
func (p *underscore) SortBy(arr []interface{}, field, order string) []interface{} {
	// 1. 判断arr的具体类型
	if len(arr) == 0 {
		// 数组为空，直接返回
		return arr
	}

	// 2. 判断元素类型并建立map
	elem := reflect.ValueOf(arr[0])
	isMap := elem.Kind() == reflect.Map

	m := map[interface{}][]interface{}{}
	for _, item := range arr {
		v := reflect.ValueOf(item)
		var key reflect.Value
		if isMap {
			key = v.MapIndex(reflect.ValueOf(field))
		} else {
			key = v.FieldByName(field)
		}
		m[key.Interface()] = append(m[key.Interface()], item)
	}

	// 3. 对field的值进行排序
	var keys []interface{}
	for k := range m {
		if v, ok := reflect.ValueOf(k).Interface().(string); ok {
			keys = append(keys, v)
		} else if v, ok := reflect.ValueOf(k).Interface().(int); ok {
			keys = append(keys, v)
		}
	}

	keys = p.Sort(keys, order).([]interface{})

	// 4. 重新建立数组
	result := make([]interface{}, 0, len(arr))
	for _, key := range keys {
		result = append(result, m[key]...)
	}

	return result
}

func (p *underscore) Sort(arr []interface{}, order string) interface{} {
	if len(arr) == 0 {
		// 数组为空，直接返回
		return arr
	}

	switch arr[0].(type) {
	case string:
		switch order {
		case "asc":
			sort.SliceStable(arr, func(i, j int) bool {
				return fmt.Sprintf("%v", arr[i]) < fmt.Sprintf("%v", arr[j])
			})
		case "desc":
			sort.SliceStable(arr, func(i, j int) bool {
				return fmt.Sprintf("%v", arr[i]) > fmt.Sprintf("%v", arr[j])
			})
		}
	case int:
		switch order {
		case "asc":
			sort.SliceStable(arr, func(i, j int) bool {
				return arr[i].(int) < arr[j].(int)
			})
		case "desc":
			sort.SliceStable(arr, func(i, j int) bool {
				return arr[i].(int) > arr[j].(int)
			})
		}
	case float64:
		switch order {
		case "asc":
			sort.SliceStable(arr, func(i, j int) bool {
				return arr[i].(float64) > arr[j].(float64)
			})
		case "desc":
			sort.SliceStable(arr, func(i, j int) bool {
				return arr[i].(float64) < arr[j].(float64)
			})
		}
	}

	return arr
}

// Keys 获取map的所有键
func (p *underscore) Keys(val interface{}, order ...bool) []string {
	keys := make([]string, 0)
	if val == nil {
		return keys
	}

	switch reflect.TypeOf(val).Kind() {
	case reflect.Map:
		mapKeys := reflect.ValueOf(val).MapKeys()
		for _, key := range mapKeys {
			keys = append(keys, key.String())
		}
	}
	if len(order) > 0 && order[0] {
		sort.SliceStable(keys, func(i, j int) bool {
			return fmt.Sprintf("%v", keys[i]) < fmt.Sprintf("%v", keys[j])
		})
	}
	return keys
}

// Values 获取map的所有值
func (p *underscore) Values(val map[string]interface{}) []interface{} {
	values := make([]interface{}, 0, len(val))

	for _, value := range val {
		values = append(values, value)
	}
	return values
}

// HasKey 判断Map是否包含Key
func (p *underscore) HasKey(val interface{}, key string) bool {
	if val == nil {
		return false
	}

	switch reflect.TypeOf(val).Kind() {
	case reflect.Map:
		mapKeys := reflect.ValueOf(val).MapKeys()
		for _, k := range mapKeys {
			if k.String() == key {
				return true
			}
		}
	}
	return false
}

// Padding 将数字填充到指定长度
func (p *underscore) Padding(num int, length int) string {
	numStr := fmt.Sprintf("%d", num)
	paddingLength := length - len(numStr)

	if paddingLength <= 0 {
		return numStr
	}

	paddedStr := strings.Repeat("0", paddingLength) + numStr
	return paddedStr
}

// 返回一个UUID
func (p *underscore) Uuid() string {
	guid, _ := uuid.NewUUID()
	return strings.ReplaceAll(guid.String(), "-", "")
}

func (p *underscore) CallMethod(methodName string, args ...interface{}) (interface{}, error) {
	underscoreType := reflect.TypeOf(Underscore)

	for i := 0; i < underscoreType.NumMethod(); i++ {
		method := underscoreType.Method(i)
		if method.Name == methodName {
			result := method.Func.Call([]reflect.Value{reflect.ValueOf(Underscore), reflect.ValueOf(args)})
			if len(result) > 0 {
				switch result[0].Kind() {
				case reflect.Int32:
					return result[0].Int(), nil
				case reflect.Bool:
					return result[0].Bool(), nil
				}
			}
			return nil, nil
		}
	}
	return nil, nil
}
