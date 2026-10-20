# India data sovereignty and compliance mapping

This is a design mapping, not legal advice, an audit opinion or certification. Applicability, current amendments, payment-system role, card scope, and sponsor-bank/NPCI obligations must be confirmed by PaySecure counsel and assessor before production.

| Decision | RBI payment data localisation | RBI PA/PG / security controls | PCI DSS 4.0.1 | NPCI UPI | India localisation |
|---|---|---|---|---|---|
| Mumbai + Hyderabad only for payment data | Store full end-to-end payment system data in India | Availability, access and audit controls | Scope every component that stores/processes/transmits account data | Confirm participant-specific geography/registration | No foreign payment-data replica |
| Single Aurora financial writer | Integrity, auditable records | Reconciliation and access controls | Restrict access; log/audit; secure configuration | Rail state and reversals per participant rules | Ledger backups remain India |
| DDB metadata / Redis cache / Kafka streams | Treat derived transaction data as payment data unless counsel says otherwise | Retention, security, resiliency | Minimize scope and sensitive authentication data | Topic/schema, timeout/retry and idempotency per specs | Keep streams, snapshots, logs in India |
| Encryption and key custody | Protect data in storage and transit | Key/access governance | Requirement 3/4/7/8/10 as applicable | HSM/certificate controls per integration | India-region KMS/HSM control path validated |
| Telemetry and support | Logs can contain payment data; control location/retention | Audit trails | Protect logs and incident response | Incident escalation contract | Inspect vendor support/telemetry paths |
| DR and incident response | Maintain availability and retrieval | BCP, security incident processes | Req 12.10 incident-response plan | Notify rail/sponsor bank per circular/contract | CERT-In scope and timeline assessed |

RBI’s 2018 circular requires full end-to-end payment system data storage in India, with limited foreign-leg handling and return within the prescribed timeline; RBI FAQs explain applicability and offshore processing boundaries. This architecture avoids offshore copies entirely. PCI DSS does not permit storing sensitive authentication data after authorization. Do not accept/store PAN, CVV, PIN or UPI PIN in this demo. CERT-In directions may require covered entities to report specified incidents within six hours; verify the entity’s status and current rules. NPCI detailed standards and sponsor-bank contracts may be restricted; request current versions and map each API behavior before implementation.

## Evidence references

- RBI payment data circular: https://rbi.org.in/Scripts/NotificationUser.aspx?Id=11244
- RBI FAQ: https://systemhealth.rbi.org.in/Scripts/FAQDisplay.aspx_Id=130(1).html
- RBI PA/PG direction: https://www.rbi.org.in/Scripts/NotificationUser.aspx?Id=11822
- RBI digital payment security controls: https://www.rbi.org.in/Scripts/BS_ViewMasDirections.aspx?id=12032
- PCI SSC DSS v4.0.1: https://blog.pcisecuritystandards.org/just-published-pci-dss-v4-0-1
- CERT-In directions: https://www.cert-in.org.in/Directions70B.jsp
- NPCI UPI circular index (obtain participant documents): https://www.npci.org.in/circulars/upi

## PCI scope note

A payment aggregator may have obligations based on actual role and card data flows. Segment CDE, tokenize through approved provider where applicable, never retain CVV/PIN, inventory service providers, validate segmentation and obtain QSA guidance. This architecture alone does not establish PCI compliance.
