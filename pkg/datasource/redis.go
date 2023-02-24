package datasource

import (
	"fmt"
	configPkg "github.com/Adesubomi/magic-ayo-api/pkg/config"
	"github.com/go-redis/redis"
)

func connectByString(conf *configPkg.Config) (*redis.Client, error) {
	options, err := redis.ParseURL(fmt.Sprintf(
		"rediss://%v:%v@%v:%v",
		conf.Redis.Username,
		conf.Redis.Password,
		conf.Redis.Host,
		conf.Redis.Port,
	))

	if err != nil {
		return nil, err
	}

	return redis.NewClient(options), nil
}

func connectByOptions(conf *configPkg.Config) (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Network: "",
		Addr: fmt.Sprintf(
			"%v:%v", // localhost:6379
			conf.Redis.Host,
			conf.Redis.Port,
		),
		Password: conf.Redis.Password,
		DB:       0,
	}), nil
}

func RedisConnection(conf *configPkg.Config) (*redis.Client, error) {
	var client *redis.Client
	var err error

	if conf.Redis.Username != "" {
		client, err = connectByString(conf)
	} else {
		client, err = connectByOptions(conf)
	}

	if err != nil {
		fmt.Printf(" ?? Could not connect to Redis because: %v\n", err)
		return nil, err
	}

	if _, err = client.Ping().Result(); err != nil {
		fmt.Printf(" ?? Could not connect to Redis because: %v\n", err)
		return client, err
	}

	fmt.Println(" ✔ Redis Connection Established")
	return client, nil
}
