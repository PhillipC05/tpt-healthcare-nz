-- Maps inbound HL7 lab site codes (MSH-4 Sending Facility / ZNZL LabSite)
-- to the tenant that owns the results, so MLLP-originated pathology results
-- are scoped to the correct tenant instead of a shared default.
CREATE TABLE IF NOT EXISTS tenant_lab_sites (
    lab_site   TEXT PRIMARY KEY,
    tenant_id  UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tenant_lab_sites_tenant ON tenant_lab_sites(tenant_id);
