# Uffizzi Controller Helm Chart

This chart installs the Kubernetes Controller for [Uffizzi](https://uffizzi.com), the continuous previews application. This is just a standard open-source Uffizzi setup.

## Configuration

### Dependencies

This chart depends upon two subcharts:

- [`ingress-nginx`](https://kubernetes.github.io/ingress-nginx/)
- [`cert-manager`](https://cert-manager.io/docs/)

### Custom Resource Definitions

By default, this Helm chart will tell `cert-manager` to create `CustomResourceDefinitions` within your Cluster. These may require special care if you have your own `cert-manager` installation somewhere else. You can disble this behavior by setting the Helm value `cert-manager.installCRDs = false`. See https://helm.sh/docs/chart_best_practices/custom_resource_definitions/

### `ClusterIssuer`: ZeroSSL

Right now this chart is configured to use the free certificate service ZeroSSL. To use this you must obtain a key from them: https://zerossl.com/documentation/acme/

Then configure the key and its ID within the `eab` values for this chart.

### Secrets

The following secrets are configurable:

- `global.uffizzi.controller.username`
- `global.uffizzi.controller.password`

- `eab.hmacKey`
- `eab.keyId`

### Environment Variables

The controller pod also has the following environment variables as Helm values:

- `podCidr`
- `env`
- `sandbox`

## More Info

See this project's main repository here: https://github.com/UffizziCloud/uffizzi_controller

And explore Uffizzi https://uffizzi.com
