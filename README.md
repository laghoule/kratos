# Kratos

## Yet another deployment tools

Kratos is simple, but with simplicity come less flexibility, so if you want a full fledge deploying tools, this is propably not for you. But, if you have some simple container, maybe a nginx with your html5 website, this may be the perfect alternative to custom Kubernetes YAML, or the build of helm templates.

### Use case

I was having a little html5 demo container that I wanted to host on my Kubernetes cluster. To deploy on the cluster I have 2 options at hand (ok, I haven't search for others options beside these two):

* Helm template packaging
* Kubernetes yaml files

These solutions are not difficult (if you are familiar with Kubernetes), but time consuming.

So Kratos is born from this use case.

## Prerequisite

* Kubernetes 1.19+ (use of Ingress V1)
* Certmanager for TLS certificates
* A working kubeconfig configuration

## Cmdline

```text
Alternative to helm for deploying simple container, without the pain of managing Kubernetes yaml templates.

Usage:
  kratos [command]

Available Commands:
  create      Deploy an application in an namespace.
  delete      Delete a deployment in a namespace.
  get         Retreive a configuration of a kratos deployment.
  help        Help about any command
  init        Create an empty configuration file.
  list        List application of managed kratos deployment.

Flags:
  -h, --help                help for kratos
  -k, --kubeconfig string   kubernetes configuration file (default "/home/user/.kube/config")

Use "kratos [command] --help" for more information about a command.
```

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
| deployment.labels | Deployment & Pod labels | no |
| deployment.replicas | Numbers of pod replicas | yes |
| deployment.port | Port to use for communication with pod | yes
| containers | List of containers in the pods | yes |
| containers.name | Name of the containers | yes |
| containers.image | Name of the Docker image | yes |
| containers.tag | Tag version of the image | yes |
| containers.resources.requests.cpu | Request this amount of CPU | no |
| containers.resources.requests.memory | Request this amount of RAM | no |
| containers.resources.limits.cpu | Max amount of CPU | no |
| containers.resources.limites.memory | Max amount of RAM | no |
| configmaps | List of configmaps | no |
| configmaps.labels | Configmaps labels | no |
| configmaps.annotations | Configmaps annotations | no |
| configmaps.name | Name of the configmap | yes |
| configmaps.mountPath | Path of the mount point in the pod | yes |
| configmaps.data | Contents of the configmap | yes |
| secrets | List of secrets | no |
| secrets.labels | Secrets labels | no |
| secrets.annotation | Secrets annotations | no |
| secrets.name | Name of the secret | yes |
| secrets.mountPath | Path of the mount point in the pod | yes |
| secrets.data | Contents of the secret | yes |
| ingress.labels | Ingress labels | no |
| ingress.annotations | Ingress annotations | no |
| ingress.ingressClass | Name of the ingressClass to use | yes |
| ingress.clusterIssuer | Name of the clusterIssuer to use | yes |
| ingress.hostnames | List of hostnames associate with this deployment | yes |


### Example of a full features configuration:

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

cronjobs:
  labels: {}
  annotations: {}
  shedule: 0 0 * * *
  retryOnError: 1
  containers:
    - name: myjobs
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
cronjobs:
  shedule: 0 0 * * *
  containers:
    - name: myjobs
      image: laghoule/crunchdata
      tag: v1.0.0
```