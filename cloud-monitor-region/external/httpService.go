package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/pkg/config"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	once       sync.Once
	httpClient *http.Client
)

func initHttpClient(config *config.HttpConfig) {
	connectTimeOut := time.Duration(config.ConnectionTimeOut) * time.Second
	rwTimeOut := time.Duration(config.ReadTimeOut+config.WriteTimeOut) * time.Second
	httpClient = &http.Client{Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:  connectTimeOut,
			Deadline: time.Now().Add(rwTimeOut),
		}).DialContext,
	},
	}
}

func HttpClient() *http.Client {
	once.Do(func() {
		initHttpClient(&config.GetConfig().HttpConfig)
	})
	return httpClient
}

func PageList(userCode string, request interface{}, url string) ([]byte, error) {
	requestObject, err := json.Marshal(request)
	if err != nil {
		logger.Logger().Errorf("auth  post object parse to json failed, err:%v\n", err)
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(requestObject)))
	req.Header.Add("userCode", userCode)
	resp, err := HttpClient().Do(req)
	logger.Logger().Debugf("http,request:%v", requestObject)
	logger.Logger().Debugf("http,response:%v", resp)
	if err != nil {
		logger.Logger().Errorf("post failed, err:%v\n", err)
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Logger().Errorf("get resp failed, err:%v\n", err)
		return nil, err
	}
	return b, nil
}
