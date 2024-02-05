package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
)

func MD5(xpath string) string {
	hash := md5.Sum([]byte(xpath))
	result := hex.EncodeToString(hash[:])
	return result
}

func UUID() string {
	data := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz@#$^&*()_+-="
	result := make([]byte, 16)
	bs := []byte(data)
	for i := 0; i < 16; i++ {
		result[i] = data[rand.Intn(len(bs))]
	}
	str := string(result)

	return str
}
