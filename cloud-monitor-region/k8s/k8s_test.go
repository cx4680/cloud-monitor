package k8s

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestK8sClient(t *testing.T) {
	InitK8s()
	dto := &forms.AlertRuleDTO{
		AlertRuleId: "x",
		TenantId:    "1",
		Region:      "",
		Zone:        "",
	}
	str := "{  \"apiVersion\": \"monitoring.coreos.com/v1\",  \"kind\": \"PrometheusRule\",  \"metadata\": {    \"name\": \"myrule\",    \"labels\": {      \"namespace\": \"product-cec-hawkeye\",      \"prometheus\": \"k8s\",      \"role\": \"alert-rules\"    }  },  \"spec\": {    \"groups\": [      {        \"name\": \"my-node-time\",        \"rules\": [          {            \"alert\": \"myClockSkewDetected\",            \"annotations\": {              \"message\": \"myClock skew detected on node-exporter {{`{{ $labels.namespace }}`}}/{{`{{ $labels.pod }}`}}. Ensure NTP is configured correctly on this host.\"            },            \"expr\": \"abs(node_timex_offset_seconds{job=\\\"node-exporter\\\"}) > 0.03\",            \"for\": \"2m\",            \"labels\": {              \"severity\": \"warning\"            }          }        ]      }    ]  }}"
	bytes := []byte(str)
	json.Unmarshal(bytes, dto)
	rule, err := CreateAlertRule(dto)
	if err != nil {

	}
	fmt.Sprintf("%v", rule)
}

func TestUpdateAlertRule(t *testing.T) {
	InitK8s()
	dto := &forms.AlertRuleDTO{
		AlertRuleId: "x",
		TenantId:    "1",
		Region:      "",
		Zone:        "",
	}
	err := UpdateAlertRule(dto)
	if err != nil {

	}
}

func TestDeleteAlertRule(t *testing.T) {
	InitK8s()
	err := DeleteAlertRule("myrule")
	if err != nil {

	}
}

func TestStr(t *testing.T) {
	fold := strings.Compare("LoCal", "local")
	print(fold)
}
