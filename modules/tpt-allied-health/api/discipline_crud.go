package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/PhillipC05/tpt-healthcare/core/db"
	"github.com/jackc/pgx/v5"
)

// The generic helpers below persist allied-health clinical records (OT,
// podiatry, speech assessments/plans/notes) using a JSONB "body" column for
// full fidelity plus indexed scalar columns for the common list filters.
// Tables are created by db/migrations/007_discipline_records.sql.

func insertDisciplineRecord(ctx context.Context, pool db.Pool, logger *slog.Logger, table, id, patientNHI, clinicianID, status, typ, category string, body any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", table, err)
	}
	_, err = pool.Exec(ctx,
		fmt.Sprintf(`INSERT INTO %s (id, patient_nhi, clinician_id, status, type, category, body, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,now(),now())`, table),
		id, nullStr(patientNHI), nullStr(clinicianID), nullStr(status), nullStr(typ), nullStr(category), payload,
	)
	if err != nil {
		return fmt.Errorf("insert %s: %w", table, err)
	}
	return nil
}

func getDisciplineRecord(ctx context.Context, pool db.Pool, table, id string) (body []byte, patientNHI string, err error) {
	err = pool.QueryRow(ctx, fmt.Sprintf(`SELECT body, patient_nhi FROM %s WHERE id=$1`, table), id).
		Scan(&body, &patientNHI)
	if err != nil {
		return nil, "", err
	}
	return body, patientNHI, nil
}

func updateDisciplineRecord(ctx context.Context, pool db.Pool, table, id, patientNHI, clinicianID, status, typ, category string, body any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", table, err)
	}
	tag, err := pool.Exec(ctx,
		fmt.Sprintf(`UPDATE %s SET patient_nhi=$1, clinician_id=$2, status=$3, type=$4, category=$5, body=$6, updated_at=now() WHERE id=$7`, table),
		nullStr(patientNHI), nullStr(clinicianID), nullStr(status), nullStr(typ), nullStr(category), payload, id,
	)
	if err != nil {
		return fmt.Errorf("update %s: %w", table, err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func deleteDisciplineRecord(ctx context.Context, pool db.Pool, table, id string) error {
	tag, err := pool.Exec(ctx, fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, table), id)
	if err != nil {
		return fmt.Errorf("delete %s: %w", table, err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// listDisciplineRecords returns the stored JSONB bodies matching the provided
// filters (only non-empty, known columns are applied).
func listDisciplineRecords(ctx context.Context, pool db.Pool, table string, filters map[string]string, limit, offset int) ([][]byte, error) {
	q := fmt.Sprintf("SELECT body FROM %s WHERE 1=1", table)
	args := []any{}
	order := 1
	for _, col := range []string{"patient_nhi", "clinician_id", "status", "type", "category"} {
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

// disciplineNotFound writes a 404 when err is pgx.ErrNoRows, otherwise a 500.
func disciplineError(w http.ResponseWriter, logger *slog.Logger, op string, err error) {
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "record not found", http.StatusNotFound)
		return
	}
	logger.Error(op, slog.Any("error", err))
	http.Error(w, "failed to process request", http.StatusInternalServerError)
}
