# 返回&异常Code

项目接口统一返回的 json 结构为（下载文件类接口除外，返回文件流）

```json
{
  "code":1, // 整型
  "data": xxx, // 任意类型，对象、数组、null、字符串、整型等
  "msg": "" // string 类型
}
```

不同的 code 值代表着不同的语义，code 值有以下几种情况：

1. 接口正常返回的情况下 code 值为1。
2. 非语义化的异常 code 值为-1，例如查询数据异常。
3. 需要给出明确提示的错误，code 返回对应的值。例如登录接口中，用户名或密码错误时需要给出一个特定的code。
   
写一个函数或方法时，约定俗成是返回结果和 error，为了保持这种使用方式不变又能满足可以返回上述的这种结构，实现了ecode v2 库。
ecode v2 同时可以用于 grpc ，即兼容 gin 和 grpc。
ecode v2 库实现了 error 接口，并拓展了相关功能。

## ecode v2库使用方法

主要是在 api 和 biz 层使用 ecode 库，在其他地方使用 errors 库即可（特殊情况除外）。接下来举个简单的示例，首先在 /app/errcode/code.go 文件（错误信息统一收口到此文件，方便管理）定义对应的错误信息：

```go
package errcode

import "github.com/sunmi-OS/gocore/v2/api/ecode"

var (
	// 用户名或密码错误
	LoginErr = ecode.NewV2(40000, "incorrect username or password")
)
```

在 biz 层的某个文件中 例如 /app/biz/user.go 中的 login 方法：

```go
package biz

import (
	"xxx/errcode"
	"xxx/dal/account"
)

var UserHandler = &User{}

type User struct{}

func (u *User) Login(username, password string) error {
	userInfo, err := account.UserHandler.GetByUserame(username)
	if err != nil {
        // 返回使用 ecode v2 库生成的 error 对象
		return errcode.LoginErr
	}
	// 省略其他逻辑
	// ...
	return nil
}
```

在 api 层对应的方法中， 例如 /app/api/user.go 中的 login 方法：

```
package api

import (
    "xxx/biz"
    "github.com/gin-gonic/gin"
    "github.com/sunmi-OS/gocore/v2/api"
)

var UserHandler = &User{}

type User struct{}

func (u *User) Login(g *gin.Context) error {
    ctx := api.NewContext(g)
    // 省略其他逻辑
    // ...
    err := biz.UserHandler.Login(username, password)
    if err != nil {
       ctx.RetJSON(nil, err)
       return
    }
    ctx.RetJSON("", nil)
}
```