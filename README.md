1. 获取是否开启加密配置

GET /api/tmp/FHNCAPIS?ajaxmethod=is_encrypt

```javascript
{
    "enable": 1
}
```

2. 获取sessionid 用于POST请求的加解密

GET /api/tmp/FHNCAPIS?ajaxmethod=get_refresh_sessionid

```javascript
{
    "sessionid": "R46yw17E0CK5z51ny34IF1mxnBr4iE6b"
}
```

3. 获取是否有新短信

GET /api/tmp/FHAPIS?ajaxmethod=get_new_sms

```javascript
{
    "new_sms_flag": "false"
}
```

4. 心跳

GET /api/tmp/heartbeat