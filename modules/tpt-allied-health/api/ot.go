// Package api implements HTTP handlers for occupational therapy services.
package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/PhillipC05/tpt-healthcare/core/consent"
	"github.com/PhillipC05/tpt-healthcare/core/db"
	"github.com/PhillipC05/tpt-healthcare/core/hpi"
	"github.com/PhillipC05/tpt-healthcare/modules/tpt-allied-health/internal/ot"
	"github.com/google/uuid"
)

// OTHandler handles occupational therapy API endpoints.
type OTHandler struct {
	hpiClient    *hpi.Client
	consentStore *consent.Store
	pool         db.Pool
	logger       *slog.Logger
}

// NewOTHandler creates a new OT handler.
func NewOTHandler(hpiClient *hpi.Client, consentStore *consent.Store, pool db.Pool, logger *slog.Logger) *OTHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return &OTHandler{hpiClient: hpiClient, consentStore: consentStore, pool: pool, logger: logger}
}

// RegisterRoutes registers OT routes.
func (h *OTHandler) RegisterRoutes(mux *http.ServeMux, protect func(http.HandlerFunc) http.Handler) {
	mux.Handle("POST /api/v1/ot/assessments", protect(h.CreateAssessment))
	mux.Handle("GET /api/v1/ot/assessments", protect(h.ListAssessments))
	mux.Handle("GET /api/v1/ot/assessments/{id}", protect(h.GetAssessment))
	mux.Handle("PUT /api/v1/ot/assessments/{id}", protect(h.UpdateAssessment))
	mux.Handle("DELETE /api/v1/ot/assessments/{id}", protect(h.DeleteAssessment))

	mux.Handle("POST /api/v1/ot/intervention-plans", protect(h.CreateInterventionPlan))
	mux.Handle("GET /api/v1/ot/intervention-plans", protect(h.ListInterventionPlans))
	mux.Handle("GET /api/v1/ot/intervention-plans/{id}", protect(h.GetInterventionPlan))
	mux.Handle("PUT /api/v1/ot/intervention-plans/{id}", protect(h.UpdateInterventionPlan))

	mux.Handle("POST /api/v1/ot/session-notes", protect(h.CreateSessionNote))
	mux.Handle("GET /api/v1/ot/session-notes", protect(h.ListSessionNotes))
	mux.Handle("GET /api/v1/ot/session-notes/{id}", protect(h.GetSessionNote))
	mux.Handle("PUT /api/v1/ot/session-notes/{id}", protect(h.UpdateSessionNote))

	mux.Handle("GET /api/v1/ot/outcome-measures", protect(h.ListOutcomeMeasures))
}

// CreateAssessment creates a new OT assessment.
func (h *OTHandler) CreateAssessment(w http.ResponseWriter, r *http.Request) {
	if !requireAPC(w, r, h.hpiClient) {
		return
	}

	var assessment ot.Assessment
	if err := json.NewDecoder(r.Body).Decode(&assessment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	assessment.ID = uuid.New().String()
	now := time.Now().UnixMilli()
	assessment.CreatedAt = now
	assessment.UpdatedAt = now

	if err := assessment.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := insertDisciplineRecord(r.Context(), h.pool, h.logger, "ot_assessments", assessment.ID, assessment.PatientNHI, assessment.ClinicianID, string(assessment.Status), string(assessment.Type), "", &assessment); err != nil {
		disciplineError(w, h.logger, "create ot assessment", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(assessment)
}

// GetAssessment retrieves an assessment by ID.
func (h *OTHandler) GetAssessment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	body, patientNHI, err := getDisciplineRecord(r.Context(), h.pool, "ot_assessments", id)
	if err != nil {
		disciplineError(w, h.logger, "get ot assessment", err)
		return
	}

	if !checkConsent(w, r, h.consentStore, patientNHI) {
		return
	}

	var assessment ot.Assessment
	if err := json.Unmarshal(body, &assessment); err != nil {
		http.Error(w, "failed to decode record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assessment)
}

// ListAssessments lists assessments with filters.
func (h *OTHandler) ListAssessments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	patientNHI := query.Get("patient_nhi")
	clinicianID := query.Get("clinician_id")
	assessmentType := query.Get("type")
	status := query.Get("status")
	limit, offset := parsePagination(r)

	if !checkConsent(w, r, h.consentStore, patientNHI) {
		return
	}

	bodies, err := listDisciplineRecords(r.Context(), h.pool, "ot_assessments", map[string]string{
		"patient_nhi":  patientNHI,
		"clinician_id": clinicianID,
		"type":         assessmentType,
		"status":       status,
	}, limit, offset)
	if err != nil {
		disciplineError(w, h.logger, "list ot assessments", err)
		return
	}

	assessments := make([]ot.Assessment, 0, len(bodies))
	for _, b := range bodies {
		var a ot.Assessment
		if json.Unmarshal(b, &a) == nil {
			assessments = append(assessments, a)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":   assessments,
		"limit":  limit,
		"offset": offset,
		"total":  len(assessments),
		"filters": map[string]string{
			"patient_nhi":  patientNHI,
			"clinician_id": clinicianID,
			"type":         assessmentType,
			"status":       status,
		},
	})
}

// UpdateAssessment updates an assessment.
func (h *OTHandler) UpdateAssessment(w http.ResponseWriter, r *http.Request) {
	if !requireAPC(w, r, h.hpiClient) {
		return
	}

	id := r.PathValue("id")

	var assessment ot.Assessment
	if err := json.NewDecoder(r.Body).Decode(&assessment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	assessment.ID = id
	assessment.UpdatedAt = time.Now().UnixMilli()

	if err := assessment.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := updateDisciplineRecord(r.Context(), h.pool, "ot_assessments", id, assessment.PatientNHI, assessment.ClinicianID, string(assessment.Status), string(assessment.Type), "", &assessment); err != nil {
		disciplineError(w, h.logger, "update ot assessment", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assessment)
}

// DeleteAssessment deletes an assessment.
func (h *OTHandler) DeleteAssessment(w http.ResponseWriter, r *http.Request) {
	if !requireAPC(w, r, h.hpiClient) {
		return
	}
	if err := deleteDisciplineRecord(r.Context(), h.pool, "ot_assessments", r.PathValue("id")); err != nil {
		disciplineError(w, h.logger, "delete ot assessment", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// CreateInterventionPlan creates a new intervention plan.
func (h *OTHandler) CreateInterventionPlan(w http.ResponseWriter, r *http.Request) {
	if !requireAPC(w, r, h.hpiClient) {
		return
	}

	var plan ot.InterventionPlan
	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	plan.ID = uuid.New().String()
	now := time.Now().UnixMilli()
	plan.CreatedAt = now
	plan.UpdatedAt = now

	if err := plan.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := insertDisciplineRecord(r.Context(), h.pool, h.logger, "ot_intervention_plans", plan.ID, plan.PatientNHI, plan.ClinicianID, string(plan.Status), "", "", &plan); err != nil {
		disciplineError(w, h.logger, "create ot intervention plan", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(plan)
}

// GetInterventionPlan retrieves an intervention plan by ID.
func (h *OTHandler) GetInterventionPlan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	body, patientNHI, err := getDisciplineRecord(r.Context(), h.pool, "ot_intervention_plans", id)
	if err != nil {
		disciplineError(w, h.logger, "get ot intervention plan", err)
		return
	}

	if !checkConsent(w, r, h.consentStore, patientNHI) {
		return
	}

	var plan ot.InterventionPlan
	if err := json.Unmarshal(body, &plan); err != nil {
		http.Error(w, "failed to decode record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

// ListInterventionPlans lists intervention plans with filters.
func (h *OTHandler) ListInterventionPlans(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	patientNHI := query.Get("patient_nhi")
	clinicianID := query.Get("clinician_id")
	status := query.Get("status")
	limit, offset := parsePagination(r)

	if !checkConsent(w, r, h.consentStore, patientNHI) {
		return
	}

	bodies, err := listDisciplineRecords(r.Context(), h.pool, "ot_intervention_plans", map[string]string{
		"patient_nhi":  patientNHI,
		"clinician_id": clinicianID,
		"status":       status,
	}, limit, offset)
	if err != nil {
		disciplineError(w, h.logger, "list ot intervention plans", err)
		return
	}

	plans := make([]ot.InterventionPlan, 0, len(bodies))
	for _, b := range bodies {
		var p ot.InterventionPlan
		if json.Unmarshal(b, &p) == nil {
			plans = append(plans, p)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":   plans,
		"limit":  limit,
		"offset": offset,
		"total":  len(plans),
		"filters": map[string]string{
			"patient_nhi":  patientNHI,
			"clinician_id": clinicianID,
			"status":       status,
		},
	})
}

// UpdateInterventionPlan updates an intervention plan.
func (h *OTHandler) UpdateInterventionPlan(w http.ResponseWriter, r *http.Request) {
	if !requireAPC(w, r, h.hpiClient) {
		return
	}

	id := r.PathValue("id")

	var plan ot.InterventionPlan
	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	plan.ID = id
	plan.UpdatedAt = time.Now().UnixMilli()

	if err := plan.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := updateDisciplineRecord(r.Context(), h.pool, "ot_intervention_plans", id, plan.PatientNHI, plan.ClinicianID, string(plan.Status), "", "", &plan); err != nil {
		disciplineError(w, h.logger, "update ot intervention plan", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

// CreateSessionNote creates a new session note.
func (h *OTHandler) CreateSessionNote(w http.ResponseWriter, r *http.Request) {
	if !requireAPC(w, r, h.hpiClient) {
		return
	}

	var note ot.SessionNote
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	note.ID = uuid.New().String()
	now := time.Now().UnixMilli()
	note.CreatedAt = now
	note.UpdatedAt = now

	if err := note.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := insertDisciplineRecord(r.Context(), h.pool, h.logger, "ot_session_notes", note.ID, note.PatientNHI, note.ClinicianID, "", "", "", &note); err != nil {
		disciplineError(w, h.logger, "create ot session note", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

// GetSessionNote retrieves a session note by ID.
func (h *OTHandler) GetSessionNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	body, patientNHI, err := getDisciplineRecord(r.Context(), h.pool, "ot_session_notes", id)
	if err != nil {
		disciplineError(w, h.logger, "get ot session note", err)
		return
	}

	if !checkConsent(w, r, h.consentStore, patientNHI) {
		return
	}

	var note ot.SessionNote
	if err := json.Unmarshal(body, &note); err != nil {
		http.Error(w, "failed to decode record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// ListSessionNotes lists session notes with filters.
func (h *OTHandler) ListSessionNotes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	patientNHI := query.Get("patient_nhi")
	interventionPlanID := query.Get("intervention_plan_id")
	limit, offset := parsePagination(r)

	if !checkConsent(w, r, h.consentStore, patientNHI) {
		return
	}

	bodies, err := listDisciplineRecords(r.Context(), h.pool, "ot_session_notes", map[string]string{
		"patient_nhi": patientNHI,
	}, limit, offset)
	if err != nil {
		disciplineError(w, h.logger, "list ot session notes", err)
		return
	}

	notes := make([]ot.SessionNote, 0, len(bodies))
	for _, b := range bodies {
		var n ot.SessionNote
		if json.Unmarshal(b, &n) == nil {
			notes = append(notes, n)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":   notes,
		"limit":  limit,
		"offset": offset,
		"total":  len(notes),
		"filters": map[string]string{
			"patient_nhi":          patientNHI,
			"intervention_plan_id": interventionPlanID,
		},
	})
}

// UpdateSessionNote updates a session note.
func (h *OTHandler) UpdateSessionNote(w http.ResponseWriter, r *http.Request) {
	if !requireAPC(w, r, h.hpiClient) {
		return
	}

	id := r.PathValue("id")

	var note ot.SessionNote
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	note.ID = id
	note.UpdatedAt = time.Now().UnixMilli()

	if err := note.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := updateDisciplineRecord(r.Context(), h.pool, "ot_session_notes", id, note.PatientNHI, note.ClinicianID, "", "", "", &note); err != nil {
		disciplineError(w, h.logger, "update ot session note", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// ListOutcomeMeasures lists standardised OT outcome measures.
func (h *OTHandler) ListOutcomeMeasures(w http.ResponseWriter, r *http.Request) {
	measures := []map[string]string{
		{"code": "COPM", "name": "Canadian Occupational Performance Measure", "domain": "performance_satisfaction"},
		{"code": "FIM", "name": "Functional Independence Measure", "domain": "function"},
		{"code": "AMPS", "name": "Assessment of Motor and Process Skills", "domain": "motor_process"},
		{"code": "MOHOST", "name": "Model of Human Occupation Screening Tool", "domain": "occupation"},
		{"code": "BARTHEL", "name": "Barthel Index", "domain": "adl"},
		{"code": "LAWTON", "name": "Lawton IADL Scale", "domain": "iadl"},
		{"code": "MMSE", "name": "Mini-Mental State Examination", "domain": "cognitive"},
		{"code": "MOCA", "name": "Montreal Cognitive Assessment", "domain": "cognitive"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measures)
}
