# TODO — TPT Healthcare NZ Testing Suite

## Phase 1-5: Core + Frontend (COMPLETED)

- [x] Test infrastructure (testutil/helpers, mock/)
- [x] Core pure-logic (rbac, auth, ddi, resilience, translate, hl7)
- [x] Core data-layer (repo, terminology, storage, nhi)
- [x] Frontend package vitest (nz-codes, api-client, offline-store, fhir-types, ui)
- [x] Frontend component tests (@tpt/ui components — Button, Card, Badge, Table, Input, Modal, PatientBanner, NHIInput, ErrorBoundary)

## Phase 6: Module Unit Tests (IN PROGRESS)

### Tier 1 — High Priority (Rich Pure Logic)

- [x] **tpt-dental** — fdi/chart_test.go (~25 tests)
  - [x] Write test file
  - [x] Run & verify all pass
- [x] **tpt-dental** — fdi/surface_test.go (~18 tests)
  - [x] Write test file
  - [x] Run & verify all pass
- [x] **tpt-dental** — procedure/codes_test.go (~10 tests)
  - [x] Write test file
  - [x] Run & verify all pass
- [x] **tpt-dental** — acc/claim_test.go (~12 tests)
  - [x] Write test file
  - [x] Run & verify all pass
- [x] **tpt-vision** — refraction/prescription_test.go (~22 tests)
  - [x] Write test file
  - [x] Run & verify all pass
- [x] **tpt-vision** — optical/dispensing_test.go (~8 tests)
  - [x] Write test file
  - [x] Run & verify all pass
- [x] **tpt-vision** — ophthalmology/exam_test.go (~5 tests)
  - [x] Write test file
  - [x] Run & verify all pass
- [x] **tpt-vision** — acc/claim_test.go (~10 tests)
  - [x] Write test file
  - [x] Run & verify all pass
- [x] **tpt-allied-health** — acc/claim_test.go (~19 tests)
  - [x] Write test file
  - [x] Run & verify all pass

### Tier 1 — Sub-discipline Tests (Repeated Patterns)

- [x] **tpt-allied-health** — speech/therapy_test.go (~10 tests)
  - [x] Write test file
  - [x] Run & verify all pass
- [x] **tpt-allied-health** — physio/treatment_test.go (~8 tests)
  - [x] Write test file
  - [x] Run & verify all pass
- [x] **tpt-allied-health** — ot/assessment_test.go (~8 tests)
  - [x] Write test file
  - [x] Run & verify all pass
- [x] **tpt-allied-health** — podiatry/care_test.go (~10 tests)
  - [x] Write test file
  - [x] Run & verify all pass

### Tier 2 — Community Health

- [x] **tpt-community-health** — homevisit/visit_test.go (~12 tests)
  - [x] Write test file
  - [x] Run & verify all pass
- [x] **tpt-community-health** — outreach/program_test.go (~12 tests)
  - [x] Write test file
  - [x] Run & verify all pass
- [x] **tpt-community-health** — districtnursing/plan_test.go (~8 tests)
  - [x] Write test file
  - [x] Run & verify all pass

### Tier 2 — Addiction

- [x] **tpt-addiction** — methadone/programme_test.go (~7 tests)
  - [x] Write test file
  - [x] Run & verify all pass

## Skipped

- [ ] **tpt-palliative** — No testable logic (pure data types, zero functions)

## Phase 7: Hospital Go-Live Gaps (Auckland City / Starship)

Gap analysis for a large adult tertiary hospital (Auckland City) and paediatric tertiary
hospital (Starship) going live on `modules/tpt-hospital` + `modules/tpt-maternal-child-health`.
See plan `lets-say-auckland-city-jolly-pinwheel.md` for full detail.

- [x] **tpt-hospital / tpt-pathology / tpt-radiology** — CPOE implemented: `clinical_orders` table (`012_cpoe.sql`), `CPOEHandler` with 7 routes + result callback, HL7 ORM^O01 dispatch, admission-linked order lifecycle
  - [x] Design order model (admissionID-linked) and wire to pathology/radiology result callback
- [x] **core/hl7** — HL7 ADT (A01/A02/A03/A08) message builders (`core/hl7/adt.go`) now wired into `tpt-hospital` admissions lifecycle (`AdmissionsHandler.sendADT`) via optional MLLP client (`Config.HL7MLLPAddr`), same pattern as CPOE's ORM dispatch
  - [x] Add ADT message builders for admit/transfer/discharge/update events
  - [x] ORM builders for lab/imaging order messages
- [x] **tpt-hospital/api/billing.go** — `deriveDRG` now delegates to the existing `core/terminology.DRGGrouper` (MDC/category rules + MCC-driven complexity/weight adjustment) instead of its own hardcoded 4-bucket switch; MCC signal comes from `additional-diagnosis` codes on the admission
  - [x] Implement real AR-DRG/WIES grouper in `core/terminology/` (grouper already existed but was unused; now wired in)
- [x] **core/fhir** — `Location` resource exists in both `core/fhir/r5/location.go` and `core/fhir/r4/location.go`; `Encounter` exists in both `r5` and `r4` (`core/fhir/r4/encounter.go`). The FHIR REST API (`interop/api/fhir.go`) is generic over `resourceType` (JSONB store), so both work through CRUD/search with no extra registration.
  - [x] Add Location resource type; add r4 Encounter for compatibility
- [x] **tpt-hospital/api/pharmacy_*.go** — eMAR now has real barcode/five-rights verification (`bedside_verification.go` `IsFiveRightsOK`) and a DB-backed S8 controlled-drug register (`pharmacy_query.go` `controlled_drug_register` table, running-balance INSERT/SELECT); IV pump/smart-infusion types (`IVPumpType`, `IVPumpStatus`) already tracked in `pharmacy_types.go`/`pharmacy_handler.go`
  - [x] Scope and implement bedside verification + S8 register
- [x] **tpt-hospital/api/icu.go** — Fluid balance charting (`AddFluidBalance`/`insertFluidBalanceEntry`/`listFluidBalanceEntries`) and EWS/PEWS scoring engine (`CalculateEWS`/`CalculatePEWS`/`insertEWSRecord`/`listEWSRecords`) are DB-backed, not stubs
  - [x] Add fluid balance charting
  - [x] Add EWS (adult) / PEWS (paediatric) scoring engine
- [x] **tpt-maternal-child-health** — NICU/PICU records now carry `HospitalAdmissionID` (`nicu.go`) and paediatric child-protection records use a `paediatric_admission_id` FK; this tooling lives in `tpt-maternal-child-health` rather than separate `tpt-hospital` files
  - [x] Wire NICU/PICU/growth-chart records to the hospital admission model
- [x] **tpt-hospital/api/icu.go** — `CalculatePaediatricDose` implements weight-based (mg/kg) dosing with max-dose caps
  - [x] Implement paediatric dosing calculator
- [x] **tpt-hospital/api/wards.go** — `PatientFlowForecast`/`buildFlowForecast` computes discharge-ETA style forecasts from average LOS + active admissions, beyond a static snapshot
  - [x] Design and implement patient-flow forecasting dashboard
- [ ] **tpt-hospital/api/admissions_*.go** — `discharge_auto.go` (`AutoPopulateDischargeSummary`, `GPTransmissionReady`) and `admissions_handler.go` build the discharge summary and check GP-transmission readiness, but the actual bundle send is still just a comment ("handled through core/gp2gp") — `core/gp2gp.go` exists but is not called
  - [x] Auto-populate discharge summary from admission/coding/pharmacy data
  - [ ] Wire the GP2GP bundle transfer call (currently unwired despite `core/gp2gp` existing)

## Phase 8: Replace Stubs & Scaffolds with Real Implementations

Repo-wide audit for placeholder/stub/scaffold code that needs replacing with real, working
implementations. Found via full sweep of `core/`, `interop/`, and all `modules/*`.
See plan `lets-say-auckland-city-jolly-pinwheel.md` for full detail.

### Critical — integration/persistence surfaces

- [x] **interop/api/fhir.go** — FHIR R4/R5 REST API is already wired to `core/repo.Store` (CRUD+Search all delegate to the real store; 3-tier fallback: explicit store → PostgresStore from pool → MemoryStore)
- [x] **interop/api/terminology.go** — `coreTermStore` adapter already bridges `core/terminology` stores to the `TermStore` interface; wired in `main.go` via config keys (`snomed_csv`, `loinc_csv`, `icd10_csv`, `nzmt_csv`)
- [x] **core/db/migrate.go** — `Migrate()` now fails loudly on non-string third argument; fixed call sites in tpt-vision, tpt-telehealth, tpt-radiology, tpt-mental-health, tpt-hospital, tpt-dental (passed `*slog.Logger` → now pass correct dir string or `""`)
- [x] **RunMigrations signatures** — Removed unused `logger *slog.Logger` parameter from 27 modules' `RunMigrations` functions and all callers

### Fully-scaffolded modules (need full persistence)

- [ ] **tpt-pharmacy** — Dispensing/claims workflow is entirely `// In production: ...` placeholders; no `db/migrate` directory; nothing persisted
- [x] **tpt-counselling** — `api/query.go` now has real `SELECT`/`INSERT INTO counselling_eap_claims` etc. via pgx; no longer in-memory
- [x] **tpt-nutrition** — `api/query.go` has real SQL persistence for food diary/meal plans/body composition (lighter coverage than other modules — worth a follow-up spot-check)
- [x] **tpt-immunisation** — `api/query.go` added; `nir.go` `Submit()` now translates to a real `core/fhir/r4.Immunization` (`translateImmunisationToR4`) and POSTs it to the NIR instead of a placeholder struct
- [x] **tpt-health-billing** — `api/query.go` (1148 lines, 33 SQL statements) now backs ACC/insurance/invoices/reconciliation; no longer hard-coded JSON
- [ ] **tpt-clinical-trials** — Real SQL migrations exist, but every handler (participants, adverse events, protocols, visits — ~43 methods) just returns HTTP 501 via `notImplemented()`
- [ ] **tpt-chiropractic** — Handlers hold everything in-memory (`internal/spine/chart.go`, `internal/xray/referral.go`); migrations now work (unblocked by migrate.go fix)
- [ ] **tpt-osteopathy** — Same pattern as tpt-chiropractic
- [ ] **tpt-acupuncture** — Same pattern as tpt-chiropractic
- [ ] **tpt-tcm** — Same pattern as tpt-chiropractic
- [ ] **tpt-massage** — Same pattern as tpt-chiropractic
- [ ] **tpt-naturopathy** — Same pattern as tpt-chiropractic

### Missing/broken migrations

- [x] **tpt-aged-care** — Handlers issue real SQL against `aged_care_plans`, `aged_care_funded_hours_allocations`, `aged_care_interrai_assessments`, `aged_care_nasc_referrals`, `aged_care_nasc_service_plans` — added `modules/tpt-aged-care/db/migrations` + `embed.go`, wired via `migrate.New(agedcaredb.Migrations, pool)` in `RunMigrations`
- [x] **tpt-blood-bank** — Real query code (crossmatch/donors/inventory) — added `modules/tpt-blood-bank/db/migrations/001_blood_bank_tables.sql` + `embed.go`, wired via `migrate.New(bloodbankdb.Migrations, pool)` in `RunMigrations`
- [x] **tpt-practice** — Real query code (rostering/settings) — tables already exist in `core/db/migrate/007_practice_management.sql`; fixed the broken `db.Migrate(ctx, pool, "")` call site to use `migrate.New(migrate.MigrationsFS, pool)`

### Partial stubs in otherwise-real modules

- [x] **tpt-allied-health** — `acc_handler.go` queries the real repo; `ot.go`, `podiatry.go`, `speech.go` "get by ID" (and full CRUD) handlers now persist and read real rows via `api/discipline_crud.go` + migration `006_discipline_records.sql` (12 tables), with consent checks.
- [x] **tpt-vision** — `acc.go` (`GetClaimFHIR`), `ophth.go` (`GetExamFHIR`), `optical.go` (`GetOrderFHIR`), `refraction.go` (`GetPrescriptionFHIR`) now query the real `fhir_resource` JSONB column from `vision_acc_claims` / `vision_ophthalmic_exams` / `vision_dispensing_orders` / `vision_prescriptions` and return it as `application/fhir+json` (404 if not found). Note: `CreateClaim`/`CreateOrder`/`CreatePrescription` still don't persist, so records must exist (via SQL or the real `CreateExam`) for the FHIR getters to return data — full Create persistence for those three is a follow-up.
- [ ] **tpt-dental** — `acc.go` (SubmitClaim + claim CRUD) and `procedure.go` (treatment-record CRUD) explicitly commented "Simplified stub", operate on in-memory data
- [ ] **tpt-doctor** — `api/pho.go` (PHO extract "would transmit" but doesn't) and `api/referrals.go` (`Send` doesn't dispatch to receiving provider's inbox)
- [x] **tpt-pathology** — `api/mllp.go` `resolveTenant` now looks up the sending lab site (MSH-4 / ZNZL LabSite) against the new `tenant_lab_sites` mapping table (`002_tenant_lab_sites.sql`); unknown senders fall back to the nil tenant with a warning instead of a silent default.
- [ ] **tpt-cardiology / tpt-rehabilitation / tpt-maternal-child-health** — `api/helpers.go` `recordAudit` isn't in the same DB transaction as the clinical write it audits (durability/consistency gap)

### Moderate stubs in core/

- [x] **core/subscription/engine.go** — `buildNotificationBundle` already emits the real subscription ID (`Subscription/` + `subscriptionID.String()`); verified no `Subscription/unknown` remains
- [x] **core/backup/scheduler.go** — `recordSuccess` is only called after a real upload (`result.SizeBytes`); the no-provider path correctly calls `recordFailure`. Removed leftover dead `bytes` no-op.
- [x] **core/nhi/nhi.go** — Already uses `fhirr4.Patient` throughout; no custom NHI types to replace

## Phase 9: Tenant Service-Line Profiles

Onboarding a new facility (e.g. a large tertiary hospital vs. a single-service clinic)
currently has no concept of "which service lines does this tenant run." Recommendation:
a service-line profile/enablement layer rather than fixed per-hospital templates, since
real facilities are combinations (e.g. one campus with both adult and paediatric wards)
that a hard-coded template can't cleanly represent. See `core/tenant` for the existing
generic `Tenant` model this would extend.

- [x] **core/tenant** — Design a service-line profile: tenant selects which service lines it runs (ED, ICU, NICU/PICU, theatre, oncology, etc.) at onboarding
  - [x] Define the service-line catalogue/schema — `core/servicelines/catalogue.go`, 16 service lines, each with modules/ward-types/triage-scale/formulary; persisted per-tenant in `tenant_service_lines` (`core/db/migrate/013_tenant_service_lines.sql`)
  - [x] Wire service-line selection to toggle relevant modules/routes per tenant — `PUT /api/v1/practice/service-lines` unions resolved modules into the existing `tenants.settings.activeModules` (additive; manually-enabled modules are preserved)
  - [x] Seed sensible defaults per service line (ward-type list, triage scale, relevant formulary subset) — exposed via `GET /api/v1/practice/service-lines` and `core/servicelines.Resolve{Modules,WardTypes,FormularySubset,TriageScales}`
  - [x] Support facilities with multiple/mixed service lines (no forking or per-site hard-coded templates) — a tenant selects any subset of the catalogue; defaults are unioned, not templated
  - [ ] Frontend: admin UI for selecting service lines during onboarding (backend/API complete; no UI wired yet)

## Final Verification

- [x] Run `go test ./modules/tpt-dental/...` — all packages pass (see middleware fix below)
- [x] Run `go test ./modules/tpt-vision/...` — all internal packages pass
- [x] Run `go test ./modules/tpt-allied-health/...` — all internal packages pass
- [x] Run `go test ./modules/tpt-community-health/...` — all packages pass (see middleware fix below)
- [x] Run `go test ./modules/tpt-addiction/...` — all internal packages pass (no test files needed for counselling/prescribing — no pure logic to test)
- [x] Run `gofmt -w ./modules/...` — applied to tpt-dental, tpt-vision, tpt-allied-health, tpt-community-health, tpt-addiction (all had alignment-only diffs)
- [x] Run `make lint` (via `golangci-lint run ./internal/...` per module) — 0 issues across all five modules

### Pre-existing build break found during verification (fixed)

`modules/tpt-dental/api/server.go` and `modules/tpt-community-health/api/server.go` failed to compile because their
`chain` middleware helper still called the old pre-refactor `core/middleware` signatures (`middleware.Auth`,
`middleware.Tenant`, `middleware.Recovery`, and single-arg `AuditWrap`/`CORS`/`RateLimit`). Fixed both to match the
current API (same pattern already used in `tpt-vision/api/server.go`): `middleware.AuditWrap(s.auditTrail)(h)`,
`auth.RequireAuth(s.auth)(h)`, `middleware.TenantExtraction()(h)`, `middleware.CORS([]string{"*"})(h)`,
`middleware.RateLimit(10, 30)(h)`, `middleware.RecoveryMiddleware()(h)`. Both modules now build and test clean.
