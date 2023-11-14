package config

// 配置yaml文件用"-"区分单词, 转为驼峰方便

var (
	Server = &ServerConfig{}
	Redis  = &RedisConfig{}
	Ldap   = &LdapConfig{}
	Mail   = &MailConfig{}
)

// ServerConfig 配置
type ServerConfig struct {
	Mode          string  `mapstructure:"mode"`
	Addr          string  `mapstructure:"addr"`
	RateLimit     float64 `mapstructure:"rate-limit"`
	RSAPublicKey  string  `mapstructure:"rsa-public-key"`
	RSAPrivateKey string  `mapstructure:"rsa-private-key"`
}

// RedisConfig 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LdapConfig 配置
type LdapConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	TLS       bool   `mapstructure:"tls"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	Base      string `mapstructure:"base"`
	Filter    string `mapstructure:"filter"`
	SizeLimit int    `mapstructure:"size-limit"`
	TimeLimit int    `mapstructure:"time-limit"`
}

// MailConfig 邮箱配置
type MailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	TLS      bool   `mapstructure:"tls"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}
