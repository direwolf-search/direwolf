package helpers

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

func RandomInt(num int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	random := r.Intn(num)

	return random
}

func GetMd5(fields ...string) string {
	buf := new(bytes.Buffer)
	buf.Grow(256)
	for _, str := range fields {
		buf.Write([]byte(str))
	}
	return fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
}
