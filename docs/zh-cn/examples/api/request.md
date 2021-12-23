# 请求参数


更具配置文件定义会生成接口对应的func，并会更具入参出参生成对应的结构体,执行bind和validator对入参进处理

默认情况下仅支持json方式传参

参数验证对生产环境屏蔽细节

[参数验证逻辑使用：go-playground/validator](https://github.com/go-playground/validator)


err := ctx.BindValidator(req)

会判断环境，生产环境不提供具体细节提升安全


```go

// GetUserInfo 获取用户信息
func GetUserInfo(g *gin.Context) {
	ctx := api.NewContext(g)
	req := new(def.GetUserInfoRequest)
	err := ctx.BindValidator(req)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Success(def.GetUserInfoResponse{})
}

...

type GetUserInfoRequest struct {
	Uid int `json:"uid" binding:"required,min=1,max=100000"` // 用户ID

}

type GetUserInfoResponse struct {
	Detail *User   `json:"detail" binding:""` // 用户详情
	List   []*User `json:"list" binding:""`   // 用户列表
}

type User struct {
	Uid  int    `json:"uid" binding:""`  // 用户ID
	Name string `json:"name" binding:""` // 用户名
}

```