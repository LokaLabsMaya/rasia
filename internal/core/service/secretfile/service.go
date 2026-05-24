package secretfile

import (
	"context"

	"github.com/redhajuanda/komon/fail"
	"github.com/redhajuanda/komon/logger"
	"github.com/redhajuanda/komon/tracer"
	"github.com/lokalabsmaya/rasia/configs"
	"github.com/lokalabsmaya/rasia/internal/core/domain"
	"github.com/lokalabsmaya/rasia/internal/core/port/inbound"
	"github.com/lokalabsmaya/rasia/internal/core/port/outbound"
)

// Service implements inbound.SecretFile.
type Service struct {
	cfg  *configs.Config
	log  logger.Logger
	repo outbound.Repository
}

var _ inbound.SecretFile = (*Service)(nil)

// NewService creates a new secret file service.
func NewService(cfg *configs.Config, log logger.Logger, repo outbound.Repository) *Service {
	return &Service{cfg: cfg, log: log, repo: repo}
}

// ListSecretFiles returns all files in a namespace.
func (s *Service) ListSecretFiles(ctx context.Context, namespaceID string) ([]*domain.SecretFile, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	files, err := s.repo.GetSecretFileRepository().ListSecretFiles(ctx, namespaceID)
	if err != nil {
		return nil, fail.Wrap(err)
	}
	return files, nil
}

// CreateSecretFile creates a new file in a namespace.
func (s *Service) CreateSecretFile(ctx context.Context, file *domain.SecretFile) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	err := s.repo.GetSecretFileRepository().CreateSecretFile(ctx, file)
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}

// DeleteSecretFile soft-deletes a file and its secrets/content.
func (s *Service) DeleteSecretFile(ctx context.Context, id string) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	_, err := s.repo.DoInTransaction(ctx, func(repo outbound.Repository) (any, error) {
		var (
			repoFile = repo.GetSecretFileRepository()
			repoSec  = repo.GetSecretRepository()
			repoFC   = repo.GetFileContentRepository()
		)

		if err := repoSec.DeleteSecretsByFileID(ctx, id); err != nil {
			return nil, err
		}
		if err := repoFC.DeleteFileContent(ctx, id); err != nil {
			return nil, err
		}
		if err := repoFile.DeleteSecretFile(ctx, id); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}
