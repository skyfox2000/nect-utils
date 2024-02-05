package file

import "os"

// File 对应的结构体
var File = &fileStruct{}

type fileStruct struct{}

func (p *fileStruct) Exists(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	return err == nil && fileInfo != nil && !fileInfo.IsDir()
}
