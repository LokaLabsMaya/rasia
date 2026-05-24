package namespace

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

// Service implements inbound.Namespace.
type Service struct {
	cfg  *configs.Config
	log  logger.Logger
	repo outbound.Repository
}

var _ inbound.Namespace = (*Service)(nil)

// NewService creates a new namespace service.
func NewService(cfg *configs.Config, log logger.Logger, repo outbound.Repository) *Service {
	return &Service{cfg: cfg, log: log, repo: repo}
}

// GetTree returns the full namespace tree.
func (s *Service) GetTree(ctx context.Context) ([]*domain.Namespace, error) {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	var (
		repoNS = s.repo.GetNamespaceRepository()
	)

	all, err := repoNS.GetAllNamespaces(ctx)
	if err != nil {
		return nil, fail.Wrap(err)
	}

	return buildTree(all), nil
}

// buildTree converts a flat list into a nested tree.
func buildTree(flat []*domain.Namespace) []*domain.Namespace {
	index := make(map[string]*domain.Namespace, len(flat))
	for _, ns := range flat {
		index[ns.ID] = ns
	}

	var roots []*domain.Namespace
	for _, ns := range flat {
		if ns.ParentID == nil {
			roots = append(roots, ns)
		} else {
			if parent, ok := index[*ns.ParentID]; ok {
				parent.Children = append(parent.Children, ns)
			}
		}
	}
	return roots
}

// CreateNamespace creates a new namespace node.
func (s *Service) CreateNamespace(ctx context.Context, ns *domain.Namespace) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	var (
		repoNS = s.repo.GetNamespaceRepository()
	)

	err := repoNS.CreateNamespace(ctx, ns)
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}

// DeleteNamespace soft-deletes a namespace and all its descendants.
func (s *Service) DeleteNamespace(ctx context.Context, id string) error {
	ctx, span := tracer.Trace(ctx)
	defer span.End()

	_, err := s.repo.DoInTransaction(ctx, func(repo outbound.Repository) (any, error) {
		var (
			repoNS   = repo.GetNamespaceRepository()
			repoFile = repo.GetSecretFileRepository()
			repoSec  = repo.GetSecretRepository()
			repoFC   = repo.GetFileContentRepository()
		)

		ids, err := repoNS.GetDescendantIDs(ctx, id)
		if err != nil {
			return nil, err
		}
		// include the target itself
		ids = append(ids, id)

		for _, nsID := range ids {
			files, err := repoFile.ListSecretFiles(ctx, nsID)
			if err != nil {
				return nil, err
			}
			for _, f := range files {
				if err := repoSec.DeleteSecretsByFileID(ctx, f.ID); err != nil {
					return nil, err
				}
				if err := repoFC.DeleteFileContent(ctx, f.ID); err != nil {
					return nil, err
				}
			}
			if err := repoFile.DeleteSecretFilesByNamespaceID(ctx, nsID); err != nil {
				return nil, err
			}
			if err := repoNS.DeleteNamespace(ctx, nsID); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		return fail.Wrap(err)
	}
	return nil
}
