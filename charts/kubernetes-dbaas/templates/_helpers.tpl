{{/*
Expand the name of the chart.
*/}}
{{- define "kubernetes-dbaas.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "kubernetes-dbaas.fullname" -}}
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
{{- define "kubernetes-dbaas.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "kubernetes-dbaas.labels" -}}
helm.sh/chart: {{ include "kubernetes-dbaas.chart" . }}
{{ include "kubernetes-dbaas.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "kubernetes-dbaas.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kubernetes-dbaas.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
DatabaseClass generator
*/}}
{{- define "kubernetes-dbaas.dbcGenerator" -}}
{{- range .Values.dbc }}
apiVersion: databaseclass.dbaas.bedag.ch/v1
kind: DatabaseClass
metadata:
  name: {{ .name }}
  labels:
    {{- include "kubernetes-dbaas.labels" $ | nindent 4 }}
spec:
  driver: {{ .driver }}
  operations:
    {{- toYaml .operations | nindent 4 }}
  secretFormat:
    {{- toYaml .secretFormat | nindent 4 }}
---
{{- end }}
{{- end }}

{{/*
DBMS endpoint Secrets generator
*/}}
{{- define "kubernetes-dbaas.dbmsSecretsGenerator" -}}
{{- range .Values.dbmsSecrets }}
apiVersion: v1
kind: Secret
metadata:
  name: {{ .name }}
  labels:
    {{- include "kubernetes-dbaas.labels" $ | nindent 4 }}
type: Opaque
stringData:
  {{- toYaml .stringData | nindent 2 }}
---
{{- end }}
{{- end }}