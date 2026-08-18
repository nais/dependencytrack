# Dependency-Track v5 migration runbook

Use this runbook to rehearse and then cut over a Dependency-Track v4
deployment to this v5 chart. Read the upstream
[v4-to-v5 migration guide](https://dependencytrack.github.io/docs/next/guides/administration/migrating-from-v4/)
before scheduling a production cutover.

## Rules

- v5 is not an in-place upgrade. It needs a new, dedicated Postgres target.
- Do not install the v5 chart against the target database before the migration
  has completed and been verified. The API server seeds data that the migrator
  must populate from v4.
- The v4 API server must be offline during a production migration. For a
  rehearsal, use a restored v4 database copy whenever possible.
- Pin exact v5 image tags or digests. Do not use `latest` or `snapshot`.
- Do not pass passwords, API keys, or the KEK on a command line or commit them
  to a values file.

## 1. Prepare the target

1. Confirm the v4 source is Dependency-Track `4.14.2` or newer.
2. Provision a new Postgres 14+ database for v5. Do not reuse the v4 database.
3. Grant the v5 database user ownership and DDL privileges for the target. Use
   the same user for the migrator and the later v5 API server.
4. Ensure the target has about twice the final v4 dataset size free for the
   temporary migration staging schema and v5 data.
5. Ensure the target supports the `pg_trgm` extension.
6. Create or retrieve the v5 secret-management KEK. Keep it stable after the
   first v5 start; changing it makes database-stored secrets unreadable.
7. Configure private connectivity so the migration environment can reach both
   the v4 source database and the new v5 target database.

For Cloud SQL, run the migrator from an approved environment with Cloud SQL
Auth Proxy connectivity to both instances. Do not expose either database
publicly just for the migration.

## 2. Validate the chart

Run these commands from the repository root before each dev deployment and
before packaging a chart release:

```bash
helm dependency build charts-v5
helm lint --strict charts-v5 --values charts-v5/test-values.yaml
helm template dependencytrack-v5 charts-v5 \
  --namespace <namespace> \
  --values charts-v5/test-values.yaml > /tmp/dependencytrack-v5.yaml
```

Review the rendered manifest:

- API server and frontend images use the intended exact version or digest.
- `app.frontend.apiBaseUrl` matches the public URL plus `/api`.
- Cloud SQL proxy has the intended instance connection name.
- API server runtime ConfigMap contains the intended OIDC and CORS values.
- The ServiceMonitor has no `metadata.namespace`; Helm installs it in the
  release namespace.
- Local file storage has one API-server replica. Use S3 or RWX storage before
  increasing `app.apiServer.web.replicaCount`.

Also compare `charts-v5/values.yaml` with the exact upstream chart version
declared in `charts-v5/Chart.yaml` before upgrading that dependency.

## 3. Rehearse the data migration

Run this first with a restored copy or snapshot of the v4 database and a fresh
v5 target. Do not point a rehearsal at a v4 database with active writers if you
need trustworthy row counts.

Use a migrator image version that exactly matches the target Dependency-Track
v5 release. Run it as an ephemeral job or approved workstation container; it
is not part of this Helm chart.

1. Keep the v5 API server offline.
2. Bootstrap the empty v5 target schema:

```bash
v4-migrator bootstrap \
  --target-url '<v5-jdbc-url>' \
  --target-user '<v5-db-user>' \
  --target-pass
```

3. Verify the empty target. Apart from the `PERMISSION` table, row counts must
   be zero:

```bash
v4-migrator verify \
  --target-url '<v5-jdbc-url>' \
  --target-user '<v5-db-user>' \
  --target-pass
```

4. Run a dry run. Choose the retention period deliberately; `0` drops all v4
   metrics and a positive number retains that many days:

```bash
v4-migrator run \
  --target-url '<v5-jdbc-url>' --target-user '<v5-db-user>' --target-pass \
  --source-url '<v4-jdbc-url>' --source-user '<v4-db-user>' --source-pass \
  --metrics-retention-days <days> \
  --dry-run
```

5. If the dry run succeeds, run the migration without `--dry-run`:

```bash
v4-migrator run \
  --target-url '<v5-jdbc-url>' --target-user '<v5-db-user>' --target-pass \
  --source-url '<v4-jdbc-url>' --source-user '<v4-db-user>' --source-pass \
  --metrics-retention-days <days>
```

6. Verify the result and inspect every reported row-count reduction and probe:

```bash
v4-migrator verify \
  --target-url '<v5-jdbc-url>' \
  --target-user '<v5-db-user>' \
  --target-pass
```

Do not proceed unless `verify` completes successfully and every reduction is
understood. Keep the staging schema until review is complete. Run
`v4-migrator cleanup` only after that review.

## 4. Deploy the v5 chart in development

Create a Fasit/Feature configuration for a new development target. The chart
must point to the freshly migrated v5 database, not the existing v4 database.

Required values include:

- `app.apiServer.image.tag` and `app.frontend.image.tag`
- `cloudSqlProxy.connectionName`
- `app.database.jdbcUrl`, `app.database.username`, and `app.database.password`
- `app.secretManagement.database.kek.value`
- `app.frontend.apiBaseUrl`
- OIDC client IDs and issuer values
- bootstrap image tag, API URL, frontend URL, and required bootstrap secrets

Install or deploy the chart only after the data migration has verified. Wait
for the API server, frontend, and bootstrap Job to complete. The bootstrap URL
for the default monolith topology is:

```text
http://dependencytrack-v5-api-server:8080/api
```

If the release name or `app.fullnameOverride` changes, render the chart and
update `bootstrap.baseUrl` to the generated API service name before deploying.

## 5. Validate the development deployment

1. Check API-server, frontend, Cloud SQL proxy, and bootstrap Job logs.
2. Open the public frontend and complete an OIDC login.
3. Confirm the frontend calls the expected `app.frontend.apiBaseUrl`.
4. Upload a known SBOM and confirm it is processed.
5. Spot-check known projects, components, findings, teams, and users from v4.
6. Confirm Prometheus discovers the release-namespace ServiceMonitor and
   scrapes the API-server management endpoint.
7. Re-enter repository, analyzer, and vulnerability-source credentials in v5.
   The migrator intentionally does not carry these encrypted values over.
8. Review and re-enable migrated notification rules. They are intentionally
   disabled by the migrator.
9. Resolve users renamed with `-CONFLICT-LDAP` or `-CONFLICT-OIDC`.
10. Wait for upstream feeds to run, then verify EPSS data has repopulated.

## 6. Production cutover

Repeat the rehearsed sequence with the following additional controls:

1. Record the exact chart, subchart, API-server, frontend, bootstrap, and
   migrator versions in the change record.
2. Take a tested backup of the v4 database and verify it can be restored.
3. Announce downtime and stop the v4 API server. Confirm it has no active
   writers before extracting data.
4. Migrate into a fresh v5 target, verify it, then deploy this chart.
5. Complete the validation list above before reopening the service.

Rollback is a traffic decision: keep v4 stopped until v5 validation succeeds,
then route traffic back to the still-intact v4 deployment if the migration or
v5 validation fails. Do not attempt to migrate v5 changes back into v4.
