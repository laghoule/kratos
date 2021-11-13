# Kratos

## Config

```yaml
deployment:
  labels: {}
  annotations: {}
  replicas: 1
  containers:
    - name: pacman
      image: laghoule/patate-poil
      tag: v1.0.1
      port: 80

configmaps:
  labels: {}
  annotations: {}
  - name: configuration.yaml
    mountPath: /etc/cfg
    data: |
      my configuration data

secrets:
  labels: {}
  annotations: {}
  - name: credentials.yaml
    mountPath: /etc/cfg
    data: |
      usename: patate
      password: poil

ingress:
  labels: {}
  annotations: {}
  ingressClass: nginx
  clusterIssuer: letsencrypt
  hostnames:
    - example.com
    - www.example.com
```

## Development map

### v1.0.0
