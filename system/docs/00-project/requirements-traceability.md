# Project requirements traceability

**Source:** the supplied 68-page Project 1B brief. Its assignment and acceptance criteria are treated as project requirements; its example facts, commands, prices and legal claims are separately checked. **Status meanings:** Implemented = artifact/code exists; Assumption = fictional topology is not evidenced; Pending customer = requires real account data, contract or decision; Not applied = this repository does not provision live cloud resources.

| Brief requirement | Deliverable | Status |
| --- | --- | --- |
| Current Mumbai single-region baseline and failure-domain analysis | `../01-current-state/current-architecture.md`, `current-state.mmd` | Assumption-labelled; the brief does not provide actual 12 service names, sizes, CIDRs, queue topics or listener rules. |
| Active-active and active-passive diagrams, analysis, 8 quantified dimensions | `../02-multi-region-design/` | Both patterns documented; quantified engineering/cost values labelled estimates or targets. |
| Aurora, DynamoDB, Redis, Kafka and S3 replication, sequence diagrams, CAP and latency | `../03-data-replication/` | Per-component strategy, timing, alarms and reconciliation. |
| DNS failover, TTL, health checks, valid YAML/JSON | `../04-dns-failover/` | Config templates are valid once customer-owned names/IDs/targets are supplied. |
| 12 mandatory scenarios, exact ten-part runbook structure, 15+ steps, CLI, timing, branching, comms and rollback | `../05-runbooks/` | Scenario-specific operator procedures are supplied, but remain templates pending client inventory, concrete resource IDs and rehearsal; do not execute as-is. External rail, key, DNS and ransomware actions are guarded. |
| FMEA with at least 20 modes and RPN | `../06-risk-cost/fmea.csv` | 24 illustrative scores; reassess with incident history and control owners. |
| Cost model with 15+ categories and Cold/Warm/Hot/Active-Active | `../06-risk-cost/cost-model.xlsx`, cost notes | Formula-driven scenario multipliers and transparent line items; prices left as editable inputs because there is no workload quote or live AWS price export. |
| Data sovereignty map for 6 categories | `../07-compliance/data-sovereignty.md` | India-only payment-data design; legal applicability requires counsel confirmation. |
| Annual plan: quarterly full failover, monthly component, weekly verification | `../08-operations/annual-drill-plan.md` | Twelve monthly events and 52 weekly checks listed. |
| Terraform 1.5+ with network modules, service designs, and remote state | `../configs/terraform/` | **Partial:** only a two-region VPC/subnet starter is implemented. EKS, database, DynamoDB, cache, messaging, DNS, monitoring and security resources are architectural specifications, not Terraform modules. Backend and provider validation remain client-account gates. |
| Kubernetes payment-api, transaction-processor and settlement-engine, probes, HPA, PDB, config/secret, network policy | `../configs/kubernetes/` | Manifests supplied; images, names, quotas and external controller values are customer inputs. |
| BCRB, ten challenge answers, comms matrix, monitoring/dashboard, chaos design, multi-cloud | `../09-review/`, `../08-operations/` | Simulated; does not represent real board or regulator approval. |
| Identify five deliberate errors | `../10-source-review/deliberate-errors-found.md` | Five high-risk claims/calculation errors are documented with confidence. The brief does not label which five it planted, so exact author intent cannot be inferred. |
| Verify overall client and PDF acceptance coverage | `client-readiness-audit.md` | Explicit status audit: distinguishes design artifacts, starter templates, simulated outputs, and evidence not yet available. |

**Client-share caveat:** “Implemented” in this matrix means an artifact exists in the repository; it does not mean a production AWS service is deployed or that an acceptance target has passed. In particular, all 12 runbooks still contain client-specific placeholders and require review and rehearsal; Terraform only provisions a network starter; RPO/RTO, availability and regulatory acceptance are unproven. See `client-readiness-audit.md` before presenting this as completed production work.

## Scenario inputs used consistently

3.2 million transactions/day; INR 500 crore/day gross transaction value; peak 1,200 TPS; 45,000 active merchants; 180 ms P99; current availability 99.92%; target 99.99%; target RPO <60 seconds and RTO <5 minutes; 42 engineers total, 8 platform/infra; INR 8 crore/year current infrastructure. These are scenario inputs, not PaySecure production telemetry.

Arithmetic checks: calendar-day average = 3,200,000 / 86,400 = 37.04 TPS. Business-hours average over 16 hours = 55.56 TPS, not the 250 TPS stated in the brief. Scenario fee-revenue estimate at 1.8% take rate = INR 9 crore/day or INR 37.5 lakh/hour. Gross transaction throughput by value is INR 20.83 crore/hour over 24 hours; it is not PaySecure revenue. 99.99% availability allows 52.56 minutes per 365-day year; 99.92% implies 420.48 minutes (7.01 hours).

