# Cloud Platform Label Pods

The service is a mutating admission controller webhook. Given the [hard char limit of 63](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set) for pod labels, we are forced to use annotations to pass users' github team data up to fluent-bit. The location of the "label" can be found on the pod at `.metadata.annotations.github_teams`. A strict hard limit does not exist on annotations, instead annotations values are limited on size, in this case [256kb](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/api/validation/objectmeta.go#L44-L67) is more than enough.

This service will be deployed to the `cloud-platfom-label-pods` namespace.

## Purpose

> Label all user pods with their github team name.

## Source code

[https://github.com/ministryofjustice/cloud-platform-terraform-label-pods](https://github.com/ministryofjustice/cloud-platform-terraform-label-pods)
