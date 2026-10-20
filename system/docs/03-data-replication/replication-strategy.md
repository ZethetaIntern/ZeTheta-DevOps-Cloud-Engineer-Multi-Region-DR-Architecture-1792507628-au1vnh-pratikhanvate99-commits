# Data replication strategy

## Authority and invariants

Aurora PostgreSQL in Mumbai is the single financial ledger writer; Hyderabad is a hot, provisioned cross-region secondary. Fence the old writer before promotion. Replication is asynchronous, so RPO <60s is an objective measured by lag alarms and exercises, not a provider guarantee. Reconcile every uncertain payment by stable merchant/idempotency key and rail reference before retrying. Do not infer authorization from an HTTP timeout.

| Component | Strategy | DR semantics / guardrails |
|---|---|---|
| Aurora PostgreSQL | Aurora Global Database, Mumbai primary, Hyderabad secondary; continuous storage replication; pre-created promotion runbook | Async typically sub-second but nonzero; measure `AuroraGlobalDBReplicationLag`; promoted secondary becomes writer; one-writer fencing and transaction reconciliation required. |
| DynamoDB | Global Tables MREC in ap-south-1 and ap-south-2 for non-ledger idempotency/merchant operational metadata | Multi-region eventual consistency; concurrent same-item writes can conflict. Deterministic ownership/conditional writes; never use it as the payment ledger. |
| ElastiCache Redis | Global Datastore primary/secondary; cache/session only | Async and potentially lossy; rebuild from durable stores; no balance, idempotency, or payment truth kept only here. Promote after database authority is established. |
| MSK Kafka | MSK Replicator Mumbai to Hyderabad; topic/ACL/config parity; consumer idempotency | Async, at-least-once replication; offsets and ordering are not a cross-region transaction. Stable event IDs, replay-safe consumers, lag monitoring, pause producers during cutover. |
| S3 | Versioned, encrypted buckets in India; CRR for approved objects; replication metrics and inventory | Replication is asynchronous; verify object replication status and KMS key access. Keep payment data within India. Backups are immutable and separately permissioned. |

## Recovery sequence

1. Detect and declare using independent synthetic checks and operator confirmation.
2. Fence primary write paths at ingress, worker deployments, and database credentials; preserve evidence.
3. Compare Aurora lag / latest durable sequence and Kafka replicated offsets; quantify possible loss.
4. BCRB incident commander accepts promotion risk and records event/time.
5. Promote Aurora in Hyderabad, verify writer endpoint and consistency; never force-promote while primary can still accept writes.
6. Enable Hyderabad consumers and API; check outbox/inbox idempotency, then Redis cache.
7. Shift DNS gradually; verify payment acceptance and rail-side reconciliation before full traffic.
8. Reconcile uncertain window and replay only deduplicated events. Keep Mumbai fenced until explicit failback.

## CAP, latency and acceptance

For the financial ledger, consistency and single-writer safety take priority during a partition; some writes may be rejected/queued. Active-active for independent stateless API/read paths does not imply dual-writer ledger. Region RTT Mumbai-Hyderabad is scenario-estimated 15-25ms; measure actual network and provider behavior. Acceptance requires p95/p99 transaction latency, replication lag, failover time, duplicate/lost event counts, and rail reconciliation evidence over repeated peak-load drills.
