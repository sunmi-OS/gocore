<div align="center">

![logo](https://file.cdn.sunmi.com/logo.png?x-oss-process=image/resize,h_200)

</div>

Gocore Web Framework
---
English | [中文](README_ZH_CN.md)

[![Go Report Card](https://goreportcard.com/badge/github.com/sunmi-OS/gocore)](https://goreportcard.com/report/github.com/sunmi-OS/gocore)
[![GoDoc](https://godoc.org/github.com/sunmi-OS/gocore/v2?status.svg)](https://pkg.go.dev/github.com/sunmi-OS/gocore/v2)
[![Release](https://img.shields.io/github/v/release/sunmi-OS/gocore.svg?style=flat-square)](https://github.com/sunmi-OS/gocore/releases)

gocore is a highly integrated development framework and provides scaffolding for generating project structure, supports api, rpc, job and other development methods, and integrates various mainstream open source libraries into best practices, and ultimately realizes simplified processes, improved efficiency, and unified specifications.

![cli](https://file.cdn.sunmi.com/gocore_cli.svg)

## Features

- integrated widely-used libraries including gin, gorm, viper, zap, offering a robust and efficient foundation for low-level operations.
- gocore scaffolding is designed to expedite project setup by automating the creation of API routes, parameter bindings, and database schemas.
- support for multiple environment-specific configuration files and is compatible with the Nacos configuration center, enabling dynamic configuration loading and hot-swapping capabilities.
- Integrated with a range of essential utilities, including signature, encryption, file processing, mail delivery, random number generation, tracing and logging.
- adopt a non-intrusive design philosophy, enabling developers to concentrate on crafting business logic without the distraction of underlying system complexities.
- integrated standard Alibaba Cloud middleware such as SLS, RocketMQ, and Nacos, simplifying the utilization of cloud services.
- out-of-the-box, greatly simplifying the project startup and development workflow.

## Getting started

### Prerequisites

- **[Go](https://go.dev/)** >= 1.18
- **[Go module](https://github.com/golang/go/wiki/Modules)**


### Getting Gocore

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```
import "github.com/sunmi-OS/gocore/v2"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the `gocore` package:

```sh
$ go get -u github.com/sunmi-OS/gocore/v2
```

### Getting Gocore scaffolding

```sh
$ go install github.com/sunmi-OS/gocore/v2/tools/gocore@latest

$ gocore --version

   __ _    ___     ___    ___    _ __    ___
  / _` |  / _ \   / __|  / _ \  | '__|  / _ \
 | (_| | | (_) | | (__  | (_) | | |    |  __/
  \__, |  \___/   \___|  \___/  |_|     \___|
  |___/

gocore version v2.0.1
```

## Quick Start

### Create project directory

```sh
$ mkdir test
$ cd test
````

### Create the gocore.yaml file to generate the project structure

```sh
$ gocore yaml create 
```

### Create the project structure after modifying the gocore.yaml file as required

```sh
$ gocore service create 
```

## Configuration file description


```yaml
service:
  projectName: demo # Project name
  version: v1.0.0 # Project version
config:
  cNacos: true # Whether nacos is used
  cRocketMQConfig: true # Whether rocketMQ is used
  cMysql: # MySQL configuration
    - name: app # Database name
      hotUpdate: false # Hot update or not
      models: # model file
        - name: user # Table name
          auto: false # Whether to automatically create the table
          fields: # Table fields, gorm rules, one in a row
            - column:id;primary_key;type:int AUTO_INCREMENT
            - column:name;type:varchar(100) NOT NULL;default:'';comment:'User name';unique_index
          comment: User information table # Table remark
  cRedis: # Redis configuration
    - name: default # Redis name
      hotUpdate: false # Hot update or not
      index:
        db0: 0 # db index
rpcEnable: false # Whether to generate the rpc service layer
httpApiEnable: true # Whether to generate interface programs
jobEnable: true # Whether to generate a resident task
httpApis:
  host: 0.0.0.0 # Listening ip address
  port: "80" # Listening port
  apis:
    - prefix: /app/user # API interface prefix
      moduleName: user # Module name
      handle: # API interface
        - name: GetUserInfo # API handler name, full path is /app/user/GetUserInfo
          method: Any
          requestParams: # Request parameters
            - name: uid # Field name
              type: int # Field type
              comment: UserID # Field remark
              validate: required,min=1,max=100000 # Validate rules
          responseParams: # Response parameters
            - name: detail  # Field name
              type: '*User'  # Field type
              comment: User detail # Field remark
              validate: ""
            - name: list
              type: '[]*User'
              comment: User list
              validate: ""
          comment: Get user information
  params:
    User:
      - name: uid
        type: int
        comment: UserID
        validate: ""
      - name: name
        type: string
        comment: Username
        validate: ""
jobs:
  - name: InitUser # One-time task or resident task method names
    comment: Initialize user information # One-time task and resident task remark
```

## Project structure

Using the three-tier architecture (HTTP service: api->biz->dal, RPC service: rpc->biz->dal)：

- api(rpc): Application Programming Interface
  - Defines interface names, validates request parameters, invokes methods in the business logic layer (biz) to process business operations, and returns response data.
  - Can only call methods in the business logic layer (biz), and is prohibited from calling methods in the data access layer (dal).
- biz: Business Logic Layer
  - The layer responsible for processing business logic, it receives parameters from the API layer, utilizes methods from the DAL layer to complete business logic processing, and returns the necessary data.
  - Prohibited from calling methods in the API layer.
- dal: Data Access Layer
  - In charge of database access; inter-method calls within this layer are prohibited.
  - Prohibited from calling methods from both the API and Biz layers.

Directory structure：

```
├── app                  // Source code
│  ├── api               // API layer, Delete this folder if no HTTP service
│  ├── rpc               // RPC layer, Delete this folder if no RPC service
│  ├── biz               // Business Logic Layer
│  ├── dal               // Data Access Layer
│  ├── middleware        // Middleware
│  ├── cmd               // Task launch entry, defining the initialization methods for each component
│  │  ├── api.go
│  │  ├── init.go
│  │  └── job.go
│  ├── conf              // Configuration files
│  │  ├── base.go        // Basic configuration
│  │  └── local.go       // Configuration file for local debugging, the local environment variable needs to be set to RUN_TIME=local
│  ├── errcode           // Define error code
│  │  └── errcode.go
│  ├── job               // Task definition entry, scheduled tasks, one-time tasks, consumer queue tasks
│  ├── param             // Definition of request and response parameter structures.
│  │  └── user.go
│  ├── pkg               // Contains your dependencies
│  │  ├── locationtools  
│  │  │  └── country.go
│  │  └── util          
│  │      └── util.go
│  ├── route             // The Routes Directory
│  │   └── routers
│  ├── go.mod
│  ├── go.mod
│  └── main.go           // main
├── .gitignore
├── CODEOWNERS       
├── Dockerfile
└── README.md
```

## Contributing

Gocore is the work of hundreds of contributors. We appreciate your help!