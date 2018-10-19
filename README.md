# cxton-sms-sdk
北京空间畅想(短信)接口开发协议Golang SDK

# 短信发送
```
c := cxtonsms.NewClient("http://xxx.xxx.xxx.xxx", "name", "passwd")
b, err := client.SendStrongUTF8(&ctxonsms.Strong{
    Dest:    strings.Split("xxxxxxxx,xxxxxxxx", ","),
    Content: "content",
})
```

# 余额查询
```
c := cxtonsms.NewClient("http://xxx.xxx.xxx.xxx", "name", "passwd")
b, err := client.GetBalanceUTF8(nil)
```
