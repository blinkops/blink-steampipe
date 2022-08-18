# Blink - Steampipe

---
Blink uses the following
[Dockerfile](https://github.com/blinkops/blink-steampipe/Dockerfile)
to generate a wrapped version of the Steampipe CLI.

Our wrapped version of the CLI comes with the following plugins pre-installed:

- [AWS](https://github.com/turbot/steampipe-plugin-aws)
- [Github](https://github.com/turbot/steampipe-plugin-github)
- [Azure](https://github.com/turbot/steampipe-plugin-azure)
- [GCP](https://github.com/turbot/steampipe-plugin-gcp)
- [Kubernetes](https://github.com/turbot/steampipe-plugin-kubernetes)

And starts with a credential generator on entrypoint, That converts environment variables into credentials for AWS, GCP
& Kubernetes.

# Credential Generators

---

## AWS

---
Passing the following environment variables:

`AWS_ACCESS_KEY_ID` & `AWS_SECRET_ACCESS_KEY` - Will use steampipe as with the provided credentials

`AWS_ROLE_ARN` & `AWS_EXTERNAL_ID` - Will try to assume role with the following information and continue using the
generated credentials.

---

## GCP

---
Passing the following environment variable:

`GOOGLE_CREDENTIALS` - Allows you to provide your `JSON` credentials as environment variable.
Extracting `project_id` and continue working with the provided credentials

---

## Kubernetes

---
Passing the following environment variables:

`KUBERNETES_API_URL` & `KUBERNETES_BEARER_TOKEN` - Allows you to provide your kubernetes cluster api url and bearer token
to create kubernetes configuration for steampipe.

---

- [Steampipe](https://github.com/turbot/steampipe/)