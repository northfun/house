package redis

import "github.com/go-redis/redis"

type Config struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func NewClient(conf *Config) (*redis.Client, error) {
	rdc := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password, // no password set
		DB:       conf.DB,       // use default DB
	})

	stt := rdc.Ping()
	return rdc, stt.Err()
}
