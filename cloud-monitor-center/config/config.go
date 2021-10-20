package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App            string    `yaml:"app"`
	Serve          Serve     `yaml:"serve"`
	DB             DB        `yaml:"db"`
	Logger         LogConfig `yaml:"logger"`
	Rocketmq       Rocketmq  `yaml:"rocketmq"`
	TenantUrl      string    `yaml:"tenantUrl"`
	HasNoticeModel bool      `yaml:"hasNoticeModel"`
}

type Serve struct {
	Debug        bool `yaml:"debug"`
	Port         int  `yaml:"port"`
	ReadTimeout  int  `yaml:"read_timeout"`
	WriteTimeout int  `yaml:"write_timeout"`
}

type DB struct {
	Dialect       string        `yaml:"dialect"`
	Url           string        `yaml:"url"`
	MaxIdleConnes int           `yaml:"max_idle_connes"`
	MaxOpenConnes int           `yaml:"max_open_connes"`
	MaxLifeTime   time.Duration `yaml:"time.Hour"`
}

type LogConfig struct {
	Debug         bool   `yaml:"debug"`
	App           string `yaml:"app"`
	DataLogPrefix string `yaml:"data_log_prefix"`
	Group         string `yaml:"group"`
	MaxSize       int    `yaml:"max_size"`
	MaxBackups    int    `yaml:"max_backups"`
	MaxAge        int    `yaml:"max_age"`
	Compress      bool   `yaml:"compress"`
	Stdout        bool   `yaml:"stdout"`
}

type Rocketmq struct {
	NameServer             string `yaml:"name-server"`
	AlertContactTopic      string `yaml:"alertContactTopic"`
	AlertContactGroupTopic string `yaml:"alertContactGroupTopic"`
	AlertContactGroup      string `yaml:"alertContactGroup"`
	AlertContactGroupGroup string `yaml:"alertContactGroupGroup"`
}

var config Config

func InitConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig() *Config {
	return &config
}
