# Authentication Failure Runbook

## Overview

This runbook covers JWT token validation and authorization failures in the Stellabill backend. The service uses JWT-based authentication with tenant isolation via the `X-Tenant-ID` header.

## Detection

### Automated Alerts
- Increased 401/403 HTTP status codes
- Auth middleware error rate spikes
- Failed login attempts (if applicable)
- Tenant mismatch errors in logs

### Manual Detection
- API responses showing auth-related errors:
  - "authorization header required"
  - "invalid or expired token"
  - "tenant mismatch"
  - "token missing subject claim"
- Application logs showing JWT parsing failures
- Monitoring dashboard showing auth failure metrics

### Impact Assessment
- **Low**: Individual user requests failing, other users unaffected
- **Medium**: Multiple users in same tenant affected
- **High**: Entire tenant unable to access service
- **Critical**: Widespread auth failures across multiple tenants

## Mitigation

### Immediate Actions (5-10 minutes)

1. **Check JWT Secret Configuration**
   ```bash
   # Verify JWT_SECRET is set and not expired
   echo $JWT_SECRET | wc -c  # Should be > 32 characters

   # Check if secret was recently rotated
   grep "JWT_SECRET" /var/log/stellarbill/backend.log | tail -10
   ```

2. **Validate Token Format**
   - Check if tokens are properly formatted Bearer tokens
   - Verify token expiration times are reasonable
   - Confirm tenant claims match X-Tenant-ID headers

3. **Check Auth Middleware Logs**
   ```bash
   # Look for patterns in auth failures
   grep "authorization header required" /var/log/stellarbill/backend.log | tail -20
   grep "tenant mismatch" /var/log/stellarbill/backend.log | tail -20
   ```

### Short-term Recovery (15-30 minutes)

1. **Token Validation Bypass** (temporary, monitored)
   - For critical endpoints, implement temporary auth bypass
   - Log all bypassed requests for security review
   - Set up monitoring for abuse

2. **JWT Secret Rotation Check**
   - If secret was rotated, ensure all services have new secret
   - Check if old tokens are still being accepted during transition
   - Verify token issuer and audience claims

3. **Tenant Configuration Review**
   - Validate tenant IDs are properly configured
   - Check for tenant header validation issues
   - Review cross-tenant access policies

## Recovery

### Full Service Restoration

1. **Verify Auth Configuration**
   ```bash
   # Test JWT validation
   curl -H "Authorization: Bearer <valid-test-token>" \
        -H "X-Tenant-ID: test-tenant" \
        http://localhost:8080/api/health

   # Test tenant isolation
   curl -H "Authorization: Bearer <valid-test-token>" \
        -H "X-Tenant-ID: wrong-tenant" \
        http://localhost:8080/api/health  # Should fail
   ```

2. **Clear Auth Caches** (if applicable)
   - Restart application to clear any cached auth state
   - Verify middleware is loading current configuration

3. **Update Client Tokens**
   - If tokens were invalidated, coordinate with frontend teams
   - Provide new tokens or refresh mechanisms
   - Communicate expected resolution time

4. **Monitor Recovery**
   - Watch 401/403 rates return to baseline
   - Verify tenant-scoped operations working
   - Check subscription access by correct callers

### Rollback Plan
- Revert to previous JWT secret if rotation caused issues
- Restore previous auth middleware version if code changes introduced bugs
- Implement gradual token migration strategy

## Observability

### Key Metrics
- Auth success rate: `auth_requests_total{status="success"}`
- Token validation errors: `jwt_validation_errors_total`
- Tenant mismatch count: `tenant_mismatch_errors_total`
- HTTP 401/403 rates: `http_requests_total{status=~"401|403"}`

### Dashboards
- [Authentication Dashboard](https://monitoring.example.com/auth-dashboard)
- [API Error Dashboard](https://monitoring.example.com/api-errors)
- [Tenant Isolation Dashboard](https://monitoring.example.com/tenant-isolation)

### Logs
- Application logs: `/var/log/stellarbill/backend.log`
- Auth middleware logs: Filter by "AuthMiddleware" component
- JWT parsing errors: Search for "invalid token" or "signature invalid"

## Post-Incident Review Checklist

### Technical Analysis
- [ ] JWT secret rotation process reviewed
- [ ] Token expiration policies validated
- [ ] Auth middleware code audited for bugs
- [ ] Tenant configuration verified

### Process Improvements
- [ ] Token rotation communication improved
- [ ] Auth failure alerting thresholds adjusted
- [ ] Client token refresh mechanisms reviewed
- [ ] Security monitoring enhanced

### Prevention Measures
- [ ] Automated token validation testing implemented
- [ ] JWT secret rotation procedures documented
- [ ] Auth failure simulation tests added
- [ ] Multi-tenant auth testing expanded

### Communication
- [ ] Incident timeline documented
- [ ] Affected users/tenants notified
- [ ] Token refresh instructions provided
- [ ] Retrospective meeting scheduled

## Security Notes

- JWT secrets never logged or exposed in error messages
- Failed auth attempts logged with IP and user agent for security monitoring
- Tenant isolation prevents unauthorized cross-tenant access
- Token blacklisting implemented for compromised credentials
- All auth operations use constant-time comparisons to prevent timing attacks

## Related Documentation

- [Authentication Middleware](../../internal/middleware/auth.go)
- [JWT Configuration](../../internal/config/config.go)
- [Multi-tenant Security](../../README.md#multi-tenant-security)
- [Security Guidelines](../../internal/worker/SECURITY.md)</content>
<parameter name="filePath">/workspaces/stellabill-backend/docs/ops/auth-failure-runbook.md