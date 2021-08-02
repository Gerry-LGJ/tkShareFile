echo @off
set GOARCH=arm
set GOOS=linux

go build -o tkShareFile-arm-linux main.go upload.go

set GOARCH=amd64
set GOOS=windows
