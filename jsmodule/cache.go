package jsmodule

import (
	"github.com/skyfox2000/nect-utils/cache"
)

var cacheModule = map[string]interface{}{
	"get": func(key string) interface{} {
		result, _ := cache.Cache.Get(key)
		return result
	},
	"mget": func(keys interface{}) map[string]interface{} {
		results := cache.Cache.MGet(keys)
		return results
	},
	"set": func(key string, data interface{}, exp *int) {
		cache.Cache.Set(key, data, exp)
	},
	"mset": func(data map[string]interface{}, exp *int) {
		cache.Cache.MSet(data, exp)
	},
	"delete": func(key string, data interface{}, exp *int) {
		cache.Cache.Delete(key)
	},
	"mdelete": func(keys []string) {
		cache.Cache.MDelete(keys)
	},
	"deleteKeys": func(filter *string) {
		cache.Cache.DeleteKeys(filter)
	},
	"keys": func(filter *string) []string {
		return cache.Cache.Keys(filter)
	},
}

func init() {
	JSModules["Cache"] = cacheModule
}
