package zip

import (
	"bytes"
	"compress/gzip"

	"github.com/skyfox2000/nect-utils/json"
)

// GZIP 对应的结构体
var GZIP = &gzipStruct{}

type gzipStruct struct{}

func (p *gzipStruct) Zip(data interface{}) (interface{}, error) {
	dataStr := json.JSON.Stringify(data)
	var buf bytes.Buffer

	// 创建 GzipWriter 进行数据压缩
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write([]byte(dataStr))
	if err != nil {
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *gzipStruct) Unzip(compressedData []byte) (string, error) {
	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
