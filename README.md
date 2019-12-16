# Cloud Platform helm charts repository (ministryofjustice.github.io/cloud-platform-helm-charts)

[Helm](https://helm.sh) repo for different internal charts which are installed on [Kubernetes](https://kubernetes.io)

Add the repo via:
```console
$ helm repo add cp-ministryofjustice https://ministryofjustice.github.io/cloud-platform-helm-charts
```
:q
## Regenerate `index.yaml`

```console
helm repo index --url  https://ministryofjustice.github.io/cloud-platform-helm-charts/ --merge index.yaml .
```
