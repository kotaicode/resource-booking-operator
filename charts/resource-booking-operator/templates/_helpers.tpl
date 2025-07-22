{{/*
Expand the name of the chart.
*/}}
{{- define "resource-booking-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "resource-booking-operator.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "resource-booking-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "resource-booking-operator.labels" -}}
helm.sh/chart: {{ include "resource-booking-operator.chart" . }}
{{ include "resource-booking-operator.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "resource-booking-operator.selectorLabels" -}}
app.kubernetes.io/name: {{ include "resource-booking-operator.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "resource-booking-operator.serviceAccountName" -}}
{{- if .Values.rbac.serviceAccount.create }}
{{- default (include "resource-booking-operator.fullname" .) .Values.rbac.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.rbac.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the namespace name
*/}}
{{- define "resource-booking-operator.namespace" -}}
{{- .Values.namespace.name | default "system" }}
{{- end }}

{{/*
Create the operator image
*/}}
{{- define "resource-booking-operator.image" -}}
{{- $registry := .Values.global.imageRegistry -}}
{{- $repository := .Values.operator.image.repository -}}
{{- $tag := .Values.operator.image.tag | default .Chart.AppVersion -}}
{{- if $registry -}}
{{- printf "%s/%s:%s" $registry $repository $tag -}}
{{- else -}}
{{- printf "%s:%s" $repository $tag -}}
{{- end -}}
{{- end }}

{{/*
Create the auth proxy image
*/}}
{{- define "resource-booking-operator.authProxyImage" -}}
{{- $registry := .Values.global.imageRegistry -}}
{{- $repository := .Values.authProxy.image.repository -}}
{{- $tag := .Values.authProxy.image.tag -}}
{{- if $registry -}}
{{- printf "%s/%s:%s" $registry $repository $tag -}}
{{- else -}}
{{- printf "%s:%s" $repository $tag -}}
{{- end -}}
{{- end }}

{{/*
Create the operator command args
*/}}
{{- define "resource-booking-operator.args" -}}
{{- $args := list "--leader-elect" -}}
{{- if .Values.authProxy.enabled -}}
{{- $args = concat $args (list "--health-probe-bind-address=:8081" "--metrics-bind-address=127.0.0.1:8080") -}}
{{- else -}}
{{- $args = concat $args (list "--health-probe-bind-address=:8081" "--metrics-bind-address=:8080") -}}
{{- end -}}
{{- join " " $args -}}
{{- end }} 