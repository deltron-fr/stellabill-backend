# Operational Runbooks

This directory contains operational runbooks for common backend incidents in the Stellabill backend service.

## Runbooks

- [Database Outages](./db-outage-runbook.md) - Handling database connectivity and availability issues
- [Authentication Failures](./auth-failure-runbook.md) - Managing JWT token and authorization problems
- [Elevated Error Rates](./elevated-errors-runbook.md) - Responding to increased error rates across the API

## Overview

These runbooks are designed for the current Go/Gin API architecture with the following components:

- **Authentication**: JWT-based with tenant isolation via `X-Tenant-ID` header
- **Database**: PostgreSQL with repository layer abstraction
- **Monitoring**: Health check endpoint at `/api/health`
- **Background Jobs**: Worker system for billing operations
- **Multi-tenancy**: Tenant-scoped data access with caller ownership checks

## General Response Process

1. **Detection** - Monitor alerts, logs, and health checks
2. **Assessment** - Gather context and impact analysis
3. **Mitigation** - Apply immediate fixes to restore service
4. **Recovery** - Full restoration and verification
5. **Post-Incident Review** - Analysis and prevention measures

## Contact Information

- **On-call Engineer**: Check current rotation in team communication channels
- **Team Lead**: Contact for escalation decisions
- **DevOps/SRE**: For infrastructure-related issues

## Tools and Access

- **Logs**: Application logs via logging service
- **Metrics**: Monitoring dashboards (links in individual runbooks)
- **Database**: Direct access for read-only queries (emergency access only)
- **Deployments**: CI/CD pipeline access for hotfixes</content>
<parameter name="filePath">/workspaces/stellabill-backend/docs/ops/README.md