package json

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/google/uuid"
	"github.com/skyfox2000/nect-utils/underscore"
)

// JSON 对应的结构体
var JSON = &jsonStruct{}

type jsonStruct struct{}

// 重新对Params根据插件数据结构赋值
func convertValue(value interface{}) interface{} {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	newValue := value
	return convertValueWithLock(newValue)
}

func convertValueWithLock(value interface{}) interface{} {
	var result interface{}
	newValue := value
	switch v := newValue.(type) {
	case map[string]interface{}:
		result = convertMap(v)
	case map[interface{}]interface{}:
		result = convertMap(v)
	case []interface{}:
		result = convertSlice(v)
	default:
		result = v
	}
	return result
}

// ConvertSlice recursively converts values in the slice
func convertSlice(slice []interface{}) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = convertValueWithLock(v)
	}
	return result
}

// ConvertMap recursively converts map[interface{}]interface{} to map[string]interface{}
func convertMap(inputMap interface{}) map[string]interface{} {
	results := make(map[string]interface{})
	newValue := inputMap
	switch m := newValue.(type) {
	case map[interface{}]interface{}:
		for key, value := range m {
			keyString, _ := key.(string)
			result := fillMap(keyString, value)
			results[keyString] = result
		}
	case map[string]interface{}:
		for key, value := range m {
			result := fillMap(key, value)
			results[key] = result
		}
	}

	return results
}

func fillMap(key string, value interface{}) interface{} {
	var result interface{}
	newValue := value
	switch v := newValue.(type) {
	case map[string]interface{}:
		result = convertValueWithLock(v)
	case map[interface{}]interface{}:
		result = convertValueWithLock(v)
	case []interface{}:
		result = convertSlice(v)
	default:
		result = v
	}
	return result
}

// ParseParams函数根据插件数据结构解析并赋值给plugParams参数，并赋值默认值
func (p *jsonStruct) ParseParams(
	plugParams any,
	actionParams, defaultParams map[string]interface{}) interface{} {
	if actionParams == nil {
		return nil
	}
	// Convert actionParams to map[string]interface{}
	convertedParams := convertValue(actionParams).(map[string]interface{})

	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	// 使用Marshal和Unmarshal处理类型转换
	b, err := json.Marshal(convertedParams)
	if err != nil {
		return err
	}

	// 将JSON数据解析到plugParams中
	json.Unmarshal(b, plugParams)

	if defaultParams != nil {
		// 使用反射循环遍历plugParams的每个参数，如果为空则使用defaultParams进行填充
		dataType := reflect.ValueOf(plugParams).Elem()

		for i := 0; i < dataType.NumField(); i++ {
			fieldName := dataType.Type().Field(i).Name
			handleField(dataType.Field(i), fieldName, defaultParams)
		}
	}

	return plugParams
}

// 判断字段类型并进行处理
func handleField(field reflect.Value, fieldName string, defaultParams map[string]interface{}) {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	switch field.Interface().(type) {
	case string, map[string]interface{}, []interface{}:
		// Map字段
		if field.Len() == 0 {
			// 空数据处理逻辑
			if value, ok := defaultParams[fieldName]; ok {
				field.Set(reflect.ValueOf(value))
			}
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		// 数字字段
		if field.Interface() == 0 {
			// 0 值处理逻辑
			if value, ok := defaultParams[fieldName]; ok {
				field.Set(reflect.ValueOf(value))
			}
		}
	default:
		// 其他类型字段的处理逻辑
	}
}

var XReplaceMutex = &sync.RWMutex{}

// 根据data对象替换dataStr内的路径
// mustHas 必须有值，否则报错
func (p *jsonStruct) ReplaceXPathValue(
	data interface{},
	dataStr string,
	mustHas bool) (string, bool) {
	XReplaceMutex.RLock()
	defer XReplaceMutex.RUnlock()
	dataStr = strings.TrimSpace(dataStr)
	resultStr := dataStr
	// 替换Key包含${}结构
	pattern := `\$([A-Z$]{1}[A-Za-z0-9_]*)(\[\"){0,2}([\p{Han}A-Za-z0-9._-]*)(\"\]){0,2}([.A-Za-z0-9_^]*)`
	reg := regexp.MustCompile(pattern)
	matches := reg.FindAllString(resultStr, -1)

	for _, match := range matches {
		matchStr := match[1:] // 去掉$符号，得到变量名
		// if p.HasXPath(matchStr) {
		// 	matchStr, _ = p.ReplaceXPathValue(data, matchStr, mustHas)
		// 	resultStr = strings.Replace(resultStr, match, "${"+matchStr+"}", -1)
		// }

		xpathValue, exists := p.GetXPathValue(data, matchStr)
		if exists {
			replacement := ""
			switch v := xpathValue.(type) {
			case map[string]interface{}, []interface{}:
				// 如果是对象或数组，将其转换为 JSON 字符串格式
				jsonValue, err := json.Marshal(v)
				if err != nil {
					// 处理错误，这里可以根据实际情况进行处理
					replacement = fmt.Sprint(xpathValue)
				} else {
					replacement = string(jsonValue)
				}
			case float64, int, int32, int64:
				// 如果是数字，转换为字符串
				replacement = fmt.Sprint(v)
			case nil:
				replacement = "null"
			default:
				// 如果是字符串，直接使用
				replacement = xpathValue.(string)
				// 处理其他类型，如 null 或 undefined，不进行替换
			}

			// 进行替换
			resultStr = strings.Replace(resultStr, "$"+matchStr, replacement, -1)
		} else if mustHas {
			// 必须有值，则报错
			return "", false
		} else {
			// 进行替换
			resultStr = strings.Replace(resultStr, "$"+matchStr, "", -1)
		}
	}

	return resultStr, true
}

var XGetMutex = &sync.RWMutex{}

// 根据XPath获取data的指定路径的数据
func (p *jsonStruct) GetXPathValue(data interface{}, xpath string) (interface{}, bool) {
	XGetMutex.RLock()
	defer XGetMutex.RUnlock()
	if data == nil {
		return nil, true
	}

	xpath = strings.TrimSpace(xpath)
	if strings.HasPrefix(xpath, "${") && strings.HasSuffix(xpath, "}") {
		xpath = xpath[2 : len(xpath)-1]
	}
	// 拆分xpath字符串，获取键和索引
	// 例如 "Git.Download[1].Path" 会被拆分为 ["Git", "Download[1]", "Path"]
	keys := strings.Split(xpath, ".")

	// 开始递归查询
	return getValueByKeys(data, keys)
}

var XSetMutex = &sync.RWMutex{}

func (p *jsonStruct) SetXPathValue(currentMap interface{}, xpath string, data interface{}) {
	XSetMutex.Lock()
	defer XSetMutex.Unlock()
	names := strings.Split(xpath, ".")
	for i, name := range names {
		if i == len(names)-1 {
			switch m := currentMap.(type) {
			case map[string]interface{}:
				m[name] = data
			case *sync.Map:
				m.Store(name, data)
			}
			break
		} else {
			switch m := currentMap.(type) {
			case map[string]interface{}:
				if m[name] == nil {
					m[name] = make(map[string]interface{})
				}
				currentMap = m[name]
			case *sync.Map:
				result, _ := m.Load(name)
				if result == nil {
					result = &sync.Map{}
					m.Store(name, result)
				}
				currentMap = result
			}
		}
	}
}

func getValueByKeys(data interface{}, keys []string) (interface{}, bool) {
	currentData := data
	for index, key := range keys {
		// 常规的key查询
		value, exists := getValueByKey(currentData, key)

		if !exists {
			// 如果当前层级不存在该键，则返回失败
			return nil, false
		}
		if index == len(keys)-1 {
			return value, true
		}

		// 如果值是一个对象，我们将其作为数据，继续查询下一层的键
		if subData, ok := value.(map[string]interface{}); ok {
			currentData = subData
		} else {
			// 到达最后一层，直接返回nil和 true
			return nil, false
		}
	}

	// 如果循环结束仍未返回结果，说明路径无效
	return nil, false
}

func getValueByKey(data interface{}, key string) (interface{}, bool) {
	// 常规的key查询
	var value interface{}
	var exists bool

	switch m := data.(type) {
	case map[string]interface{}:
		value, exists = m[key]
	case map[string]string:
		value, exists = m[key]
	case *sync.Map:
		value, exists = m.Load(key)
	}
	// 判断key的格式是否为 "Field[0]" 或 "Field["0"]"
	if strings.HasSuffix(key, "]") {
		// 将key分割为两部分
		parts := strings.Split(key, "[")
		firstPart := parts[0]
		indexPart := parts[1][:len(parts[1])-1] // 去掉开头的"["和结尾的"]"

		var keyValue interface{}
		// 获取对应的值
		switch m := data.(type) {
		case map[string]interface{}:
			keyValue, exists = m[firstPart]
		case map[string]string:
			keyValue, exists = m[firstPart]
		case *sync.Map:
			keyValue, exists = m.Load(firstPart)
		}
		if exists {
			// 如果值是一个数组，则根据索引获取对应的值
			switch v := keyValue.(type) {
			case []interface{}:
				index, err := strconv.Atoi(indexPart)
				if err == nil {
					if index >= 0 && index < len(v) {
						result := v[index]
						return result, true
					}
				}
			case map[string]interface{}:
				resultKey := strings.ReplaceAll(indexPart, "\"", "")
				result, ok := v[resultKey]
				return result, ok
			case sync.Map:
				resultKey := strings.ReplaceAll(indexPart, "\"", "")
				result, ok := v.Load(resultKey)
				return result, ok
			}
		}
		return nil, false
	}

	return value, exists
}

// 是否匹配，主键名，主键Key，子属性
func (p *jsonStruct) IsXPath(key string) (bool, string, string, string) {
	pattern := `\$([A-Z$]{1}[A-Za-z0-9_]*)[\[\"]{0,2}([\p{Han}\$A-Za-z0-9._-]*)[\"\]]{0,2}([.A-Za-z0-9_]*)`
	reg := regexp.MustCompile(pattern)

	// 查找第一个匹配项
	matches := reg.FindAllStringSubmatch(key, -1)
	if len(matches) > 0 {
		if strings.HasPrefix(matches[0][2], ".") {
			return matches[0][0] == key, matches[0][1], "", matches[0][2]
		}
		return matches[0][0] == key, matches[0][1], matches[0][2], matches[0][3]
	}
	return false, "", "", ""
}

func (p *jsonStruct) HasXPath(key string) bool {
	pattern := `\$([A-Z$]{1}[A-Za-z0-9_]*)[\[\"]{0,2}([\p{Han}\$A-Za-z0-9._-]*)[\"\]]{0,2}([.A-Za-z0-9_]*)`
	reg := regexp.MustCompile(pattern)

	// 查找第一个匹配项
	match1 := reg.FindString(key)
	if match1 != "" {
		return true
	}

	pattern = `\$[A-Z$]{0,1}[\$A-Za-z0-9.]*`
	reg = regexp.MustCompile(pattern)

	// 查找第一个匹配项
	match2 := reg.FindString(key)

	return match2 != ""
}

func (p *jsonStruct) Stringify(data interface{}) string {
	if data == nil {
		return ""
	}
	jsonData := data
	if !underscore.Underscore.IsObject(jsonData) {
		return data.(string)
	}

	convertedParams := convertValue(data)
	jsonStr, err := json.Marshal(convertedParams)
	if err != nil {
		// 处理转换失败的情况
		return err.Error()
	}
	return string(jsonStr)
}

func (p *jsonStruct) Parse(data interface{}) (interface{}, bool) {
	// 检查data是否为[]byte类型
	if dataBytes, ok := data.([]byte); ok {
		var result interface{}
		err := json.Unmarshal(dataBytes, &result)
		if err != nil {
			// 处理解析失败的情况
			return err.Error(), false
		}
		return result, true
	}

	// 检查data是否为string类型
	if dataStr, ok := data.(string); ok {
		var result interface{}
		err := json.Unmarshal([]byte(dataStr), &result)
		if err != nil {
			// 处理解析失败的情况
			return err.Error(), false
		}
		return result, true
	}

	// 如果data既不是string也不是[]byte类型，返回原始数据和false
	return data, false
}

func (p *jsonStruct) Clone(data interface{}) interface{} {
	if data == nil || !underscore.Underscore.IsObject(data) {
		return data
	}

	result, _ := p.Parse(p.Stringify(data))
	return result
}

func (p *jsonStruct) ToJSON(inputStr string, data map[string]interface{}) (interface{}, error) {
	// 复杂表达式，通过VM获取
	newVm := goja.New()

	for k, v := range data {
		newVm.Set("$"+k, v)
	}

	prepareCode := fmt.Sprintf("(function(){\nconst result=%s; return result;\n})();", inputStr)

	uuid := uuid.New().String()
	uuid = strings.Replace(uuid, "-", "", -1)
	prog, err := goja.Compile(uuid+".js", prepareCode, false)
	if err != nil {
		fmt.Println("compile error:", err.Error())
		return nil, err
	}
	exprResult, err := newVm.RunProgram(prog)

	if err != nil {
		return nil, err
	}

	if exprResult == nil {
		return exprResult, nil
	}

	return exprResult.Export(), nil
}

func (p *jsonStruct) Log(data interface{}, depth int, arrayLimit int) interface{} {
	if data == nil {
		return string("<nil>")
	}
	// convertedParams := convertValue(data)
	result := p.limitJSONDepth(data, depth, arrayLimit)

	bytes, _ := json.MarshalIndent(result, "", "  ")
	message := json.RawMessage(bytes)

	return string(message)
}

// limitJSONDepth limits the depth of the JSON output to the specified depth.
// Arrays will only print the first few elements as specified by arrayLimit.
func (p *jsonStruct) limitJSONDepth(data interface{}, depth int, arrayLimit int) interface{} {
	switch v := data.(type) {
	case []interface{}:
		limit := arrayLimit
		if len(v) < arrayLimit {
			limit = len(v)
		}
		slice := make([]interface{}, limit)
		for i := 0; i < limit; i++ {
			slice[i] = p.limitJSONDepth(v[i], depth, arrayLimit)
		}
		if len(v) > limit {
			slice = append(slice, "...")
		}
		return slice
	case []map[string]interface{}:
		limit := arrayLimit
		if len(v) < arrayLimit {
			limit = len(v)
		}
		slice := make([]interface{}, limit)
		for i := 0; i < limit; i++ {
			slice[i] = p.limitJSONDepth(v[i], depth, arrayLimit)
		}
		if len(v) > limit {
			slice = append(slice, "...")
		}
		return slice
	case map[string]interface{}:
		resultMap := make(map[string]interface{})
		for k, val := range v {
			if depth > -1 {
				resultMap[k] = p.limitJSONDepth(val, depth-1, arrayLimit)
			} else {
				switch val.(type) {
				case []interface{}:
				case map[string]interface{}:
					return "..."
				default:
					return data
				}
			}
		}
		return resultMap
	default:
		return v
	}
}
