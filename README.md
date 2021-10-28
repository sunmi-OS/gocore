<a href="https://sunmi.com"><img height="180" src="https://file.cdn.sunmi.com/gocore-logo.png"></a>

# Integrated Development Framework



---

## Installation && Usage

```bash
go get github.com/sunmi-OS/gocor/v2
```

```go
import (
	"github.com/sunmi-OS/gocore/v2/xxxxx"
)
```

### Supported Go versions

- Golang > 1.13
- [Go module](https://github.com/golang/go/wiki/Modules)

---

## Examples
- 简单工程创建
```
mkdir vsim #创建工程文件夹
cd vsim
gocore conf create #创建yaml配置文件模板gocore.yaml
gocroe service create #修改gocore.yaml模板之后,根据yaml文件创建工程项目
```

- 工程创建时导入已有数据库
```
mkdir vsim #创建工程文件夹
cd vsim
gocore conf create #创建yaml配置文件模板gocore.yaml
gocore mysql create_yaml #创建连接数据库的配置文件模板mysql.yaml
gocroe mysql add #修改mysql.yaml之后,连接数据库将字段合并到gocore.yaml
gocroe service create #修改gocore.yaml模板之后,根据yaml文件创建工程项目
```
