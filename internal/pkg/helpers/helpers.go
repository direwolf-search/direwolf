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

// ErrorBuilder builds error from error message and additional fields.
// Error message must be a first element in argument list.
func ErrorBuilder(fields ...interface{}) error {
	var formatString = ""
	for i := 0; i < len(fields); i++ {
		formatString += " %v"
	}
	return fmt.Errorf(formatString, fields...)
}
