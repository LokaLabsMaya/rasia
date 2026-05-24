package inbound

//go:generate mockgen -source=export.go -destination=../../../mocks/inbound/mock_export.go -package=mocks_inbound

import (
	"context"

	"github.com/lokalabsmaya/rasia/internal/core/domain"
)

type Export interface {
	// ExportByPath resolves a slash-separated path and returns all decrypted files under that namespace.
	ExportByPath(ctx context.Context, path string) ([]*domain.ExportFile, error)
}
