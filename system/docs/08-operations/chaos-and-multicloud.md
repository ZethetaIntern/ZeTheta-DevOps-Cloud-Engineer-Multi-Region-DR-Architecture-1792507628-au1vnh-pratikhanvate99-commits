# Chaos engineering and multi-cloud option

## Controlled experiments

Start in isolated non-production with synthetic payments and explicit blast-radius limits. Experiments: AZ workload eviction; Aurora replica lag observation; Kafka broker interruption; cross-region link loss; DNS resolver cache delay; Redis promotion; capacity surge; third-party timeout; cert rotation. Each experiment defines hypothesis, abort condition, owner, telemetry, restore path and evidence. Abort on any real payment traffic, unexpected data exposure, duplicate/lost synthetic transaction, or unplanned customer impact. Never inject faults into NPCI/production endpoints without written authorization.

## Multi-cloud position

Keep a portable container/application and PostgreSQL logical export strategy as a later resilience option, not the primary <5 minute recovery path. A second cloud adds network, identity, key custody, Kafka semantics, data residency, operational and certification complexity. Test an offline immutable backup restore and periodic cross-cloud recovery only after India locality and PCI scope are proven. Do not claim cross-cloud active-active ledger consistency. The two AWS India regions remain the designed recovery pair.
