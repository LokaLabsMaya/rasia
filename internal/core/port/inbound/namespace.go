package inbound

//go:generate mockgen -source=namespace.go -destination=../../../mocks/inbound/mock_namespace.go -package=mocks_inbound

import (
	"context"

	"github.com/lokalabsmaya/rasia/internal/core/domain"
)

type Namespace interface {
	// GetTree returns the full namespace tree.
	GetTree(ctx context.Context) ([]*domain.Namespace, error)
	// CreateNamespace creates a new namespace node.
	CreateNamespace(ctx context.Context, ns *domain.Namespace) error
	// DeleteNamespace soft-deletes a namespace and all its descendants.
	DeleteNamespace(ctx context.Context, id string) error
}
