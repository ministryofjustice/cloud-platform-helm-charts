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
2) Artifacts inside *gh-pages* branch. The artifacts used by helm are tarballs (`tgz`) and they are built using `helm package` command

If we want to include a new helm chart in this repo we should add its source code into the master branch followed by the artifacts into *gh-pages* branch. 

Last step is to regenerate the index (`index.yaml`), which is the file indexed when we do `helm repo update`.

## Regenerating the `index.yaml`

Regenerate `index.yaml`

```console
helm repo index --url  https://ministryofjustice.github.io/cloud-platform-helm-charts/ --merge index.yaml .
```

## More documentation

- [Official Helm Documentation](https://helm.sh/docs/)
