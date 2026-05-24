package inbound

//go:generate mockgen -source=secret_file.go -destination=../../../mocks/inbound/mock_secret_file.go -package=mocks_inbound

import (
	"context"

	"github.com/lokalabsmaya/rasia/internal/core/domain"
)

type SecretFile interface {
	// ListSecretFiles returns all files in a namespace.
	ListSecretFiles(ctx context.Context, namespaceID string) ([]*domain.SecretFile, error)
	// CreateSecretFile creates a new file in a namespace.
	CreateSecretFile(ctx context.Context, file *domain.SecretFile) error
	// DeleteSecretFile soft-deletes a file and its secrets/content.
	DeleteSecretFile(ctx context.Context, id string) error
}
