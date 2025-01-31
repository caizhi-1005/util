package common

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
	"sync"
	"time"
)

// IsEmpty 判读数据是否为空
func IsEmpty(a interface{}) bool {
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

var Locker = make(map[string]*sync.RWMutex)

// Lock 锁
func Lock(index string) {
	for {
		_, ok := Locker[index]
		if !ok {
			Locker[index] = &sync.RWMutex{}
			break
		}
		//100ms轮训一次状态
		time.Sleep(100 * time.Millisecond)
	}

	Locker[index].Lock()
}

// Unlock 解锁
func Unlock(index string) {
	Locker[index].Unlock()
	//删除使用过的锁，避免map无限增加
	delete(Locker, index)
}

// getMd5String 生成32位md5字串
func getMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateUniqueId 生成Guid字串
func GenerateUniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return getMd5String(base64.URLEncoding.EncodeToString(b))
}

// StructToMap
func StructToMap(obj interface{}) map[string]interface{} {
	ty := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < ty.NumField(); i++ {
		data[ty.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

// IsUrlErr 判断错误是否为路径错误
func IsUrlErr(err error) bool {
	errStr := fmt.Sprintf("%T", err)
	if errStr == "*url.Error" {
		return true
	}
	return false
}

// IsUpper 判断是否是大写字母
func IsUpper(b byte) bool {
	return 'A' <= b && b <= 'Z'
}

// IsLower 判断是否是小写字母
func IsLower(b byte) bool {
	return 'a' <= b && b <= 'z'
}

// IsDigit 判断是否是数字
func IsDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

// ToLower 转换为小写字母
func ToLower(b byte) byte {
	if IsUpper(b) {
		return b - 'A' + 'a'
	}
	return b
}

// DefaultIfNil checks if the value is nil, if true returns the default value otherwise the original
func DefaultIfNil(value interface{}, defaultValue interface{}) interface{} {
	if value != nil {
		return value
	}
	return defaultValue
}

// FirstNonNil returns the first non nil parameter
func FirstNonNil(values ...interface{}) interface{} {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}
