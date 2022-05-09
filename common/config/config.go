package config

import (
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

type AppConfig struct {
	App        string       `yaml:"app"`
	Serve      Serve        `yaml:"serve"`
	Db         DB           `yaml:"db"`
	Logger     LogConfig    `yaml:"logger"`
	HttpConfig HttpConfig   `yaml:"http"`
	Rocketmq   Rocketmq     `yaml:"rocketmq"`
	Prometheus Prometheus   `yaml:"prometheus"`
	Common     CommonConfig `yaml:"common"`
	Redis      RedisConfig  `yaml:"redis"`
	Iam        IamConfig    `yaml:"iam"`
}

type CommonConfig struct {
	Env                    string `yaml:"env"`
	EnvType                string `yaml:"envType"`
	Nk                     string `yaml:"nk"`
	TenantUrl              string `yaml:"tenantUrl"`
	SmsCenterPath          string `yaml:"smsCenterPath"`
	ActivateInformationUrl string `yaml:"activateInformationUrl"`
	RegionName             string `yaml:"regionName"`
	EcsInnerGateway        string `yaml:"ecs-inner-gateway"`
	IsSingleRegion         bool   `yaml:"isSingleRegion"`
	MsgUrl                 string `yaml:"msgUrl"`
}

type Serve struct {
	Debug        bool `yaml:"debug"`
	Port         int  `yaml:"port"`
	ReadTimeout  int  `yaml:"read_timeout"`
	WriteTimeout int  `yaml:"write_timeout"`
}

type DB struct {
	Dialect       string        `yaml:"dialect"`
	Username      string        `yaml:"username"`
	Password      string        `yaml:"password"`
	Url           string        `yaml:"url"`
	MaxIdleConnes int           `yaml:"max_idle_connes"`
	MaxOpenConnes int           `yaml:"max_open_connes"`
	MaxLifeTime   time.Duration `yaml:"time.Hour"`
}

type LogConfig struct {
	Debug         bool   `yaml:"debug"`
	DataLogPrefix string `yaml:"data_log_prefix"`
	ServiceName   string `yaml:"service_name"`
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
	Site   string `yaml:"site"`
	Region string `yaml:"region"`
	Log    string `yaml:"log"`
}

var Cfg = defaultAppConfig()

func defaultAppConfig() AppConfig {
	return AppConfig{
		HttpConfig: HttpConfig{
			ConnectionTimeOut: 3,
			ReadTimeOut:       3,
			WriteTimeOut:      3,
		},
		Logger: LogConfig{
			DataLogPrefix: "./logs/",
			ServiceName:   "cloud-monitor-region",
		},
		Rocketmq: Rocketmq{
			NameServer: "127.0.0.1:9876",
		},
		Common: CommonConfig{
			Env:                    "local",
			Nk:                     "",
			TenantUrl:              "",
			SmsCenterPath:          "",
			ActivateInformationUrl: "",
			RegionName:             uuid.New().String(),
			EcsInnerGateway:        "",
		},
	}
}

func InitConfig(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if e := yaml.Unmarshal(data, &Cfg); e != nil {
		return e
	}
	return nil
}
