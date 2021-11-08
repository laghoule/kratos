# Kratos

## Config

```yaml
name: myApp
namespace: myNamespace
deployment:
  replicas: 1
service:
  port: 80
ingress:
  ingressClass: nginx
  clusterIssuer: letsencrypt
  hostnames:
    - example.com
    - www.example.com
  port: 80
```

## Development map

### v1.0.0
