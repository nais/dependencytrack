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
- `apiServerEnv.*` and `frontendEnv.*` are rendered into ConfigMaps and loaded
  through `extraEnvFrom`; the upstream chart injects `app.frontend.apiBaseUrl`
  directly, so it remains the single source of truth for the frontend API URL.
- `bootstrap.*` retains the repository bootstrap Job and runs it with
  `DEPENDENCYTRACK_V5=true`. v5 replaced the v4 config properties for
  vulnerability sources and analyzers with extension configs that are only
  writable through the v2 API. The Job stores tokens as managed secrets and
  configures the NVD, GitHub and OSV data sources and the Trivy and OSS Index
  analyzers. It enables Trivy OS package scanning, which v4 enabled by default
  and v5 does not. Validate the Job against the exact v5 release before the
  production cutover.
- File storage defaults to S3 on Google Cloud Storage (`storage.googleapis.com`)
  with HMAC keys, so the API server can run more than one replica. Local
  `ReadWriteOnce` storage remains available as a fallback with one replica.
- Metrics are always enabled by the official v5 chart and are exposed on its
  management port. This wrapper installs its ServiceMonitor in the release
  namespace, so it is not a Feature-configurable setting.

## Validation

Run `helm dependency update charts-v5`, then render with a non-secret values
file containing all values required by `runtime-config.yaml` and `bootstrap.yaml`.

For the development rehearsal and production cutover procedure, see
[`RUNBOOK.md`](RUNBOOK.md).
