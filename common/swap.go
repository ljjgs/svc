package common

import "encoding/json"

// 将 src 结构体对象中的字段值拷贝到 dest 结构体对象中对应的同名字段中。
func SwapTo(src interface{}, dest interface{}) error {
	marshal, err := json.Marshal(src)

	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, dest)

	if err != nil {
		return err
	}
	return nil
}
