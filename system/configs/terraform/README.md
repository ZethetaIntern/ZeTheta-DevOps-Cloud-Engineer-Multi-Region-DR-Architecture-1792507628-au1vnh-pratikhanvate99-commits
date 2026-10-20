# Terraform baseline

This minimal Terraform starter creates two India-region VPCs and private subnet sets only. It does **not** create EKS, Aurora, DynamoDB, Redis, MSK, Route 53, security controls, or a production-ready network path; those require client CIDRs, topology, quotas, policies and cost review. The architecture docs specify service settings, replication and controls. The required production Terraform modules for compute, databases, DynamoDB, cache, messaging, DNS, monitoring and security remain to be implemented and wired after account discovery. Nothing has been applied to AWS.

## Remote state bootstrap

The root configuration does not embed a backend block because bucket ownership and names are client-specific. Create the state bucket and controls through a separately reviewed bootstrap, copy `backend.hcl.example` to a protected local file, replace the bucket value, and initialize with `terraform init -backend-config=<protected-file>`. Use an India-region bucket, versioning, encryption, access logging, least-privilege roles and state locking. Never commit the populated backend file or Terraform state. This example uses S3 lockfiles and requires Terraform 1.10+ for `use_lockfile`; Terraform 1.5 is the project syntax baseline, not the remote-locking feature baseline.

Run only in a sandbox account after supplying approved CIDRs and reviewing provider availability:

```sh
terraform init
terraform fmt -check -recursive
terraform validate
terraform plan
```

Never apply to production without change approval. No AWS resource was provisioned for this project.
