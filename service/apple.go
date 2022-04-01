package service

import (
	"errors"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/tttlkkkl/asc-go/asc"
)

// AppleCmdArgs apple 关联的命令行参数
type AppleCmdArgs struct {
	// Init 是否从 api 拉取数据并覆盖本地设置
	Init     bool
	BundleID string
	Platform string
	Auth     AppStoreAuth
}

type AppStoreAuth struct {
	// Key ID for the given private key, described in App Store Connect
	KeyID string `toml:"KeyID"`
	// Issuer ID for the App Store Connect team
	IssuerID string `toml:"IssuerID"`
	// A duration value for the lifetime of a token. App Store Connect does not accept a token with a lifetime of longer than 20 minutes
	ExpiryDuration time.Duration `toml:"ExpiryDuration"`
	// The bytes of the PKCS#8 private key created on App Store Connect. Keep this key safe as you can only download it once.
	PrivateKeyFile string `toml:"PrivateKeyFile"`
	privateKeyData []byte
}

// MetaVersionLocalization 多语言化版本信息
type MetaVersionLocalization struct {
	Description     string `toml:"description"`
	Keywords        string `toml:"keywords"`
	MarketingURL    string `toml:"marketing_url"`
	PromotionalText string `toml:"promotional_text"`
	SupportURL      string `toml:"support_url"`
}

// AppleMetaData apple 应用元数据
type AppleMetaData struct {
	MetaVersionLocalization []MetaVersionLocalization
}

// AppleService apple
type AppleService struct {
	Service
	client    *asc.Client
	AppleArgs *AppleCmdArgs
}

func initConfig(args *AppleCmdArgs) (err error) {
	args.Auth.KeyID = viper.GetString("apple.KeyID")
	args.Auth.IssuerID = viper.GetString("apple.IssuerID")
	args.Auth.ExpiryDuration = viper.GetDuration("apple.ExpiryDuration")
	args.Auth.PrivateKeyFile = viper.GetString("apple.PrivateKeyFile")
	if args.Auth.KeyID == "" {
		return errors.New("missing KeyID")
	}
	if args.Auth.IssuerID == "" {
		return errors.New("missing IssuerID")
	}
	if args.Auth.ExpiryDuration > (20 * time.Minute) {
		return errors.New("App Store Connect does not accept a token with a lifetime of longer than 20 minutes")
	}
	if args.Auth.ExpiryDuration <= 0 {
		args.Auth.ExpiryDuration = 20 * time.Minute
	}
	if args.Auth.PrivateKeyFile == "" {
		return errors.New("missing PrivateKeyFile")
	}
	args.Auth.privateKeyData, err = os.ReadFile(args.Auth.PrivateKeyFile)
	if err != nil {
		return err
	}
	return nil
}

func NewAppleService(cmdArgs *CmdArgs, appleArgs *AppleCmdArgs) (*AppleService, error) {
	if err := initConfig(appleArgs); err != nil {
		return nil, err
	}
	auth, err := asc.NewTokenConfig(appleArgs.Auth.KeyID, appleArgs.Auth.IssuerID, appleArgs.Auth.ExpiryDuration, appleArgs.Auth.privateKeyData)
	if err != nil {
		return nil, err
	}
	return &AppleService{Service: Service{CmdArgs: cmdArgs}, client: asc.NewClient(auth.Client()), AppleArgs: appleArgs}, nil
}
