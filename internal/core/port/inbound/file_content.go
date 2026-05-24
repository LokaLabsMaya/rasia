package inbound

//go:generate mockgen -source=file_content.go -destination=../../../mocks/inbound/mock_file_content.go -package=mocks_inbound

import "context"

type FileContent interface {
	// GetContent decrypts and returns the raw file content.
	GetContent(ctx context.Context, fileID string) (string, error)
	// SaveContent encrypts and upserts the file content.
	SaveContent(ctx context.Context, fileID, content string) error
}
