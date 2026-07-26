-- tpt-dental: ACC dental claims and treatment records.
-- Stored with indexed scalar columns for common filters plus a JSONB "body"
-- column holding the full record for round-trip fidelity.

CREATE TABLE IF NOT EXISTS dental_acc_claims (
    id           TEXT PRIMARY KEY,
    patient_nhi  TEXT NOT NULL,
    clinician_id TEXT NOT NULL,
    status       TEXT,
    body         JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_dental_claims_patient ON dental_acc_claims(patient_nhi);
CREATE INDEX IF NOT EXISTS idx_dental_claims_clinician ON dental_acc_claims(clinician_id);
CREATE INDEX IF NOT EXISTS idx_dental_claims_status ON dental_acc_claims(status);

CREATE TABLE IF NOT EXISTS dental_treatment_records (
    id           TEXT PRIMARY KEY,
    patient_nhi  TEXT NOT NULL,
    clinician_id TEXT NOT NULL,
    status       TEXT,
    body         JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_dental_records_patient ON dental_treatment_records(patient_nhi);
CREATE INDEX IF NOT EXISTS idx_dental_records_clinician ON dental_treatment_records(clinician_id);
