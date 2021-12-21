{{- define "alertManagerConfig" -}}
{
    "apiVersion": "monitoring.coreos.com/v1alpha1",
    "kind": "AlertmanagerConfig",
    "metadata":{
      "name": {{- .Name }},
      "namespace": "product-cec-hawkeye"
    },
    "spec":{
      "route":{
        "receiver": "webhook-inner-op-{{ .Name }}",
        "routes":[
        {{ range .Router}}
          {
            "receiver": "webhook-alert-{{ $.Name }}",
            "matchers":[
              {
                "name": "app",
                "value": "hawkeye"
              }
              {{ range $key, $value := .Matchers }}
              ,
              {
                  "name": "{{$key}}",
                  "value": "{{$value}}"
              }
              {{end}}
            ],
            "continue": false,
            "repeatInterval": "{{.RepeatInterval}}"
          },
        {{end}}
        ],
        "repeatInterval": "3h"
      },
      "receivers":[
        {
          "name": "webhook-alert-{{ .Name }}",
          "webhookConfigs":[
            {
              "sendResolved": true,
              "url": "http://region-web-provider-svc.product-cec-hawkeye/inner/alertRecord/insert",
              "maxAlerts": 1000
            }
          ]
        },
        {
          "name":  "webhook-inner-op-{{ .Name }}",
          "webhookConfigs":[
            {
              "sendResolved": true,
              "url": "http://op.other.com/alert",
              "maxAlerts": 1000
            }
          ]
        }
      ]
    }
}

{{ end }}
