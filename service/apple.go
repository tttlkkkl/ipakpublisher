package service

import (
	"errors"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/tttlkkkl/asc-go/asc"
	"github.com/tttlkkkl/ipakpublisher/com"
)

// AppleCmdArgs apple 关联的命令行参数
type AppleCmdArgs struct {
	// Init 是否从 api 拉取数据并覆盖本地设置
	Init bool
	// 是否将本地元数据同步到api
	Sync     bool
	BundleID string
	Platform string
	Version  string
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

// AppleService apple
type AppleService struct {
	Service
	client    *asc.Client
	AppleArgs *AppleCmdArgs
	SubDir    string
}

func initAppleConfig(args *AppleCmdArgs) (err error) {
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
	if err := initAppleConfig(appleArgs); err != nil {
		return nil, err
	}
	auth, err := asc.NewTokenConfig(appleArgs.Auth.KeyID, appleArgs.Auth.IssuerID, appleArgs.Auth.ExpiryDuration, appleArgs.Auth.privateKeyData)
	if err != nil {
		return nil, err
	}
	c := &AppleService{
		Service:   Service{CmdArgs: cmdArgs},
		client:    asc.NewClient(auth.Client()),
		AppleArgs: appleArgs,
		SubDir:    cmdArgs.Workdir.SubPath(com.SubDirApple),
	}
	// c.client.SetHTTPDebug(true)
	return c, nil
}
