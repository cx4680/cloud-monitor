package actuator

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"os"
	"strings"
)

type envObject struct {
	Configuration     map[string]interface{} `json:"configuration"`
	SystemEnvironment map[string]string      `json:"systemEnvironment"`
}

func Env() string {
	env := new(envObject)
	env.SystemEnvironment = make(map[string]string)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		env.SystemEnvironment[pair[0]] = pair[1]
	}

	return jsonutil.ToString(env)
}
