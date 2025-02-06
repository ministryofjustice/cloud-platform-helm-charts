# cloud-platform-github-teams-filter

## WIP

This service exists to help mitigate issues with SAML authentication to AWS Console for CP user readonly access.

The attribute used to permit users to view their tagged resources consists of the set of team names a user's github account is associated with. There is a hard limit of 256 characters in this attribute, so users in many teams / long team names find authentication breaks if this limit is exceeded.

## Solution

The teams filter service works by: 

- listening for requests containing a `:` separated list of teams ie `:team1:team2:team3`
- queries the CP cluster's rolebinding objects and generated a deduplicated set of all teams that are "registered" across all namespaces
- removes from the input list any teams which are not present in any cluster rolebindings
- returns the filtered string

The service is to be called within our auth0 AWS SSO action, and requires an api key.

## Source code

https://github.com/ministryofjustice/cloud-platform-github-teams-filter


