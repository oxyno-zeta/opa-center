package partitions

import "text/template"

const opaCfgTemplateString = `
services:
  opacenter-{{ .Partition.Name }}:
    url: {{ .ServiceURL }}

decision_logs:
  service: opacenter-{{ .Partition.Name }}
  partition_name: {{ .Partition.ID }}
  reporting:
    min_delay_seconds: 30
    max_delay_seconds: 60

status:
  service: opacenter-{{ .Partition.Name }}
  partition_name: {{ .Partition.ID }}

`

func loadOpaCfgTemplate() (*template.Template, error) {
	return template.New("service-tpl").Parse(opaCfgTemplateString)
}
