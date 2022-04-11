[中文说明](README-zh.md)
## ipakpublisher cli
This program can be limited to help you edit the application information on the Apple App store. Help you upload and publish apps to Google play.

It will work well in CI / CD automation pipeline.

Realize the command line operation of [App Store Connect API](https://developer.apple.com/documentation/ipakpublisherapi).
Realize the command line operation of [Google Play Developer API](https://developers.google.com/android-publisher/api-ref/rest).

## Usage
### set api auth info
- use toml config file,like:
```toml
[apple]
# Key ID for the given private key, described in App Store Connect
KeyID="xxxx"
# Issuer ID for the App Store Connect team
IssuerID="xxxx"
# A duration value for the lifetime of a token. App Store Connect does not accept a token with a lifetime of longer than 20 minutes
ExpiryDuration ="20m"
# The bytes of the PKCS#8 private key created on App Store Connect. Keep this key safe as you can only download it once.
PrivateKeyFile ="/path/xxxx.p8"
[google]
application_credentials_file=""
```
- Using environment variables
```bash
 export AT_APPLE_KEYID=xxx
 export AT_APPLE_ISSUERID=yyy
 export AT_APPLE_PRIVATEKEYFILE=/path
 export AT_APPLE_EXPIRYDURATION=20s
 export AT_APPLE_EXPIRYDURATION=20s
 export AT_GOOGLE_APPLICATION_CREDENTIALS_FILE=./xxx
```
## Usage
### Global
```bash
Usage:
  ipakpublisher [command]

Available Commands:
  aab         push the aab file
  apk         push the apk file
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  ipa         Perform IPA related operations

Flags:
  -c, --config string    config file for app store connect api auth (default is $HOME/.ipakpublisher.toml)
  -h, --help             help for ipakpublisher
  -w, --workdir string   workdir path (default is $(pwd))
```
### ipa
```bash
Perform IPA related operations.

Usage:
  ipakpublisher ipa [flags]

Flags:
  -b, --bundle-id string   Specify the ipa bundleid.
  -h, --help               help for ipa
  -i, --init               If you specify this option, metadata is downloaded from the API and an attempt is made to overwrite the local record.When this operation is completed, the program terminates.
  -p, --platform string    Specify the platform. (default "IOS")
  -s, --sync               If you specify this option, To simultaneously transfer the local data to the remote location through API.
  -v, --version string     the app version string,e.g. 1.0.0

Global Flags:
  -c, --config string    config file for app store connect api auth (default is $HOME/.ipakpublisher.toml)
  -w, --workdir string   workdir path (default is $(pwd))
```
### apk
```bash
Use this command when you want to publish an app in APK format

Usage:
  ipakpublisher apk [flags]

Flags:
  -f, --file string           the app file path.
  -h, --help                  help for apk
  -i, --init                  If you specify this option, A multilingual change log configuration template will be generated locally. Due to API limitations, this method will not pull data from the remote.
      --notes-file string     This option loads a text content as the update log.
  -n, --package-name string   app package name
  -p, --platform string       Specify the platform. (default "android")
  -s, --sync                  If you specify this option, To simultaneously transfer the local data to the remote location through API.
      --timeout string        Set upload timeout (default "2m")
  -t, --track strings         If this option is set.
                                it will be published to these tracks after the application upload is completed.
                                You can specify multiple times.
                                The optional values are:
                                alpha
                                beta
                                internal
                                production

Global Flags:
  -c, --config string    config file for app store connect api auth (default is $HOME/.ipakpublisher.toml)
  -w, --workdir string   workdir path (default is $(pwd))
```
### aab
```bash
Use this command when you want to publish an app in AAB format

Usage:
  ipakpublisher aab [flags]

Flags:
  -f, --file string           the app file path.
  -h, --help                  help for aab
  -i, --init                  If you specify this option, A multilingual change log configuration template will be generated locally. Due to API limitations, this method will not pull data from the remote.
      --notes-file string     This option loads a text content as the update log.
  -n, --package-name string   app package name
  -p, --platform string       Specify the platform. (default "android")
  -s, --sync                  If you specify this option, To simultaneously transfer the local data to the remote location through API.
      --timeout string        Set upload timeout (default "2m")
  -t, --track strings         If this option is set.
                                it will be published to these tracks after the application upload is completed.
                                You can specify multiple times.
                                The optional values are:
                                alpha
                                beta
                                internal
                                production

Global Flags:
  -c, --config string    config file for app store connect api auth (default is $HOME/.ipakpublisher.toml)
  -w, --workdir string   workdir path (default is $(pwd))
```
