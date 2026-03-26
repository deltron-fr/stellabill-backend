# Elevated Error Rates Runbook

## Overview

This runbook covers increased error rates across the Stellabill backend API. The service includes health checks, subscription management, plan listings, and background billing operations.

## Detection

### Automated Alerts
- HTTP 5xx error rate exceeds threshold (e.g., >5% of requests)
- Health check endpoint failures: `GET /api/health` returns non-200
- Background worker job failure rate spikes
- API latency increases beyond acceptable limits

### Manual Detection
- API responses showing application errors:
  - "internal error"
  - "subscription not found"
  - "forbidden"
  - "billing parse error"
- Application logs showing panic traces or unhandled errors
- Monitoring dashboard showing error rate trends

### Impact Assessment
- **Low**: Isolated endpoint errors, service mostly operational
- **Medium**: Multiple endpoints affected, degraded user experience
- **High**: Core functionality impaired, billing operations failing
- **Critical**: Service unavailable, widespread customer impact

## Mitigation

### Immediate Actions (5-10 minutes)

1. **Check Application Health**
   ```bash
   # Health check
   curl http://localhost:8080/api/health

   # Application process status
   systemctl status stellarbill-backend
   ```

2. **Review Recent Logs**
   ```bash
   # Check for error patterns
   grep "ERROR" /var/log/stellarbill/backend.log | tail -20

   # Look for panic traces
   grep "panic" /var/log/stellarbill/backend.log | tail -10
   ```

3. **Circuit Breaker Check**
   - Verify if circuit breaker is tripping for external services
   - Check HTTP client error rates

4. **Resource Usage Check**
   ```bash
   # Memory and CPU usage
   top -p $(pgrep stellarbill-backend)

   # Disk space
   df -h

   # Network connectivity
   ping -c 3 database-host
   ```

### Short-term Recovery (15-30 minutes)

1. **Application Restart**
   ```bash
   # Graceful restart
   systemctl restart stellarbill-backend

   # Check startup logs
   journalctl -u stellarbill-backend -n 20
   ```

2. **Load Balancer Check**
   - Verify load balancer health checks passing
   - Check if instances are being marked unhealthy
   - Review load distribution across instances

3. **Dependency Health**
   - Test database connectivity
   - Check external service availability
   - Verify network connectivity to dependencies

4. **Rate Limiting Review**
   - Check if rate limits are being hit
   - Review request patterns for abuse
   - Adjust limits if configured too conservatively

## Recovery

### Full Service Restoration

1. **Verify All Endpoints**
   ```bash
   # Test all API endpoints
   curl -H "Authorization: Bearer <test-token>" \
        -H "X-Tenant-ID: test-tenant" \
        http://localhost:8080/api/health

   curl -H "Authorization: Bearer <test-token>" \
        -H "X-Tenant-ID: test-tenant" \
        http://localhost:8080/api/plans

   curl -H "Authorization: Bearer <test-token>" \
        -H "X-Tenant-ID: test-tenant" \
        http://localhost:8080/api/subscriptions
   ```

2. **Background Worker Verification**
   ```bash
   # Check worker status
   systemctl status stellarbill-worker

   # Verify job processing
   curl http://localhost:8080/internal/worker/status  # If available
   ```

3. **Performance Monitoring**
   - Watch latency return to baseline
   - Monitor error rates decreasing
   - Verify resource usage stabilizing

4. **Load Testing** (if needed)
   - Run basic load test to verify stability
   - Monitor for memory leaks or performance degradation

### Rollback Plan
- Revert to previous application version if code changes caused issues
- Restore from backup if data corruption suspected
- Implement feature flags to disable problematic functionality

## Observability

### Key Metrics
- HTTP error rates: `http_requests_total{status=~"4..|5.."}`
- API latency: `http_request_duration_seconds`
- Worker job failures: `worker_jobs_failed_total`
- Resource usage: `process_cpu_usage`, `process_memory_usage`

### Dashboards
- [API Performance Dashboard](https://monitoring.example.com/api-performance)
- [Error Rate Dashboard](https://monitoring.example.com/error-rates)
- [Worker Health Dashboard](https://monitoring.example.com/worker-health)
- [Infrastructure Dashboard](https://monitoring.example.com/infrastructure)

### Logs
- Application logs: `/var/log/stellarbill/backend.log`
- Worker logs: `/var/log/stellarbill/worker.log`
- System logs: `/var/log/syslog`
- Access logs: `/var/log/stellarbill/access.log`

## Post-Incident Review Checklist

### Technical Analysis
- [ ] Error logs analyzed for root cause patterns
- [ ] Code review of error-handling paths
- [ ] Performance profiling conducted
- [ ] Resource limits and configurations reviewed

### Process Improvements
- [ ] Alert thresholds calibrated based on incident
- [ ] Monitoring dashboards enhanced
- [ ] Error tracking and reporting improved
- [ ] Automated testing expanded

### Prevention Measures
- [ ] Circuit breakers implemented for external dependencies
- [ ] Graceful degradation strategies added
- [ ] Load testing automated in CI/CD
- [ ] Error budget monitoring implemented

### Communication
- [ ] Incident timeline documented
- [ ] Customer impact quantified and communicated
- [ ] Service level objectives reviewed
- [ ] Incident report published internally

## Security Notes

- Error messages sanitized to prevent information leakage
- Sensitive data never included in error responses
- Authentication failures logged securely without exposing credentials
- Rate limiting prevents abuse during error conditions
- Audit logs maintained for security incident investigation

## Related Documentation

- [Error Handling](../../internal/service/errors.go)
- [Health Check Implementation](../../internal/handlers/health.go)
- [Circuit Breaker](../../internal/httpclient/circuitbreaker.go)
- [Worker Error Handling](../../internal/worker/executor.go)</content>
<parameter name="filePath">/workspaces/stellabill-backend/docs/ops/elevated-errors-runbook.md