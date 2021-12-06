package cbr

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type QueryParam struct {
	TenantId   string
	VaultId    string
	VaultName  string
	Status     string
	PageNumber string
	PageSize   string
}

func PageList(form *QueryParam) (*QueryPageVO, error) {
	var ecsInnerGateway = config.GetCommonConfig().EcsInnerGateway
	path := ecsInnerGateway + "/noauth/backup/vault/pageList"
	jsonStr, _ := json.Marshal(form)
	resp, err := http.Post(path, "application/json", strings.NewReader(string(jsonStr)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	var cbrQueryPageVO QueryPageVO
	err = json.Unmarshal(result, &cbrQueryPageVO)
	if err != nil {
		return nil, err
	}
	return &cbrQueryPageVO, nil
}

type QueryPageVO struct {
	Code        string
	Message     string
	Total_count int
	Data        []InfoBean
}

type InfoBean struct {
	VaultId      string
	TenantId     string
	Name         string
	Type         string
	Status       string
	Region       string
	Zone         string
	Capacity     string
	UsedCapacity string
	CreatedAt    string
	UpdatedAt    string
}
