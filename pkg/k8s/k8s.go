package k8s

import (
	"bytes"
	c "code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/inhibit"
	"context"
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
	"text/template"
)

const (
	group     = "monitoring.coreos.com"
	version   = "v1"
	namespace = "product-cec-hawkeye"
	plural    = "prometheusrules"
)

var client dynamic.Interface
var resource *schema.GroupVersionResource
var alertManagerResource *schema.GroupVersionResource

func InitK8s() error {
	var config *rest.Config
	if strings.EqualFold(c.Cfg.Common.Env, "local") {
		cfg, err := clientcmd.BuildConfigFromFlags("", "k8s-config-local.yml")
		if err != nil {
			return err
		}
		config = cfg
	} else {
		cfg, err := rest.InClusterConfig()
		if err != nil {
			return err
		}
		config = cfg
	}

	clientSet, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	client = clientSet
	resource = &schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: plural,
	}
	alertManagerResource = &schema.GroupVersionResource{
		Group:    group,
		Version:  "v1alpha1",
		Resource: "alertmanagerconfigs",
	}
	return nil
}

func DeleteAlertRule(alertRuleId string) error {
	err := client.Resource(*resource).Namespace(namespace).Delete(context.TODO(), alertRuleId, metav1.DeleteOptions{})
	if err != nil {
		logger.Logger().Error("调用api删除规则失败, name=", alertRuleId, err)
		return errors.NewBusinessErrorCode(errors.DeleteError, "调用api删除规则失败")
	}
	return nil
}

func alertRuleToObject(alertRuleDTO *AlertRuleDTO) *map[string]interface{} {
	result := map[string]interface{}{}
	result["apiVersion"] = "monitoring.coreos.com/v1"
	result["kind"] = "PrometheusRule"
	labels := map[string]string{}
	labels["prometheus"] = "k8s"
	labels["role"] = "alert-rules"
	labels["namespace"] = namespace
	metadata := map[string]interface{}{}
	metadata["labels"] = &labels
	metadata["name"] = alertRuleDTO.TenantId
	metadata["namespace"] = namespace
	result["metadata"] = &metadata
	spec := map[string]interface{}{}
	var array = make([]*map[string]interface{}, len(alertRuleDTO.SpecGroupsList))
	for index, specGroup := range alertRuleDTO.SpecGroupsList {
		group := map[string]interface{}{}
		group["name"] = specGroup.Name
		var ruleList = make([]*map[string]interface{}, len(specGroup.AlertList))
		for index, alert := range specGroup.AlertList {
			rule := map[string]interface{}{}
			rule["alert"] = alert.Alert
			rule["expr"] = alert.Expr
			rule["for"] = alert.ForTime
			rule["labels"] = alert.Labels
			annotations := map[string]interface{}{}
			annotations["summary"] = alert.Summary
			annotations["description"] = alert.Description
			rule["annotations"] = &annotations
			ruleList[index] = &rule
		}
		group["rules"] = &ruleList
		array[index] = &group
	}
	spec["groups"] = &array
	result["spec"] = &spec
	return &result
}

func ApplyAlertRule(alertRuleDTO *AlertRuleDTO) error {
	rules := alertRuleToObject(alertRuleDTO)
	requestObj, err2 := json.Marshal(rules)
	if err2 != nil {
		logger.Logger().Errorf("调用api更新规则失败%v", err2)
		return errors.NewBusinessErrorCode(errors.ApplyFail, "调用api更新规则失败")
	}
	_, err := client.Resource(*resource).Namespace(namespace).Patch(context.TODO(), alertRuleDTO.TenantId, types.ApplyPatchType, requestObj, metav1.ApplyOptions{FieldManager: "application/apply-patch", Force: true}.ToPatchOptions())
	if err != nil {
		return err
	}
	return nil
}

func ApplyAlertManagerConfig(cfg AlertManagerConfig) error {
	var buf bytes.Buffer
	var err error
	logger.Logger().Infof("alertmanagercfg : %+v", cfg)
	templates, err := template.ParseFiles("templates/alert_manager_config.tpl")
	if err != nil {
		logger.Logger().Errorf(err.Error())
		return err
	}
	err = templates.ExecuteTemplate(&buf, "alertManagerConfig", cfg)
	if err != nil {
		logger.Logger().Errorf(err.Error())
		return err
	}
	_, err = client.Resource(*alertManagerResource).
		Namespace(namespace).
		Patch(context.TODO(), cfg.Name, types.ApplyPatchType,
			[]byte(buf.String()), metav1.ApplyOptions{FieldManager: "application/apply-patch", Force: true}.ToPatchOptions())
	if err != nil {
		logger.Logger().Errorf(err.Error())
		return err
	}
	return nil
}

func DeleteAlertManagerConfig(configName string) error {
	err := client.Resource(*alertManagerResource).Namespace(namespace).Delete(context.TODO(), configName, metav1.DeleteOptions{})
	if err != nil {
		logger.Logger().Warn("调用api删除AlertManagerConfig失败, name=", configName, err)
		return errors.NewBusinessErrorCode(errors.DeleteError, "调用api删除AlertManagerConfig失败")
	}
	return nil
}

const LevelInhibitName = "hawkeye-level-inhibit"

func ApplyInhibitRules(levels []uint8) error {
	rules := inhibit.BuildRules(levels)
	tpl, err := template.ParseFiles("templates/inhibit_rules.tpl")
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = tpl.ExecuteTemplate(&buf, "inhibitRulesConfig", struct {
		Name      string
		Namespace string
		Rules     []inhibit.InhibitRule
	}{
		Name:      LevelInhibitName,
		Namespace: namespace,
		Rules:     rules,
	})
	if err != nil {
		logger.Logger().Errorf(err.Error())
		return err
	}
	_, err = client.Resource(*alertManagerResource).
		Namespace(namespace).
		Patch(context.TODO(), LevelInhibitName, types.ApplyPatchType,
			[]byte(buf.String()), metav1.ApplyOptions{FieldManager: "application/apply-patch", Force: true}.ToPatchOptions())
	if err != nil {
		logger.Logger().Errorf(err.Error())
		return err
	}
	return nil
}
