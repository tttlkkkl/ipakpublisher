package service

import (
	"os"
	"time"

	"github.com/tttlkkkl/ipakpublisher/com"
	androidpublisher "google.golang.org/api/androidpublisher/v3"
)

// GoogleCmdArgs google 命令行操作
type GoogleCmdArgs struct {
	FilPath              string
	PackageName          string
	Track                []string
	ReleaseNotes         string
	ReleaseNotesFilePath string
	Timeout              string
}

// App 应用程序包定义
type App struct {
	packageName  string
	file         *os.File
	fileInfo     os.FileInfo
	track        []com.Track
	releaseNotes string
	appEdit      AppEdit
	timeout      time.Duration
	action       Action
}

// AppEdit app 暂存数据
type AppEdit struct {
	id          string
	versionCode int64
}

// Action 数据操作
type Action struct {
	service *androidpublisher.Service
	edit    *androidpublisher.EditsService
}

// GoogleService google service
type GoogleService struct {
	Service
	GoogleArgs *GoogleCmdArgs
	service    *androidpublisher.Service
	edit       *androidpublisher.EditsService
}

func initGoogleConfig(args *GoogleCmdArgs) (err error) {
	return nil
}
func NewGoogleService(cmdArgs *CmdArgs, googleArgs *GoogleCmdArgs) (*GoogleService, error) {
	c := GoogleService{}
	return &c, nil
}
