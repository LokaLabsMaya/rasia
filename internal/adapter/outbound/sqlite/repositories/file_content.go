package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/redhajuanda/komon/fail"
	"github.com/redhajuanda/komon/tracer"
	"github.com/lokalabsmaya/rasia/internal/core/domain"
)

// fileContentRepository implements repositories.FileContent for SQLite.
type fileContentRepository struct {
	db querier
}

// NewFileContentRepository creates a new SQLite FileContentRepository.
func NewFileContentRepository(db querier) *fileContentRepository {
	return &fileContentRepository{db: db}
}

// GetFileContent returns the content record for a file.
func (r *fileContentRepository) GetFileContent(ctx context.Context, fileID string) (*domain.FileContent, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	var fc domain.FileContent
	err := r.db.QueryRowContext(ctx,
		`SELECT id, file_id, content, updated_at, deleted_at
		 FROM file_contents WHERE file_id = ? AND deleted_at = 0`, fileID).
		Scan(&fc.ID, &fc.FileID, &fc.Content, &fc.UpdatedAt, &fc.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fail.Wrap(err).WithFailure(fail.ErrNotFound)
		}
		return nil, fail.Wrap(err)
	}
	return &fc, nil
}

// UpsertFileContent inserts or replaces the content record.
func (r *fileContentRepository) UpsertFileContent(ctx context.Context, fc *domain.FileContent) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO file_contents (id, file_id, content, deleted_at)
		 VALUES (?, ?, ?, 0)
		 ON CONFLICT(file_id, deleted_at) DO UPDATE SET content = excluded.content, updated_at = datetime('now')`,
		fc.ID, fc.FileID, fc.Content)
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}

// DeleteFileContent soft-deletes the content for a file.
func (r *fileContentRepository) DeleteFileContent(ctx context.Context, fileID string) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	_, err := r.db.ExecContext(ctx,
		`UPDATE file_contents SET deleted_at = strftime('%s','now') WHERE file_id = ? AND deleted_at = 0`, fileID)
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}

// GetFileContentsByFileIDs returns content records for multiple file IDs.
func (r *fileContentRepository) GetFileContentsByFileIDs(ctx context.Context, fileIDs []string) ([]*domain.FileContent, error) {
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
		`SELECT id, file_id, content, updated_at, deleted_at
		 FROM file_contents WHERE file_id IN (`+placeholders+`) AND deleted_at = 0`,
		args...)
	if err != nil {
		return nil, fail.Wrap(err)
	}
	defer rows.Close()

	var result []*domain.FileContent
	for rows.Next() {
		var fc domain.FileContent
		if err := rows.Scan(&fc.ID, &fc.FileID, &fc.Content, &fc.UpdatedAt, &fc.DeletedAt); err != nil {
			return nil, fail.Wrap(err)
		}
		result = append(result, &fc)
	}
	if err := rows.Err(); err != nil {
		return nil, fail.Wrap(err)
	}
	return result, nil
}
