package config

var config Config

type Config struct {
	AuthSdkConfig *AuthSdkConfig `json:"authSdkConfig"`
	LogConfig     *LogConfig     `json:"logConfig"`
}
type AuthSdkConfig struct {
	AuthRequestSite   string `json:"authRequestSite"`
	ResourceName      string `json:"ResourceName"`
	ConnectionTimeOut int    `json:"ConnectionTimeOut"`
	ReadTimeOut       int    `json:"ReadTimeOut"`
	WriteTimeOut      int    `json:"WriteTimeOut"`
	RegionId          string `json:"RegionId"`
}

type LogConfig struct {
	Debug      bool   `json:"debug"`
	Directory  string `json:"directory"`
	MaxSize    int    `json:"maxSize"`
	MaxBackups int    `json:"maxBackups"`
	MaxAge     int    `json:"maxAge"`
	Compress   bool   `json:"compress"`
	Stdout     bool   `json:"stdout"`
}

func InitConfig(identityUrl, regionId, logDir string, v ...string) {
	config = Config{
		AuthSdkConfig: &AuthSdkConfig{
			AuthRequestSite:   identityUrl,
			ResourceName:      "cecloud",
			RegionId:          regionId,
			ConnectionTimeOut: 10,
			ReadTimeOut:       10,
			WriteTimeOut:      5,
		},
		LogConfig: &LogConfig{
			Debug:     false,
			Directory: logDir,
			Stdout:    true,
		},
	}
}

func GetConfig() *Config {
	return &config
}
