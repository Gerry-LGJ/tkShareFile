echo @off

:: 进入工程根目录，并创建.\pkg\tkShareFile_windows_amd64_latest文件夹
cd ..
md pkg\tkShareFile_windows_amd64_latest
md pkg\tkShareFile_windows_amd64_latest\www

:: 编译
set GOARCH=amd64
set GOOS=windows
go build -o pkg\tkShareFile_windows_amd64_latest\tkShareFile.exe main.go upload.go
set GOARCH=amd64
set GOOS=windows

:: 将编译后的文件复制到pkg\tkShareFile_windows_amd64_latest
:: /e复制目录和子目录，包括空的 /y禁止提示以确认改写一个现存目标文件 /d复制那些源时间比目标时间新的文件
xcopy www pkg\tkShareFile_windows_amd64_latest\www /eyd
