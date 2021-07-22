package conf

type GoCore struct {
	Service  Service
	Config   Config
	HttpApis HttpApi
	CronJobs []CronJob
	Jobs     []Job
	Models   []Model
}

type Service struct {
	ProjectName string // 项目名称
	Version     string // 项目版本
}

// HttpApi 路由拼接规则 /public/v1/项目名/模块名/接口名
// TODO: swagger.json导入
type HttpApi struct {
	Enable bool           // 是否开启HttpApi
	Host   string         // 地址
	Port   string         // 端口
	Apis   map[string]Api // key 模块名
}

type Api struct {
	Handle Handle // 路由配置
}

type Handle struct {
	Name   string  // 接口名
	Method string  // Get Post Any
	Params []Param // 参数列表
}

type Param struct {
	Name     string
	Required string
	Type     string
	Title    string
	Params   []Param
}

type CronJob struct {
	Enable bool   // 是否开启 CronJob 默认不开启
	Spec   string // 定时规则
}

type Job struct {
	Enable bool   // 是否开启 Job 任务
	Name   string // 任务名称
	Usage  string // 任务描述
}

type Config struct {
	CNacos Nacos
	CMysql []Mysql
	CRedis []Redis
}

type Nacos struct {
	Enable bool // 是否开启 Nacos 默认开启
}

type Mysql struct {
	Name   string // Mysql名称，默认default
	Models []Model
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
	Index    bool   // 是否开启所有
}

type Redis struct {
	Name  string         // Redis名称，默认default
	Index map[string]int // index和Key的映射关系
}
