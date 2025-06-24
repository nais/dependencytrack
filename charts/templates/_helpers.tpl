{{/*
Expand the name of the chart.
*/}}
{{- define "dependencytrack.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}


{{- define "dependencytrack-backend.name" -}}
{{- .Chart.Name }}-backend
{{- end }}

{{- define "dependencytrack-frontend.name" -}}
{{- .Chart.Name }}-frontend
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this ().
If release name contains chart name it will be used as a full name.
*/}}
{{- define "dependencytrack.fullname" -}}
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
{{- define "dependencytrack.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "dependencytrack-backend.labels" -}}
helm.sh/chart: {{ include "dependencytrack.chart" . }}
{{ include "dependencytrack-backend.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "dependencytrack-frontend.labels" -}}
helm.sh/chart: {{ include "dependencytrack.chart" . }}
{{ include "dependencytrack-frontend.selectorLabels" . }}
{{- if .Chart.Version }}
app.kubernetes.io/version: {{ .Chart.Version | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels for backend
*/}}
{{- define "dependencytrack-backend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "dependencytrack.name" . }}
app.kubernetes.io/instance: {{ include "dependencytrack-backend.name" . }}
{{- end }}

{{/*
Selector labels for frontend
*/}}
{{- define "dependencytrack-frontend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "dependencytrack.name" . }}
app.kubernetes.io/instance: {{ include "dependencytrack-frontend.name" . }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "dependencytrack.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "dependencytrack.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}