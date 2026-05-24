package repositories

//go:generate mockgen -source=secret_file.go -destination=../../../../mocks/outbound/repositories/mock_secret_file_repository.go -package=mocks_outbound_repositories

import (
	"context"

	"github.com/lokalabsmaya/rasia/internal/core/domain"
)

type SecretFile interface {
	// GetSecretFileByID returns a file by its ID.
	GetSecretFileByID(ctx context.Context, id string) (*domain.SecretFile, error)
	// ListSecretFiles returns all non-deleted files in a namespace.
	ListSecretFiles(ctx context.Context, namespaceID string) ([]*domain.SecretFile, error)
	// ListSecretFilesByNamespaceIDs returns files for multiple namespace IDs.
	ListSecretFilesByNamespaceIDs(ctx context.Context, namespaceIDs []string) ([]*domain.SecretFile, error)
	// CreateSecretFile inserts a new file.
	CreateSecretFile(ctx context.Context, file *domain.SecretFile) error
	// DeleteSecretFile soft-deletes a file.
	DeleteSecretFile(ctx context.Context, id string) error
	// DeleteSecretFilesByNamespaceID soft-deletes all files in a namespace.
	DeleteSecretFilesByNamespaceID(ctx context.Context, namespaceID string) error
}
