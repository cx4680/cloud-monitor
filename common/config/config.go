package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

type Config struct {
	app        string       `yaml:"app"`
	serve      Serve        `yaml:"serve"`
	db         DB           `yaml:"db"`
	logger     LogConfig    `yaml:"logger"`
	httpConfig HttpConfig   `yaml:"http"`
	rocketmq   Rocketmq     `yaml:"rocketmq"`
	prometheus Prometheus   `yaml:"prometheus"`
	common     CommonConfig `yaml:"common"`
	redis      RedisConfig  `yaml:"redis"`
	iam        IamConfig    `yaml:"iam"`
}

type CommonConfig struct {
	Env                   string `yaml:"env"`
	Nk                    string `yaml:"nk"`
	TenantUrl             string `yaml:"tenantUrl"`
	SmsCenterPath         string `yaml:"smsCenterPath"`
	CertifyInformationUrl string `yaml:"certifyInformationUrl"`
	HawkeyeCenterPath     string `yaml:"hawkeyeCenterPath"`
	MsgIsOpen             string `yaml:"msgIsOpen"`
	MsgChannel            string `yaml:"msgChannel"`
	RegionName            string `yaml:"regionName"`
	EcsInnerGateway       string `yaml:"ecs-inner-gateway"`
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
		httpConfig: HttpConfig{
			ConnectionTimeOut: 3,
			ReadTimeOut:       3,
			WriteTimeOut:      3,
		},
		logger: LogConfig{
			DataLogPrefix: "../logs/",
			//TODO group
			Group: "cloud-monitor-region",
		},
		rocketmq: Rocketmq{
			NameServer: "127.0.0.1:9876",
		},
		common: CommonConfig{
			Env:                   "local",
			Nk:                    "",
			TenantUrl:             "",
			SmsCenterPath:         "",
			CertifyInformationUrl: "",
			HawkeyeCenterPath:     "",
			MsgIsOpen:             MsgOpen,
			MsgChannel:            MsgChannelEmail,
			RegionName:            "local",
			EcsInnerGateway:       "",
		},
	}
}

func InitConfig(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &cfg)
}

func GetCommonConfig() CommonConfig {
	return cfg.common
}
func GetServeConfig() Serve {
	return cfg.serve
}
func GetDbConfig() DB {
	return cfg.db
}
func GetLogConfig() LogConfig {
	return cfg.logger
}

func GetRedisConfig() RedisConfig {
	return cfg.redis
}

//TODO 写在tools中
func GetHttpConfig() HttpConfig {
	return cfg.httpConfig
}

func GetRocketmqConfig() Rocketmq {
	return cfg.rocketmq
}

func GetPrometheusConfig() Prometheus {
	return cfg.prometheus
}

func GetIamConfig() IamConfig {
	return cfg.iam
}
