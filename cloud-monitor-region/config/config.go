package config

type K8sConfig struct {
	Group     string
	Version   string
	Namespace string
	Plural    string
	ApiUrl    string
}

var k8sConfig = &K8sConfig{
	Group:     "monitoring.coreos.com",
	Version:   "v1",
	Namespace: "product-cec-hawkeye",
	Plural:    "prometheusrules",
}

func initK8s() {
	/*config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}*/

}
