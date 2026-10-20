# RB-12: Cascading microservice/fraud failure

**Owner:** Incident Commander. **Severity:** SEV-1. **Objective:** restore safe service within 5 minutes where the scenario permits; preserve RPO <60s as a measured target, not a guarantee. **Prerequisites:** current contact tree, customer account IDs, role ARNs, resource inventory, tested break-glass access. **Roles:** Incident Commander (IC) owns declaration and timeline; Platform Lead executes cloud actions; Data Lead validates writer/replication; Security Lead handles evidence and keys; Payments Lead owns rail state/reconciliation; Communications Lead issues approved notices.

## 1. Trigger and declaration

dependency latency, rejection anomalies, queue growth. Declare when two independent indicators confirm service impact or when safety/integrity is at risk. Record UTC declaration and detection times.

## 2. Roles and communications

IC controls command and go/no-go. Platform Lead runs AWS actions; Data Lead owns database/stream correctness; Security Lead owns incident evidence; Payments Lead owns rail-side state; Communications Lead approves customer, bank, NPCI and regulator notices with Legal/Compliance. Use out-of-band channel if corporate identity is suspect.

## 3. Safety checks

No dual ledger writer. Require two-person review for promotion, key rotation and destructive recovery. Confirm exact account/region/resource before commands. Commands below are templates: replace all `REPLACE_*` tokens; never paste secrets. Obtain legal/rail approval for external notification.

## 4. Timed action plan

| Elapsed | Owner | Action |
|---|---|---|
| 0-1 min | IC | Declare, page roles, freeze unsafe retries |
| 1-2 min | Data + Security | Capture evidence, writer identity, lag/checkpoints |
| 2-3 min | Platform | disable noncritical fanout, apply safe degraded mode and independent fraud fallback |
| 3-4 min | Payments + Platform | Canary, rail status inquiry, controlled traffic |
| 4-5 min | IC | Objective gate, customer status and next decision |
| 5+ min | All | Reconcile, sustain, or rollback to safe state |

## 5. Operator steps

1. Incident commander opens a severity-1 incident for **Cascading microservice/fraud failure** and records UTC start time.
2. On-call verifies the signal from two independent sources: dependency latency, rejection anomalies, queue growth.
3. Payments lead freezes risky automated retries and identifies affected transaction window.
4. Security lead preserves CloudTrail, application, database and provider evidence with immutable retention.
5. Platform lead pages primary and secondary on-call; establish out-of-band bridge.
6. Data lead captures writer identity, replication lag, consumer offsets and last known good checkpoint.
7. IC confirms customer impact, region, merchant scope and payment uncertainty; announce internal status.
8. Execute this scenario action: disable noncritical fanout, apply safe degraded mode and independent fraud fallback.
9. Verify all write-capable services have exactly one authorized region/owner.
10. Apply protective traffic controls; preserve merchant idempotency and request correlation.
11. Validate API synthetic flow and database read/write canary using nonfinancial test tenant.
12. Reconcile uncertain requests against durable records and rail/provider status; prohibit blind retries.
13. Ramp traffic in 5%, 25%, 50%, 100% steps, holding at least five minutes at each gate.
14. At every gate check error rate, p99 latency, Aurora replication/writer status, Kafka lag and fraud controls.
15. Payments lead signs off duplicate, missing, delayed and settlement exception counts.
16. IC confirms recovery objective timestamps and customer/partner notice status.
17. Keep the failed domain fenced; do not fail back until root cause and data parity are proved.
18. Security and service owners review temporary permissions, routes, scaling and circuit breakers.
19. Record event IDs and remediation owners; preserve all relevant artifacts.
20. Within five business days, issue post-incident review, corrective actions and drill updates.

## 6. Commands and branches

```bash
aws sts get-caller-identity
aws rds describe-global-clusters --region ap-south-1 --global-cluster-identifier REPLACE_GLOBAL_CLUSTER
aws cloudwatch get-metric-statistics --region ap-south-1 --namespace AWS/RDS --metric-name AuroraGlobalDBReplicationLag --start-time "$(date -u -d '10 minutes ago' +%FT%TZ)" --end-time "$(date -u +%FT%TZ)" --period 60 --statistics Maximum --dimensions Name=DBClusterIdentifier,Value=REPLACE_CLUSTER
aws route53 list-resource-record-sets --hosted-zone-id REPLACE_ZONE
aws kafka list-clusters-v2 --region ap-south-2
```

**Branch A:** if the primary is confirmed fenced and DR data is within approved loss bound, IC authorizes promotion. **Branch B:** if primary may still accept writes or lag is unknown, keep writes fenced/queued and escalate; do not promote. **Branch C:** if DR health/capacity fails canary, revert DNS/traffic to the sole safe writer if available or remain in controlled outage. These commands are read-only examples except explicitly approved promotion/change commands; use reviewed account-specific command from change record.

## 7. Validation and recovery gates

{guard}. Require exactly one writer; no unresolved offset skip; canary success; p99 within the approved SLO; error rate below 0.1% for 5 minutes; observed replication and event reconciliation; fraud, audit, logging and encryption controls healthy. Actual scenario objective is confirmed by timestamps from declaration to safe processing, not DNS change alone.

## 8. Rollback / failback

Stop traffic ramp and return to the last verified safe state if duplicate/unknown outcomes rise, writer count >1, fraud bypass occurs, or integrity checks fail. Keep failed region isolated. Failback is a separate planned change: restore/reseed, verify parity, fence DR writes, reverse replication, execute canary, and shift traffic gradually with IC approval. Never reverse roles while both sides are writable.

## 9. Communications

Internal updates every 15 minutes during SEV-1. Payments/Legal decide partner and regulator obligations using applicable contracts and law; log recipient, timestamp, approved facts and next update. Customer notice must explain impact and safe retry guidance without exposing security details.

## 10. Closure and post-incident review

Close only after payments reconciliation, single writer proof, stable service for 30 minutes, stakeholder updates, and action owners. Within five business days document timeline, RPO/RTO observed, loss/duplicate counts, evidence links, what failed, corrective actions, owner/date, and next exercise.
