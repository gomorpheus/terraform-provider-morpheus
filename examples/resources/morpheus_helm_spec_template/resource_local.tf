resource "morpheus_helm_spec_template" "tfexample_helm_spec_template_local" {
  name         = "tf-helm-spec-example-local"
  source_type  = "local"
  spec_content = <<TFEOF
apiVersion: v1
kind: Service
metadata:
name: {{ template "fullname" . }}
labels:
    chart: "{{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}"
spec:
type: {{ .Values.service.type }}
ports:
- port: {{ .Values.service.externalPort }}
    targetPort: {{ .Values.service.internalPort }}
    protocol: TCP
    name: {{ .Values.service.name }}
selector:
    app: {{ template "fullname" . }}
TFEOF
}