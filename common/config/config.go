package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

type Config struct {
	App        string     `yaml:"app"`
	Serve      Serve      `yaml:"serve"`
	DB         DB         `yaml:"db"`
	Logger     LogConfig  `yaml:"logger"`
	HttpConfig HttpConfig `yaml:"http"`
	Nk         string     `yaml:"nk"`
	Rocketmq   Rocketmq   `yaml:"rocketmq"`
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

type HttpConfig struct {
	ConnectionTimeOut int `yaml:"connection_time_out"`
	ReadTimeOut       int `yaml:"read_time_out"`
	WriteTimeOut      int `yaml:"write_time_out"`
}

type Rocketmq struct {
	NameServer        string `yaml:"name-server"`
	AlertContactTopic string `yaml:"alertContactTopic"`
	AlertContactGroup string `yaml:"alertContactGroup"`
}

var config Config = defaultAuthSdkConfig()

func defaultAuthSdkConfig() Config {
	return Config{
		HttpConfig: HttpConfig{
			ConnectionTimeOut: 3,
			ReadTimeOut:       3,
			WriteTimeOut:      3,
		},
		Logger: LogConfig{
			DataLogPrefix: "../logs/",
			Group:         "cloud-monitor-region",
		}}
}

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
