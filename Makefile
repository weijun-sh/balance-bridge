all:
	env GO111MODULE=on GOPROXY=https://goproxy.io go build -o balance *.go
