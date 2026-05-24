package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/redhajuanda/komon/fail"
	"github.com/redhajuanda/komon/tracer"
	"github.com/lokalabsmaya/rasia/internal/core/domain"
	"github.com/lokalabsmaya/rasia/shared/failure"
)

// secretRepository implements repositories.Secret for SQLite.
type secretRepository struct {
	db querier
}

// NewSecretRepository creates a new SQLite SecretRepository.
func NewSecretRepository(db querier) *secretRepository {
	return &secretRepository{db: db}
}

// GetSecretByID returns a secret by its ID.
func (r *secretRepository) GetSecretByID(ctx context.Context, id string) (*domain.Secret, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	var s domain.Secret
	err := r.db.QueryRowContext(ctx,
		`SELECT id, file_id, key_name, value_enc, created_at, updated_at, deleted_at
		 FROM secrets WHERE id = ? AND deleted_at = 0`, id).
		Scan(&s.ID, &s.FileID, &s.KeyName, &s.ValueEnc, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fail.Wrap(err).WithFailure(failure.ErrSecretNotFound)
		}
		return nil, fail.Wrap(err)
	}
	return &s, nil
}

// ListSecrets returns all non-deleted secrets in a file.
func (r *secretRepository) ListSecrets(ctx context.Context, fileID string) ([]*domain.Secret, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, file_id, key_name, value_enc, created_at, updated_at, deleted_at
		 FROM secrets WHERE file_id = ? AND deleted_at = 0 ORDER BY key_name`, fileID)
	if err != nil {
		return nil, fail.Wrap(err)
	}
	defer rows.Close()

	var result []*domain.Secret
	for rows.Next() {
		var s domain.Secret
		if err := rows.Scan(&s.ID, &s.FileID, &s.KeyName, &s.ValueEnc, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt); err != nil {
			return nil, fail.Wrap(err)
		}
		result = append(result, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, fail.Wrap(err)
	}
	return result, nil
}

// ListSecretsByFileIDs returns secrets for multiple file IDs.
func (r *secretRepository) ListSecretsByFileIDs(ctx context.Context, fileIDs []string) ([]*domain.Secret, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	if len(fileIDs) == 0 {
		return nil, nil
	}

	placeholders := strings.Repeat("?,", len(fileIDs))
	placeholders = placeholders[:len(placeholders)-1]

	args := make([]any, len(fileIDs))
	for i, id := range fileIDs {
		args[i] = id
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, file_id, key_name, value_enc, created_at, updated_at, deleted_at
		 FROM secrets WHERE file_id IN (`+placeholders+`) AND deleted_at = 0 ORDER BY file_id, key_name`,
		args...)
	if err != nil {
		return nil, fail.Wrap(err)
	}
	defer rows.Close()

	var result []*domain.Secret
	for rows.Next() {
		var s domain.Secret
		if err := rows.Scan(&s.ID, &s.FileID, &s.KeyName, &s.ValueEnc, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt); err != nil {
			return nil, fail.Wrap(err)
		}
		result = append(result, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, fail.Wrap(err)
	}
	return result, nil
}

// CreateSecret inserts a new secret.
func (r *secretRepository) CreateSecret(ctx context.Context, secret *domain.Secret) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO secrets (id, file_id, key_name, value_enc) VALUES (?, ?, ?, ?)`,
		secret.ID, secret.FileID, secret.KeyName, secret.ValueEnc)
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}

// UpdateSecret updates an existing secret's value.
func (r *secretRepository) UpdateSecret(ctx context.Context, secret *domain.Secret) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	_, err := r.db.ExecContext(ctx,
		`UPDATE secrets SET value_enc = ? WHERE id = ? AND deleted_at = 0`,
		secret.ValueEnc, secret.ID)
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}

// DeleteSecret soft-deletes a secret.
func (r *secretRepository) DeleteSecret(ctx context.Context, id string) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	_, err := r.db.ExecContext(ctx,
		`UPDATE secrets SET deleted_at = strftime('%s','now') WHERE id = ? AND deleted_at = 0`, id)
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}

// DeleteSecretsByFileID soft-deletes all secrets in a file.
func (r *secretRepository) DeleteSecretsByFileID(ctx context.Context, fileID string) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	_, err := r.db.ExecContext(ctx,
		`UPDATE secrets SET deleted_at = strftime('%s','now') WHERE file_id = ? AND deleted_at = 0`, fileID)
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}
