package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
)

func IsBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String, reflect.Slice:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

func StringToUint32(str string) (uint32, error) {
	var num uint32
	num64, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}
	num = uint32(num64)
	return num, nil
}

func SplitSuffix(str, substr string) (string, string) {
	if !strings.Contains(str, substr) {
		return "", ""
	}
	vals := strings.Split(str, substr)
	return vals[0], vals[1]
}

func RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	suffix := make([]rune, n)
	for i := range suffix {
		suffix[i] = letters[rand.Intn(len(letters))]
	}
	return string(suffix)
}

func SHA1Hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}
