# DependencyTrack

## Trivy

Note!

PVCs does not get deleted when the pod is deleted. There is no way to override this behavior in the helm chart.
This has to be done manually.

## Encryption key

The chart creates `<release>-secret-key` during installation and keeps it on
upgrades. Do not delete or replace this Secret in an established environment:
Dependency-Track uses it to decrypt encrypted configuration values.

If a newly created environment has an incorrectly generated key and contains
no encrypted configuration data, delete only `<release>-secret-key` before a
new Helm install. Confirm that the mounted `secret.key` contains 32 bytes
before allowing the bootstrap Job to update encrypted properties.
