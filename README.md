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
go build ./...
```