# Monitoring and alert catalog

| Signal | Warning | Critical / action |
|---|---|---|
| Aurora Global replication lag | >20s 3 min | >45s 1 min: freeze promotion readiness and page Data Lead |
| API success/5xx | 5xx >0.1% 5 min | >1% 1 min: IC page |
| P99 latency | >180ms 5 min | >250ms 2 min: shed noncritical load |
| Kafka replicator/consumer lag | >20s | >45s: page streaming owner; no offset skip |
| Outbox oldest event | >10s | >30s: page publisher owner |
| Redis replica lag | >20s | >45s: promote only after authority established |
| DNS synthetic | one probe failure | 3 regions failing: IC declaration review |
| Cert expiry | 30/14/7 days | <72h: security page |
| DDB replication conflict/error | baseline anomaly | sustained errors: metadata degraded mode |
| Capacity | >70% sustained | >85%: scale / traffic gates |

Dashboard panels: request rate/errors/latency by region and merchant cohort; payment state transition counts; outbox age; Kafka broker/replicator/consumer lag; Aurora role/lag/connections; Dynamo errors; Redis lag/evictions; DNS health; fraud and rail adapter outcomes; active alerts and change markers. Label data sources and time zones. Page on actionable symptoms, not derived vanity metrics.
