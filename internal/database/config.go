package database

type DBConfig struct {
	Host     string
	User     string
	DBName   string
	Port     int
	Password string
	SslMode  string
}
