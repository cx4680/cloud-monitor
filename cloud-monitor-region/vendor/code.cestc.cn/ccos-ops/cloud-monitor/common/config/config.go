package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

type Config struct {
	App        string       `yaml:"app"`
	Serve      Serve        `yaml:"serve"`
	DB         DB           `yaml:"db"`
	Logger     LogConfig    `yaml:"logger"`
	HttpConfig HttpConfig   `yaml:"http"`
	Rocketmq   Rocketmq     `yaml:"rocketmq"`
	Prometheus Prometheus   `yaml:"prometheus"`
	Common     CommonConfig `yaml:"common"`
	Redis      RedisConfig  `yaml:"redis"`
	Iam        IamConfig    `yaml:"iam"`
}

type CommonConfig struct {
	Env               string `yaml:"env"`
	Nk                string `yaml:"nk"`
	TenantUrl         string `yaml:"tenantUrl"`
	SmsCenterPath     string `yaml:"smsCenterPath"`
	HawkeyeCenterPath string `yaml:"hawkeyeCenterPath"`
	HasNoticeModel    bool   `yaml:"hasNoticeModel"`
	RegionName        string `yaml:"regionName"`
	EcsInnerGateway   string `yaml:"ecs-inner-gateway"`
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
	NameServer string `yaml:"name-server"`
}

type Prometheus struct {
	Url        string `yaml:"url"`
	Query      string `yaml:"query"`
	QueryRange string `yaml:"queryRange"`
}

type RedisConfig struct {
	Addr     string
	Password string
}

type IamConfig struct {
	Site   string
	Region string
	Log    string
}

var cfg = defaultAuthSdkConfig()

func defaultAuthSdkConfig() Config {
	return Config{
		HttpConfig: HttpConfig{
			ConnectionTimeOut: 3,
			ReadTimeOut:       3,
			WriteTimeOut:      3,
		},
		Logger: LogConfig{
			DataLogPrefix: "../logs/",
			//TODO group
			Group: "cloud-monitor-region",
		},
		Rocketmq: Rocketmq{
			NameServer: "127.0.0.1:9876",
		},
		Common: CommonConfig{
			Env:               "local",
			Nk:                "",
			TenantUrl:         "",
			SmsCenterPath:     "",
			HawkeyeCenterPath: "",
			HasNoticeModel:    false,
			RegionName:        "local",
			EcsInnerGateway:   "",
		},
	}
}

func InitConfig(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}
	//TODO Read from envs

	return nil
}

func GetCommonConfig() CommonConfig {
	return cfg.Common
}
func GetServeConfig() Serve {
	return cfg.Serve
}
func GetDbConfig() DB {
	return cfg.DB
}
func GetLogConfig() LogConfig {
	return cfg.Logger
}

func GetRedisConfig() RedisConfig {
	return cfg.Redis
}

//TODO 写在tools中
func GetHttpConfig() HttpConfig {
	return cfg.HttpConfig
}

func GetRocketmqConfig() Rocketmq {
	return cfg.Rocketmq
}

func GetPrometheusConfig() Prometheus {
	return cfg.Prometheus
}

func GetIamConfig() IamConfig {
	return cfg.Iam
}
