package service

import (
	"os"
	"path/filepath"

	"github.com/tttlkkkl/ipakpublisher/com"
)

// CmdLineArgs 命令行参数
type CmdArgs struct {
	// Workdir 工作目录
	Workdir  com.Workdir
	RootPath string
	// ConfigFile 配置文件
	ConfigFile string
}

type Service struct {
	CmdArgs *CmdArgs
}

func InitConfig(args *CmdArgs) (err error) {
	if args.RootPath == "" {
		args.RootPath, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	args.Workdir = com.Workdir(args.RootPath)
	return nil
}

func (s *Service) initDir(sub com.SubDir, subDir ...string) error {
	dir := s.CmdArgs.Workdir.SubPath(sub)
	dir = filepath.Join(dir, filepath.Join(subDir...))
	_, err := os.Stat(dir)
	// 目录不存在则创建目录
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
