package service

import (
	"encoding/json"
	"os"
)

// WriteMetaDataToLocal 覆盖写入本地 json 文件
func WriteMetaDataToLocal(file string, data interface{}) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	// 清空文件内容
	if err := f.Truncate(0); err != nil {
		return err
	}
	// 从头开始写
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}
	// json.MarshalIndent(&data, prefix string, indent string)
	jr := json.NewEncoder(f)
	jr.SetIndent("", "\t")
	return jr.Encode(&data)
}

// ReadMetaDataFromLocal 读取本地 json 文件
func ReadMetaDataFromLocal(file string, data any) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	// 从头开始写
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}
	jr := json.NewDecoder(f)
	return jr.Decode(data)
}
