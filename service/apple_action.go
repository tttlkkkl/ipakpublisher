package service

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/tttlkkkl/asc-go/asc"
	"github.com/tttlkkkl/ipakpublisher/com"
)

// IpaArgs ipa 处理子命令

// Ipa ipa 命令处理
func (s *AppleService) Ipa() error {
	ctx, cl := context.WithTimeout(context.Background(), s.AppleArgs.Auth.ExpiryDuration)
	defer cl()
	if s.AppleArgs.BundleID == "" {
		return errors.New("missing BundleID")
	}
	if s.AppleArgs.Init {
		if err := s.initLocal(ctx); err != nil {
			return err
		}
	}
	// 遍历工作目录
	// filepath.Walk(s.CmdArgs.Workdir.SubPath(com.SubDirApple), func(path string, info fs.FileInfo, err error) error {
	// 	return nil
	// })
	return nil
}

// 从 api 下载元数据并尝试初始化本地数据
func (s *AppleService) initLocal(ctx context.Context) error {
	// 初始化目录
	if err := s.initDir(com.SubDirApple, s.AppleArgs.Platform); err != nil {
		return err
	}
	subDir := s.CmdArgs.Workdir.SubPath(com.SubDirApple)
	// 获取应用信息
	apps, _, err := s.client.Apps.ListApps(ctx, &asc.ListAppsQuery{
		FilterBundleID: []string{s.AppleArgs.BundleID},
	})
	var app asc.App
	if err != nil {
		return err
	} else if len(apps.Data) == 1 {
		app = apps.Data[0]
		fileName := filepath.Join(subDir, s.AppleArgs.Platform, "app.json")
		if err := WriteMetaDataToLocal(fileName, app); err != nil {
			return err
		}
	} else {
		return errors.New("can not get app info from api")
	}
	// get versions
	versions, _, err := s.client.Apps.ListAppStoreVersionsForApp(ctx, app.ID, &asc.ListAppStoreVersionsQuery{
		FilterVersionString: []string{"0.0.1"},
		FilterPlatform:      []string{s.AppleArgs.Platform},
	})
	var version asc.AppStoreVersion
	if err != nil {
		return err
	} else if len(versions.Data) == 1 {
		version = versions.Data[0]
		fileName := filepath.Join(subDir, s.AppleArgs.Platform, "version.json")
		if err := WriteMetaDataToLocal(fileName, version); err != nil {
			return err
		}
	} else {
		return errors.New("can not get version info from api")
	}
	// Get all localizations for the version
	localizations, _, err := s.client.Apps.ListLocalizationsForAppStoreVersion(ctx, version.ID, nil)

	for _, v := range localizations.Data {
		fileName := filepath.Join(subDir, s.AppleArgs.Platform, *v.Attributes.Locale+".json")
		if err := WriteMetaDataToLocal(fileName, v); err != nil {
			return err
		}
	}
	return nil
}
