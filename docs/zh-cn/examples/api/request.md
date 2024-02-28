# 请求参数

根据配置文件定义会生成接口对应的 handler，并会根据入参出参生成对应的结构体,执行 bind 和 validator 对入参进行校验。 默认情况下仅支持 json 方式传参。

在项目根目录下的 /app/param 目录创建出入参结构体：

```
// 入参
type GetUserInfoRequest struct {
	Uid int `json:"uid" binding:"required,min=1,max=100000"` // 用户ID

}
// 出参
type GetUserInfoResponse struct {
	Detail *User   `json:"detail" binding:""` // 用户详情
	List   []*User `json:"list" binding:""`   // 用户列表
}
type User struct {
	Uid  int    `json:"uid" binding:""`  // 用户ID
	Name string `json:"name" binding:""` // 用户名
}
```
在 handler 里面校验参数：
```go

// GetUserInfo 获取用户信息
func GetUserInfo(g *gin.Context) {
	ctx := api.NewContext(g)
	
	// 自动校验参数
	req := new(param.GetUserInfoRequest)
	err := ctx.BindValidator(req)
	if err != nil {
		ctx.Error(err)
		return
	}
	
	ctx.Success(def.GetUserInfoResponse{})
}
```