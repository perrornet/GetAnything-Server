package utils

import (
	"crypto/md5"
	"encoding/hex"
)

type _md5 struct {
	text []byte
}

func NewMd5(text interface{}) _md5 {
	if v, ok := text.(string); ok {
		return _md5{text: []byte(v)}
	} else {
		return _md5{text: text.([]byte)}
	}
}

func (m _md5) Encrypt() string {
	d := md5.New()
	d.Write(m.text)
	return hex.EncodeToString(d.Sum(nil))
}
