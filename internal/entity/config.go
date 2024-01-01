package entity

type Config struct {
	HTTP HTTPConfig
	DB   DBConfig
}

type HTTPConfig struct {
	IP   string
	Port string
}

type DBConfig struct {
	MySQL MySQLConfig
	Redis RedisConfig
}

type MySQLConfig struct {
	Driver string
	DBName string
	Host   string
	Port   string
	User   string
	Pass   string
}

type RedisConfig struct {
	Host string
	Pass string
	DB   string
}
