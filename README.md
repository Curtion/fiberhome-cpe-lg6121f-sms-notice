# 说明

自动读取烽火 5G CPE 内 SIM 卡短信，并通过 bark 通知到手机。

要求程序运行在能够访问烽火 5G CPE 的网络环境中。

如果你手动登录 CPE，程序会自动被登出，此时程序会默认暂停监听 10（可手动指定）分钟，然后再次尝试监听。


# 使用方法

`./fh.exe -u [username] -p [password] -b [barkKey] -s [sound] -url http://192.168.8.1 -w 10`

- `-u` 烽火 5G CPE 登录用户名
  - 如：`admin`
- `-p` 烽火 5G CPE 登录密码
  - 如：`u7A1g`
- `-b` bark 通知 key
  - 如：`Rr2YbnXjGzqNB73ccTo2`（单设备端）
  - 如：`Rr2YbnXjGzqNB73ccTo2,Tr3Yb6XjUz9NB73ccFg1`（多设备端用英文逗号分隔）
- `-s` bark 通知铃声名称
  - 可选项，默认为：`healthnotification`
- `-url` 烽火 5G CPE 地址
  - 可选项，默认为：`http://192.168.8.1`
- `-w` 被动下线后距离下次开始监听间隔，单位分钟
  - 可选项，默认为：`10`

# 编译

```shell
# Windows PowerShell
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o ./bin/fh_sms_notice-v0.0.3-darwin-amd64 main.go aes.go fiberhome.go
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o ./bin/fh_sms_notice-v0.0.3-darwin-arm64 main.go aes.go fiberhome.go
$env:GOOS="linux"; $env:GOARCH="386"; go build -o ./bin/fh_sms_notice-v0.0.3-linux-386 main.go aes.go fiberhome.go
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o ./bin/fh_sms_notice-v0.0.3-linux-amd64 main.go aes.go fiberhome.go
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -o ./bin/fh_sms_notice-v0.0.3-linux-arm64 main.go aes.go fiberhome.go
$env:GOOS="windows"; $env:GOARCH="386"; go build -o ./bin/fh_sms_notice-v0.0.3-windows-386.exe main.go aes.go fiberhome.go
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o ./bin/fh_sms_notice-v0.0.3-windows-amd64.exe main.go aes.go fiberhome.go
```

```shell
# Linux/macOS Bash
GOOS=darwin GOARCH=amd64 go build -o ./bin/fh_sms_notice-v0.0.3-darwin-amd64 main.go aes.go fiberhome.go
GOOS=darwin GOARCH=arm64 go build -o ./bin/fh_sms_notice-v0.0.3-darwin-arm64 main.go aes.go fiberhome.go
GOOS=linux GOARCH=386 go build -o ./bin/fh_sms_notice-v0.0.3-linux-386 main.go aes.go fiberhome.go
GOOS=linux GOARCH=amd64 go build -o ./bin/fh_sms_notice-v0.0.3-linux-amd64 main.go aes.go fiberhome.go
GOOS=linux GOARCH=arm64 go build -o ./bin/fh_sms_notice-v0.0.3-linux-arm64 main.go aes.go fiberhome.go
GOOS=windows GOARCH=386 go build -o ./bin/fh_sms_notice-v0.0.3-windows-386.exe main.go aes.go fiberhome.go
GOOS=windows GOARCH=amd64 go build -o ./bin/fh_sms_notice-v0.0.3-windows-amd64.exe main.go aes.go fiberhome.go
```