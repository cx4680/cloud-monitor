package templates

import (
	"bytes"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/k8s"
	"log"
	"testing"
	"text/template"
)

func Test_tpl(t *testing.T) {
	templates, err := template.ParseFiles("alert_manager_config.tpl")
	if err != nil {
		log.Println(err)
		return
	}
	var buf bytes.Buffer
	param := k8s.AlertManagerConfig{
		Name: "jim-config",
		Router: []k8s.Router{{
			Matchers:       map[string]string{"app": "hawkeye", "ruleId": "123"},
			RepeatInterval: "5m",
		}},
	}
	err = templates.ExecuteTemplate(&buf, "alertManagerConfig", param)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("ret:" + buf.String())
	log.Println("--------------")

}

func Test_tpl_expression(t *testing.T) {
	//s := "{{Calc .OSTYPE}}"
	s := "{{eq .OSTYPE \"windows\"  }}"
	m := map[string]string{"OSTYPE": "windows"}
	var buf bytes.Buffer
	temp, _ := template.New("exp").Parse(s)

	if err := temp.Execute(&buf, m); err != nil {
		log.Println(err)
	}

	log.Println(buf.String())
	//log.Println(strconv.ParseBool(buf.String()))
}
