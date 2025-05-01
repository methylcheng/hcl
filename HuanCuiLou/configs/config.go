package configs

import "time"

// MySQLConfig 定义 MySQL 配置结构体
type MySQLConfig struct {
	DSN string
}

// JwtConfig 定义 JWT 配置结构体
type JwtConfig struct {
	SecretKey                  string        //密钥
	AccessTokenExpireDuration  time.Duration //accessToken过期时间
	RefreshTokenExpireDuration time.Duration //refreshToken过期时间
}

// CodeConfig 定义验证码配置结构体
type CodeConfig struct {
	ExpireDuration time.Duration // 定时删除过期验证码的间隔时间
}

// Config 定义配置结构体
type Config struct {
	MySQL MySQLConfig
	Jwt   JwtConfig
	Code  CodeConfig
}

// GetConfig 获取配置实例
func GetConfig() Config {
	const (
		host     = "116.196.120.25"
		port     = "3307"
		user     = "HuanCuiLou"
		password = "HuanCuiLou"
		dbname   = "huancuilou"
	)

	return Config{
		MySQL: MySQLConfig{
			DSN: user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname +
				"?charset=utf8mb4&parseTime=True&loc=Local",
		},
		Jwt: JwtConfig{
			SecretKey:                  "huancuilou",
			AccessTokenExpireDuration:  time.Hour,
			RefreshTokenExpireDuration: time.Hour * 7 * 24,
		},
		Code: CodeConfig{
			ExpireDuration: time.Minute * 2,
		},
	}
}
