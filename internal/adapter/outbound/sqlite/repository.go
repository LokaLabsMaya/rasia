package sqlite

import (
	"context"
	"database/sql"

	"github.com/redhajuanda/komon/fail"
	implRepo "github.com/lokalabsmaya/rasia/internal/adapter/outbound/sqlite/repositories"
	"github.com/lokalabsmaya/rasia/internal/core/port/outbound"
	portRepo "github.com/lokalabsmaya/rasia/internal/core/port/outbound/repositories"
)

// sqliteRepository implements outbound.Repository backed by SQLite.
type sqliteRepository struct {
	db *sql.DB
	tx *sql.Tx

	namespaceRepo   portRepo.Namespace
	secretFileRepo  portRepo.SecretFile
	secretRepo      portRepo.Secret
	fileContentRepo portRepo.FileContent
}

// Interface compliance check
var _ outbound.Repository = (*sqliteRepository)(nil)

// NewRepository creates a new SQLite-backed Repository.
func NewRepository(db *DB) *sqliteRepository {
	return &sqliteRepository{
		db:              db.Client,
		namespaceRepo:   implRepo.NewNamespaceRepository(db.Client),
		secretFileRepo:  implRepo.NewSecretFileRepository(db.Client),
		secretRepo:      implRepo.NewSecretRepository(db.Client),
		fileContentRepo: implRepo.NewFileContentRepository(db.Client),
	}
}

// GetNamespaceRepository returns the NamespaceRepository instance.
func (r *sqliteRepository) GetNamespaceRepository() portRepo.Namespace {
	if r.tx != nil {
		return implRepo.NewNamespaceRepository(r.tx)
	}
	return r.namespaceRepo
}

// GetSecretFileRepository returns the SecretFileRepository instance.
func (r *sqliteRepository) GetSecretFileRepository() portRepo.SecretFile {
	if r.tx != nil {
		return implRepo.NewSecretFileRepository(r.tx)
	}
	return r.secretFileRepo
}

// GetSecretRepository returns the SecretRepository instance.
func (r *sqliteRepository) GetSecretRepository() portRepo.Secret {
	if r.tx != nil {
		return implRepo.NewSecretRepository(r.tx)
	}
	return r.secretRepo
}

// GetFileContentRepository returns the FileContentRepository instance.
func (r *sqliteRepository) GetFileContentRepository() portRepo.FileContent {
	if r.tx != nil {
		return implRepo.NewFileContentRepository(r.tx)
	}
	return r.fileContentRepo
}

// DoInTransaction executes fn inside a SQLite transaction.
func (r *sqliteRepository) DoInTransaction(ctx context.Context, fn func(repo outbound.Repository) (any, error)) (out any, err error) {
	if r.tx != nil {
		return fn(r)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fail.Wrapf(err, "failed to begin transaction")
	}

	txRepo := &sqliteRepository{
		db: r.db,
		tx: tx,
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fail.Wrapf(commitErr, "failed to commit transaction")
		}
	}()

	out, err = fn(txRepo)
	if err != nil {
		return nil, fail.Wrapf(err, "failed to execute function in transaction")
	}
	return
}
