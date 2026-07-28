{{- define "backend-api.fullname" -}}
{{ .Release.Name }}
{{- end -}}

{{- define "backend-api.labels" -}}
app.kubernetes.io/name: backend-api
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "backend-api.selectorLabels" -}}
app.kubernetes.io/name: backend-api
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
