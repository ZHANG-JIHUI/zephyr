package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}
