# Database Outage Runbook

## Overview

This runbook covers database connectivity and availability issues in the Stellabill backend. The service uses PostgreSQL with a repository layer abstraction for data access.

## Detection

### Automated Alerts
- Health check failures: `GET /api/health` returns non-200 status
- Database connection pool exhaustion
- Increased error rates on endpoints: `/api/subscriptions`, `/api/plans`
- Worker job failures (billing operations stuck in pending state)

### Manual Detection
- API responses showing database-related errors:
  - "connection refused"
  - "connection timeout"
  - "server closed the connection unexpectedly"
- Background worker logs showing DB connection failures
- Monitoring dashboard showing DB connection metrics

### Impact Assessment
- **Low**: Single read queries failing, service remains partially operational
- **Medium**: Write operations failing, subscription modifications blocked
- **High**: Complete DB unavailability, all data-dependent endpoints failing
- **Critical**: Background billing jobs failing, potential revenue impact

## Mitigation

### Immediate Actions (5-10 minutes)

1. **Check Database Status**
   ```bash
   # Check if PostgreSQL is running
   systemctl status postgresql

   # Check database connectivity
   psql -h localhost -U stellarbill -d stellarbill -c "SELECT 1;"
   ```

2. **Restart Application** (if connection pool issues)
   ```bash
   # Graceful restart to clear connection pools
   systemctl restart stellarbill-backend
   ```

3. **Enable Read-Only Mode** (if partial DB access)
   - Set environment variable: `DB_READONLY=true`
   - This allows read operations while blocking writes

### Short-term Recovery (15-30 minutes)

1. **Database Restart**
   ```bash
   # If PostgreSQL is down
   systemctl restart postgresql

   # Check logs for startup issues
   journalctl -u postgresql -n 50
   ```

2. **Connection Pool Reset**
   - Application automatically handles pool recovery on restart
   - Monitor connection pool metrics post-restart

3. **Failover Check** (if using replication)
   - Verify primary database is healthy
   - Check replica lag if applicable

## Recovery

### Full Service Restoration

1. **Verify Database Health**
   ```sql
   -- Run basic health checks
   SELECT version();
   SELECT COUNT(*) FROM subscriptions;
   SELECT COUNT(*) FROM plans;
   ```

2. **Test API Endpoints**
   ```bash
   # Test health endpoint
   curl -H "Authorization: Bearer <test-token>" \
        -H "X-Tenant-ID: test-tenant" \
        http://localhost:8080/api/health

   # Test data endpoints
   curl -H "Authorization: Bearer <test-token>" \
        -H "X-Tenant-ID: test-tenant" \
        http://localhost:8080/api/subscriptions
   ```

3. **Restart Background Workers**
   ```bash
   # Ensure billing workers are processing jobs
   systemctl restart stellarbill-worker
   ```

4. **Monitor Recovery**
   - Watch error rates return to baseline
   - Verify billing job queue is processing
   - Check tenant isolation is working

### Rollback Plan
- If database corruption suspected, restore from backup
- Coordinate with DBA team for point-in-time recovery
- Notify affected tenants of data restoration timeline

## Observability

### Key Metrics
- Database connection count: `db_connections_active`
- Query latency: `db_query_duration_seconds`
- Error rate by endpoint: `http_requests_total{status=~"5.."}`
- Worker job success rate: `worker_jobs_completed_total`

### Dashboards
- [Database Performance Dashboard](https://monitoring.example.com/db-performance)
- [API Health Dashboard](https://monitoring.example.com/api-health)
- [Worker Jobs Dashboard](https://monitoring.example.com/worker-jobs)

### Logs
- Application logs: `/var/log/stellarbill/backend.log`
- Database logs: `/var/log/postgresql/postgresql-*.log`
- Worker logs: `/var/log/stellarbill/worker.log`

## Post-Incident Review Checklist

### Technical Analysis
- [ ] Root cause identified (network, disk space, configuration, etc.)
- [ ] Database logs reviewed for error patterns
- [ ] Connection pool settings validated
- [ ] Backup/recovery procedures tested

### Process Improvements
- [ ] Alert thresholds reviewed and adjusted
- [ ] Monitoring coverage gaps identified
- [ ] Runbook accuracy verified
- [ ] On-call rotation feedback collected

### Prevention Measures
- [ ] Database maintenance windows scheduled
- [ ] Connection pool monitoring enhanced
- [ ] Automated failover testing implemented
- [ ] Capacity planning updated based on incident

### Communication
- [ ] Incident timeline documented
- [ ] Customer impact assessed and communicated
- [ ] Follow-up actions assigned with owners
- [ ] Retrospective meeting scheduled

## Security Notes

- Database credentials never logged in application logs
- Direct database access limited to authorized personnel only
- All database queries use parameterized statements to prevent SQL injection
- Tenant isolation enforced at repository layer prevents cross-tenant data access

## Related Documentation

- [Database Configuration](../config/database.md)
- [Health Check Implementation](../../internal/handlers/health.go)
- [Repository Layer](../../internal/repository/)</content>
<parameter name="filePath">/workspaces/stellabill-backend/docs/ops/db-outage-runbook.md