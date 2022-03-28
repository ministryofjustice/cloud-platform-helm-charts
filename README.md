# Cloud Platform helm charts repository

Cloud Platform [Helm](https://helm.sh) repository to store the internal helm charts. 

The main advantage is directly related to maintainability: we don’t need to keep an eye ay service, just in source code.

## How to use this repository 

Adding a helm repository is as easy as:

```console
$ helm repo add cloud-platform https://ministryofjustice.github.io/cloud-platform-helm-charts
```

## CRUD for Charts

Keep in mind two things:

1) Source code of helm charts are inside the *master* branch, 
2) Artifacts inside *gh-pages* branch which are the ones publicly distributed and available in url [https://ministryofjustice.github.io/cloud-platform-helm-charts](https://ministryofjustice.github.io/cloud-platform-helm-charts)

If you want to include a new helm chart of amend existing helm chart to create new artifacts(with latest chart version) and publish the artifacts follow the below steps

- Build the artifacts tarballs(`tgz`) using the below command
`helm package <chart-name>`

- Regenerate `index.yaml`

```console
helm repo index --url  https://ministryofjustice.github.io/cloud-platform-helm-charts/ --merge index.yaml .
```

- Raise a PR with all the changes

Once the PR is merged to *main* branch, the github action `sync-branches.yml`  will push the artifacts into *gh-pages* branch. The regenerated index (`index.html`) can then used to upgrade the chart by doing `helm repo update`

## More documentation

- [Official Helm Documentation](https://helm.sh/docs/)
