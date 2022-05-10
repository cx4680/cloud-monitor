package model

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
)

type AlarmItem struct {
	Id               int        `gorm:"id" json:"id"`
	RuleBizId        string     `gorm:"rule_biz_id" json:"rule_biz_id"`
	MetricCode       string     `gorm:"metric_code" json:"metric_code"`
	TriggerCondition *Condition `gorm:"trigger_condition" json:"trigger_condition"` // 条件表达式
	Level            uint8      `gorm:"level" json:"level"`
	SilencesTime     string     `gorm:"silences_time" json:"silences_time"` // 告警间隔
}

func (*AlarmItem) TableName() string {
	return "t_alarm_item"
}

type Condition struct {
	MetricName         string  `json:"metricName"`
	MetricCode         string  `json:"metricCode"`
	Period             int     `json:"period"`
	Times              int     `json:"times"`
	Statistics         string  `json:"statistics"`
	ComparisonOperator string  `json:"comparisonOperator"`
	Threshold          float64 `json:"threshold"`
	Unit               string  `json:"unit"`
	Labels             string  `json:"labels"`
}

func (c *Condition) Value() (driver.Value, error) {
	bs, err := json.Marshal(c)
	return string(bs), errors.WithStack(err)
}
func (c *Condition) Scan(v interface{}) error {
	var err error
	switch vt := v.(type) {
	case string:
		err = json.Unmarshal([]byte(vt), &c)
	case []byte:
		err = json.Unmarshal(vt, &c)
	default:
		return errors.New("rule condition 转换错误")
	}
	return err
}
