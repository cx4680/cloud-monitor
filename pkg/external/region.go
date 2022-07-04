package external

type RegionService struct {
}

type RegionResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (service *RegionService) GetRegionList(tenantId string) ([]string, error) {
	//todo
	//respStr, err := httputil.HttpGet(url + param.Params)
	//if err != nil {
	//	return nil, err
	//}
	//var resp BmsResponse
	//jsonutil.ToObject(respStr, &resp)
	return nil, nil
}
