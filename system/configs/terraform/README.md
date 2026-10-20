# Terraform baseline

This minimal Terraform starter creates two India-region VPCs and private subnet sets only. It does **not** create EKS, Aurora, DynamoDB, Redis, MSK, Route 53, security controls, or a production-ready network path; those require client CIDRs, topology, quotas, policies and cost review. The architecture docs specify service settings, replication and controls. Remote backend bootstrapping is intentionally not included; configure an approved encrypted S3 state backend and locking before any `init`.

Run only in a sandbox account after supplying approved CIDRs and reviewing provider availability:

```sh
terraform init
terraform fmt -check -recursive
terraform validate
terraform plan
```

Never apply to production without change approval. No AWS resource was provisioned for this project.