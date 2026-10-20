# Annual resilience and DR exercise plan

Plan year starts January; shift months to the approved business calendar. Four full failovers occur quarterly; twelve component-level drills monthly; 52 weekly verification cycles (one each week). All exercises use non-production accounts or explicitly approved production game days with rollback controls.

| Month | Component drill | Quarterly full failover | Evidence / exit gate |
|---|---|---|---|
| Jan | Aurora lag and restore | — | PITR restore, checksums, measured lag |
| Feb | Kafka broker and replicator lag | — | replay, duplicate count, offset parity |
| Mar | Redis Global Datastore promotion | Q1 Mumbai-to-Hyderabad | RTO/RPO, reconciliation, sign-off |
| Apr | DynamoDB conflict and recovery | — | conditional write conflict evidence |
| May | DNS and certificate rotation | — | resolver convergence by network |
| Jun | AZ loss and capacity | Q2 failover | 1,560 TPS target and safe degraded mode |
| Jul | KMS/key compromise tabletop | — | access isolation and evidence custody |
| Aug | DDoS/peak load | — | 150k RPS defensive test in approved scope |
| Sep | Ransomware immutable restore | Q3 failover | clean-room recovery and evidence |
| Oct | NPCI/third-party outage simulation | — | status inquiry, no blind retry |
| Nov | Fraud/microservice cascade | — | controls remain effective in degraded mode |
| Dec | Full end-to-end component review | Q4 failback and failover | two-way procedure, lessons/actions |

Weekly checks (52 entries in the calendar): verify backups/restore evidence, KMS alarms, IAM drift, certificates, DNS synthetics, Aurora lag, DDB replication, Redis status, Kafka lag, outbox backlog, quotas/capacity, alert routing, contacts, immutable logs and open actions. Rotate focus weekly; record owner, date, result and ticket. Do not mark a metric green without evidence.

Quarterly failover gates: single writer proof, RPO under 60s observed, safe transaction handling, full Hyderabad capacity, recovery within five minutes, rail reconciliation, locality proof, 24-hour monitoring, rollback available. Any missed objective creates a severity-rated remediation with owner/date and retest.
