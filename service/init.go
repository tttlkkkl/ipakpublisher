package service

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tttlkkkl/ipakpublisher/com"
)

// InitLocalWorkDir 初始化本地工作目录
func (s *AppleService) InitLocalWorkDir() error {
	iosDir := filepath.Join(s.args.WorkDir, com.SubDirIos.String())
	_, err := os.Stat(iosDir)
	// 目录不存在则创建目录
	if os.IsNotExist(err) {
		if err = os.MkdirAll(iosDir, os.ModePerm); err != nil {
			return err
		}
	} else {
		return err
	}
	err = filepath.Walk(iosDir, func(path string, info fs.FileInfo, err error) error {
		fmt.Println("------>", path)
		return nil
	})
	return err
}
