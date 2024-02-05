package jsmodule

import (
	"github.com/skyfox2000/nect-utils/zip"
)

var zipModule = map[string]interface{}{
	"zip": func(obj interface{}) string {
		result, _ := zip.GZIP.Zip(obj)
		return result.(string)
	},
	"unzip": func(compressedData interface{}) interface{} {
		result, _ := zip.GZIP.Unzip(compressedData.([]byte))
		return result
	},
}

func init() {
	JSModules["ZIP"] = zipModule
}
