package k8s

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"context"
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
)

const (
	group     = "monitoring.coreos.com"
	version   = "v1"
	namespace = "product-cec-hawkeye"
	plural    = "prometheusrules"
)

var client dynamic.Interface
var resource *schema.GroupVersionResource

func InitK8s() {
	var config *rest.Config
	if strings.EqualFold(os.Getenv("active"), "local") {
		cfg, err := clientcmd.BuildConfigFromFlags("", "D:\\dev-go\\cloud-monitor\\cloud-monitor-region\\k8s\\config.yml")
		if err != nil {
			panic(err.Error())
		}
		config = cfg
	} else {
		cfg, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		config = cfg
	}

	clientSet, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	client = clientSet
	resource = &schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: plural,
	}
}

func CreateAlertRule(alertRuleDTO *forms.AlertRuleDTO) (*forms.AlertRuleDTO, error) {
	requestObj := alertRuleToObject(alertRuleDTO)
	rules := &unstructured.Unstructured{
		Object: requestObj,
	}
	create, err := client.Resource(*resource).Namespace(namespace).Create(context.TODO(), rules, metav1.CreateOptions{})
	if err != nil {
		logger.Logger().Errorf("调用api创建规格失败 %+v", err)
		return nil, errors.NewBussinessError(2, "调用api创建规格失败")
	}
	alertRuleDTO.AlertRuleId = alertRuleDTO.TenantId
	fmt.Sprintf("%v", create)
	return alertRuleDTO, nil
}

func UpdateAlertRule(alertRuleDTO *forms.AlertRuleDTO) error {
	rules := alertRuleToObject(alertRuleDTO)
	requestObj, err2 := json.Marshal(rules)
	if err2 != nil {
		logger.Logger().Errorf("调用api更新规则失败%v", err2)
		return errors.NewBussinessError(3, "调用api更新规则失败")
	}
	_, err := client.Resource(*resource).Namespace(namespace).Patch(context.TODO(), alertRuleDTO.AlertRuleId, types.MergePatchType, requestObj, metav1.PatchOptions{})
	if err != nil {
		logger.Logger().Errorf("调用api更新规则失败%v", err)
		return errors.NewBussinessError(3, "调用api更新规则失败")
	}
	return nil
}

func DeleteAlertRule(alertRuleId string) error {
	err := client.Resource(*resource).Namespace(namespace).Delete(context.TODO(), alertRuleId, metav1.DeleteOptions{})
	if err != nil {
		logger.Logger().Errorf("调用api删除规则失败")
		return errors.NewBussinessError(3, "调用api删除规则失败")
	}
	return nil
}

func alertRuleToObject(alertRuleDTO *forms.AlertRuleDTO) map[string]interface{} {
	result := map[string]interface{}{}
	result["apiVersion"] = "monitoring.coreos.com/v1"
	result["kind"] = "PrometheusRule"
	labels := map[string]string{}
	labels["prometheus"] = "k8s"
	labels["role"] = "alert-rules"
	labels["namespace"] = namespace
	metadata := map[string]interface{}{}
	metadata["labels"] = labels
	metadata["name"] = alertRuleDTO.TenantId
	metadata["namespace"] = namespace
	result["metadata"] = metadata
	spec := map[string]interface{}{}
	var array = make([]map[string]interface{}, len(alertRuleDTO.SpecGroupsList))
	for index, specGroup := range alertRuleDTO.SpecGroupsList {
		group := map[string]interface{}{}
		group["name"] = specGroup.Name
		var ruleList = make([]map[string]interface{}, len(specGroup.AlertList))
		for index, alert := range specGroup.AlertList {
			rule := map[string]interface{}{}
			rule["alert"] = alert.Alert
			rule["expr"] = alert.Expr
			rule["for"] = alert.ForTime
			rule["labels"] = alert.Labels
			annotations := map[string]interface{}{}
			annotations["summary"] = alert.Summary
			annotations["description"] = alert.Description
			rule["annotations"] = annotations
			ruleList[index] = rule
		}
		group["rules"] = ruleList
		array[index] = group
	}
	spec["groups"] = array
	result["spec"] = spec
	return result
}
