build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o  ipakpublisher-linux-amd64 -work $(shell pwd)/
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o  ipakpublisher-windows-amd64.exe -work $(shell pwd)/
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o  ipakpublisher-darwin-amd64 -work $(shell pwd)/
run:
	go build -o ipakpublisher
	./ipakpublisher -c $(shell pwd)/../config.toml -w $(shell pwd)/.workdir ipa -s -b 'com.dex.slope' -v '0.0.1'
runw:
	go build -o ipakpublisher
	./ipakpublisher -c $(shell pwd)/../config.toml -w $(shell pwd)/.workdir ipa -i -b 'com.zly.wallet' -v '2.1.0'