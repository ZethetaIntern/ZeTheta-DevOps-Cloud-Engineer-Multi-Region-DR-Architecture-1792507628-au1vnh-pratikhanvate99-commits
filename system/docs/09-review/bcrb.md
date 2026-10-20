# Simulated Business Continuity Review Board

**This is a prepared simulation, not a real approval.** The five roles are: Chief Risk Officer (chair), CTO, Head of Payments, CISO, and CFO. Q3 2026 target date has passed as of this document’s creation (20 October 2026); re-baseline dates and report actual measured readiness.

## Recommendation

Approve a gated implementation plan for Mumbai primary / Hyderabad hot standby with one Aurora writer and India-only payment data. Do not declare 99.99%, <60-second RPO or <5-minute RTO achieved until witnessed drills meet objective with reconciled payment evidence and sustained capacity. Fund detailed discovery, nonproduction build and certification work; production cutover requires separate change approval.

## Board challenge questions and responses

1. **Can 5-minute RTO be guaranteed?** No provider feature guarantees end-to-end recovery; it is an acceptance objective. Exercise must measure detection, decision, promotion, application, DNS and rail readiness.
2. **Why not dual-writer active-active?** Ledger safety and uncertain UPI outcomes require a clear financial authority. Async replicas cannot resolve split-brain safely.
3. **What is the RPO under a partition?** Single-writer policy may reject writes; measured Aurora and stream lag define actual exposure. Unknown is not zero.
4. **Are all payment records in India?** Design keeps payment data, replicas, backups, telemetry and keys in India; vendor telemetry and support paths need validation.
5. **Will failover create duplicate UPI transactions?** Idempotency and status inquiry/reconciliation are required; no retry on unknown outcome until rail status is established.
6. **How much does it cost?** Workbook is a planning sensitivity based on current INR 8 crore baseline multipliers, not a cloud quote. Obtain AWS calculator/export and staffing quotes.
7. **How is ransomware recovery safe?** Isolated clean-room restore, immutable evidence, separate identities and security sign-off precede connectivity.
8. **What has been validated?** Local simulator/static artifacts only; no AWS deployment, NPCI integration, performance test or live failover has been executed.
9. **What does PCI cover?** Applicability depends on actual card flows, PA role, scope and assessor determination; this is not an attestation.
10. **What blocks production?** Customer inventory, contracts, quotas, capacity pricing, threat model, legal mapping, restore drills, regional failover evidence and board change approval.

## Decision log

Requested decision: authorize discovery/nonproduction implementation and resource estimates. Production approval: pending evidence. Owners: CTO for platform, Head of Payments for rail integration, CISO for security, CRO for risk acceptance, CFO for funding. Re-review after two successful quarterly full failovers.
