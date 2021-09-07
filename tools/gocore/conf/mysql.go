package conf

type MysqlDb struct {
	Host     string   `yaml:"host"`
	Port     string   `yaml:"port"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
	Database string   `yaml:"database"`
	Tables   []string `yaml:"tables"`
}
