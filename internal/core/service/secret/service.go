package secret

import (
	"context"

	"github.com/redhajuanda/komon/fail"
	"github.com/redhajuanda/komon/logger"
	"github.com/redhajuanda/komon/tracer"
	"github.com/lokalabsmaya/rasia/configs"
	"github.com/lokalabsmaya/rasia/internal/core/domain"
	"github.com/lokalabsmaya/rasia/internal/core/port/inbound"
	"github.com/lokalabsmaya/rasia/internal/core/port/outbound"
	"github.com/lokalabsmaya/rasia/shared/utils"
)

// Service implements inbound.Secret.
type Service struct {
	cfg    *configs.Config
	log    logger.Logger
	repo   outbound.Repository
	encKey []byte
}

var _ inbound.Secret = (*Service)(nil)

// NewService creates a new secret service.
func NewService(cfg *configs.Config, log logger.Logger, repo outbound.Repository) *Service {
	return &Service{
		cfg:    cfg,
		log:    log,
		repo:   repo,
		encKey: utils.DeriveKey(cfg.Crypto.Key),
	}
}

// ListSecrets returns all secrets in a file with values masked.
func (s *Service) ListSecrets(ctx context.Context, fileID string) ([]*domain.Secret, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	secrets, err := s.repo.GetSecretRepository().ListSecrets(ctx, fileID)
	if err != nil {
		return nil, fail.Wrap(err)
	}

	// mask values
	for _, sec := range secrets {
		sec.ValueEnc = "***"
	}
	return secrets, nil
}

// AddSecret encrypts and stores a new key-value pair.
func (s *Service) AddSecret(ctx context.Context, secret *domain.Secret) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	enc, err := utils.Encrypt(s.encKey, secret.ValueEnc)
	if err != nil {
		return fail.Wrap(err)
	}
	secret.ValueEnc = enc

	if err := s.repo.GetSecretRepository().CreateSecret(ctx, secret); err != nil {
		return fail.Wrap(err)
	}
	return nil
}

// UpdateSecret encrypts and updates an existing key-value pair.
func (s *Service) UpdateSecret(ctx context.Context, secret *domain.Secret) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	enc, err := utils.Encrypt(s.encKey, secret.ValueEnc)
	if err != nil {
		return fail.Wrap(err)
	}
	secret.ValueEnc = enc

	if err := s.repo.GetSecretRepository().UpdateSecret(ctx, secret); err != nil {
		return fail.Wrap(err)
	}
	return nil
}

// RevealSecret decrypts and returns the plaintext value.
func (s *Service) RevealSecret(ctx context.Context, id string) (string, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	sec, err := s.repo.GetSecretRepository().GetSecretByID(ctx, id)
	if err != nil {
		return "", fail.Wrap(err)
	}

	plain, err := utils.Decrypt(s.encKey, sec.ValueEnc)
	if err != nil {
		return "", fail.Wrap(err)
	}
	return plain, nil
}

// DeleteSecret soft-deletes a secret.
func (s *Service) DeleteSecret(ctx context.Context, id string) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	if err := s.repo.GetSecretRepository().DeleteSecret(ctx, id); err != nil {
		return fail.Wrap(err)
	}
	return nil
}
