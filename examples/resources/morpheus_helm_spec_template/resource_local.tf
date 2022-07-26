resource "morpheus_helm_spec_template" "tfexample_helm_spec_template_local" {
  name         = "tf-helm-spec-example-local"
  source_type  = "local"
  spec_content = <<TFEOF
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
TFEOF
}