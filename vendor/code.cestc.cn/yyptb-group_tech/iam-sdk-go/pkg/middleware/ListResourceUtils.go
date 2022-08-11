package middleware

import (
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/handler/authhttp"
)

// GetContextListResources 获取资源列表
func GetContextListResources(listResource *authhttp.ListResource) map[string]map[string][]string {

	if listResource == nil {
		return nil
	}

	if listResource.Resources == nil {
		return nil
	}

	return listResource.Resources
}

// GetContextListResourceOfProduct 根据产品获取资源ID列表（通过资源类型分组）
func GetContextListResourceOfProduct(listResource *authhttp.ListResource, product string) map[string][]string {
	resources := GetContextListResources(listResource)

	if resources == nil {
		return nil
	}

	return resources[product]
}

// GetContextListResourceOfResourceType 根据产品和资源类型获取资源ID列表
func GetContextListResourceOfResourceType(listResource *authhttp.ListResource, product string, resourceType string) []string {

	productResources := GetContextListResourceOfProduct(listResource, product)

	if productResources == nil {
		return nil
	}

	return productResources[resourceType]
}
