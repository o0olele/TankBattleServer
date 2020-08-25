package common

import (
	"crypto/md5"
	"io"
)

const salt string = "cbh8a932gnvf"

func GetToken(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return string(h.Sum([]byte(salt)))
}

func CheckToken(s, token string) bool {
	tok := GetToken(s)
	if tok == token {
		return true
	}
	return false
}
