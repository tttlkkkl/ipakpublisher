package service

import (
	"errors"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/tttlkkkl/asc-go/asc"
)

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

// CmdLineArgs 命令行参数
type CmdLineArgs struct {
	// WorkDir 工作目录
	WorkDir string
	// ConfigFile 配置文件
	ConfigFile string
	Auth       AppStoreAuth
}

// VersionLocalizationArgs 多语言化版本信息
type VersionLocalizationArgs struct {
	Description     string
	Keywords        string
	MarketingURL    string
	PromotionalText string
	SupportURL      string
}

// AppleService apple
type AppleService struct {
	client *asc.Client
	args   *CmdLineArgs
}

func InitConfig(args *CmdLineArgs) (err error) {
	if args.WorkDir == "" {
		args.WorkDir, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	args.Auth.KeyID = viper.GetString("KeyID")
	args.Auth.IssuerID = viper.GetString("IssuerID")
	args.Auth.ExpiryDuration = viper.GetDuration("ExpiryDuration")
	args.Auth.PrivateKeyFile = viper.GetString("PrivateKeyFile")
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

func NewAppleService(args *CmdLineArgs) (*AppleService, error) {
	auth, err := asc.NewTokenConfig(args.Auth.KeyID, args.Auth.IssuerID, args.Auth.ExpiryDuration, args.Auth.privateKeyData)
	if err != nil {
		return nil, err
	}
	// ctx, cl := context.WithTimeout(context.Background(), args.Auth.ExpiryDuration)
	return &AppleService{client: asc.NewClient(auth.Client()), args: args}, nil
}
