package form

type AlertRecordGroupLabels struct {
	Alertname string
}

type AlertRecordCommonLabels struct {
	Alertname string
	Severity  string
	Team      string
}

type AlertRecordCommonAnnotations struct {
	Summary string
}

type AlertBeanLabelsBean struct {
	Alertname                           string
	Beta_kubernetes_io_arch             string
	Beta_kubernetes_io_fluentd_ds_ready string
	Beta_kubernetes_io_os               string
	Instance                            string
	Job                                 string
	Kubernetes_io_arch                  string
	Kubernetes_io_hostname              string
	Kubernetes_io_os                    string
	Severity                            string
	Team                                string
}

type AlertBeanAnnotationsBean struct {
	Description string
	Summary     string
}

type AlertRecordAlertsBean struct {
	Status       string
	Labels       *AlertBeanLabelsBean
	Annotations  *AlertBeanAnnotationsBean
	StartsAt     string
	EndsAt       string
	GeneratorURL string
}

type InnerAlertRecordAddForm struct {
	Receiver          string
	Status            string
	GroupLabels       *AlertRecordGroupLabels
	CommonLabels      *AlertRecordCommonLabels
	commonAnnotations *AlertRecordCommonAnnotations
	ExternalURL       string
	Version           string
	GroupKey          string
	Alerts            []*AlertRecordAlertsBean
}
