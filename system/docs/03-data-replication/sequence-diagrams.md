# Replication and cutover sequence diagrams

## Normal write and asynchronous replication

```mermaid
sequenceDiagram
 participant M as Merchant
 participant A as API Mumbai
 participant DB as Aurora Mumbai writer
 participant O as Outbox publisher
 participant K as MSK Mumbai
 participant R as MSK Replicator
 participant H as Hyderabad consumers
 M->>A: POST with idempotency key
 A->>DB: transaction: payment + outbox
 DB-->>A: commit
 A-->>M: accepted, stable payment ID
 O->>K: publish event ID
 K-->>O: ack then mark outbox published
 R-->>H: async replicated event (at least once)
 H->>H: dedupe event ID before side effect
```

## Region failover

```mermaid
sequenceDiagram
 participant IC as Incident commander
 participant DNS as Route 53
 participant M as Mumbai
 participant H as Hyderabad
 participant DB as Aurora Global DB
 IC->>M: fence writes and workers
 IC->>DB: inspect lag and last commit evidence
 IC->>DB: promote Hyderabad after fencing
 DB-->>IC: writer endpoint active
 IC->>H: enable consumers, validate idempotency
 IC->>DNS: shift weighted/failover records
 DNS-->>H: traffic ramps under health gates
 IC->>H: reconcile uncertain payment window
```

## Aurora PostgreSQL cross-region replication and promotion

```mermaid
sequenceDiagram
 participant API as Payment API Mumbai
 participant P as Aurora Mumbai writer
 participant G as Aurora Global Database
 participant S as Aurora Hyderabad secondary
 participant IC as Incident Commander
 API->>P: Commit payment + outbox in one transaction
 P->>G: Asynchronous storage replication
 G-->>S: Replicate cluster storage changes
 IC->>P: Fence primary write paths
 IC->>S: Inspect lag, writer role, and last durable commit
 IC->>S: Promote only after primary fencing is proven
 S-->>IC: Hyderabad becomes writable primary
 IC->>API: Switch database endpoint and run canary
```

## DynamoDB Global Tables replication

```mermaid
sequenceDiagram
 participant App as Mumbai service
 participant M as DynamoDB Mumbai replica
 participant GT as Global Tables replication
 participant H as DynamoDB Hyderabad replica
 participant DR as Hyderabad service
 App->>M: Conditional put/update with idempotency key
 M-->>App: Local write acknowledgement
 M->>GT: Replicate item asynchronously
 GT-->>H: Apply replicated item (eventual consistency)
 DR->>H: Read local replica after regional cutover
 DR->>H: Conditional write using deterministic merchant ownership
```

Use Global Tables only for operational metadata in this design. Concurrent same-item writes can resolve with last-writer-wins semantics; never use the table as the authoritative payment ledger or assume a write acknowledged in one region is immediately readable in another.

## ElastiCache Redis Global Datastore recovery

```mermaid
sequenceDiagram
 participant API as Mumbai cache client
 participant P as Redis primary
 participant GD as Global Datastore replication
 participant S as Hyderabad secondary
 participant IC as Incident Commander
 API->>P: Cache lookup or rebuildable session data
 P->>GD: Asynchronous replication
 GD-->>S: Replicate cache state
 IC->>S: Check replication lag and primary isolation
 IC->>S: Promote secondary after database authority decision
 S-->>IC: Hyderabad cache becomes writable
 IC->>API: Reconnect; tolerate misses and rebuild from durable stores
```

Redis is a cache/session optimization only. Never use it as the sole copy of balance, idempotency, payment status, or transaction ownership. Cache loss is recovered from durable systems; expired or absent keys are normal during promotion.

## Kafka MSK Replicator recovery

```mermaid
sequenceDiagram
 participant P as Mumbai producer/outbox
 participant K as MSK Mumbai
 participant R as MSK Replicator
 participant H as MSK Hyderabad
 participant C as Hyderabad consumer
 P->>K: Publish keyed event with stable event_id
 K-->>P: Acknowledge source record
 K->>R: Replicate topic record asynchronously
 R-->>H: Copy record (at least once)
 Note over H,C: Consumer group offsets are independent across clusters
 C->>H: Read from approved recovery offset
 C->>C: Deduplicate event_id in durable inbox
 C-->>H: Commit offset only after idempotent side effect
```

Replicator does not provide atomic cross-region producer acknowledgement or consumer offset transfer. Define the replay point from durable outbox/inbox/event IDs and reconcile the ambiguous window. Never skip an offset merely to make lag appear green.

## S3 cross-region replication

```mermaid
sequenceDiagram
 participant App as Mumbai service
 participant B1 as India source bucket
 participant CRR as S3 replication
 participant B2 as India destination bucket
 participant IC as Recovery operator
 App->>B1: Write encrypted, versioned object
 B1-->>App: Source write acknowledged
 B1->>CRR: Queue object version for replication
 CRR-->>B2: Replicate object and metadata asynchronously
 IC->>B2: Verify replication status, version, checksum, and KMS access
 IC->>B2: Restore only verified object versions into clean recovery path
```

Ordinary S3 CRR is asynchronous. If a contractual replication-time objective is required, validate eligibility and configure S3 Replication Time Control separately; do not describe RTC timing as a default CRR guarantee. Keep both buckets and KMS keys in approved India regions.

## Split-brain authority and fencing

```mermaid
sequenceDiagram
 participant M as Mumbai writer
 participant H as Hyderabad secondary
 participant W as Independent fencing witness
 participant IC as Incident Commander
 IC->>M: Observe writer role and commit position
 IC->>H: Observe role, lag, and local write state
 Note over M,H: Inter-region partition; both local health checks may pass
 IC->>W: Verify out-of-band authority and fencing path
 alt Mumbai can be conclusively fenced
  IC->>M: Fence writes and credentials
  IC->>H: Promote only after fence verification
 else Fence cannot be proven
  IC->>M: Keep writes disabled/queued where reachable
  IC->>H: Do not promote; escalate and preserve single-writer safety
 end
```

DynamoDB MREC, Redis Global Datastore, Kafka Replicator, and S3 CRR all copy asynchronously and independently. There is no distributed atomic commit across these services. Treat each as a recoverable projection; reconcile from authoritative Aurora plus rail settlement records. Set measured lag SLOs per stream and alert before 60 seconds; do not assume service default replication timing guarantees the business RPO.
