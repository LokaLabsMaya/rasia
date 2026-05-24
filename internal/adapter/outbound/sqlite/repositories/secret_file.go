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

// secretFileRepository implements repositories.SecretFile for SQLite.
type secretFileRepository struct {
	db querier
}

// NewSecretFileRepository creates a new SQLite SecretFileRepository.
func NewSecretFileRepository(db querier) *secretFileRepository {
	return &secretFileRepository{db: db}
}

// GetSecretFileByID returns a file by its ID.
func (r *secretFileRepository) GetSecretFileByID(ctx context.Context, id string) (*domain.SecretFile, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	var f domain.SecretFile
	err := r.db.QueryRowContext(ctx,
		`SELECT id, namespace_id, name, ext, created_at, updated_at, deleted_at
		 FROM secret_files WHERE id = ? AND deleted_at = 0`, id).
		Scan(&f.ID, &f.NamespaceID, &f.Name, &f.Ext, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fail.Wrap(err).WithFailure(failure.ErrSecretFileNotFound)
		}
		return nil, fail.Wrap(err)
	}
	return &f, nil
}

// ListSecretFiles returns all non-deleted files in a namespace.
func (r *secretFileRepository) ListSecretFiles(ctx context.Context, namespaceID string) ([]*domain.SecretFile, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, namespace_id, name, ext, created_at, updated_at, deleted_at
		 FROM secret_files WHERE namespace_id = ? AND deleted_at = 0 ORDER BY name`, namespaceID)
	if err != nil {
		return nil, fail.Wrap(err)
	}
	defer rows.Close()

	var result []*domain.SecretFile
	for rows.Next() {
		var f domain.SecretFile
		if err := rows.Scan(&f.ID, &f.NamespaceID, &f.Name, &f.Ext, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt); err != nil {
			return nil, fail.Wrap(err)
		}
		result = append(result, &f)
	}
	if err := rows.Err(); err != nil {
		return nil, fail.Wrap(err)
	}
	return result, nil
}

// ListSecretFilesByNamespaceIDs returns files for multiple namespace IDs.
func (r *secretFileRepository) ListSecretFilesByNamespaceIDs(ctx context.Context, namespaceIDs []string) ([]*domain.SecretFile, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	if len(namespaceIDs) == 0 {
		return nil, nil
	}

	placeholders := strings.Repeat("?,", len(namespaceIDs))
	placeholders = placeholders[:len(placeholders)-1]

	args := make([]any, len(namespaceIDs))
	for i, id := range namespaceIDs {
		args[i] = id
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, namespace_id, name, ext, created_at, updated_at, deleted_at
		 FROM secret_files WHERE namespace_id IN (`+placeholders+`) AND deleted_at = 0 ORDER BY namespace_id, name`,
		args...)
	if err != nil {
		return nil, fail.Wrap(err)
	}
	defer rows.Close()

	var result []*domain.SecretFile
	for rows.Next() {
		var f domain.SecretFile
		if err := rows.Scan(&f.ID, &f.NamespaceID, &f.Name, &f.Ext, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt); err != nil {
			return nil, fail.Wrap(err)
		}
		result = append(result, &f)
	}
	if err := rows.Err(); err != nil {
		return nil, fail.Wrap(err)
	}
	return result, nil
}

// CreateSecretFile inserts a new file.
func (r *secretFileRepository) CreateSecretFile(ctx context.Context, file *domain.SecretFile) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO secret_files (id, namespace_id, name, ext) VALUES (?, ?, ?, ?)`,
		file.ID, file.NamespaceID, file.Name, file.Ext)
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}

// DeleteSecretFile soft-deletes a file.
func (r *secretFileRepository) DeleteSecretFile(ctx context.Context, id string) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	_, err := r.db.ExecContext(ctx,
		`UPDATE secret_files SET deleted_at = strftime('%s','now') WHERE id = ? AND deleted_at = 0`, id)
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}

// DeleteSecretFilesByNamespaceID soft-deletes all files in a namespace.
func (r *secretFileRepository) DeleteSecretFilesByNamespaceID(ctx context.Context, namespaceID string) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	_, err := r.db.ExecContext(ctx,
		`UPDATE secret_files SET deleted_at = strftime('%s','now') WHERE namespace_id = ? AND deleted_at = 0`, namespaceID)
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}
