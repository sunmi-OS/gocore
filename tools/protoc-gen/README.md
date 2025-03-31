# protoc-gen-go-gin
根据proto生成sunmi gin controller

## 安装
```bash
# protoc
brew install protobuf

# protoc-gen-go (生成 {package}.pb.go)
wget 'http://qiniu.brightguo.com/sunmi/protoc-gen-go_mac13.0' -O protoc-gen-go && chmod +x protoc-gen-go && mv protoc-gen-go $(go env GOPATH)/bin

# protoc-gen-go-gin（生成{package}_http_client.pb.go {package}_http_server.pb.go {package}_json.pb.go）
go install github.com/sunmi-OS/gocore/v2/tools/protoc-gen/cmd/protoc-gen-go-gin@latest

# protoc-gen-go-errors (生成{package}_ecode.pb.go)
go install github.com/sunmi-OS/gocore/v2/tools/protoc-gen/cmd/protoc-gen-go-errors@latest

# protoc-gen-go-openapiv2 生成{package}.swagger.json）
go install github.com/sunmi-OS/gocore/v2/tools/protoc-gen/cmd/protoc-gen-openapi@latest
  
# 拷贝third_party目录(protoc-gen-go-gin和protoc-gen-validate会用到)
git clone https://github.com/guoming0000/protoc-gen-go-gin.git
cp -r protoc-gen-go-gin/third_party $(go env GOPATH)/pkg/mod/github.com/guoming0000/
```

## 使用
```bash
protoc -I. -I ./third_party --go-gin_out=./ --go_out=./  api/article.proto

protoc --go-errors_out=./ api/article_error.proto

```

proto demo
https://github.com/guoming0000/protoc-gen-go-gin/blob/main/api/article/article.proto

生成后的文件demo
https://github.com/guoming0000/protoc-gen-go-gin/blob/main/api/article/

server使用demo：
https://github.com/guoming0000/protoc-gen-go-gin/blob/main/api/article/articleserver_test.go

client使用demo：
https://github.com/guoming0000/protoc-gen-go-gin/blob/main/api/article/client.go

## struct转proto定义方法
和chatgpt聊天，话术如下：
```bash
请把这个golang struct转换为proto形式:
type GetArticlesReq struct {
	Title string `json:"title,omitempty"`
	Page  int32  `json:"page,omitempty"`
	// 字段名使用小写下划线的风格，例如 string status_code = 1
	PageSize int32                    `json:"page_size,omitempty"`
	AuthorId int32                    `json:"author_id,omitempty"`
	Email    string                   `json:"email,omitempty"`
	Name     string                   `json:"name,omitempty"`
	Home     *GetArticlesReq_Location `json:"home,omitempty"`
}
```
## TODO
- [x] 支持proto-gen-sm-error-go
- [x] 支持自定义错误码，支持通过注释生成错误码字符串
- [x] 支持gin形式的binding参数校验方法
- [x] 支持生成swagger
- [x] 支持配置请求头
- [ ] 支持枚举
- [ ] 支持配置非标准json注解
- [ ] 支持通过yapi生成proto3
