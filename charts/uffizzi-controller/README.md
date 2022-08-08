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

## More Info

See this project's main repository here: https://github.com/UffizziCloud/uffizzi_controller

And explore Uffizzi https://uffizzi.com
