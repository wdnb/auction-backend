package utils

import (
	"encoding/json"
	"errors"
	"fmt"
)

func ToJSONBytes(data interface{}) ([]byte, error) {
	// 将数据编码为JSON格式的字节数组
	bytes, err := json.Marshal(data)
	if err != nil {
		errMsg := fmt.Sprintf("failed to convert data to bytes: %v", err)
		return nil, errors.New(errMsg)
	}

	return bytes, nil
}

func FromJSONBytes(bytes []byte, data interface{}) error {
	err := json.Unmarshal(bytes, data)
	if err != nil {
		errMsg := fmt.Sprintf("failed to convert JSON bytes to data: %v", err)
		return errors.New(errMsg)
	}
	return nil
}
