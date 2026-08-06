# Dependency-Track v5 chart

This chart wraps the official `dependency-track` Helm chart v2.x from
`https://dependencytrack.github.io/helm-charts` and retains NAIS-specific Cloud
SQL proxy and bootstrap configuration.

It is intentionally separate from `../charts`: Dependency-Track v4 cannot be
upgraded in-place to v5. Migrate the database while v4 is offline, uninstall
the v4 release, then install this chart.

## Design

- `app.*` values are forwarded to the official v5 chart.
- `cloudSqlProxy.connectionName` configures the Cloud SQL Auth Proxy sidecar.
- `apiServerEnv.*` and `frontendEnv.*` are rendered into a ConfigMap and loaded
  through `extraEnvFrom`; this keeps Fasit-facing values readable.
- `bootstrap.*` retains the repository bootstrap Job. Validate its API calls
  against the exact Dependency-Track v5 release before production cutover.
- Local file storage uses `ReadWriteOnce`; keep `app.apiServer.web.replicaCount`
  at `1` unless you switch to S3 or an RWX-capable StorageClass.

## Validation

Run `helm dependency update charts-v5`, then render with a non-secret values
file containing all values required by `runtime-config.yaml` and `bootstrap.yaml`.
