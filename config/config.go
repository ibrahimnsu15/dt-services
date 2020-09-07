package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var Conf *Config

func InitConf() {
	Conf = new(Config)
	Conf.setEnv()
	if err := Conf.setupViper(); err != nil {
		panic(err)
	}
	Conf.setupConfig()
}

type Config struct {
	Env          string
	ConsulHost   string
	ConsulPort   string
	MongoRC      []string
	GithubToken  string
	GithubAppKey []byte
}

func (c *Config) setEnv() {
	env := os.Getenv("ENV")
	if env == DEV {
		c.Env = DEV
		c.ConsulHost = CONSUL_HOST_CLUSTER
		c.ConsulPort = CONSUL_PORT_CLUSTER
	} else if env == STAGE {
		c.Env = STAGE
		c.ConsulHost = CONSUL_HOST_CLUSTER
		c.ConsulPort = CONSUL_PORT_CLUSTER
	} else if env == PROD {
		c.Env = PROD
		c.ConsulHost = CONSUL_HOST_CLUSTER
		c.ConsulPort = CONSUL_PORT_CLUSTER
	} else {
		c.Env = WORKSTATION
		c.ConsulHost = CONSUL_HOST_DEV
		c.ConsulPort = CONSUL_PORT_DEV
	}
}

func (c *Config) setupViper() error {
	consulUrl := fmt.Sprintf("%s:%s", c.ConsulHost, c.ConsulPort)
	if err := viper.AddRemoteProvider("consul", consulUrl, CONSUL_KV); err != nil {
		return err
	}
	viper.SetConfigType("json")
	if err := viper.ReadRemoteConfig(); err != nil {
		return err
	}

	return nil
}

func (c *Config) setupConfig() {
	if c.GetEnv() == WORKSTATION {
		c.MongoRC = MONGOHOSTS_WORKSTATION
	} else {
		mongoRCI := viper.Get("mongo_rc").([]interface{})
		var mongoRC []string
		for _, mongoHost := range mongoRCI {
			mongoRC = append(mongoRC, mongoHost.(string))
		}

		c.MongoRC = mongoRC
	}

	c.GithubAppKey = []byte(viper.Get("github_app_key").(string))
}

func (c *Config) GetEnv() string {
	return c.Env
}
