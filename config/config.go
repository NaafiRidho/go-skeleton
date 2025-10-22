package config

import (
	"github.com/sirupsen/logrus"
	"os"
	"user-service/common/utils"
)

var Config AppConfig

type AppConfig struct {
	Port                  int      `json:"port"`
	AppName               string   `json:"appName"`
	AppKey                string   `json:"appKey"`
	SignatureKey          string   `json:"signatureKey"`
	Database              Database `json:"database"`
	RateLimiterMaxRequest float64  `json:"rateLimiterMaxRequest"`
	RateLimiterTimeSecond int      `json:"rateLimiterTimeSecond"`
	JwtSecretKey          string   `json:"jwtSecret"`
	JwtExpireTime         int      `json:"jwtExpireTime"`
}

type Database struct {
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	Name                  string `json:"name"`
	UserName              string `json:"username"`
	Password              string `json:"password"`
	MaxOpenConnections    int    `json:"maxOpenConnections"`
	MaxLifeTimeConnection int    `json:"maxLifeTimeConnection"`
	MaxIdleConnections    int    `json:"maxIdleConnections"`
	MaxIdleTime           int    `json:"maxIdleTime"`
}

func Init() {
	err := utils.BindFromJson(&Config, "config.json", ".")
	if err != nil {
		logrus.Infof("Failed to bind config from file: %v", err)
		err = utils.BindFromConsul(&Config, os.Getenv(os.Getenv("CONSUL_HTTP_URL")), os.Getenv("CONSUL_HTTP_KEY"))
		if err != nil {
			panic(err)
		}
	}
}
