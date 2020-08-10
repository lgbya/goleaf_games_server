package tool

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

//md5加密
func Md5(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//获取随机数
func RandNum(num int) int{
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
}