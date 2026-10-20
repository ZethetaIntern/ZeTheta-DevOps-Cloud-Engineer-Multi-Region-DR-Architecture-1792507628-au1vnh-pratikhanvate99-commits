# PaySecure multi-region payment system

This repository is an implementation starter for the architecture in `../deliverables/PaySecure_Multi_Region_DR_Architecture_Client_Design.pdf`. It contains a locally runnable payment-flow simulator, AWS infrastructure definitions, Kubernetes manifests, and operator commands. It does not connect to NPCI, an issuer, an acquirer, or a real payment network. It does not accept or store PAN, CVV, PIN, UPI PIN, or other payment credentials. The demo is not a PCI DSS attestation or production payment service.

## What is implemented

- Go HTTP API with merchant-scoped idempotency and a PostgreSQL transactional outbox.
- Kafka outbox publisher and idempotent demo processor; duplicate delivery does not create a second payment result.
- Local Docker Compose stack: API, PostgreSQL, Kafka in KRaft mode.
- Kubernetes deployment templates for API, processor and settlement workers, HTTP liveness/readiness probes, disruption budgets, API autoscaling, network policy and zone spreading.
- Terraform starter for a two-region AWS India network foundation. Managed data services are specified in architecture docs and require client-specific implementation.
- 12 scenario-specific runbooks for failover, failback, incident response and recovery verification.

The templates still require client-specific image digests, secret-provider integration, quotas, network rules and partner configuration. Worker readiness probes verify the database; Kafka consumer lag and health must also be monitored externally. None of these manifests has been applied to a client cluster.

## Local demonstration

Requirements: Docker Compose v2 and Go 1.23+ (only needed to build outside Docker).

```powershell
# If local port 8080 is occupied, choose another host port first:
$env:PAYSECURE_HOST_PORT = '18080'
docker compose up --build
```

Create a simulated payment (amounts are paise):

```powershell
$key = [guid]::NewGuid().ToString()
$port = if ($env:PAYSECURE_HOST_PORT) { $env:PAYSECURE_HOST_PORT } else { '8080' }
Invoke-RestMethod -Method Post -Uri http://localhost:$port/v1/payments `
  -Headers @{ 'Idempotency-Key' = $key; 'X-Merchant-ID' = 'merchant-demo-001'; 'X-API-Key' = 'local-demo-api-key' } `
  -ContentType 'application/json' `
  -Body '{"amount_minor":159900,"currency":"INR"}'
```

The local compose profile enables only `merchant-demo-001` with a local-only API key. Repeat the same request with the same key to receive the same payment ID. Poll `GET /v1/payments/{id}`; the demo worker moves it from `PENDING` to `SIMULATED_APPROVED`. This is a software simulation only. `GET /health/live` and `/health/ready` are available for orchestration.

Run the local API checks with `./scripts/smoke-test.ps1`. Stop with `docker compose down`. Local volume data is retained; remove it only when intentionally resetting the local demo.

## AWS deployment path

1. Use a dedicated AWS account, Terraform state bucket with locking, and separate deployment/operator roles. Keep credentials out of this repository.
2. Review `configs/terraform/terraform.tfvars.example`; populate account-specific CIDRs, hostnames and approved capacity. Check current service quotas, regional feature availability, PCI scope and NPCI/sponsor-bank endpoint registration.
3. Initialize and review the plan from `configs/terraform/`:

   ```powershell
   terraform init
   terraform validate
   terraform plan -out paysecure.tfplan
   ```

4. Human change approval is required before `terraform apply`. This configuration provisions billable resources. It is not applied by this project.
5. Build and push a pinned image digest to ECR, deploy the Kubernetes manifests, and populate Secrets Manager / Kubernetes secrets through the organization’s approved secret pipeline.
6. Configure MSK Replicator according to `docs/03-data-replication/replication-strategy.md`; source/target IAM and private connectivity must be designed and tested separately.
7. Execute the rehearsal checklist in `docs/05-runbooks/` in a non-production account. Promote no database or region based only on an alarm.

### Important system boundaries

- Aurora PostgreSQL is the authoritative payment record. DynamoDB is provisioned for non-ledger idempotency/operational metadata only; conflict resolution must not determine financial truth.
- Aurora Global Database, DynamoDB MREC, Redis Global Datastore and MSK Replicator replicate asynchronously. The RPO objective is an acceptance gate, not a service guarantee.
- The demo API models acceptance into a queue and a simulated processor result. A real payment integration requires certified provider adapters, HSM/key custody, legal and compliance review, NPCI/bank certification, reconciliation, fraud controls and 24x7 operations.
- Terraform state, logs, snapshots, backups, secrets and Kafka topics must remain in approved India regions for payment data. Review every telemetry and vendor support path.
- Do not use `terraform destroy` against a production workspace. No destroy workflow is included.

## Architecture and operations package

See `docs/00-project/requirements-traceability.md` for the requirement matrix; `docs/05-runbooks/` contains 12 scenario runbooks; `docs/06-risk-cost/cost-model.xlsx` is a planning sensitivity; `docs/07-compliance/` maps locality and controls; `docs/09-review/bcrb.md` contains the simulated review.

## Layout

```text
app/                  Go API, outbox publisher and demo processor
db/                   PostgreSQL schema and migrations
configs/terraform/    Two-region network starter
configs/kubernetes/  Workload and policy templates
docs/                 Architecture, replication, runbooks, compliance, operations
```
