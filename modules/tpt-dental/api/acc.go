package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/PhillipC05/tpt-healthcare/core/audit"
	"github.com/PhillipC05/tpt-healthcare/core/db"
	"github.com/PhillipC05/tpt-healthcare/core/encryption"
	"github.com/google/uuid"

	dentalacc "github.com/PhillipC05/tpt-healthcare/modules/tpt-dental/internal/acc"
)

// ACCHandler handles ACC dental claim operations.
type ACCHandler struct {
	pool       db.Pool
	enc        *encryption.Cipher
	auditTrail *audit.Trail
	logger     *slog.Logger
}

// ListClaims returns ACC dental claims, optionally filtered by patient or status.
func (h *ACCHandler) ListClaims(w http.ResponseWriter, r *http.Request) {
	patientNHI := r.URL.Query().Get("patient_nhi")
	status := r.URL.Query().Get("status")
	limit, offset := parsePagination(r)

	bodies, err := jsonbList(r.Context(), h.pool, "dental_acc_claims", map[string]string{
		"patient_nhi": patientNHI,
		"status":      status,
	}, limit, offset)
	if err != nil {
		jsonbError(w, h.logger, "list dental claims", err)
		return
	}

	claims := make([]dentalacc.DentalClaim, 0, len(bodies))
	for _, b := range bodies {
		var c dentalacc.DentalClaim
		if json.Unmarshal(b, &c) == nil {
			claims = append(claims, c)
		}
	}
	writeJSON(w, http.StatusOK, claims)
}

// CreateClaim creates a new ACC dental claim in draft status.
func (h *ACCHandler) CreateClaim(w http.ResponseWriter, r *http.Request) {
	var claim dentalacc.DentalClaim
	if err := decodeJSON(r, &claim); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{
			Code: "INVALID_JSON", Message: fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	claim.ID = uuid.New().String()
	claim.Status = dentalacc.ClaimDraft
	claim.CreatedAt = time.Now().UTC()
	claim.UpdatedAt = claim.CreatedAt

	// Run validation.
	result := claim.Validate()
	if !result.Valid {
		writeJSON(w, http.StatusUnprocessableEntity, apiError{
			Code:    "VALIDATION_FAILED",
			Message: "Claim failed validation",
			Details: result.Errors,
		})
		return
	}

	if err := jsonbInsert(r.Context(), h.pool, h.logger, "dental_acc_claims", claim.ID, claim.PatientNHI, claim.ProviderHPI, string(claim.Status), &claim); err != nil {
		jsonbError(w, h.logger, "create dental claim", err)
		return
	}

	h.logger.Info("ACC dental claim created",
		slog.String("claim_id", claim.ID),
		slog.String("patient_nhi", claim.PatientNHI),
		slog.Int("teeth", len(claim.Teeth)))

	writeJSON(w, http.StatusCreated, claim)
}

// GetClaim returns details for a specific ACC dental claim.
func (h *ACCHandler) GetClaim(w http.ResponseWriter, r *http.Request) {
	claimID := r.PathValue("claimId")
	if claimID == "" {
		writeJSON(w, http.StatusBadRequest, apiError{
			Code: "MISSING_CLAIM_ID", Message: "Claim ID is required",
		})
		return
	}

	body, _, err := jsonbGet(r.Context(), h.pool, "dental_acc_claims", claimID)
	if err != nil {
		jsonbError(w, h.logger, "get dental claim", err)
		return
	}
	var claim dentalacc.DentalClaim
	if err := json.Unmarshal(body, &claim); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Code: "DB_ERROR", Message: "failed to decode claim"})
		return
	}
	writeJSON(w, http.StatusOK, claim)
}

// UpdateClaim updates fields on an existing draft claim.
func (h *ACCHandler) UpdateClaim(w http.ResponseWriter, r *http.Request) {
	claimID := r.PathValue("claimId")
	if claimID == "" {
		writeJSON(w, http.StatusBadRequest, apiError{
			Code: "MISSING_CLAIM_ID", Message: "Claim ID is required",
		})
		return
	}

	var claim dentalacc.DentalClaim
	if err := decodeJSON(r, &claim); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{
			Code: "INVALID_JSON", Message: fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	claim.ID = claimID
	claim.UpdatedAt = time.Now().UTC()

	if err := jsonbUpdate(r.Context(), h.pool, "dental_acc_claims", claimID, claim.PatientNHI, claim.ProviderHPI, string(claim.Status), &claim); err != nil {
		jsonbError(w, h.logger, "update dental claim", err)
		return
	}

	h.logger.Info("ACC dental claim updated",
		slog.String("claim_id", claimID))

	writeJSON(w, http.StatusOK, claim)
}

// SubmitClaim submits a draft claim to ACC for processing.
func (h *ACCHandler) SubmitClaim(w http.ResponseWriter, r *http.Request) {
	claimID := r.PathValue("claimId")
	if claimID == "" {
		writeJSON(w, http.StatusBadRequest, apiError{
			Code: "MISSING_CLAIM_ID", Message: "Claim ID is required",
		})
		return
	}

	// Load the draft claim, mark it submitted, and persist. Actual transmission
	// to ACC (core/acc.Client) is performed by the claims worker; here we record
	// the submission state and an ACC claim number for downstream correlation.
	body, _, err := jsonbGet(r.Context(), h.pool, "dental_acc_claims", claimID)
	if err != nil {
		jsonbError(w, h.logger, "submit dental claim", err)
		return
	}
	var claim dentalacc.DentalClaim
	if err := json.Unmarshal(body, &claim); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Code: "DB_ERROR", Message: "failed to decode claim"})
		return
	}

	claim.Status = dentalacc.ClaimSubmitted
	claim.ACCClaimNumber = "ACC-" + uuid.New().String()[:8]
	claim.UpdatedAt = time.Now().UTC()

	if err := jsonbUpdate(r.Context(), h.pool, "dental_acc_claims", claimID, claim.PatientNHI, claim.ProviderHPI, string(claim.Status), &claim); err != nil {
		jsonbError(w, h.logger, "submit dental claim persist", err)
		return
	}

	h.logger.Info("ACC dental claim submitted",
		slog.String("claim_id", claimID),
		slog.String("acc_claim_number", claim.ACCClaimNumber))

	writeJSON(w, http.StatusOK, map[string]any{
		"status":         "submitted",
		"claimId":        claimID,
		"accClaimNumber": claim.ACCClaimNumber,
	})
}

// CheckStatus polls ACC for the current status of a submitted claim.
func (h *ACCHandler) CheckStatus(w http.ResponseWriter, r *http.Request) {
	claimID := r.PathValue("claimId")
	if claimID == "" {
		writeJSON(w, http.StatusBadRequest, apiError{
			Code: "MISSING_CLAIM_ID", Message: "Claim ID is required",
		})
		return
	}

	// Simplified stub — real implementation polls core/acc.Client.Poll().
	writeJSON(w, http.StatusOK, map[string]string{
		"claimId": claimID,
		"status":  "pending",
		"message": "Awaiting ACC adjudication",
	})
}

// ValidateClaim pre-validates a claim payload without saving it.
func (h *ACCHandler) ValidateClaim(w http.ResponseWriter, r *http.Request) {
	var claim dentalacc.DentalClaim
	if err := decodeJSON(r, &claim); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{
			Code: "INVALID_JSON", Message: fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	result := claim.Validate()
	writeJSON(w, http.StatusOK, result)
}

// InjuryTypes returns the list of recognised ACC dental injury classifications.
func (h *ACCHandler) InjuryTypes(w http.ResponseWriter, r *http.Request) {
	types := dentalacc.DentalInjuryTypes()
	writeJSON(w, http.StatusOK, types)
}
