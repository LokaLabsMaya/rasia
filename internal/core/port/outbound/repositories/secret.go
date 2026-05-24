package repositories

//go:generate mockgen -source=secret.go -destination=../../../../mocks/outbound/repositories/mock_secret_repository.go -package=mocks_outbound_repositories

import (
	"context"

	"github.com/lokalabsmaya/rasia/internal/core/domain"
)

type Secret interface {
	// GetSecretByID returns a secret by its ID.
	GetSecretByID(ctx context.Context, id string) (*domain.Secret, error)
	// ListSecrets returns all non-deleted secrets in a file.
	ListSecrets(ctx context.Context, fileID string) ([]*domain.Secret, error)
	// ListSecretsByFileIDs returns secrets for multiple file IDs.
	ListSecretsByFileIDs(ctx context.Context, fileIDs []string) ([]*domain.Secret, error)
	// CreateSecret inserts a new secret.
	CreateSecret(ctx context.Context, secret *domain.Secret) error
	// UpdateSecret updates an existing secret's value.
	UpdateSecret(ctx context.Context, secret *domain.Secret) error
	// DeleteSecret soft-deletes a secret.
	DeleteSecret(ctx context.Context, id string) error
	// DeleteSecretsByFileID soft-deletes all secrets in a file.
	DeleteSecretsByFileID(ctx context.Context, fileID string) error
}
