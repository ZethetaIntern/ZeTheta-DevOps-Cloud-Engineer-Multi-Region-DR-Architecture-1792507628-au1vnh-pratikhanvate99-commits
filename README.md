# ZeTheta DevOps & Cloud Engineer — Multi-Region DR Architecture for Payment Systems

Client review package for the fictional PaySecure Gateway Private Limited scenario. It includes architecture designs and a local simulator; production deployment and regulatory acceptance are not complete. The assignment brief is treated as the specification; factual and legal claims in the source have been independently qualified in the project notes.

## Deliverables

- [Client architecture overview PDF](deliverables/PaySecure_Multi_Region_DR_Architecture_Client_Design.pdf)
- [Implementation starter and operations package](system/README.md)
- [Requirements traceability](system/docs/00-project/requirements-traceability.md)
- [Client-readiness and brief coverage audit](system/docs/00-project/client-readiness-audit.md)
- [Architecture and active-active comparison](system/docs/02-multi-region-design/comparison-matrix.md)
- [Replication strategy](system/docs/03-data-replication/replication-strategy.md)
- [12 disaster recovery runbooks](system/docs/05-runbooks/)
- [Cost model workbook](system/docs/06-risk-cost/cost-model.xlsx)
- [India locality and compliance mapping](system/docs/07-compliance/data-sovereignty.md)
- [BCRB simulation](system/docs/09-review/bcrb.md)

## Important status

The architecture recommends a hot active-passive Mumbai/Hyderabad design with a single financial ledger writer. The scenario targets of 99.99% availability, RPO under 60 seconds and RTO under five minutes have **not** been proven by production drills. Q3 2026 has passed as of 20 October 2026; the client should approve a re-baselined program schedule. No live AWS resources, NPCI connectivity, card processing, production failover or performance test were executed. Terraform and Kubernetes files are starter templates requiring account-specific review and completion.

The supplied source brief contains several arithmetic and regulatory/technical claims that should not be carried forward unqualified. See `system/docs/10-source-review/` and the requirement traceability matrix.

The repository is suitable to share as a **design package and local demonstrator**. It is not a production-ready DR deployment or compliance attestation. See the [client-readiness audit](system/docs/00-project/client-readiness-audit.md) for coverage, evidence gaps, and the recommended client-facing GitHub description. The GitHub About field itself still needs an owner login to update.
