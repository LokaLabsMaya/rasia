package filecontent

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
	"github.com/oklog/ulid/v2"
)

// Service implements inbound.FileContent.
type Service struct {
	cfg    *configs.Config
	log    logger.Logger
	repo   outbound.Repository
	encKey []byte
}

var _ inbound.FileContent = (*Service)(nil)

// NewService creates a new file content service.
func NewService(cfg *configs.Config, log logger.Logger, repo outbound.Repository) *Service {
	return &Service{
		cfg:    cfg,
		log:    log,
		repo:   repo,
		encKey: utils.DeriveKey(cfg.Crypto.Key),
	}
}

// GetContent decrypts and returns the raw file content.
func (s *Service) GetContent(ctx context.Context, fileID string) (string, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	fc, err := s.repo.GetFileContentRepository().GetFileContent(ctx, fileID)
	if err != nil {
		return "", fail.Wrap(err)
	}

	plain, err := utils.Decrypt(s.encKey, fc.Content)
	if err != nil {
		return "", fail.Wrap(err)
	}
	return plain, nil
}

// SaveContent encrypts and upserts the file content.
func (s *Service) SaveContent(ctx context.Context, fileID, content string) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	enc, err := utils.Encrypt(s.encKey, content)
	if err != nil {
		return fail.Wrap(err)
	}

	fc := &domain.FileContent{
		ID:      ulid.Make().String(),
		FileID:  fileID,
		Content: enc,
	}

	if err := s.repo.GetFileContentRepository().UpsertFileContent(ctx, fc); err != nil {
		return fail.Wrap(err)
	}
	return nil
}
