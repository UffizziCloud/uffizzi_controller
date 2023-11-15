# Uffizzi Controller Helm Chart

This chart installs the Kubernetes Controller for [Uffizzi](https://uffizzi.com), the continuous previews application. This is just a standard open-source Uffizzi setup.

### Dependencies

This chart depends upon two subcharts:

- [`ingress-nginx`](https://kubernetes.github.io/ingress-nginx/)
- [`cert-manager`](https://cert-manager.io/docs/)

### Custom Resource Definitions

By default, this Helm chart will tell `cert-manager` to create `CustomResourceDefinitions` within your Cluster. These may require special care if you have your own `cert-manager` installation somewhere else. You can disble this behavior by setting the Helm value `cert-manager.installCRDs = false`. See https://helm.sh/docs/chart_best_practices/custom_resource_definitions/

## Configuration

### `ClusterIssuer`: Let's Encrypt

Right now this chart is configured to use the free certificate service Let's Encrypt. To use this you must specify the value `cert-email`.

### `ClusterIssuer`: ZeroSSL

You may also use the free certificate service ZeroSSL. To use this you must obtain a key from them: https://zerossl.com/documentation/acme/

Then configure the key and its ID within the `zerossl.eab` values for this chart. You'll also want to specify the value `clusterIssuer` as `zerossl` instead of `letsencrypt`.

### Controller Ingress

This chart configures a Kubernetes `Ingress` for the Controller. For it to successfully obtain a certificate, specify its DNS hostname as the value `ingress.hostname`.

### Secrets

The following secrets are configurable:

- `global.uffizzi.controller.username`
- `global.uffizzi.controller.password`

- `zerossl.eab.hmacKey`
- `zerossl.eab.keyId`

### Environment Variables

The controller pod also has the following environment variables as Helm values:

- `env` - (Doesn't do much yet.)
- `sandbox` - Enable `nodeSelector` and `taint` options for gVisor on GKE.

## Installation

If this is your first time using Helm, consult their documentation: https://helm.sh/docs/intro/quickstart/

Begin by adding our Helm repository:

```
helm repo add uffizzi-controller https://uffizzicloud.github.io/uffizzi_controller/
```

Then install the lastest version as a new release using the values you specified earlier. We recommend isolating Uffizzi in its own Namespace.

```
helm install uc uffizzi-controller/uffizzi-controller --values myvals.yaml --namespace uffizzi-controller --create-namespace
```

If you encounter any errors here, tell us about them in [our Slack](https://join.slack.com/t/uffizzi/shared_invite/zt-ffr4o3x0-J~0yVT6qgFV~wmGm19Ux9A)!

You should then see the release is installed:
```
helm list --namespace uffizzi-controller
```

### Troubleshooting

When installing this chart, you may see errors like this:
```
clusterroles.rbac.authorization.k8s.io "my-uffizzi-controller-flux-default-source-controller-helmchart" already exists
```

This happens when more than one resource within a dependency chart (in this case `flux`) has a very long name truncated into the same name as another resource. To avoid this, use shorter release names as in the example above.

## More Info

See this project's main repository here: https://github.com/UffizziCloud/uffizzi_controller

And explore Uffizzi https://uffizzi.com
