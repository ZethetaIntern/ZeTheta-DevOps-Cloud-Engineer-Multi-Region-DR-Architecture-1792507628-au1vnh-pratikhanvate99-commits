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

DynamoDB MREC, Redis Global Datastore, Kafka Replicator, and S3 CRR all copy asynchronously and independently. There is no distributed atomic commit across these services. Treat each as a recoverable projection; reconcile from authoritative Aurora plus rail settlement records. Set measured lag SLOs per stream and alert before 60 seconds; do not assume service default replication timing guarantees the business RPO.
