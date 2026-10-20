# Active-passive design: hot standby for sub-five-minute recovery

## Decision and acceptance boundary

Use Mumbai (`ap-south-1`) as the normal payment-write region and Hyderabad (`ap-south-2`) as a hot standby. The passive region continuously runs the full network/security boundary, EKS control plane, core worker capacity, database instances, Kafka brokers, cache replica and operational tooling. Maintain sufficient ready capacity for the measured 1,200 TPS peak plus a 30% engineering headroom target (1,560 TPS). Autoscaling is for burst margin, not the primary recovery mechanism. This is a hot standby, not a low-capacity warm standby: the brief's tier table gives warm standby a typical 15-60 minute RTO, which cannot substantiate a sub-five-minute objective.

Only Mumbai is authorized to originate payment, ledger and settlement writes in steady state. The Hyderabad copy is replication receiver/read-only for stores whose service supports it. A region change is a controlled promotion, with explicit write fencing and one operator-approved authority. Do not rely on DNS health alone to promote data services. The system prefers consistency for payment records during an inter-region partition: stop new writes in the region that cannot prove authority and return safe retryable responses. Merchant idempotency keys and partner transaction references make uncertain outcomes recoverable; they do not make an external bank authorization atomic with the database.

## Region choice and latency budget

The source brief nominates Hyderabad and quotes 15-25 ms Mumbai-Hyderabad RTT. Treat that as an initial planning range, not an SLA or live measurement. During discovery, measure p50/p95/p99 TCP/TLS and database replication lag from actual VPCs at peak and settlement windows. The profile P99 is 180 ms. A simple upper-bound budget using 25 ms synchronous inter-region round trip would leave at most 95 ms from the brief's claimed 300 ms UPI ceiling (`300 - 180 - 25`), before any remote commit, queueing, tail amplification or external rail time. That arithmetic does not prove end-to-end feasibility. The 300 ms ceiling itself is not verified as a universal NPCI rule and is not used as a compliance acceptance threshold.

Aurora Global Database, DynamoDB MREC, Redis Global Datastore, MSK Replicator and S3 CRR are asynchronous patterns. Do not state that cross-region replication is synchronous or has zero data loss. Use the RPO as a measured gate, alert and business admission decision. At 1,200 TPS, 60 seconds is up to 72,000 request attempts crossing the theoretical exposure window; at 250 TPS it is 15,000. Those are attempts, not unique committed payments or an INR loss estimate. Financial impact must be reconciled from durable event IDs and partner responses.

## Topology and write path

1. Route 53 health checks and ARC routing controls observe independent regional endpoints. Each endpoint includes app, dependency and synthetic-write health; a shallow TCP or `/live` check cannot trigger regional failover by itself.
2. Mumbai WAF and NLB/ALB accept merchant HTTPS requests. EKS services validate API identity, merchant status, amount, currency and merchant-scoped idempotency key.
3. A PostgreSQL transaction writes the payment intent and a transactional outbox row together. An outbox publisher writes immutable event IDs to Kafka. Consumers are idempotent and persist their consumer offsets/dedup state with business effects.
4. The authoritative payment/ledger state is Aurora PostgreSQL. DynamoDB Global Tables are limited to non-ledger operational metadata that can tolerate the selected consistency behavior. Redis contains rebuildable cache entries only. Kafka carries events and is not treated as the payment ledger. S3 stores encrypted backups/audit exports, not the only current transaction copy.
5. Mumbai replicates required data to Hyderabad. Hyderabad has matching application image digests, schema compatibility, endpoint secrets, firewall rules, certificates, IAM permissions, Kafka topics/ACLs, monitoring and test capacity.
6. On failover, declare the event, freeze deploys, fence the old writer, inspect lag and external transaction status, promote Aurora, promote Redis only if used, validate Kafka target offsets, then open a controlled traffic canary. Increase traffic only after a payment, reversal/status inquiry and duplicate idempotency probe pass.

## Failover phases and timing budget

The five-minute objective is measured from first missed critical-journey SLI to successful recovery of the critical payment journey. Target checkpoints: detection 0-30 s; declaration/fencing 30-60 s; database promotion and endpoint refresh 60-180 s; stream/cache activation and synthetic validations 180-240 s; traffic ramp and business verification 240-300 s. This is an acceptance budget, not a vendor timing promise. Aurora unplanned promotion may take minutes; any run failing 300 seconds is recorded as an objective miss. If the organization wants the customer RTO to exclude partner network outages, it must publish both PaySecure platform RTO and end-to-end merchant RTO.

Do not auto-failback. Rebuild Mumbai from the new authoritative Hyderabad writer, verify catch-up and configuration parity, perform a planned zero-data-loss switchover when both sites are healthy, then restore Hyderabad standby. Keep settlement and finance control owners in the go/no-go.

## Data sovereignty and tertiary event

Keep payment message/transaction data, copies, logs containing payment records, database backups, Kafka payloads, operational exports and recovery data in Indian regions. No payment data is routed to Singapore. For simultaneous failure of both Indian regions, the safe default is fail closed for new payment writes, protect unknown outcomes, preserve encrypted India-held backups, communicate, and restore when an Indian region is available. A third-region compute-only fallback is not a payment-processing DR site unless legal/compliance validates exactly what data and operations are permitted; it is excluded from this design.

## Operational and cost implications

This design requires the Hyderabad hot data plane and sufficient ready EKS/MSK capacity, so its cost resembles the brief's Hot Standby tier (1.6-1.8x its current annual infrastructure baseline) rather than Warm Standby. The planning midpoint 1.7x on INR 8 crore is INR 13.6 crore/year, an incremental INR 5.6 crore/year. The cost workbook leaves all service rates editable and this multiplier is only scenario arithmetic. A rightsized lower floor may cost less but cannot be described as meeting 1,560 TPS or five-minute recovery until measured.

Operationally, eight platform engineers can operate this only with a narrow service catalog, one-button evidence capture, quarterly full regional drills, monthly component drills, weekly automated readiness, named application and finance owners, and 24x7 on-call. It reduces conflict-resolution logic compared with active-active but retains high-quality incident-command, partner and reconciliation responsibilities.

## Go-live acceptance gates

- Both regions deploy from identical reviewed IaC and have approved quotas, image versions, network policies and regional secrets.
- Controlled failovers demonstrate RTO <300 s for the defined synthetic journey and RPO <60 s for authoritative payment/ledger state across repeated runs, including test at 1,200 TPS.
- A separate overload drill sustains 1,560 TPS with noncritical load shed and no duplicate authorization/settlement effects.
- Tests prove old-region write fencing even with asymmetric connectivity, delayed control-plane APIs and stale DNS clients.
- Operations validates recovery of backups, keys, TLS, partner circuits, status inquiry, reconciliation and out-of-band communications.
- Compliance, QSA, Legal, sponsor bank and NPCI counterparties confirm applicable controls and notices.

