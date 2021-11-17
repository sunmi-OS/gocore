快速开始
---

通过gocore工具可以快速生成开发骨架，框架会提供一个gocore.yaml文件来管理动态生成配置，开发者可以通过对yaml文件进行修改来定义cmd、api、job、cron、model、config以及中间件等 。

特性：
- 当创建项目同时会执行mod init、mod tidy、fmt、goimports来保障项目符合Golang标准
- 在数据表结构创建支持连接mysql反向生成model结构
- @TODO 未来将支持从swagger导入和导出swagger功能

创建一个示例项目
```bash
# 创建工程文件夹
> export PROJECT_NAME=demo
> mkdir PROJECT_NAME
> cd PROJECT_NAME

# 创建yaml配置文件模板gocore.yaml
> gocore conf create 
...
Welcome to GoCore, Configuration file has been generated.

# 修改gocore.yaml模板之后,根据yaml文件创建工程项目
> gocore service create 

   __ _    ___     ___    ___    _ __    ___
  / _` |  / _ \   / __|  / _ \  | '__|  / _ \
 | (_| | | (_) | | (__  | (_) | | |    |  __/
  \__, |  \___/   \___|  \___/  |_|     \___|
  |___/

Run go mod init.
[11/11] Initialize the Request return parameters... 100% [========================================]   
Run go mod tidy .
Run go fmt .
goimports -l -w .
Welcome to GoCore, the project has been initialized.


# 下次迭代增加新的接口或数据表更新代码
> gocore service create 

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
> gocore mysql add 

# 修改gocore.yaml模板之后,根据yaml文件创建工程项目
> gocore service create 
```
