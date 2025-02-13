package model

type UserRegistrationData struct {
	FirstName  string
	MiddleName string
	LastName   string
	Login      string
	Password   string
}

type Config struct {
	DBConfig    DBConfig
	AppConfig   AppConfig
	JWTConfig   JWTConfig
	RedisConfig RedisConfig
}

type AppConfig struct {
	Host string
}
type DBConfig struct {
	DBConnection string
}

type JWTConfig struct {
	SecretKey     string
	TokenDuration int
}

type RedisConfig struct {
	Servers  string
	Password string
	TTL      int64
}
