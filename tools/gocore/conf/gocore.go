package conf

type GoCore struct {
	Service       Service
	Config        Config
	NacosEnable   bool // 是否开启 Nacos 默认开启
	HttpApiEnable bool // 是否开启HttpApi
	CronJobEnable bool // 是否开启 CronJob 默认不开启
	JobEnable     bool // 是否开启 Job 任务
	HttpApis      HttpApi
	CronJobs      []CronJob
	Jobs          []Job
}

type Service struct {
	ProjectName string // 项目名称
	Version     string // 项目版本
}

// HttpApi 路由拼接规则 /public/v1/项目名/模块名/接口名
// TODO: swagger.json导入
type HttpApi struct {
	Host string // 地址
	Port string // 端口
	Apis []Api
}

type Api struct {
	ModuleName string // 模块名
	Handle     []Handle
}

type Handle struct {
	Name           string  // 接口名
	Method         string  // Get Post Any
	RequestParams  []Param // 请求参数列表
	ResponseParams []Param // 返回参数列表
}

type Param struct {
	Name     string
	Required bool
	Type     string
	Title    string
	Params   []Param
}

type CronJob struct {
	Spec string // 定时规则
	Job  Job
}

type Job struct {
	Name  string // 任务名称
	Usage string // 任务描述
}

type Config struct {
	CNacos Nacos
	CMysql []Mysql
	CRedis []Redis
}

type Nacos struct {
	Env         bool
	EmailConfig bool
}

type Mysql struct {
	Name      string // Mysql名称，默认default
	HotUpdate bool
	Models    []Model
}

// Model TODO：支持建表SQL导入
type Model struct {
	Name   string  // 表名
	Auto   bool    // 是否自动创建表结构
	Fields []Field // 字段列表
}

type Field struct {
	Name     string // 字段名
	GormRule string // Gorm规则
	Index    bool   // 是否开启索引
}

type Redis struct {
	Name      string // Redis名称，默认default
	HotUpdate bool
	Index     map[string]int // index和Key的映射关系
}

func GetGocoreConfig() *GoCore {
	return &GoCore{
		Service: Service{
			ProjectName: "app",
			Version:     "v1.0.0",
		},
		Config: Config{
			CNacos: Nacos{
				EmailConfig: false,
			},
			CMysql: []Mysql{
				{
					Name: "default",
					Models: []Model{
						{
							Name: "user",
							Auto: false,
							Fields: []Field{
								{
									Name:     "Id",
									GormRule: "primary_key;type:BIGINT AUTO_INCREMENT",
									Index:    false,
								},
								{
									Name:     "Name",
									GormRule: "type:varchar(55) NOT NULL;default:'';comment:'用户名'",
									Index:    false,
								},
							},
						},
					},
				},
			},
			CRedis: []Redis{
				{
					Name: "default",
					Index: map[string]int{
						"default": 0,
					},
				},
			},
		},
		NacosEnable:   true,
		HttpApiEnable: true,
		CronJobEnable: false,
		JobEnable:     false,
		HttpApis: HttpApi{
			Host: "0.0.0.0",
			Port: "80",
			Apis: []Api{
				{
					ModuleName: "User",
					Handle: []Handle{
						{
							Name:   "GetUserInfo",
							Method: "Any",
							RequestParams: []Param{
								{
									Name:     "uid",
									Required: true,
									Type:     "int",
									Title:    "用户ID",
									Params:   nil,
								},
							},
							ResponseParams: []Param{
								{
									Name:     "code",
									Required: true,
									Type:     "int",
									Title:    "响应名",
									Params:   nil,
								},
								{
									Name:     "data",
									Required: true,
									Type:     "object",
									Title:    "返回内容",
									Params: []Param{
										{
											Name:     "uid",
											Required: true,
											Type:     "int",
											Title:    "用户ID",
											Params:   nil,
										},
										{
											Name:     "name",
											Required: true,
											Type:     "string",
											Title:    "用户名",
											Params:   nil,
										},
									},
								},
								{
									Name:     "msg",
									Required: true,
									Type:     "string",
									Title:    "返回异常文本描述",
									Params:   nil,
								},
							},
						},
					},
				},
			},
		},
		CronJobs: []CronJob{
			{
				Spec: "@every 30m",
				Job: Job{
					Name:  "SyncUser",
					Usage: "同步用户",
				},
			},
		},
		Jobs: []Job{
			{
				Name:  "InitUser",
				Usage: "初始化默认用户",
			},
		},
	}
}
