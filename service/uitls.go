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
	// json.MarshalIndent(&data, prefix string, indent string)
	jr := json.NewEncoder(f)
	jr.SetIndent("", "\t")
	return jr.Encode(&data)
}
