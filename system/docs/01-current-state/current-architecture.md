# Current-state architecture and risk baseline

## Evidence boundary

The source brief only establishes that the fictional PaySecure platform runs in AWS Mumbai (`ap-south-1`), handles 3.2 million transactions and INR 500 crore gross payment value per day, serves about 45,000 merchants, sees up to 1,200 TPS and has 180 ms P99 latency and 99.92% reported availability. It does not name the twelve microservices, instance sizes, CIDR ranges, AZ placement, database engine/version, broker/topic counts, load-balancer listeners, network policy, or current backup/monitoring configuration. The diagram therefore shows logical components from the brief and explicitly marks implementation details as inventory required; it must not be read as an as-built attestation.

## Working as-is model

Merchant APIs arrive over HTTPS at a public edge protected by WAF and load balancing. Payment orchestration validates merchant credentials, idempotency, amount and payment-method routing. A transaction-processing service calls an external rail adapter (UPI/NPCI, acquiring bank, card scheme, net banking, wallet or BNPL). A durable payment/ledger store records the request and outcome. Kafka carries asynchronous status, settlement, webhook, audit and notification events. Settlement reconciles PaySecure records with bank and network reports. Merchant configuration, operational state, cache and observability are separate logical domains. All are presumed to run in the single Mumbai region because the brief gives no contrary evidence.

For implementation planning, the twelve service slots are provisionally named `merchant-api`, `payment-api`, `payment-orchestrator`, `transaction-processor`, `upi-adapter`, `card-adapter`, `banking-adapter`, `fraud-service`, `token-service`, `settlement-engine`, `merchant-webhook`, and `notification-service`. These names are placeholders for discovery, not claims about the existing application. Replace them with an observed service catalog before migration.

## Illustrative failure domains and single points of failure

The confirmed architectural single point is regional scope: the brief says there is one AWS region and no automated regional failover. Any regional control-plane or data-plane event can affect ingress, compute, primary persistence, event processing and monitoring together. Even within Mumbai, the number of AZs and placement are unknown. A single writer may be a single logical authority even when Aurora is Multi-AZ. External rails, merchant DNS/resolvers, certificate renewal, CI/CD, secrets, identity, KMS authorization, observability alert delivery, NAT/egress, and staff access are additional dependencies to inventory; none should be marked redundant without evidence.

The current availability baseline implies approximately 7 hours of annual downtime. The target of 99.99% permits approximately 52.6 minutes annually, so the service needs about an eightfold reduction in unavailability. For 3.2m daily transactions, each minute at calendar-day average represents roughly 2,222 requests; at the stated festival peak it may represent 72,000 transaction attempts. These are throughput impact estimates, not unique lost payments or financial loss. Unknown outcomes must be reconciled before retries.

## Required discovery before production migration

1. Export AWS Config/resource inventory, subnet/CIDR/AZ map, EKS namespace/deployment inventory, ALB listener/rule map, security groups, WAF rules, IAM roles, KMS keys, DNS records and certificates.
2. For each data store, record engine and version, recovery model, transaction boundary, RPO/RTO evidence, encryption, retention and restore test date.
3. For Kafka, inventory brokers, topics, partitions, replication factor, retention, ACLs, schema registry, consumer groups, peak throughput and replay tolerance.
4. Map each payment method and partner endpoint, callback allowlist, mTLS certificate, participant certification, status inquiry, reversal and reconciliation process.
5. Trace all personal/payment data, logs, traces, backups, support sessions and vendor telemetry to verify locality and PCI boundary.
6. Establish service-level indicators from observed traffic, not scenario estimates: success rate, p95/p99, unknown outcomes, queue age, settlement breaks and merchant cohorts.

The current-state diagram in `current-state.mmd` uses dashed, labeled inventory placeholders to show unknowns. Discovery evidence must replace those placeholders before the drawing is called “as built.”

