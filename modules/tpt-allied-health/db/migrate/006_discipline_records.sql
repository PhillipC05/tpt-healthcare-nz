-- Allied health discipline clinical records (OT, podiatry, speech).
-- Each entity is stored with indexed scalar columns for common list filters
-- plus a JSONB "body" column holding the full record for round-trip fidelity.

CREATE TABLE IF NOT EXISTS ot_assessments (
    id           TEXT PRIMARY KEY,
    patient_nhi  TEXT NOT NULL,
    clinician_id TEXT NOT NULL,
    status       TEXT,
    type         TEXT,
    category     TEXT,
    body         JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_ot_assessments_patient ON ot_assessments(patient_nhi);
CREATE INDEX IF NOT EXISTS idx_ot_assessments_clinician ON ot_assessments(clinician_id);

CREATE TABLE IF NOT EXISTS ot_intervention_plans (
    id           TEXT PRIMARY KEY,
    patient_nhi  TEXT NOT NULL,
    clinician_id TEXT NOT NULL,
    status       TEXT,
    type         TEXT,
    category     TEXT,
    body         JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_ot_plans_patient ON ot_intervention_plans(patient_nhi);
CREATE INDEX IF NOT EXISTS idx_ot_plans_clinician ON ot_intervention_plans(clinician_id);

CREATE TABLE IF NOT EXISTS ot_session_notes (
    id           TEXT PRIMARY KEY,
    patient_nhi  TEXT NOT NULL,
    clinician_id TEXT NOT NULL,
    status       TEXT,
    type         TEXT,
    category     TEXT,
    body         JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_ot_notes_patient ON ot_session_notes(patient_nhi);
CREATE INDEX IF NOT EXISTS idx_ot_notes_clinician ON ot_session_notes(clinician_id);

CREATE TABLE IF NOT EXISTS podiatry_assessments (
    id           TEXT PRIMARY KEY,
    patient_nhi  TEXT NOT NULL,
    clinician_id TEXT NOT NULL,
    status       TEXT,
    type         TEXT,
    category     TEXT,
    body         JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_podiatry_assessments_patient ON podiatry_assessments(patient_nhi);
CREATE INDEX IF NOT EXISTS idx_podiatry_assessments_clinician ON podiatry_assessments(clinician_id);

CREATE TABLE IF NOT EXISTS podiatry_treatment_plans (
    id           TEXT PRIMARY KEY,
    patient_nhi  TEXT NOT NULL,
    clinician_id TEXT NOT NULL,
    status       TEXT,
    type         TEXT,
    category     TEXT,
    body         JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_podiatry_plans_patient ON podiatry_treatment_plans(patient_nhi);
CREATE INDEX IF NOT EXISTS idx_podiatry_plans_clinician ON podiatry_treatment_plans(clinician_id);

CREATE TABLE IF NOT EXISTS podiatry_session_notes (
    id           TEXT PRIMARY KEY,
    patient_nhi  TEXT NOT NULL,
    clinician_id TEXT NOT NULL,
    status       TEXT,
    type         TEXT,
    category     TEXT,
    body         JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_podiatry_notes_patient ON podiatry_session_notes(patient_nhi);
CREATE INDEX IF NOT EXISTS idx_podiatry_notes_clinician ON podiatry_session_notes(clinician_id);

CREATE TABLE IF NOT EXISTS podiatry_wound_assessments (
    id           TEXT PRIMARY KEY,
    patient_nhi  TEXT NOT NULL,
    clinician_id TEXT NOT NULL,
    status       TEXT,
    type         TEXT,
    category     TEXT,
    body         JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_podiatry_wound_patient ON podiatry_wound_assessments(patient_nhi);
CREATE INDEX IF NOT EXISTS idx_podiatry_wound_clinician ON podiatry_wound_assessments(clinician_id);

CREATE TABLE IF NOT EXISTS speech_assessments (
    id           TEXT PRIMARY KEY,
    patient_nhi  TEXT NOT NULL,
    clinician_id TEXT NOT NULL,
    status       TEXT,
    type         TEXT,
    category     TEXT,
    body         JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_speech_assessments_patient ON speech_assessments(patient_nhi);
CREATE INDEX IF NOT EXISTS idx_speech_assessments_clinician ON speech_assessments(clinician_id);

CREATE TABLE IF NOT EXISTS speech_therapy_plans (
    id           TEXT PRIMARY KEY,
    patient_nhi  TEXT NOT NULL,
    clinician_id TEXT NOT NULL,
    status       TEXT,
    type         TEXT,
    category     TEXT,
    body         JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_speech_plans_patient ON speech_therapy_plans(patient_nhi);
CREATE INDEX IF NOT EXISTS idx_speech_plans_clinician ON speech_therapy_plans(clinician_id);

CREATE TABLE IF NOT EXISTS speech_session_notes (
    id           TEXT PRIMARY KEY,
    patient_nhi  TEXT NOT NULL,
    clinician_id TEXT NOT NULL,
    status       TEXT,
    type         TEXT,
    category     TEXT,
    body         JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_speech_notes_patient ON speech_session_notes(patient_nhi);
CREATE INDEX IF NOT EXISTS idx_speech_notes_clinician ON speech_session_notes(clinician_id);

CREATE TABLE IF NOT EXISTS speech_swallowing_assessments (
    id           TEXT PRIMARY KEY,
    patient_nhi  TEXT NOT NULL,
    clinician_id TEXT NOT NULL,
    status       TEXT,
    type         TEXT,
    category     TEXT,
    body         JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_speech_swallow_patient ON speech_swallowing_assessments(patient_nhi);
CREATE INDEX IF NOT EXISTS idx_speech_swallow_clinician ON speech_swallowing_assessments(clinician_id);
