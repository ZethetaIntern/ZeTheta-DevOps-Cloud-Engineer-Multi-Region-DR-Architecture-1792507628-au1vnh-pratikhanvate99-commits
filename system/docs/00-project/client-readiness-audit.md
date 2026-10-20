# Client-readiness and brief coverage audit

**Audit date:** 20 October 2026  
**Scope:** supplied Project 1B PDF, original user project description, and the contents of this repository at `main`.  
**Purpose:** distinguish completed artifacts from design intent, starter code, and work that requires the fictional client’s AWS accounts, approved vendors, operational evidence, or regulator/partner inputs.

## Executive assessment

The repository is a useful architecture and local payment-flow simulation package. It is **not a complete production DR system**, a compliance certification, or evidence that 99.99% availability, RPO under one minute, or RTO under five minutes has been achieved. The headline targets remain unproven, and the PDF brief’s Q3 2026 target date has passed. Re-baseline the delivery plan with the client before representing this as a completed implementation.

The recommended design is hot active-passive across Mumbai and Hyderabad with one financial ledger writer. Both active-active and active-passive patterns are described. The repository does not create or operate the regional AWS services, connect to NPCI/banks, run production drills, or certify any workload against a payment rail.

## Coverage against the original project description

| Requested outcome | Repository evidence | Assessment |
| --- | --- | --- |
| Multi-region DR architecture across AWS India regions | Architecture docs and diagrams describe Mumbai/Hyderabad; Terraform contains a two-region network starter | **Partial.** No deployed regions or full service topology. Confirm regional service availability, quotas, data paths, and partner connectivity in the client account. |
| Evaluate active-active vs active-passive | Comparison matrix and separate design docs | **Designed.** Cost, RTO/RPO, and capacity values are planning estimates or objectives; workload benchmarks and service quotes are absent. Hot active-passive is recommended. |
| Aurora PostgreSQL, DynamoDB, ElastiCache Redis, Kafka replication | Per-service strategy and sequence diagrams | **Designed, not implemented or measured.** No replication resources, lag measurements, failover proof, or reconciliation evidence in AWS. |
| 12 production-quality DR runbooks | Twelve scenario-named runbooks, each with the requested section structure and steps | **Templates, not ready for execution.** They repeat a common procedure and include `REPLACE_*` identifiers, generic commands, and client-dependent approvals. The scenario-specific triggers, blast-radius checks, tested command output, named roles, and rehearsed timing need client completion. |
| Scenarios include region failure, split-brain, Kafka partition, and peak-load failover | Region failure and Kafka failure have dedicated files; network partition and peak DDoS/traffic conditions are covered by related runbooks | **Partial.** Split-brain is discussed as a safety concern, but there is no dedicated split-brain runbook. Peak-load failover is not a clearly isolated scenario with measured capacity proof. |
| Map each design decision to RBI, PCI DSS v4.0, NPCI standards, India localisation | Sovereignty/compliance notes and component controls | **Partial.** This is an engineering mapping, not a clause-by-clause legal control assessment or QSA/NPCI/sponsor-bank approval. Legal applicability and current partner standards must be confirmed by counsel and the client’s payment participants. |
| 99.99% availability, RPO <1 minute, RTO <5 minutes by Q3 2026 | Targets and design gates are documented | **Not demonstrated.** No production deployment, failover drill, measured replication lag, or service availability record is supplied. Q3 2026 has passed as of this audit date. |
| Present to a five-stakeholder simulated BCRB | Simulated BCRB document | **Simulated artifact only.** No actual board meeting, stakeholder approval, or signed decision record. |

## Coverage against the supplied PDF acceptance criteria

| PDF requirement | Evidence in repository | Status / remaining evidence |
| --- | --- | --- |
| Current-state Mumbai baseline and failure-domain analysis | Current-state narrative and Mermaid diagram | **Assumption-labelled.** The source provides scenario inputs, not verified service inventory, network ranges, sizes, Kafka topics, or dependency map. Replace with discovery evidence. |
| Active-active and active-passive architecture, diagrams, quantified comparison across at least eight dimensions | Two design documents, diagrams, and comparison matrix with ten dimensions | **Design complete at planning level.** Estimates need benchmark, AWS Pricing Calculator inputs, actual topology, and owner-approved decisions. |
| Aurora, DynamoDB, Redis, Kafka and S3 replication strategies and sequence diagrams | Replication strategy and sequence diagrams | **Design present.** No production configuration, lag/consistency tests, recovery-point evidence, or replay exercise. |
| DNS failover design and YAML/JSON examples | Route 53 design plus health-check YAML and record JSON | **Templates present.** Substitute and validate real zones, records, health-check IDs, IAM, and routing policies; test resolver/client cache behaviour. |
| Twelve required incident scenarios | Twelve runbook files | **Scenario list present; operational readiness incomplete.** Each remains a template and must be tailored, command-reviewed, risk-approved, and rehearsed. The brief’s split-brain and peak-load cases need explicit standalone procedures and acceptance checks. |
| Ten-section runbook format, 15+ steps, commands, timing, branch points, owners, verification, rollback, communications, PIR | Runbook files have ten sections and 20 numbered steps | **Format present; not production-quality evidence.** Placeholder IDs and generic repeated content make them unsafe to execute without a client-specific tabletop and command-by-command review. |
| FMEA with at least 20 failure modes and RPN | FMEA CSV with 24 rows | **Present as illustrative analysis.** Confirm scoring definitions, owners, mitigations, and residual risk against client incidents and controls. |
| Cost workbook with four tiers and at least 15 service categories | Cost model XLSX | **Planning workbook present.** Editable inputs/multipliers are not an AWS quote; validate formulas, workload dimensions, commercial discounts, support, network, logging, and India-region pricing. |
| Annual drill plan: quarterly full failover, monthly component drills, weekly checks | Annual drill plan | **Calendar/checklist present.** It is not evidence the drills occurred or met objectives. Add accountable names, evidence storage, pass/fail criteria, and dates agreed with client. |
| Terraform 1.5+ modules for network, compute, database, DynamoDB, cache, messaging, DNS, monitoring, security and remote state | Terraform root with network module and two-region VPC/subnet starter | **Partial.** Required service modules and remote backend are not implemented. Terraform validation was not confirmed because the AWS provider failed to respond in the prior local check. No plan/apply against a client account. |
| Kubernetes payment API, transaction processor, settlement engine; probes, HPA, PDB, ConfigMap/Secret, NetworkPolicy | Kubernetes workload and policy template | **Partial.** Starter manifests use customer-supplied images/settings; probe and runtime readiness coverage is incomplete for all workloads. Secrets and quotas require deployment-specific setup. |
| Monitoring dashboards, alerting, communications, chaos and multi-cloud considerations | Operations docs for monitoring, comms, chaos/multi-cloud | **Design/checklist level.** No installed dashboards, alert routing, chaos execution, secondary-cloud deployment, or tested escalation evidence. |
| India localisation and mapping to RBI, PCI DSS v4.0, NPCI standards | Compliance/data sovereignty notes | **Partial.** Requires applicability review against the legal entity, current directions, actual payment data flows, PCI scope, sponsor-bank/NPCI agreements, and audit/QSA evidence. |
| Ten BCRB challenge questions and five stakeholder roles | Simulated BCRB document | **Simulated.** Not a real approval. Obtain named attendees, decisions, actions, dissent, and sign-off. |
| Identify five deliberate errors in the source | Source review identifies five material arithmetic/technical/legal claims needing correction | **Independent review provided, author intent unknown.** The source does not identify which errors it deliberately planted; do not claim these are confirmed planted errors. |
| Local software demonstration | Go API, PostgreSQL/Kafka Compose stack, smoke script and unit tests | **Local simulation only.** Prior `go test ./...`, `go vet ./...`, Compose configuration check, and PowerShell smoke test passed. No real payment, NPCI, card, or AWS path is exercised. |

## Client-facing project description

Recommended GitHub About description:

> Reference architecture and local payment-flow simulator for a fictional Indian payment aggregator, with Mumbai–Hyderabad DR designs, replication strategies, runbook templates, compliance notes, and a Terraform network starter. Production RPO/RTO and regulatory compliance are not yet validated.

The public GitHub About field still contains a truncated copy of the original assignment prompt. The unauthenticated repository page does not expose an edit control; update this field while signed into the repository owner account. The repository README already includes a clear scope and limitations statement.

## Client share recommendation

The repository can be shared as a **design package and demonstrator**, provided the client is told that the service targets and production deployment are not yet validated. Do not describe it as a fully implemented or compliant production DR system. Before operational use, close at least these gates:

1. Replace assumed topology and every placeholder with approved client inventory and account identifiers; review every runbook command with a second operator.
2. Implement missing Terraform modules, remote state, IAM/security controls, managed services, DNS, observability, and deployment automation; validate plans in a non-production account.
3. Define and test split-brain fencing, ledger consistency, Kafka/DynamoDB/Redis/Aurora recovery, reconciliation and peak-load capacity with business-safe test data.
4. Execute witnessed component and full-region exercises; capture measured RPO/RTO, transaction loss/duplication, recovery timestamps, and signed acceptance against SLOs.
5. Obtain client legal/compliance, QSA, sponsor-bank/NPCI and security approvals for the applicable obligations and actual data flows.
6. Re-baseline the program schedule because the brief’s Q3 2026 target date has elapsed.

## Source-data arithmetic checks

These scenario calculations are not production telemetry: 3.2 million transactions/day averages 37.04 TPS across 24 hours and 55.56 TPS across a 16-hour business window; it does not yield the brief’s stated 250 TPS without another stated time window. 99.99% annual availability permits about 52.56 minutes of downtime in a 365-day year; 99.92% permits about 420.48 minutes (7.01 hours). INR 500 crore/day at an assumed 1.8% fee is INR 9 crore/day of gross fee estimate, not transaction value or verified revenue.

See [requirements traceability](requirements-traceability.md), [source review](../10-source-review/deliberate-errors-found.md), and the [repository README](../../../README.md).
