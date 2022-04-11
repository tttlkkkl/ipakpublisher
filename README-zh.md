## ipakpublisher cli
本应用程序通过苹果 App store api 和 google play api 实现 app 的管理工作。主要用于自动化流水线。
苹果的不具备文件上传功能，仍然需要通过苹果官方提供的命令行工具进行上传。然后通过本程序完成提审元数据填充，有多个国家地区支持时特别有用。之后仍然需要手动提审。
google 的可以上传 apk 、aab 包并发送到对应的渠道。如 alpha、beta、internal、production。
apple api 参考：[App Store Connect API](https://developer.apple.com/documentation/ipakpublisherapi).
google api 参考：[Google Play Developer API](https://developers.google.com/android-publisher/api-ref/rest).

## 安装
## 使用
### 配置
- 需要提供 api 访问授权配置，配置文件格式如下：
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
- 使用命令行参数`-c`来指定配置文件绝对路径。
- 或者将配置信息写入 `$HOME/.ipakpublisher.toml` 中。
### 工作目录
工作目录中保存本程序需要填充的元数据信息，默认为当前 shell 工作目录。
- 可通过指定 `-i` 选项来初始化本地工作目录。
- apple 将从 api 拉取数据在本地初始化工作目录数据。在使用过程中使用本参数也可用线上元数据覆盖本地元数据。
- google 无法从 api 拉取元数据信息，将生成一个发布日志固定文件模板。

## 使用示例
```bash
# 拉取应用商店元数据，并初始化本地工作目录
ipakpublisher -c $(pwd)/../config.toml -w $(pwd)/.workdir ipa -i -b 'com.zly.ss' -v '2.1.0'
# 从本地工作目录读取应用元数据信息并同步到远程
ipakpublisher -c $(pwd)/../config.toml -w $(pwd)/.workdir ipa -s -b 'com.dex.ss' -v '0.0.1'
# 初始化 android 本地模板
ipakpublisher -c $(pwd)/../config.toml -w $(pwd)/.workdir aab -i
# 上传文件并向目标渠道发布版本
ipakpublisher -c $(pwd)/../config.toml -w $(pwd)/.workdir aab -s -n 'com.wd.ss' -f /Users/m/Downloads/xx.apk -t internal
```