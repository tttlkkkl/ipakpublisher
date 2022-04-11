package service

import (
	"context"
	"io/fs"
	"path/filepath"
	"regexp"

	"github.com/tttlkkkl/ipakpublisher/com"
	androidpublisher "google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/googleapi"
)

type GoogleVersionLocalization struct {
	ReleaseNotes string `json:"releaseNotes"`
}

// 尝试拉取远程数据
func (s *GoogleService) InitLocal(ctx context.Context) error {
	if err := s.initDir(com.SubDirGoogle, s.GoogleArgs.Platform, "loc"); err != nil {
		return err
	}
	// google play api 貌似没有提供获取 app 信息的接口
	// 所以此处尝试生成一个模版
	return WriteMetaDataToLocal(s.getAppLocalizationFile("en-GB"), &GoogleVersionLocalization{ReleaseNotes: "the change log"})
}

// Do 开发处理
// ext aab or apk
func (s *GoogleService) Do(ext string) (err error) {

	ctx, cl := context.WithTimeout(context.Background(), s.App.timeout)
	defer cl()
	// 尝试初始化本地路径
	if s.GoogleArgs.Init {
		return s.InitLocal(ctx)
	}
	s.App.service, err = androidpublisher.NewService(ctx)
	if err != nil {
		return err
	}
	s.App.edit = androidpublisher.NewEditsService(s.App.service)
	appEdit, err := s.App.edit.Insert(s.GoogleArgs.PackageName, &androidpublisher.AppEdit{}).Do()
	if err != nil {
		return err
	}
	s.App.editID = appEdit.Id
	defer func() {
		if err != nil {
			com.Log.Error("An error occurred during execution. All submissions will be rolled back")
			err = s.App.edit.Delete(s.GoogleArgs.PackageName, s.App.editID).Do()
		} else {
			_, err = s.App.edit.Commit(s.GoogleArgs.PackageName, s.App.editID).Do()
		}
	}()
	// 上传
	if ext == "aab" {
		if err = s.aabUpload(); err != nil {
			com.Log.Error(err)
			return err
		}
	} else {
		if err = s.apkUpload(); err != nil {
			com.Log.Error(err)
			return err
		}
	}
	// 发布
	if s.GoogleArgs.Sync {
		for _, v := range s.GoogleArgs.Track {
			if err = s.release(com.Track(v)); err != nil {
				com.Log.Error(err)
				return err
			}
		}
	}
	return nil
}

// apkUpload app 上传
func (s *GoogleService) apkUpload() error {
	apkEdit := androidpublisher.NewEditsApksService(s.App.service)
	rs, err := apkEdit.Upload(s.GoogleArgs.PackageName, s.App.editID).Media(s.App.file, googleapi.ContentType("application/octet-stream")).Do()
	if err != nil {
		return err
	}
	com.Log.Infof("Apk upload completed, package version code is:%d", rs.VersionCode)
	s.App.versionCode = rs.VersionCode
	return nil
}

// aabUpload app 上传
func (s *GoogleService) aabUpload() error {
	aabEdit := androidpublisher.NewEditsBundlesService(s.App.service)
	call := aabEdit.Upload(s.GoogleArgs.PackageName, s.App.editID)
	// 文件大小大于 100M 时要加这个选项
	if s.App.fileInfo.Size() > 104857600 {
		call.AckBundleInstallationWarning(true)
	}
	rs, err := call.Media(s.App.file, googleapi.ContentType("application/octet-stream")).Do()
	if err != nil {
		return err
	}
	com.Log.Infof("Aab upload completed, package version code is:%d", rs.VersionCode)
	s.App.versionCode = rs.VersionCode
	return nil
}

// release 发布到渠道
func (s *GoogleService) release(trackString com.Track) error {
	track := androidpublisher.NewEditsTracksService(s.App.service)

	var releaseNotes []*androidpublisher.LocalizedText
	filepath.Walk(s.getAppLocalizationDir(), func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			var loc string
			rg, _ := regexp.Compile(`(.*)\.json$`)
			rs := rg.FindStringSubmatch(info.Name())
			if len(rs) == 2 {
				loc = rs[1]
			}
			var locData GoogleVersionLocalization
			if err := ReadMetaDataFromLocal(path, &locData); err != nil {
				return err
			}
			d := &androidpublisher.LocalizedText{
				Language: loc,
				Text:     locData.ReleaseNotes,
			}
			releaseNotes = append(releaseNotes, d)
		}
		return err
	})
	tk := androidpublisher.Track{
		Releases: []*androidpublisher.TrackRelease{
			{
				VersionCodes: googleapi.Int64s{s.App.versionCode},
				Status:       com.StatusCompleted.String(),
				ReleaseNotes: releaseNotes,
			},
		},
	}
	_, err := track.Update(s.GoogleArgs.PackageName, s.App.editID, trackString.String(), &tk).Do()
	if err != nil {
		return err
	}
	return nil
}

// func (s *GoogleService) getAppFile() string {
// 	return filepath.Join(s.SubDir, s.GoogleArgs.Platform, "app.json")
// }

// func (s *GoogleService) getAppVersionFile() string {
// 	return filepath.Join(s.SubDir, s.GoogleArgs.Platform, "version.json")
// }

func (s *GoogleService) getAppLocalizationDir() string {
	return filepath.Join(s.SubDir, s.GoogleArgs.Platform, "loc")
}

func (s *GoogleService) getAppLocalizationFile(loc string) string {
	return filepath.Join(s.getAppLocalizationDir(), loc+".json")
}
