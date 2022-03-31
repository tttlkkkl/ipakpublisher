## connect cli
Realize the command line operation of connect OpenAPI. 

## apple store api

This is a generated Go project for Apple's [App Store Connect API](https://developer.apple.com/documentation/ipakpublisherapi).

Apple provided an OpenAPI
[downloable specification](https://developer.apple.com/sample-code/app-store-connect/app-store-connect-openapi-specification.zip) which allows use to generate client api codes.

Client code generation was done with [OpenAPI Generator](https://github.com/OpenAPITools/openapi-generator)

The generation process is as follows：
- Download the API description file and unzip it .
- Then run the following command:
```bash
docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli:latest-release generate -i /local/app_store_connect_api_1.8_openapi.json -g go -o /local/openapi
```

## Usage
### set api auth info
- use toml config file,like:
```toml
# Key ID for the given private key, described in App Store Connect
KeyID="xxx"
# Issuer ID for the App Store Connect team
IssuerID="ccc"
# A duration value for the lifetime of a token. App Store Connect does not accept a token with a lifetime of longer than 20 minutes
ExpiryDuration ="20m"
# The bytes of the PKCS#8 private key created on App Store Connect. Keep this key safe as you can only download it once.
PrivateKeyFile ="/filepath"
```
- Using environment variables
```bash
 export AT_KEYID=xxx
 export AT_ISSUERID=yyy
 export AT_PRIVATEKEYFILE=/path
 export AT_EXPIRYDURATION=20s
```