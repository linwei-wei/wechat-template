
BINARY_NAME=wechat_tpl_msg

build:
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}_windows_amd64.exe main.go
	GOARCH=arm64 GOOS=windows go build -o ${BINARY_NAME}_windows_arm64.exe main.go