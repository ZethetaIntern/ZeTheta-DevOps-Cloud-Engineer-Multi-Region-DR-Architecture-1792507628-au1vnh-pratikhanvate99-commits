# DNS and traffic failover

Use Route 53 health checks against region-local synthetic endpoints that validate API health, database writer role, queue, and rail-adapter readiness without creating financial transactions. Configure a failover alias for primary and secondary. TTL 30 seconds is a planning value; resolver caching and client connection pools mean actual convergence can exceed TTL. Keep the Hyderabad endpoint warm and fully capacity-tested. Use staged weighted routing for planned events; region-failure action requires manual declaration and writer fencing. Never infer database promotion from DNS health.

Replace example domains, hosted-zone IDs, target IDs and health-check IDs. Route 53 health check addresses must be publicly reachable; for private APIs use CloudWatch alarms and calculated health checks / controlled automation. Test DNS from multiple resolvers and merchant networks.
