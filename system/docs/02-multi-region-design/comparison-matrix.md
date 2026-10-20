# DR pattern comparison matrix

The numbers below are planning values tied to the fictional brief and must be validated in a workload benchmark and AWS Pricing Calculator. “Active-active” means both regions serve traffic with a single-home writer per merchant aggregate; uncontrolled dual-write is rejected.

| Dimension | Active-active | Hot active-passive | Basis / decision |
| --- | --- | --- | --- |
| RTO | Edge routing can be near-zero for stateless paths; payment writes still need ownership fencing and can take 1-5 min | 1-5 min objective; elapsed DB/service/partner path must be proven | Brief tier table says hot standby 1-5 min; DNS alone is not full transaction RTO. |
| RPO | Seconds or less for a strong owner protocol; async replicas can still lag | Under 60 sec target based on measured async replication lag | Neither mode has a universal zero-loss claim. |
| Cost | 1.8-2.2x; INR 14.4-17.6 crore/year using INR 8cr baseline | 1.6-1.8x; INR 12.8-14.4 crore/year | Brief multipliers; workbook line items are editable. |
| Complexity | Very high: merchant ownership, epochs, global idempotency, split-brain and settlement merge | High: promotion, fencing, DNS/client tail, stream replay and reconciliation | Eight-person infra team favors one writer / hot standby. |
| Data consistency | Strong per merchant only with home-region serialization; cross-region reads may be stale | One writer; secondary is async and stale until promotion | Aurora/MSK/Redis cross-region copies asynchronous. |
| Latency impact | Local reads and stateless requests can improve; remote owner writes add at least an RTT (15-25 ms planning band from brief) | No cross-region synchronous hop on normal Mumbai write path | Measure actual p99; 180ms baseline leaves 120ms to unverified 300ms ceiling. |
| Operational burden | Two active traffic planes, traffic fairness, skew management, owner transfer and conflict remediation | Two maintained regions, drills, explicit promotion and failback | Hot standby reduces dual-write operations but not 24x7 readiness. |
| Regulatory approval | More evidence for split-brain prevention, locality and partner certification | Simpler single-authority narrative; still needs BCP, testing and audit evidence | Obtain entity-specific RBI/NPCI/QSA confirmation. |
| Peak capacity | Both regions share normal load; either region must support full peak 1,200 TPS + 30% headroom target | Hyderabad must be provisioned to full 1,560 TPS test target while idle | Autoscaling is not counted as immediate failover capacity. |
| Integrity blast radius | Potentially both regions and both write sets during partition | Replication can copy corruption, but one authoritative writer narrows conflict scope | Isolated immutable backups and clean restore needed in both. |

