package repositories

//go:generate mockgen -source=namespace.go -destination=../../../../mocks/outbound/repositories/mock_namespace_repository.go -package=mocks_outbound_repositories

import (
	"context"

	"github.com/lokalabsmaya/rasia/internal/core/domain"
)

type Namespace interface {
	// GetAllNamespaces returns all non-deleted namespaces.
	GetAllNamespaces(ctx context.Context) ([]*domain.Namespace, error)
	// GetNamespaceByID returns a namespace by its ID.
	GetNamespaceByID(ctx context.Context, id string) (*domain.Namespace, error)
	// GetNamespaceByParentAndName returns a namespace matching parent+name.
	GetNamespaceByParentAndName(ctx context.Context, parentID *string, name string) (*domain.Namespace, error)
	// CreateNamespace inserts a new namespace.
	CreateNamespace(ctx context.Context, ns *domain.Namespace) error
	// DeleteNamespace soft-deletes a namespace by ID.
	DeleteNamespace(ctx context.Context, id string) error
	// GetDescendantIDs returns IDs of all descendants of a namespace.
	GetDescendantIDs(ctx context.Context, id string) ([]string, error)
}
