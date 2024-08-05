# go-sweet-email
# 基于 go 语言的邮件发送工具
## 适用版本 go 1.20
## 使用方法
1. 引入包
```text
go get github.com/PurpleScorpion/go-sweet-email
```

2. 基本使用
```text
    1. 注册邮箱基本信息
	sweetEmail.SetPort(465) // 端口
	sweetEmail.SetHost("smtp.163.com") // 邮箱服务器
	sweetEmail.SetUserName("xxxxx@163.com") // 发件者邮箱
	sweetEmail.SetPassword("XXXXXXXXX") // 发件者邮箱识别码(注意不是登录密码)
	sweetEmail.SetEmailName("你的名字") // 发件人名字,最终会和你的邮箱一起显示 例如: 刻晴<keqing@163.com>
	// 打开日志 - 非必要
	sweetEmail.OpenLog()
	// 注册邮箱服务 - 必须
    sweetEmail.Register()
    var to = []string{"123456@qq.com", "888888@qq.com"}
	sweetEmail.SendEmail(to, "test", "测试邮件发送")
```
3. 自动注册
```text
若想使用自动注册, 则需要满足以下2点要求
1. 需使用go-sweet框架编写的项目

github地址: https://github.com/PurpleScorpion/go-sweet


2. yaml格式必须为以下写法

sweet:
  email:
    port: 25
    host: smtp.qq.com
    username: 123456@qq.com
    password: xtjqjqjqjqjqjqjq
    emailname: Sweet

代码示例:
    sweetEmail.AutoRegister(true)
	var to = []string{"123456@qq.com", "888888@qq.com"}
	sweetEmail.SendEmail(to, "test", "测试邮件发送")
```