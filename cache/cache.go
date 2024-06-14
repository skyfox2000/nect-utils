package cache

import (
	"regexp"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/skyfox2000/nect-utils/encrypt"
	"github.com/skyfox2000/nect-utils/json"
	"github.com/skyfox2000/nect-utils/underscore"
)

var cacheStorage = cache.New(5*time.Minute, 10*time.Minute)

// Cache 对应的结构体
var Cache = &cacheStruct{}
var DefaultExpiration = 600

type cacheStruct struct{}

// 获取缓存数据
func (p *cacheStruct) Get(key string) (interface{}, bool) {
	key = strings.TrimSpace(key)
	result, ok := cacheStorage.Get(key)
	if ok {
		result = json.JSON.Clone(result)
		return result, true
	} else {
		return nil, false
	}
}

// 批量获取缓存数据，支持*号和数组形式
func (p *cacheStruct) MGet(keys interface{}) map[string]interface{} {
	results := make(map[string]interface{})
	switch t := keys.(type) {
	case []string:
		for _, k := range t {
			result, ok := p.Get(k)
			if ok {
				results[k] = result
			}
		}
		return results
	case string:
		cacheKeys := p.Keys(&t)
		for _, k := range cacheKeys {
			result, ok := p.Get(k)
			if ok {
				results[k] = result
			}
		}
	}
	return results
}

// / 根据*号获取符合的Keys
func (p *cacheStruct) Keys(filter *string) []string {
	keys := underscore.Underscore.Keys(cacheStorage.Items())
	if filter == nil || *filter == "" {
		return keys
	}
	pattern := regexp.QuoteMeta(*filter)
	pattern = strings.ReplaceAll(pattern, "\\*", ".*")
	reg := regexp.MustCompile(pattern)
	subKeys := make([]string, 0)
	for _, key := range keys {
		if reg.MatchString(key) {
			subKeys = append(subKeys, key)
		}
	}
	return subKeys
}

func (p *cacheStruct) Set(key string, data interface{}, exp *int) {
	var d time.Duration
	if exp != nil {
		d = time.Duration(*exp) * time.Second
	}
	key = strings.TrimSpace(key)
	cacheStorage.Set(key, data, d)
}

func (p *cacheStruct) MSet(data map[string]interface{}, exp *int) {
	for key, value := range data {
		p.Set(key, value, exp)
	}
}

func (p *cacheStruct) Delete(key string) {
	cacheStorage.Delete(key)
}

func (p *cacheStruct) MDelete(keys []string) {
	for _, v := range keys {
		p.Delete(v)
	}
}

func (p *cacheStruct) DeleteKeys(filter *string) {
	keys := p.Keys(filter)
	p.MDelete(keys)
}

// 脚本, 数据
func (p *cacheStruct) HashKey(jscode string, data map[string]interface{}) string {
	jscode = strings.TrimSpace(jscode)
	keys := underscore.Underscore.Keys(data, true)
	sb := strings.Builder{}
	if len(jscode) > 100 {
		jscode = encrypt.MD5(jscode)
	}

	for _, key := range keys {
		sb.WriteString(data[key].(string))
	}

	dataStr := sb.String()
	if len(dataStr) > 100 {
		dataStr = encrypt.MD5(dataStr)
	}
	result := encrypt.MD5(jscode + strings.Join(keys, "") + dataStr)
	return result
}
