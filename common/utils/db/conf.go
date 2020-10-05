package db

import "fmt"

type Config struct {
	// host=myhost port=myport user=gorm dbname=gorm password=mypassword
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	User     string `yaml:"user"`
	Dbname   string `yaml:"dbname"`
	Password string `yaml:"password"`
}

func (c *Config) ConnectInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Dbname, c.Password)
}

func DefaultConf() Config {
	return Config{
		// Host:     "172.17.0.1",
		Host:     "127.0.0.1",
		Port:     15432,
		User:     "postgres",
		Dbname:   "house",
		Password: "root",
	}
}
