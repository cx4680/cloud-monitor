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
	Ecs        Ecs          `yaml:"ecs"`
	Common     CommonConfig `yaml:"common"`
}

type CommonConfig struct {
	Nk                string `yaml:"nk"`
	TenantUrl         string `yaml:"tenantUrl"`
	SmsCenterPath     string `yaml:"smsCenterPath"`
	HawkeyeCenterPath string `yaml:"hawkeyeCenterPath"`
	HasNoticeModel    bool   `yaml:"hasNoticeModel"`
	RegionName        string `yaml:"regionName"`
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
	BrokerAddr        string `yaml:"broker-addr"`
	AlertContactTopic string `yaml:"alertContactTopic"`
	AlertContactGroup string `yaml:"alertContactGroup"`
	RuleTopic         string `yaml:"ruleTopic"`
	RecordTopic       string `yaml:"recordTopic""`
}

type Prometheus struct {
	Url        string `yaml:"url"`
	Query      string `yaml:"query"`
	QueryRange string `yaml:"queryRange"`
}

type Ecs struct {
	InnerGateway string `yaml:"inner-gateway"`
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

func InitConfig(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	return nil
}

func GetCommonConfig() CommonConfig {
	return config.Common
}
func GetServeConfig() Serve {
	return config.Serve
}
func GetDbConfig() DB {
	return config.DB
}
func GetLogConfig() LogConfig {
	return config.Logger
}

//TODO 写在tools中
func GetHttpConfig() HttpConfig {
	return config.HttpConfig
}

func GetRocketmqConfig() Rocketmq {
	return config.Rocketmq
}

func GetPrometheusConfig() Prometheus {
	return config.Prometheus
}

func GetEcsConfig() Ecs {
	return config.Ecs
}
