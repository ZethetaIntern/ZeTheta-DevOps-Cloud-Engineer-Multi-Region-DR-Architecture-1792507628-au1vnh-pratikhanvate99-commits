# Source brief: five high-risk claims corrected

The 68-page assignment asks for five deliberate errors but does not identify them. These are five independently checkable material errors or unsupported claims; author intent cannot be established.

1. **Business-hour average TPS:** 3.2m / (16×3600) = 55.56 TPS, not the brief’s 250 TPS. The 250 figure needs a different explicit measurement window.
2. **Universal RBI recovery target:** the brief attributes a universal 4-hour RTO / near-zero RPO to RBI. Applicability depends on regulated entity and direction; no universal payment-aggregator 4-hour rule is evidenced here. Confirm current entity-specific directions with counsel.
3. **UPI failover / latency claim:** 99.95% availability, automatic five-minute failover and a universal 300ms transaction ceiling are not established as general NPCI UPI public requirements in the cited public circular index. Obtain participant/sponsor-bank technical specs and contracts.
4. **PCI 72-hour incident reporting:** PCI DSS 4.0.1 requires an incident response plan and notification obligations as applicable; it does not create a blanket 72-hour deadline for every incident. CERT-In directions impose a separate six-hour reporting obligation for covered incidents/entities; counsel determines scope.
5. **S3 replication SLA conflation:** a 15-minute replication SLA belongs to S3 Replication Time Control for eligible configurations, not ordinary CRR by default. Standard CRR is asynchronous; monitor replication status and select RTC if that SLA is required.

Other arithmetic cross-checks: 99.92% yields 7.01 hours/year downtime and 99.99% yields 52.56 minutes/year. INR 500 crore/day at an assumed 1.8% take rate implies INR 9 crore/day gross fee estimate; this is a scenario calculation, not known company revenue.
