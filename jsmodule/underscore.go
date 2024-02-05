package jsmodule

import (
	"reflect"
	"strings"

	"github.com/skyfox2000/nect-utils/underscore"
)

func registerUnderscore() map[string]interface{} {
	var Underscore = underscore.Underscore
	underscore := make(map[string]interface{})
	underscoreType := reflect.TypeOf(Underscore)

	for i := 0; i < underscoreType.NumMethod(); i++ {
		method := underscoreType.Method(i)
		name := method.Name
		methodName := strings.ToLower(name[0:1]) + name[1:]
		underscore[methodName] = reflect.ValueOf(Underscore).MethodByName(name).Interface()
	}

	return underscore
}

func init() {
	JSModules["underscore"] = registerUnderscore()
}
