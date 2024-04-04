# dependencytrack

[Dependencytrack](https://dependencytrack.org/) is a tool for scanning dependencies for vulnerabilities.

NAIS oversees the maintenance of this repository, which includes:

* helm charts
* pre-installed & post-install bootstrap configuration
* Stateful set for persistence and no down-time deployments

## Dependencytrack client

The repository houses a client for dependencytrack, facilitating its usage for implementations. Please feel free to
expand the client interface with additional functionality as required. To update to the most recent version of the
client, execute:

`go get -u github.com/nais/dependencytrack@HEAD`

## Local development

### Prerequisites

First `asdf install`

Run your local instance of dependencytrack with the following command:
`make compose`

A docker-compose file is available for local development. Please duplicate .env.sample as .env and provide the necessary
values. The users.yaml file allows for the creation of pre-installed users for automated testing purposes. You can
access the [dependencytrack UI](http://localhost:9000).

The API is available at [depdendencytrack-api](http://localhost:9001).

### Swagger

The [swagger UI](http://localhost:9002) is bundled with the docker-compose file.

## Verifying the image and its contents

The image is signed "keylessly" (is that a word?) using Sigstore cosign. To verify its authenticity run

```
cosign verify \
--certificate-identity "https://github.com/nais/depedencytrack/.github/workflows/main.yaml@refs/heads/main" \
--certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
europe-north1-docker.pkg.dev/nais-io/nais/images/dependencytrack@sha256:<shasum>
```

The images are also attested with SBOMs in the CycloneDX format. You can verify these by running

```
cosign verify-attestation --type cyclonedx \
--certificate-identity "https://github.com/nais/depedencytrack/.github/workflows/main.yaml@refs/heads/main" \
--certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
europe-north1-docker.pkg.dev/nais-io/nais/images/dependencytrack@sha256:<shasum>
```

## License

nais/Dependencytrack is licensed under the MIT License, see [LICENSE.md](/LICENSE.md).
