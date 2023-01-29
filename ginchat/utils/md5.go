package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	b := h.Sum(nil)
	return hex.EncodeToString(b)
}

// md5加密后转大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// 随机数加密
func MakePassword(plainpwd, slt string) string {
	return Md5Encode(plainpwd + slt)
}

func ValidatePwd(plainpwd, slt, password string) bool {
	return Md5Encode(plainpwd+slt) == password
}
