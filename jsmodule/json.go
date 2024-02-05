package jsmodule

import (
	"github.com/skyfox2000/nect-utils/json"
)

var jsonModule = map[string]interface{}{
	"stringify": func(obj interface{}) string {
		JSON := json.JSON
		return JSON.Stringify(obj)
	},
	"getXPathValue": func(obj interface{}, xpath string) interface{} {
		JSON := json.JSON
		result, _ := JSON.GetXPathValue(obj, xpath)
		return result
	},
	"parse": func(jsonStr interface{}) interface{} {
		JSON := json.JSON
		result, ok := JSON.Parse(jsonStr)
		if !ok {
			Logger.Error("failed to parse json: " + jsonStr.(string))
		}
		return result
	},
	"clone": func(obj interface{}) interface{} {
		JSON := json.JSON
		return JSON.Clone(obj)
	},
}

func init() {
	JSModules["JSON"] = jsonModule
}
