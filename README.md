# Dependencytrack

[Dependencytrack](https://dependencytrack.org/) is a tool for scanning dependencies for vulnerabilities.

## Overview

NAIS maintains this repository, which includes:

- 📦 Helm charts for Kubernetes deployments
- 🔧 Pre-install and post-install bootstrap configuration
- 🗄️ StatefulSet configuration for persistence and zero-downtime deployments

## Table of Contents

- [Dependencytrack Client](#dependencytrack-client)
- [Local Development](#local-development)
- [Image Verification](#image-verification)
- [License](#license)

## Dependencytrack Client

This repository includes a Go client for Dependencytrack, making it easy to integrate with your implementations. The client is designed to be extensible—feel free to expand the interface with additional functionality as needed.

### Using the Client

To use the client in your projects, import it as follows:

```bash
go get -u github.com/nais/dependencytrack@HEAD
```

## Local Development

### Prerequisites

Install the required tools:

```bash
asdf install
```

### Getting Started

1. **Start Dependencytrack locally:**

   ```bash
   make compose
   ```

2. **Configure environment variables:**

   - Copy `.env.sample` to `.env`
   - Fill in the necessary configuration values

3. **Access the services:**

   - **Dependencytrack UI:** http://localhost:9000
   - **API:** http://localhost:9001
   - **Swagger UI:** http://localhost:9002

### Test BOM schema validation locally

The BOM schema validation workaround can be configured in Kubernetes and locally.

#### Helm / Feature configuration

Use the Helm value:

```yaml
bootstrap:
  bomValidation:
    disabled: true
```

This value is also exposed in `charts/Feature.yaml` as:

```yaml
bootstrap.bomValidation.disabled
```

Semantics:

- Unset: leave the Dependency-Track setting unchanged
- `true`: set `bom.validation.mode=DISABLED`
- `false`: set `bom.validation.mode=ENABLED`

#### Docker Compose

To test the CycloneDX 1.7 workaround locally with Docker Compose, set the bootstrap environment variable:

```yaml
# - BOM_VALIDATION_DISABLED=true
```

by uncommenting it in the `bootstrap` service in `docker-compose.yaml`, or by adding it manually, then start the stack again:

```bash
make compose
```

Notes:

- Omit `BOM_VALIDATION_DISABLED` to leave the Dependency-Track setting unchanged
- Set `BOM_VALIDATION_DISABLED=true` to set `bom.validation.mode=DISABLED`
- Set `BOM_VALIDATION_DISABLED=false` to set `bom.validation.mode=ENABLED`

### Test Users

The `users.yaml` file contains pre-configured users for automated testing. You can modify this file to add or update test users as needed.

## Image Verification

The container images are signed using Sigstore cosign for authenticity verification.

### Verify Image Signature

```bash
cosign verify \
  --certificate-identity "https://github.com/nais/depedencytrack/.github/workflows/main.yaml@refs/heads/main" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  europe-north1-docker.pkg.dev/nais-io/nais/images/dependencytrack@sha256:<shasum>
```

### Verify Image Attestation (SBOM)

Images are attested with CycloneDX SBOMs. To verify the attestation:

```bash
cosign verify-attestation --type cyclonedx \
  --certificate-identity "https://github.com/nais/depedencytrack/.github/workflows/main.yaml@refs/heads/main" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  europe-north1-docker.pkg.dev/nais-io/nais/images/dependencytrack@sha256:<shasum>
```

## License

nais/Dependencytrack is licensed under the MIT License. See [LICENSE.md](/LICENSE.md) for details.
