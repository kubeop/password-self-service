package config

// 配置yaml文件用"-"区分单词, 转为驼峰方便

var Setting Config

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Redis   RedisConfig   `mapstructure:"redis"`
	Ldap    LdapConfig    `mapstructure:"ldap"`
	Cron    CronConfig    `mapstructure:"cron"`
	Channel ChannelConfig `mapstructure:"channel"`
}

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

// CronConfig 定时任务配置
type CronConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Schedule string `mapstructure:"schedule"`
}

type ChannelConfig struct {
	PlatformUrl    string     `mapstructure:"platform-url"`
	VerifyChannel  string     `mapstructure:"verify-channel"`
	ExpiredChannel string     `mapstructure:"expired-channel"`
	Mail           MailConfig `mapstructure:"mail"`
	AliyunSms      AliyunSms  `mapstructure:"aliyunsms"`
	TencentSms     TencentSms `mapstructure:"tencentsms"`
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

type AliyunSms struct {
	AccessKeyId         string `mapstructure:"access-key-id"`
	AccessKeySecret     string `mapstructure:"access-key-secret"`
	SignName            string `mapstructure:"sign-name"`
	TemplateCodeVerify  string `mapstructure:"template-code-verify"`
	TemplateCodeExpired string `mapstructure:"template-code-expired"`
}

type TencentSms struct {
	AccessKeyId         string `mapstructure:"access-key-id"`
	AccessKeySecret     string `mapstructure:"access-key-secret"`
	SignName            string `mapstructure:"sign-name"`
	Region              string `mapstructure:"region"`
	AppId               string `mapstructure:"app-id"`
	TemplateCodeVerify  string `mapstructure:"template-code-verify"`
	TemplateCodeExpired string `mapstructure:"template-code-expired"`
}
