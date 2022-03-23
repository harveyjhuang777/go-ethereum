package model

type IOpsConfig interface {
	GetOpsMySQLConfig() MySQLOps
}

type IAppConfig interface {
	GetLogConfig() LogConfig
	GetGinConfig() GinConfig
}

type OpsConfig struct {
	MySQLOps MySQLOps `mapstructure:"mysql_ops"`
}

type MySQLOps struct {
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Debug    bool   `mapstructure:"debug"`
}

type AppConfig struct {
	LogConfig LogConfig `mapstructure:"log_config"`
	GinConfig GinConfig `mapstructure:"gin_config"`
}

type LogConfig struct {
	Name  string `mapstructure:"name"`
	Env   string `mapstructure:"env"`
	Level string `mapstructure:"level"`
}

type GinConfig struct {
	Address string `mapstructure:"address"`
}
