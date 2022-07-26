package form

type PrometheusRequest struct {
	TenantId   string `form:"tenantId"`
	Name       string `form:"name"`
	Instance   string `form:"instance"`
	TopNum     int    `form:"topNum"`
	Time       string `form:"time"`       //瞬时数据查询参数 时间戳
	Start      int    `form:"start"`      //开始时间 时间戳
	End        int    `form:"end"`        //结束时间 时间戳
	Step       int    `form:"step"`       //步长 时间戳
	Statistics string `form:"statistics"` //聚合函数 sum(求和)  min(最小值)  max (最大值)  avg (平均值)  stddev (标准差)  stdvar (标准差异)  count (计数)
	Scope      string `form:"scope"`      //聚合范围 1s(1秒)  1m(1分钟)  1h(1小时)  1d(1天)  1w(1周)  1y(1年)
	Pid        string `form:"pid"`
}

type PrometheusResponse struct {
	Status string          `json:"status"`
	Data   *PrometheusData `json:"data"`
}

type PrometheusData struct {
	ResultType string              `json:"resultType"`
	Result     []*PrometheusResult `json:"result"`
}

type PrometheusResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
	Values [][]interface{}   `json:"values"`
}

type PrometheusValue struct {
	Time  string `json:"time"`
	Value string `json:"value"`
}

type PrometheusAxis struct {
	XAxis []string            `json:"xaxis"`
	YAxis map[string][]string `json:"yaxis"`
}

type PrometheusInstance struct {
	Instance string `json:"instance"`
	Value    string `json:"value"`
}

type NetworkData struct {
	TimeAxis  []string `json:"timeAxis"`
	ValueAxis []string `json:"ValueAxis"`
}

type ProcessData struct {
	Pid       string `json:"pid"`
	Name      string `json:"name"`
	CmdLine   string `json:"cmdLine"`
	Cpu       string `json:"cpu"`
	Memory    string `json:"memory"`
	Openfiles string `json:"openfiles"`
	Time      string `json:"time"`
}
