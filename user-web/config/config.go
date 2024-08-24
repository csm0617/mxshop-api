package config

/*
user-grpc服务配置
*/
type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

/*
jwt配置
*/
type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

/*
*
redis配置
*/
type RedisConfig struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	Expire int    `mapstructure:"expire"`
}

/*
*
consul配置
*/
type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

/*
当前user-web服务配置
*/
type ServerConfig struct {
	Name        string        `mapstructure:"name"`
	Host        string        `mapstructure:"host"`
	Port        int           `mapstructure:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user-srv"`
	JWTInfo     JWTConfig     `mapstructure:"jwt"`
	RedisInfo   RedisConfig   `mapstructure:"redis"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul"`
}
