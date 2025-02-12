# Cloud Platform Github Teams Filter

This service exists to help mitigate issues with SAML authentication to AWS Console for Cloud Platform user readonly access.

## Problem

The SAML attribute used to permit users to view their tagged resources consists of the set of team names a user's github account is associated with. There is a hard limit of 256 characters in this attribute, so users in many teams / long team names find authentication breaks if this limit is exceeded. See [this runbook entry](https://runbooks.cloud-platform.service.justice.gov.uk/debugging-aws-console-access.html#debugging-aws-console-read-only-access-issues) for more details on this issue.

## Solution

In order to alleviate this issue, the teams filter service works by: 

- listening for requests containing a `:` separated list of teams ie `:team1:team2:team3`
- queries the CP cluster's rolebinding objects and generates a deduplicated set of all teams that are "registered" across all namespaces
- removes from the input list any teams which are not present in any cluster rolebindings
- returns the filtered string

In doing so, we are removing any github teams which have no relevance for our read-only console service.

The service is to be called within our auth0 AWS SSO action, and requires an api key. This can be found in the Kubernetes Secret 'github-teams-filter-secret`.

## Source code

https://github.com/ministryofjustice/cloud-platform-github-teams-filter


## Usage

If you want to check the service in action, try something like:

```
curl https://filter-teams.apps.cloud-platform.service.justice.gov.uk/filter-teams \                                                
     -H "X-API-Key: {api-key value}" -H "Content-Type: application/json" \
     -d '{"teams": ":badteam:webops:test1:test2:dps-tech"}'
```




