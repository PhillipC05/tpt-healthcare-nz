package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/PhillipC05/tpt-healthcare/core/db"
	"github.com/jackc/pgx/v5"
)

// parsePagination reads limit and offset from query params with defaults.
func parsePagination(r *http.Request) (limit, offset int) {
	limit = 50
	offset = 0
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		if l > 200 {
			l = 200
		}
		limit = l
	}
	if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && o >= 0 {
		offset = o
	}
	return
}

// JSONB-backed persistence helpers for dental clinical records. Tables are
// created by db/migrate/002_dental_claims_and_records.sql.

// nullStr returns nil for empty strings so nullable columns store NULL.
func nullStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func jsonbInsert(ctx context.Context, pool db.Pool, logger *slog.Logger, table, id, patientNHI, clinicianID, status string, body any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", table, err)
	}
	_, err = pool.Exec(ctx,
		fmt.Sprintf(`INSERT INTO %s (id, patient_nhi, clinician_id, status, body, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,now(),now())`, table),
		id, nullStr(patientNHI), nullStr(clinicianID), nullStr(status), payload,
	)
	if err != nil {
		return fmt.Errorf("insert %s: %w", table, err)
	}
	return nil
}

func jsonbGet(ctx context.Context, pool db.Pool, table, id string) (body []byte, patientNHI string, err error) {
	err = pool.QueryRow(ctx, fmt.Sprintf(`SELECT body, patient_nhi FROM %s WHERE id=$1`, table), id).
		Scan(&body, &patientNHI)
	if err != nil {
		return nil, "", err
	}
	return body, patientNHI, nil
}

func jsonbUpdate(ctx context.Context, pool db.Pool, table, id, patientNHI, clinicianID, status string, body any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", table, err)
	}
	tag, err := pool.Exec(ctx,
		fmt.Sprintf(`UPDATE %s SET patient_nhi=$1, clinician_id=$2, status=$3, body=$4, updated_at=now() WHERE id=$5`, table),
		nullStr(patientNHI), nullStr(clinicianID), nullStr(status), payload, id,
	)
	if err != nil {
		return fmt.Errorf("update %s: %w", table, err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func jsonbList(ctx context.Context, pool db.Pool, table string, filters map[string]string, limit, offset int) ([][]byte, error) {
	q := fmt.Sprintf("SELECT body FROM %s WHERE 1=1", table)
	args := []any{}
	order := 1
	for _, col := range []string{"patient_nhi", "clinician_id", "status"} {
		if v, ok := filters[col]; ok && v != "" {
			q += fmt.Sprintf(" AND %s=$%d", col, order)
			args = append(args, v)
			order++
		}
	}
	q += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", order, order+1)
	args = append(args, limit, offset)

	rows, err := pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list %s: %w", table, err)
	}
	defer rows.Close()

	out := make([][]byte, 0)
	for rows.Next() {
		var b []byte
		if err := rows.Scan(&b); err != nil {
			continue
		}
		out = append(out, b)
	}
	return out, nil
}

// jsonbError writes a 404 on pgx.ErrNoRows, otherwise a 500.
func jsonbError(w http.ResponseWriter, logger *slog.Logger, op string, err error) {
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSON(w, http.StatusNotFound, apiError{Code: "NOT_FOUND", Message: "record not found"})
		return
	}
	logger.Error(op, slog.Any("error", err))
	writeJSON(w, http.StatusInternalServerError, apiError{Code: "DB_ERROR", Message: "failed to process request"})
}
