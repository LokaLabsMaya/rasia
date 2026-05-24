package export

import (
	"context"
	"strings"

	"github.com/redhajuanda/komon/fail"
	"github.com/redhajuanda/komon/logger"
	"github.com/redhajuanda/komon/tracer"
	"github.com/lokalabsmaya/rasia/configs"
	"github.com/lokalabsmaya/rasia/internal/core/domain"
	"github.com/lokalabsmaya/rasia/internal/core/port/inbound"
	"github.com/lokalabsmaya/rasia/internal/core/port/outbound"
	"github.com/lokalabsmaya/rasia/shared/failure"
	"github.com/lokalabsmaya/rasia/shared/utils"
)

// Service implements inbound.Export.
type Service struct {
	cfg    *configs.Config
	log    logger.Logger
	repo   outbound.Repository
	encKey []byte
}

var _ inbound.Export = (*Service)(nil)

// NewService creates a new export service.
func NewService(cfg *configs.Config, log logger.Logger, repo outbound.Repository) *Service {
	return &Service{
		cfg:    cfg,
		log:    log,
		repo:   repo,
		encKey: utils.DeriveKey(cfg.Crypto.Key),
	}
}

// ExportByPath resolves a slash-separated path and returns all decrypted files under that namespace.
func (s *Service) ExportByPath(ctx context.Context, path string) ([]*domain.ExportFile, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	var (
		repoNS   = s.repo.GetNamespaceRepository()
		repoFile = s.repo.GetSecretFileRepository()
		repoSec  = s.repo.GetSecretRepository()
		repoFC   = s.repo.GetFileContentRepository()
	)

	// resolve path segments to a namespace ID
	segments := splitPath(path)
	if len(segments) == 0 {
		return nil, fail.New("path is required").WithFailure(fail.ErrBadRequest)
	}

	var parentID *string
	for _, seg := range segments {
		ns, err := repoNS.GetNamespaceByParentAndName(ctx, parentID, seg)
		if err != nil {
			return nil, fail.Wrap(err).WithFailure(failure.ErrNamespaceNotFound)
		}
		id := ns.ID
		parentID = &id
	}

	// parentID now holds the resolved namespace ID
	namespaceID := *parentID

	// get descendant namespace IDs (include self)
	descendantIDs, err := repoNS.GetDescendantIDs(ctx, namespaceID)
	if err != nil {
		return nil, fail.Wrap(err)
	}
	allIDs := append([]string{namespaceID}, descendantIDs...)

	// load all files for these namespaces
	files, err := repoFile.ListSecretFilesByNamespaceIDs(ctx, allIDs)
	if err != nil {
		return nil, fail.Wrap(err)
	}
	if len(files) == 0 {
		return []*domain.ExportFile{}, nil
	}

	// collect file IDs
	fileIDs := make([]string, 0, len(files))
	for _, f := range files {
		fileIDs = append(fileIDs, f.ID)
	}

	// load all secrets and contents in bulk
	secrets, err := repoSec.ListSecretsByFileIDs(ctx, fileIDs)
	if err != nil {
		return nil, fail.Wrap(err)
	}
	contents, err := repoFC.GetFileContentsByFileIDs(ctx, fileIDs)
	if err != nil {
		return nil, fail.Wrap(err)
	}

	// index secrets and contents by file ID
	secretsByFile := make(map[string][]*domain.Secret, len(secrets))
	for _, sec := range secrets {
		secretsByFile[sec.FileID] = append(secretsByFile[sec.FileID], sec)
	}
	contentByFile := make(map[string]*domain.FileContent, len(contents))
	for _, fc := range contents {
		contentByFile[fc.FileID] = fc
	}

	// build export result
	result := make([]*domain.ExportFile, 0, len(files))
	for _, f := range files {
		ef := &domain.ExportFile{
			Name: f.Name,
			Ext:  f.Ext,
		}

		if f.Ext == domain.SecretFileExtEnv {
			ef.Secrets = make(map[string]string)
			for _, sec := range secretsByFile[f.ID] {
				plain, err := utils.Decrypt(s.encKey, sec.ValueEnc)
				if err != nil {
					return nil, fail.Wrap(err)
				}
				ef.Secrets[sec.KeyName] = plain
			}
		} else {
			if fc, ok := contentByFile[f.ID]; ok {
				plain, err := utils.Decrypt(s.encKey, fc.Content)
				if err != nil {
					return nil, fail.Wrap(err)
				}
				ef.Content = plain
			}
		}

		result = append(result, ef)
	}

	return result, nil
}

// splitPath splits a slash-separated path into non-empty segments.
func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}
