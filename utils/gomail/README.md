

配置文件如下

```toml
[email]
host="smtp.exmail.qq.com"
port=465
username="xxx@xxxxx.com"
password="XXXXXXXXX"
```


使用方式如下

```go
gomail.SendEmail(
	"wenzhenxi@sunmi.com",     // 发送给谁
	"service@sunmi.com",       // 发送者的邮箱
	"SUNMI",                   // 发送者的名称
	"SUNMI激活邮件",            // 邮件主题
	"URL:xxxxxxxxx",           //  邮件内容
)
```


