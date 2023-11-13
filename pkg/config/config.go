package config

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
	RateLimit     float64 `mapstructure:"rate_limit"`
	RSAPublicKey  string  `mapstructure:"rsa_public_key"`
	RSAPrivateKey string  `mapstructure:"rsa_private_key"`
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
	BaseDn    string `mapstructure:"base_dn"`
	Filter    string `mapstructure:"filter"`
	SizeLimit int    `mapstructure:"size_limit"`
	TimeLimit int    `mapstructure:"time_limit"`
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
