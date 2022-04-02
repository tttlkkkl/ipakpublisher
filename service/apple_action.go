package service

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"regexp"

	"github.com/tttlkkkl/asc-go/asc"
	"github.com/tttlkkkl/ipakpublisher/com"
)

// AppleAppInfo 自定义 app info
type AppleAppInfo struct {
	ID string `json:"id"`
	asc.AppAttributes
}

// AppleVersionInfo 自定义 app 版本信息
type AppleVersionInfo struct {
	ID string `json:"id"`
	asc.AppStoreVersionAttributes
}

// AppleAppStoreVersionLocalization 自定义本地化信息
type AppleVersionLocalization struct {
	asc.AppStoreVersionLocalizationAttributes
}

// Ipa ipa 命令处理
func (s *AppleService) Ipa() error {
	ctx, cl := context.WithTimeout(context.Background(), s.AppleArgs.Auth.ExpiryDuration)
	defer cl()
	subDir := s.CmdArgs.Workdir.SubPath(com.SubDirApple)
	// 初始化
	if s.AppleArgs.Init {
		if s.AppleArgs.BundleID == "" {
			return errors.New("missing BundleID")
		}
		if s.AppleArgs.Version == "" {
			return errors.New("You may need to specify the version number currently in use.")
		}
		if err := s.initLocal(ctx, subDir); err != nil {
			return err
		}
		// 执行初始化操作后不能再执行其他操作
		return nil
	}
	if s.AppleArgs.Sync {
		if err := s.syncMetadata(ctx, subDir); err != nil {
			return err
		}
	}
	return nil
}

// 从 api 下载元数据并尝试初始化本地数据
func (s *AppleService) initLocal(ctx context.Context, subDir string) error {
	// 初始化目录
	if err := s.initDir(com.SubDirApple, s.AppleArgs.Platform, "loc"); err != nil {
		return err
	}
	var fileName string
	// 获取应用信息
	apps, _, err := s.client.Apps.ListApps(ctx, &asc.ListAppsQuery{
		FilterBundleID: []string{s.AppleArgs.BundleID},
	})
	var app asc.App
	if err != nil {
		return err
	} else if len(apps.Data) == 1 {
		app = apps.Data[0]
		fileName = s.getAppFile()
		if err := WriteMetaDataToLocal(fileName, AppleAppInfo{
			ID:            app.ID,
			AppAttributes: *app.Attributes,
		}); err != nil {
			return err
		}
	} else {
		return errors.New("can not get app info from api")
	}
	// get versions
	versions, _, err := s.client.Apps.ListAppStoreVersionsForApp(ctx, app.ID, &asc.ListAppStoreVersionsQuery{
		FilterVersionString: []string{s.AppleArgs.Version},
		FilterPlatform:      []string{s.AppleArgs.Platform},
	})
	if err != nil {
		return err
	}
	var version asc.AppStoreVersion
	if err != nil {
		return err
	} else if len(versions.Data) == 1 {
		version = versions.Data[0]
		fileName = s.getAppVersionFile()
		if err := WriteMetaDataToLocal(fileName, AppleVersionInfo{
			ID:                        version.ID,
			AppStoreVersionAttributes: *version.Attributes,
		}); err != nil {
			return err
		}
	} else {
		return errors.New("can not get version info from api")
	}
	// Get all localizations for the version
	localizations, _, err := s.client.Apps.ListLocalizationsForAppStoreVersion(ctx, version.ID, nil)
	for _, v := range localizations.Data {
		fileName := s.getAppLocalizationFile(*v.Attributes.Locale)
		// 避免语言信息被写入本地文件，方便复制和移植
		v.Attributes.Locale = nil
		if err := WriteMetaDataToLocal(fileName, AppleVersionLocalization{
			AppStoreVersionLocalizationAttributes: *v.Attributes,
		}); err != nil {
			return err
		}
	}
	return nil
}

// 同步本地元数据
func (s *AppleService) syncMetadata(ctx context.Context, subDir string) error {
	// 尝试获取本地 app 信息
	var app AppleAppInfo
	if err := ReadMetaDataFromLocal(s.getAppFile(), &app); err != nil {
		return err
	}
	if app.ID == "" {
		return errors.New("If the app ID is missing, you can try to specify the -i option to initialize the local storage")
	}
	// 尝试获取本地版本信息
	var version AppleVersionInfo
	var err error
	if err := ReadMetaDataFromLocal(s.getAppVersionFile(), &version); err != nil {
		return err
	}
	if _, _, err := s.client.Apps.UpdateAppStoreVersion(ctx, version.ID, &asc.AppStoreVersionUpdateRequestAttributes{
		Copyright:    version.Copyright,
		Downloadable: version.Downloadable,
		// EarliestReleaseDate :version.CreatedDate,
		ReleaseType: version.ReleaseType,
		// UsesIDFA:      version.UsesIDFA,
		VersionString: version.VersionString,
	}, nil); err != nil {
		return err
	}
	// 遍历工作目录
	err = filepath.Walk(s.getAppLocalizationDir(), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			var loc string
			rg, _ := regexp.Compile(`(.*)\.json$`)
			rs := rg.FindStringSubmatch(info.Name())
			if len(rs) == 2 {
				loc = rs[1]
			}
			var locData AppleVersionLocalization
			if err := ReadMetaDataFromLocal(path, &locData); err != nil {
				return err
			}
			// 从远程读取列表
			locMp, err := s.getLocalizationMap(ctx, version.ID)
			if err != nil {
				return err
			}
			// 存在则更新，否则创建
			if locLineV, ok := locMp[loc]; !ok {
				data := asc.AppStoreVersionLocalizationCreateRequestAttributes{
					Description:     locData.Description,
					Keywords:        locData.Keywords,
					Locale:          loc,
					MarketingURL:    locData.MarketingURL,
					PromotionalText: locData.PromotionalText,
					SupportURL:      locData.SupportURL,
				}
				// 只有从线上读到这个字段才能执行创建
				if locLineV.IsWhatNew {
					data.WhatsNew = locData.WhatsNew
				}
				if _, _, err := s.client.Apps.CreateAppStoreVersionLocalization(ctx, data, version.ID); err != nil {
					return err
				}
			} else {
				data := asc.AppStoreVersionLocalizationUpdateRequestAttributes{
					Description:     locData.Description,
					Keywords:        locData.Keywords,
					MarketingURL:    locData.MarketingURL,
					PromotionalText: locData.PromotionalText,
					SupportURL:      locData.SupportURL,
				}
				// 只有从线上读到这个字段才能执行更新
				if locLineV.IsWhatNew {
					data.WhatsNew = locData.WhatsNew
				}
				if _, _, err := s.client.Apps.UpdateAppStoreVersionLocalization(ctx, locLineV.VersionData.ID, &data); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

// VersionOnline 从线上读取版本信息
type VersionOnline struct {
	VersionData asc.AppStoreVersionLocalization
	IsWhatNew   bool
}
type VersionOnlineMap map[string]VersionOnline

// 获取线上
func (s *AppleService) getLocalizationMap(ctx context.Context, versionID string) (VersionOnlineMap, error) {
	var rs = make(VersionOnlineMap)
	localizations, _, err := s.client.Apps.ListLocalizationsForAppStoreVersion(ctx, versionID, nil)
	if err != nil {
		return rs, err
	}
	for _, v := range localizations.Data {
		k := *v.Attributes.Locale
		val := VersionOnline{
			VersionData: v,
			IsWhatNew:   false,
		}
		if v.Attributes.WhatsNew != nil {
			val.IsWhatNew = true
		}
		rs[k] = val
	}
	return rs, err
}

func (s *AppleService) getAppFile() string {
	return filepath.Join(s.SubDir, s.AppleArgs.Platform, "app.json")
}

func (s *AppleService) getAppVersionFile() string {
	return filepath.Join(s.SubDir, s.AppleArgs.Platform, "version.json")
}

func (s *AppleService) getAppLocalizationDir() string {
	return filepath.Join(s.SubDir, s.AppleArgs.Platform, "loc")
}

func (s *AppleService) getAppLocalizationFile(loc string) string {
	return filepath.Join(s.getAppLocalizationDir(), loc+".json")
}
