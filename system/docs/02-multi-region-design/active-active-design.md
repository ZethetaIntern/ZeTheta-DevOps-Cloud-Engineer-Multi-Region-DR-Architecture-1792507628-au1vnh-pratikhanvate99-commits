# Active-active design: evaluated alternative

## Proposed pattern

An active-active option would run public APIs, stateless orchestration, EKS services and read paths in Mumbai and Hyderabad, with Route 53 latency or weighted routing. Compute, WAF, certificates, secrets, monitoring, CI/CD artifacts and downstream connectors would be provisioned in both regions. The choice for stateful writes is not indiscriminate multi-writer replication. Partition merchants/payment aggregates into a single home region using a versioned, strongly controlled merchant-to-region directory. The owning region serializes each merchant's ledger mutations. Requests arriving in the other region either proxy to the owner over private links or return a bounded retry response. On owner-region outage, fence the former owner, promote the replica after lag evaluation and update the directory epoch before accepting writes.

This provides active-active edge and compute capacity while preserving one write authority per payment aggregate. It is more accurate to call it partitioned active-active with home-region ownership. An unsharded two-writer payment ledger is rejected because network partitions can produce conflicting balances, duplicate payment/settlement actions and ambiguous idempotency outcomes. DynamoDB MREC's last-writer behavior is not a substitute for ledger serialization. DynamoDB MRSC has topology and supported-region constraints that must be checked against the current AWS feature matrix; no topology is assumed to make Aurora, Kafka, Redis and external rails one distributed transaction.

## Request and failure handling

The edge selects an owning region from an authoritative merchant directory and sends each merchant's ledger-affecting requests to that region. The idempotency key is merchant-scoped and stored in the same transactional boundary as payment intent/outbox. Both regions can serve static reads and non-authoritative catalog/configuration reads from local replicas with freshness labels. A successful response only follows durable commit in the owning ledger. If the directory is unavailable or the owner is unreachable, do not guess a new owner. Return a safe retryable response, preserve partner status references and execute a fenced ownership transfer. A network partition is treated separately from regional destruction: maintain the incumbent writer when it has authority, isolate the other site, and fail closed rather than allowing both sides to process the same merchant's account.

Settlement aggregates are partitioned by the same ownership key. Every event carries region and epoch metadata, immutable event ID, payment ID, idempotency key, partner reference and schema version. Reconcilers detect duplicate effects and ledger gaps across both Kafka clusters and partner reports. Failback requires a higher ownership epoch only after reconciliation and a planned switchover.

## Quantified trade-off (scenario assumptions)

With measured Mumbai-Hyderabad RTT in the brief's 15-25 ms planning band, remote owner writes add at least one inter-region coordination trip if the design requires remote authority confirmation; actual p99 adds network, queue and service time and cannot be inferred from RTT alone. A partitioned ownership model routes roughly one half of merchants to each region only if merchant load is uniform; 45,000 merchants and value/TPS are not known to be uniform, so capacity must be tested with skewed hot merchants and failover of the largest shard. When one region fails, the survivor must support the entire 1,200 TPS peak and test headroom of 1,560 TPS, not half the load.

The brief's own planning multipliers place active-active at 1.8-2.2x current spend: on the INR 8 crore baseline, INR 14.4-17.6 crore annually. That excludes unpriced engineering, partner certification, fees and potential overprovisioning. Compare against hot standby 1.6-1.8x (INR 12.8-14.4 crore); midpoints differ by INR 1.6 crore/year (2.0x vs 1.7x), a planning estimate only.

## When to reconsider

Reconsider active-active only if measured business loss from the failover window justifies additional cost, PaySecure can operate per-merchant ownership and reconciliation at scale, external rail partners certify dual-region endpoints, consistency/partition tests prove fencing, and the QSA/compliance review accepts the expanded CDE and control surface. Start with stateless read traffic and synthetic traffic distribution; do not enable dual writes by a feature flag without a formal data migration and rollback plan.

