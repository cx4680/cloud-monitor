package form

type AlarmRecordGroupLabels struct {
	Alertname string
}

type AlarmRecordCommonLabels struct {
	Alertname string
	Severity  string
	Team      string
}

type AlarmRecordCommonAnnotations struct {
	Summary string
}

type AlarmBeanLabelsBean struct {
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

type AlarmBeanAnnotationsBean struct {
	Description string
	Summary     string
}

type AlarmRecordAlertsBean struct {
	RequestId    string
	Status       string
	Labels       *AlarmBeanLabelsBean
	Annotations  *AlarmBeanAnnotationsBean
	StartsAt     string
	EndsAt       string
	GeneratorURL string
}

type InnerAlarmRecordAddForm struct {
	Receiver          string
	Status            string
	GroupLabels       *AlarmRecordGroupLabels
	CommonLabels      *AlarmRecordCommonLabels
	commonAnnotations *AlarmRecordCommonAnnotations
	ExternalURL       string
	Version           string
	GroupKey          string
	Alerts            []*AlarmRecordAlertsBean
}
