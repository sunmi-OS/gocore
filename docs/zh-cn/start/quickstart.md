快速开始
---

创建一个示例项目
```bash
# 创建工程文件夹
> mkdir test
> cd test

# 创建yaml配置文件模板gocore.yaml
> gocore conf create 
...
Welcome to GoCore, Configuration file has been generated.

# 修改gocore.yaml模板之后,根据yaml文件创建工程项目
> gocroe service create 

# 下次迭代增加新的接口或数据表更新代码
> gocroe service create 

```

工程创建时导入已有数据库
```bash
# 创建工程文件夹
> mkdir test 
> cd test

# 创建yaml配置文件模板gocore.yaml
> gocore conf create 

# 创建连接数据库的配置文件模板mysql.yaml
> gocore mysql create_yaml 

# 修改mysql.yaml之后,连接数据库将字段合并到gocore.yaml
> gocroe mysql add 

# 修改gocore.yaml模板之后,根据yaml文件创建工程项目
> gocroe service create 
```
