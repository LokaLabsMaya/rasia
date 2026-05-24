package inbound

//go:generate mockgen -source=secret.go -destination=../../../mocks/inbound/mock_secret.go -package=mocks_inbound

import (
	"context"

	"github.com/lokalabsmaya/rasia/internal/core/domain"
)

type Secret interface {
	// ListSecrets returns all secrets in a file (values masked).
	ListSecrets(ctx context.Context, fileID string) ([]*domain.Secret, error)
	// AddSecret encrypts and stores a new key-value pair.
	AddSecret(ctx context.Context, secret *domain.Secret) error
	// UpdateSecret encrypts and updates an existing key-value pair.
	UpdateSecret(ctx context.Context, secret *domain.Secret) error
	// RevealSecret decrypts and returns the plaintext value.
	RevealSecret(ctx context.Context, id string) (string, error)
	// DeleteSecret soft-deletes a secret.
	DeleteSecret(ctx context.Context, id string) error
}
