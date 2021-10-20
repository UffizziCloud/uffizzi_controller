# Uffizzi Cloud Resource Controller
This application connects to a Kubernetes (k8s) Cluster to provision Uffizzi users' workloads on their behalf.
While it provides a documented REST API for anyone to use, it's most valuable when used with [uffizzi_app](https://github.com/UffizziCloud/uffizzi_app). Learn more at <https://uffizzi.com>

# Design
The Uffizzi Continuous Previews Engine empowers development teams to conduct feature-level pre-merge testing by automatically deploying branches of application repositories for full-stack and microservices applications based on user designated triggers.  Uffizzi makes these on-demand test environments available for review by key stakeholders (QA, Peer review, Product designer, Product manager, end users, etc.) at a secure Preview URL. The on-demand test environments provisioned by Uffizzi have a purpose driven life cycle and follow the Continous Previews methodology - <https://cpmanifesto.org> - <https://github.com/UffizziCloud/Continuous_Previews_Manifesto>

Uffizzi's implementation leverages several components as well as public cloud resources, including a Kubernetes Cluster. This controller is a supporting service for `uffizzi_app` and works in conjunction with redis and postgres to provide the CP capabilty. 

This controller runs within the Cluster and accepts authenticated instructions from other Uffizzi components.
It then specifies Resources within the Cluster's Kubernetes control API.

This controller acts as a smart and secure proxy for uffizzi_app and is designed to restrict required access to the k8s cluster.  It is implemented in Golang to leverage the best officially-supported Kubernetes API client.

The controller is required as a uffizzi_app supporting service and serves these purposes:
1. Communicate deployment instructions via native Golang API client to the designated Kubernetes cluster(s) from the Uffizzi interface
2. Provide Kubernetes cluster information back to the Uffizzi interface
3. Support restricted and secure connection between the Uffizzi interface and the Kubernetes cluster

## Example story: New Preview Deployment
- `main()` loop is within `cmd/controller/controller.go`, which calls `setup()` and handles exits. This initializes `global` settings and the `sentry` logging, connects to the database, initializes the Kubernetes clients, and starts the HTTP server listening.
- An HTTP request for a new Deployment arrives and is handled within `internal/http/handlers.go`.  The request contains the new Deployment integer ID.
- The HTTP handler uses the ID as an argument to call the `ApplyDeployment` function within `internal/domain/deployment.go`. This takes a series of steps:
  - It then calls several methods from `internal/kuber/client.go`, which creates Kubernetes specifications for each k8s resource (Namespace, Deployment, NetworkPolicy, Service, etc.) and publishes them to the Cluster one at a time.
    - This function should return an IP address or hostname, which is added to the `data` for this Deployment's `state`.
- Any errors are then handled and returned to the HTTP client.

# Dependencies
This controller specifies custom Resources managed by popular open-source controllers:
- [cert-manager](https://cert-manager.io/)
- [ingress-nginx](https://kubernetes.github.io/ingress-nginx/)

You'll want these installed within the Cluster managed by this controller.

# Configuration

## Environment Variables
You can specify these within `credentials/variables.env` for use with `docker-compose` and our `Makefile`.
Some of these may have defaults within `configs/settings.yml`.

- `ENV` - Which deployment environment we're currently running within.  Default: `development`
- `CONTROLLER_LOGIN` - The username to HTTP Basic Authentication
- `CONTROLLER_PASSWORD` - The password to HTTP Basic Authentication
- `CONTROLLER_NAMESPACE_NAME_PREFIX` - Prefix for Namespaces provisioned. Default: `deployment`
- `CERT_MANAGER_CLUSTER_ISSUER` - The issuer for signing certificates. Possible values:
    - `letsencrypt` (used by default)
    - `zerossl`
- `POD_CIDR` - IP range to allowlist within `NetworkPolicy`. Default: `10.24.0.0/14`
- `POOL_MACHINE_TOTAL_CPU_MILLICORES` - Node resource to divide for Pods. Default: 2000
- `POOL_MACHINE_TOTAL_MEMORY_BYTES` - Node recourse to divide for Pods. Default: 17179869184
- `DEFAULT_AUTOSCALING_CPU_THRESHOLD` - Default: 75
- `DEFAULT_AUTOSCALING_CPU_THRESHOLD_EPSILON` - Default: 8
- `AUTOSCALING_MAX_PERFORMANCE_REPLICAS` - Horizontal Pod Autoscaler configuration. Default: 10
- `AUTOSCALING_MIN_PERFORMANCE_REPLICAS` - Horizontal Pod Autoscaler configuration. Default: 1
- `AUTOSCALING_MAX_ENTERPRISE_REPLICAS` - Horizontal Pod Autoscaler configuration. Default: 30
- `AUTOSCALING_MIN_ENTERPRISE_REPLICAS` - Horizontal Pod Autoscaler configuration. Default: 3
- `STARTUP_PROBE_DELAY_SECONDS` - Startup Probe configuration. Default: 10
- `STARTUP_PROBE_FAILURE_THRESHOLD` - Startup Probe configuration. Default: 80
- `STARTUP_PROBE_PERIOD_SECONDS` - Startup Probe configuration. Default: 15
- `EPHEMERAL_STORAGE_COEFFICIENT` - `LimitRange` configuration. Default: 1.9

## Kubernetes API Server Connection
This process expects to be provided a Kubernetes Service Account
within a Kubernetes cluster. You can emulate this with these four
pieces of configuration:

- `KUBERNETES_SERVICE_HOST` - Hostname (or IP) of the k8s API service
- `KUBERNETES_SERVICE_PORT` - TCP port number of the k8s API service (usually `443`.)
- `/var/run/secrets/kubernetes.io/serviceaccount/token` - Authentication token
- `/var/run/secrets/kubernetes.io/serviceaccount/ca.crt` - k8s API Server's x509 host certificate

Once you're configured to connect to your cluster (using `kubectl` et al)
then you can get the value for these two environment variables from the output of
`kubectl cluster-info`.

Add those two environment variables to `credentials/variables.env`.

The authentication token must come from the cluster's cloud provider, e.g.
`gcloud config config-helper --format="value(credential.access_token)"`

The server certificate must also come from the cluster's cloud provider, e.g.
`gcloud container clusters describe uffizzi-pro-production-gke --zone us-central1-c --project uffizzi-pro-production-gke --format="value(masterAuth.clusterCaCertificate)" | base64 --decode`

You should write these two values to `credentials/token` and `credentials/ca.crt`
and the `make` commands and `docker-compose` will copy them for you.

# Shell
While developing, we most often run the controller within a shell on our workstations.
`docker-compose` will set up this shell and mount the current working directory within the container so you can use other editors from outside.
To login into docker container just run:
```shell script
make shell
```

All commands in this "Shell" section should be run inside this shell.

## Compile
After making any desired changes, compile the controller:
```shell script
go install ./cmd/controller/...
```

## Execute
```shell script
/go/bin/controller
```

## Test Connection to Cluster
Once you've configured access to your k8s Cluster (see above), you can test `kubectl` within the shell:
```shell script
kubectl --token=`cat /var/run/secrets/kubernetes.io/serviceaccount/token` --certificate-authority=/var/run/secrets/kubernetes.io/serviceaccount/ca.crt get nodes
```

## Tests, Linters

In docker shell:
```
make test
make lint
make fix_lint
```

# External Testing
Once the controller is running on your workstation, you can make HTTP requests to it from outside of the shell.

## Ping controller

```shell script
curl localhost:8080 \
  --user "${CONTROLLER_LOGIN}:${CONTROLLER_PASSWORD}"
```

## Remove all workload from existing environment
This will remove the specified Preview's Namespace and all other Resources.
```shell script
curl -X POST localhost:8080/clean \
     --user "${CONTROLLER_LOGIN}:${CONTROLLER_PASSWORD}" \
     -H "Content-Type: application/json" \
     -d '{ "environment_id": 1 }'
```

## Online API Documentation
Available at <http://localhost:8080/docs/>

# Installation within a Cluster
Functional usage within a Kubernetes Cluster is beyond the scope of this document. For more, [join us on Slack](https://uffizzi.slack.com/join/shared_invite/zt-ffr4o3x0-J~0yVT6qgFV~wmGm19Ux9A#/shared-invite/email) or contact us at <info@uffizzi.com>.

That said, we've included a Kubernetes manifest to help you get started at `infrastructure/controller.yaml`.
Review it and change relevant variables before applying this manifest.
You'll also need to install and configure the dependencies identified near the top of this document.
