# Kratos

[![Go Report Card](https://goreportcard.com/badge/github.com/laghoule/kratos)](https://goreportcard.com/report/github.com/laghoule/kratos)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=laghoule_kratos&metric=coverage)](https://sonarcloud.io/summary/new_code?id=laghoule_kratos)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=laghoule_kratos&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=laghoule_kratos)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=laghoule_kratos&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=laghoule_kratos)

## Yep another deployment tools

Kratos is simple, but with simplicity come less flexibility, so if you want a full fledge deploying tools, this is propably not for you. But, if you have some simple container, maybe a nginx with your html5 website, this may be the perfect alternative to custom Kubernetes YAML, or the build of helm templates.

### Use case

I had a little html5 demo container that I wanted to host on my Kubernetes cluster. To deploy on my cluster I have 2 options at hand (ok, I haven't search for others options beside these two):

* Helm template packaging
* Kubernetes yaml files

These solutions are not difficult (if you are familiar with Kubernetes), but time consuming (and yes boring).

So Kratos is born from this use case.

## Status

Under heavy developpment, **use at your own risk**.

Contribution are welcome.

## Prerequisite

* Kubernetes 1.19+
* Golang 1.18+ for building
* Certmanager for TLS certificates
* A working kubeconfig configuration

## Cmdline

```text
Alternative to helm for deploying simple container, without the pain of managing Kubernetes yaml templates.

Usage:
  kratos [command]

Available Commands:
  create      Deploy an application.
  delete      Delete an application.
  get         Retreive a configuration of deployed application.
  help        Help about any command.
  init        Create an empty configuration file.
  list        List applications.
  update      Update an application.
  version     Show version of kratos.

Flags:
  -h, --help                help for kratos
  -k, --kubeconfig string   kubernetes configuration file (default "/home/user/.kube/config")

Use "kratos [command] --help" for more information about a command.
```

### Security restriction

Kratos don't support deployment of container running as `root`. If you deploy a container running as `root`, it will fail to start, with this error:

```text
Error: container has runAsNonRoot and image will run as root (pod: "pacman-7dc78bcb9c-hwjhp_static(11b8ab59-0a7f-45ca-87fd-3c3348e9fc7f)", container: pacman)
```

Kratos don't mount the Kubernetes service account `token` in the containers. These `token` are useful only to application who need access to Kubernetes API.

podSpec configurations:

* RunAsNonRoot is `true`

* AutomountServiceAccountToken is `false`

### Initialize a configuration file

You need to create a configuration file, you can build it from scratch, or use the `init` command:

```bash
kratos init --name myappconfig.yaml
```

You can now just `fill` the configuration with your own config values. It's very similar to helm values.yaml file.

### Deploy your application

With your configuration ready, you can now deploy to Kubernetes:

```bash
kratos create --name myapp --namespace mynamespace --config myappconfig.yaml
```

## Config

### Values definition

| Values | Descriptions | Mandatory |
|--------|-------------|---------|
| common.labels| Labels common to all Kubernetes objects | no |
| common.annotations | Annotation common to all Kubernetes objects | no |
| deployment.labels | Deployment & pod labels | no |
| deployment.annotations | Deployment & pod annotations | no |
| deployment.replicas | Numbers of pod replicas | yes |
| deployment.port | Port to use for communication with pod | yes |
| deployment.containers | List of containers in the pods | yes |
| deployment.containers.name | Name of the containers | yes |
| deployment.containers.image | Name of the Docker image | yes |
| deployment.containers.tag | Tag version of the image | yes |
| deployment.containers.resources.requests.cpu | Request this amount of CPU | no |
| deployment.containers.resources.requests.memory | Request this amount of RAM | no |
| deployment.containers.resources.limits.cpu | Max amount of CPU | no |
| deployment.containers.resources.limites.memory | Max amount of RAM | no |
| deployment.containers.health | Healthcheck configuration for the container | no |
| deployment.containers.health.live | Liveness check | no |
| deployment.containers.health.live.probe | URI to use for the check | yes |
| deployment.containers.health.live.port | Port of the container | yes |
| deployment.containers.health.live.initialDelaySeconds | Delay the check for x seconds at startup | no |
| deployment.containers.health.live.periodSeconds | Time between check in second | no |
| deployment.containers.health.ready | Readyness check | no |
| deployment.containers.health.ready.probe | URI to use for the check | yes |
| deployment.containers.health.ready.port |  Port of the container | yes |
| deployment.containers.health.ready.initialDelaySeconds | Delay the check for x seconds at startup | no |
| deployment.containers.health.ready.periodSeconds | Time between check in second| no |
| deployment.ingress.labels | Ingress labels | no |
| deployment.ingress.annotations | Ingress annotations | no |
| deployment.ingress.ingressClass | Name of the ingressClass to use | yes |
| deployment.ingress.clusterIssuer | Name of the clusterIssuer to use | yes |
| deployment.ingress.hostnames | List of hostnames associate with this deployment | yes |
| cronjobs.labels | Cronjobs labels | no |
| cronjobs.annotations | Cronjobs annotation | no |
| cronjobs.schedule | Cronjobs schedule definition | yes |
| cronjobs.retry | Number of retry if jobs fail | no |
| cronjobs.container | container definition | yes |
| cronjobs.container.name | Name of the containers | yes |
| cronjobs.container.image | Name of the Docker image | yes |
| cronjobs.container.tag | Tag version of the image | yes |
| cronjobs.container.resources.requests.cpu | Request this amount of CPU | no |
| cronjobs.container.resources.requests.memory | Request this amount of RAM | no |
| cronjobs.container.resources.limits.cpu | Max amount of CPU | no |
| cronjobs.container.resources.limites.memory | Max amount of RAM | no |
| configmaps | List of configmaps | no |
| configmaps.labels | Configmaps labels | no |
| configmaps.annotations | Configmaps annotations | no |
| configmaps.files | List of configmaps files | yes |
| configmaps.files.name | Name of the configmaps | yes |
| configmaps.files.mount | How to use the configmaps | yes |
| configmaps.files.mount.path | Path of the mount point in the pod | yes |
| configmaps.files.mount.exposedTo | List of containers to expose the configmap | no |
| configmaps.data | Contents of the configmap | yes |
| secrets | List of secrets | no |
| secrets.labels | Secrets labels | no |
| secrets.annotation | Secrets annotations | no |
| secrets.files | List of secrets files | yes |
| secrets.files.name | Name of the secret | yes |
| secrets.files.mount | How to use the secret | yes |
| secrets.files.mount.path | Path of the mount point in the pod | yes |
| secrets.files.mount.exposedTo | List of containers to expose the secret | no |
| secrets.files.data | Contents of the secret | yes |

### Example of a full features configuration

```yaml
common:
  labels: {}
  annotations: {}

deployment:
  labels: {}
  annotations: {}
  replicas: 1
  port: 80
  containers:
    - name: pacman
      image: laghoule/patate-poil
      tag: v1.0.1
      resources:
        requests:
          cpu: 25m
          memory: 32Mi
        limits:
          cpu: 50m
          memory: 64Mi
      health:
        live:
          probe: /isLive
          initialDelaySeconds: 3
          periodSeconds: 3
        ready:
          probe: /isReady
          initialDelaySeconds: 3
          periodSeconds: 3
  ingress:
    labels: {}
    annotations: {}
    ingressClass: nginx
    clusterIssuer: letsencrypt
    hostnames:
      - example.com
      - www.example.com

cronjob:
  labels: {}
  annotations: {}
  schedule: 0 0 * * *
  retry: 1
  container:
    name: myjobs
    image: laghoule/crunchdata
    tag: v1.0.0
    resources:
      requests:
        cpu: 25m
        memory: 32Mi
      limits:
        cpu: 50m
        memory: 64Mi

configmaps:
  labels: {}
  annotations: {}
  files:
    - name: configuration.yaml
      mount:
        path: /etc/cfg
        exposedTo:
          - pacman
          - myjobs
      data: |
        my configuration data

secrets:
  labels: {}
  annotations: {}
  files:
    - name: credentials.yaml
      mount:
        path: /etc/cfg
        exposedTo:
          - myjobs
      data: |
        usename: patate
        password: poil
```

### Example of a minimal configuration

#### Deployment

```yaml
deployment:
  replicas: 1
  port: 80
  containers:
    - name: pacman
      image: laghoule/patate-poil
      tag: v1.0.1
  ingress:
    ingressClass: nginx
    clusterIssuer: letsencrypt
    hostnames:
      - example.com
      - www.example.com
```

#### Cronjobs

```yaml
cronjob:
  schedule: 0 0 * * *
  container:
    name: myjobs
    image: laghoule/crunchdata
    tag: v1.0.0
```
