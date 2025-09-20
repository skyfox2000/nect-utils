package jsmodule

import (
	"github.com/skyfox2000/nect-utils/encrypt"
)

var crytoModule = map[string]interface{}{
	"md5": func(data string) string {
		return encrypt.MD5(data)
	},
}

func init() {
	JSModules["crypto"] = crytoModule
}
