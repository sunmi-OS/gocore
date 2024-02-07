package conf

import (
	"os"
	"strings"
)

type GoCore struct {
	Service       Service `yaml:"service"`
	Config        Config  `yaml:"config"`
	HttpApiEnable bool    `yaml:"httpApiEnable"` // 是否开启HttpApi，开启后生成接口表示层
	RPCEnable     bool    `yaml:"rpcEnable"`     // 是否开启 rpc server,开启后生成rpc服务器表示层
	JobEnable     bool    `yaml:"jobEnable"`     // 是否开启 Job 任务
	HttpApis      HttpApi `yaml:"httpApis"`
	Jobs          []Job   `yaml:"jobs"`
}

type Service struct {
	ProjectName string `yaml:"projectName"` // 项目名称
	Version     string `yaml:"version"`     // 项目版本
}

// HttpApi 路由拼接规则 /public/v1/项目名/模块名/接口名
// TODO: swagger.json导入
type HttpApi struct {
	Host   string             `yaml:"host"` // 地址
	Port   string             `yaml:"port"` // 端口
	Apis   []Api              `yaml:"apis"`
	Params map[string][]Param `yaml:"params"`
}

type Api struct {
	Prefix     string   `yaml:"prefix"`     //路由前缀
	ModuleName string   `yaml:"moduleName"` // 模块名
	Handle     []Handle `yaml:"handle"`
}

type Handle struct {
	Name           string  `yaml:"name"`           // 接口名
	Method         string  `yaml:"method"`         // Get Post Any
	RequestParams  []Param `yaml:"requestParams"`  // 请求参数列表
	ResponseParams []Param `yaml:"responseParams"` // 返回参数列表
	Comment        string  `yaml:"comment"`        //接口描述
}

type Param struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Comment  string `yaml:"comment"`
	Validate string `yaml:"validate"`
}

type Job struct {
	Name    string `yaml:"name"` // 任务名称
	Comment string `yaml:"comment"`
}

type Config struct {
	CNacos          bool    `yaml:"cNacos"`
	CRocketMQConfig bool    `yaml:"cRocketMQConfig"`
	CMysql          []Mysql `yaml:"cMysql"`
	CRedis          []Redis `yaml:"cRedis"`
}

type Mysql struct {
	Name      string  `yaml:"name"` // Mysql名称，默认default
	HotUpdate bool    `yaml:"hotUpdate"`
	Models    []Model `yaml:"models"`
}

type Model struct {
	Name    string   `yaml:"name"`   // 表名
	Auto    bool     `yaml:"auto"`   // 是否自动创建表结构
	Fields  []string `yaml:"fields"` // 字段列表
	Comment string   `yaml:"comment"`
}

type Redis struct {
	Name      string         `yaml:"name"` // Redis名称，默认default
	HotUpdate bool           `yaml:"hotUpdate"`
	Index     map[string]int `yaml:"index"` // index和Key的映射关系
}

func GetGocoreConfig() *GoCore {
	projectName := "demo"
	// 获取当前目录名称
	pwd, err := os.Getwd()
	if err != nil {
		panic("failed to get the rooted path name corresponding to the current directory")
	}
	path := strings.Replace(pwd, "\\", "/", -1)
	arr := strings.Split(path, "/")
	if len(arr) > 1 {
		projectName = arr[len(arr)-1]
	}
	return &GoCore{
		Service: Service{
			ProjectName: projectName,
			Version:     "v1.0.0",
		},
		Config: Config{
			CNacos:          true,
			CRocketMQConfig: true,
			CMysql: []Mysql{
				{
					Name: "demo",
					Models: []Model{
						{
							Name: "user",
							Auto: false,
							Fields: []string{
								"column:id;primary_key;type:int AUTO_INCREMENT",
								"column:name;type:varchar(100) NOT NULL;default:'';comment:'用户名';unique_index",
							},
							Comment: "用户表",
						},
					},
				},
			},
			CRedis: []Redis{
				{
					Name: "demo",
					Index: map[string]int{
						"db0": 0,
					},
				},
			},
		},
		RPCEnable:     false,
		HttpApiEnable: true,
		JobEnable:     true,
		HttpApis: HttpApi{
			Host: "0.0.0.0",
			Port: "80",
			//Params: map[string][]Param{
			//	"User": {
			//		{
			//			Name:    "id",
			//			Type:    "int64",
			//			Comment: "用户ID",
			//		},
			//		{
			//			Name:    "name",
			//			Type:    "string",
			//			Comment: "用户名",
			//		},
			//	},
			//},
			Apis: []Api{
				{
					ModuleName: "user",
					Prefix:     "/app/user",
					Handle: []Handle{
						{
							Name:    "getUserInfo",
							Method:  "POST",
							Comment: "获取用户信息",
							RequestParams: []Param{
								{
									Name:     "id",
									Type:     "int64",
									Comment:  "用户ID",
									Validate: "required,min=1,max=100000",
								},
							},
							ResponseParams: []Param{
								{
									Name:    "id",
									Type:    "int64",
									Comment: "用户ID",
								},
								{
									Name:    "name",
									Type:    "string",
									Comment: "用户名",
								},
							},
						},
					},
				},
			},
		},
		Jobs: []Job{
			{
				Name:    "InitUser",
				Comment: "初始化默认用户",
			},
		},
	}
}
