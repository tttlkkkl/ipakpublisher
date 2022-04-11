build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o  ipakpublisher-linux-amd64 -work $(shell pwd)/
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o  ipakpublisher-windows-amd64.exe -work $(shell pwd)/
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o  ipakpublisher-darwin-amd64 -work $(shell pwd)/
ipas:
	go build -o ipakpublisher
	./ipakpublisher -c $(shell pwd)/../config.toml -w $(shell pwd)/.workdir ipa -s -b 'com.dex.ss' -v '0.0.1'
ipai:
	go build -o ipakpublisher
	./ipakpublisher -c $(shell pwd)/../config.toml -w $(shell pwd)/.workdir ipa -i -b 'com.zly.ss' -v '2.1.0'
aabi:
	go build -o ipakpublisher
	./ipakpublisher -c $(shell pwd)/../config.toml -w $(shell pwd)/.workdir aab -i
apki:
	go build -o ipakpublisher
	./ipakpublisher -c $(shell pwd)/../config.toml -w $(shell pwd)/.workdir aab -i -n 'com.wd.ss' -f /Users/m/Downloads/xx.apk -t internal
aabs:
	go build -o ipakpublisher
	./ipakpublisher -c $(shell pwd)/../config.toml -w $(shell pwd)/.workdir aab -s -n 'com.wd.ss' -f /Users/m/Downloads/xx.aab -t internal
apks:
	go build -o ipakpublisher
	./ipakpublisher -c $(shell pwd)/../config.toml -w $(shell pwd)/.workdir aab -s -n 'com.wd.ss' -f /Users/m/Downloads/xx.apk -t internal