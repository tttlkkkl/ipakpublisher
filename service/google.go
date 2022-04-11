package service

import (
	"errors"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/tttlkkkl/ipakpublisher/com"
	androidpublisher "google.golang.org/api/androidpublisher/v3"
)

// GoogleCmdArgs google 命令行操作
type GoogleCmdArgs struct {
	Platform             string
	FilPath              string
	PackageName          string
	Track                []string
	ReleaseNotesFilePath string
	Timeout              string
	// Init 是否从 api 拉取数据并覆盖本地设置
	Init bool
	// 是否将本地元数据同步到api
	Sync bool
}

// App 应用程序包定义
type App struct {
	timeout     time.Duration
	editID      string
	versionCode int64
	service     *androidpublisher.Service
	edit        *androidpublisher.EditsService
	file        *os.File
	fileInfo    os.FileInfo
}

// GoogleService google service
type GoogleService struct {
	Service
	GoogleArgs *GoogleCmdArgs
	App        App
	SubDir     string
}

func (s *GoogleService) initGoogleConfig() (err error) {
	// 如果是初始化不需要进行额外的参数验证
	if s.GoogleArgs.Init {
		return nil
	}
	if key := viper.GetString("google.application_credentials_file"); key == "" {
		return errors.New("you need set application_credentials_file")
	} else {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", key)
	}
	if s.GoogleArgs.PackageName == "" {
		return errors.New("Package name must be provided")
	}
	if s.GoogleArgs.FilPath == "" {
		return errors.New("the app file must be provided")
	}
	f, err := os.Open(s.GoogleArgs.FilPath)
	if err != nil {
		return err
	}
	s.App.file = f
	fInfo, err := os.Stat(s.GoogleArgs.FilPath)
	if err != nil {
		return err
	}
	s.App.fileInfo = fInfo
	if s.GoogleArgs.Timeout != "" {
		if du, err := time.ParseDuration(s.GoogleArgs.Timeout); err != nil {
			return err
		} else {
			s.App.timeout = du
		}
	}
	if s.App.timeout <= 0 {
		s.App.timeout = 2 * time.Minute
	}
	return nil
}
func NewGoogleService(cmdArgs *CmdArgs, googleArgs *GoogleCmdArgs) (*GoogleService, error) {
	c := GoogleService{
		Service: Service{
			CmdArgs: cmdArgs,
		},
		GoogleArgs: googleArgs,
		App:        App{},
		SubDir:     cmdArgs.Workdir.SubPath(com.SubDirGoogle),
	}
	if err := c.initGoogleConfig(); err != nil {
		return nil, err
	}
	return &c, nil
}
