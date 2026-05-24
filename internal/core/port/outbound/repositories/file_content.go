package repositories

//go:generate mockgen -source=file_content.go -destination=../../../../mocks/outbound/repositories/mock_file_content_repository.go -package=mocks_outbound_repositories

import (
	"context"

	"github.com/lokalabsmaya/rasia/internal/core/domain"
)

type FileContent interface {
	// GetFileContent returns the content record for a file.
	GetFileContent(ctx context.Context, fileID string) (*domain.FileContent, error)
	// UpsertFileContent inserts or updates the content record.
	UpsertFileContent(ctx context.Context, fc *domain.FileContent) error
	// DeleteFileContent soft-deletes the content for a file.
	DeleteFileContent(ctx context.Context, fileID string) error
	// GetFileContentsByFileIDs returns content records for multiple file IDs.
	GetFileContentsByFileIDs(ctx context.Context, fileIDs []string) ([]*domain.FileContent, error)
}
