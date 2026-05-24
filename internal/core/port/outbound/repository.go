package outbound

//go:generate mockgen -source=repository.go -destination=../../../mocks/outbound/mock_repository.go -package=mocks_outbound

import (
	"context"

	"github.com/lokalabsmaya/rasia/internal/core/port/outbound/repositories"
)

type Repository interface {
	// DoInTransaction executes a function in a transaction
	DoInTransaction(ctx context.Context, fn func(repo Repository) (any, error)) (any, error)
	// GetNamespaceRepository returns the NamespaceRepository instance
	GetNamespaceRepository() repositories.Namespace
	// GetSecretFileRepository returns the SecretFileRepository instance
	GetSecretFileRepository() repositories.SecretFile
	// GetSecretRepository returns the SecretRepository instance
	GetSecretRepository() repositories.Secret
	// GetFileContentRepository returns the FileContentRepository instance
	GetFileContentRepository() repositories.FileContent
}
