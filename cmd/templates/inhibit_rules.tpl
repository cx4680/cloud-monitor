{{- define "inhibitRulesConfig" -}}
{
    "apiVersion": "monitoring.coreos.com/v1alpha1",
    "kind": "AlertmanagerConfig",
    "metadata":{
      "name": {{ .Name }},
      "namespace": {{.Namespace}}
    },
    "spec":{
        "inhibitRules": [
            {{- range .Rules }}
            {
            sourceMatch: [
               {{- range $v:= .SourceMatchers}}
                {
                    "name": "{{.Name}}",
                    "value": "{{.Value}}",
                    "regex": {{.Regex}}
                },
               {{- end}}
               ],
              targetMatch: [
              {{- range $v:= .TargetMatchers}}
                 {
                     "name": "{{.Name}}",
                     "value": "{{.Value}}",
                     "regex": {{.Regex}}
                 },
              {{- end}}
              ],
              equal: [
               {{- range $v:= .Equal}}
               "{{ $v }}",
               {{- end}}
             ]
            },
            {{- end}}
        ]
    }
}
{{ end }}