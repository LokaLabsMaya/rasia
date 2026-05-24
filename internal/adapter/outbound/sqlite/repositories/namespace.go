package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/redhajuanda/komon/fail"
	"github.com/redhajuanda/komon/tracer"
	"github.com/lokalabsmaya/rasia/internal/core/domain"
	"github.com/lokalabsmaya/rasia/shared/failure"
)

// namespaceRepository implements repositories.Namespace for SQLite.
type namespaceRepository struct {
	db querier
}

// NewNamespaceRepository creates a new SQLite NamespaceRepository.
func NewNamespaceRepository(db querier) *namespaceRepository {
	return &namespaceRepository{db: db}
}

// GetAllNamespaces returns all non-deleted namespaces.
func (r *namespaceRepository) GetAllNamespaces(ctx context.Context) ([]*domain.Namespace, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, parent_id, name, created_at, updated_at, deleted_at
		 FROM namespaces WHERE deleted_at = 0 ORDER BY parent_id, name`)
	if err != nil {
		return nil, fail.Wrap(err)
	}
	defer rows.Close()

	var result []*domain.Namespace
	for rows.Next() {
		var ns domain.Namespace
		if err := rows.Scan(&ns.ID, &ns.ParentID, &ns.Name, &ns.CreatedAt, &ns.UpdatedAt, &ns.DeletedAt); err != nil {
			return nil, fail.Wrap(err)
		}
		result = append(result, &ns)
	}
	if err := rows.Err(); err != nil {
		return nil, fail.Wrap(err)
	}
	return result, nil
}

// GetNamespaceByID returns a namespace by its ID.
func (r *namespaceRepository) GetNamespaceByID(ctx context.Context, id string) (*domain.Namespace, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	var ns domain.Namespace
	err := r.db.QueryRowContext(ctx,
		`SELECT id, parent_id, name, created_at, updated_at, deleted_at
		 FROM namespaces WHERE id = ? AND deleted_at = 0`, id).
		Scan(&ns.ID, &ns.ParentID, &ns.Name, &ns.CreatedAt, &ns.UpdatedAt, &ns.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fail.Wrap(err).WithFailure(failure.ErrNamespaceNotFound)
		}
		return nil, fail.Wrap(err)
	}
	return &ns, nil
}

// GetNamespaceByParentAndName returns a namespace matching parent+name.
func (r *namespaceRepository) GetNamespaceByParentAndName(ctx context.Context, parentID *string, name string) (*domain.Namespace, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	var (
		ns  domain.Namespace
		row *sql.Row
	)
	if parentID == nil {
		row = r.db.QueryRowContext(ctx,
			`SELECT id, parent_id, name, created_at, updated_at, deleted_at
			 FROM namespaces WHERE parent_id IS NULL AND name = ? AND deleted_at = 0`, name)
	} else {
		row = r.db.QueryRowContext(ctx,
			`SELECT id, parent_id, name, created_at, updated_at, deleted_at
			 FROM namespaces WHERE parent_id = ? AND name = ? AND deleted_at = 0`, *parentID, name)
	}

	err := row.Scan(&ns.ID, &ns.ParentID, &ns.Name, &ns.CreatedAt, &ns.UpdatedAt, &ns.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fail.Wrap(err).WithFailure(failure.ErrNamespaceNotFound)
		}
		return nil, fail.Wrap(err)
	}
	return &ns, nil
}

// CreateNamespace inserts a new namespace.
func (r *namespaceRepository) CreateNamespace(ctx context.Context, ns *domain.Namespace) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO namespaces (id, parent_id, name) VALUES (?, ?, ?)`,
		ns.ID, ns.ParentID, ns.Name)
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}

// DeleteNamespace soft-deletes a namespace by ID.
func (r *namespaceRepository) DeleteNamespace(ctx context.Context, id string) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	_, err := r.db.ExecContext(ctx,
		`UPDATE namespaces SET deleted_at = strftime('%s','now') WHERE id = ? AND deleted_at = 0`, id)
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}

// GetDescendantIDs returns IDs of all descendants using a recursive CTE.
func (r *namespaceRepository) GetDescendantIDs(ctx context.Context, id string) ([]string, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	rows, err := r.db.QueryContext(ctx, `
		WITH RECURSIVE descendants AS (
			SELECT id FROM namespaces WHERE parent_id = ? AND deleted_at = 0
			UNION ALL
			SELECT n.id FROM namespaces n
			INNER JOIN descendants d ON n.parent_id = d.id
			WHERE n.deleted_at = 0
		)
		SELECT id FROM descendants`, id)
	if err != nil {
		return nil, fail.Wrap(err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fail.Wrap(err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fail.Wrap(err)
	}
	return ids, nil
}
