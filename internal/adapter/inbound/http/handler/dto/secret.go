package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"github.com/lokalabsmaya/rasia/internal/core/domain"
)

type ReqListSecrets struct {
	FileID string `uri:"file_id" validate:"required"`
}

func (r *ReqListSecrets) Validate() error {
	return validator.New().Struct(r)
}

type ReqAddSecret struct {
	FileID   string `uri:"file_id" validate:"required" swaggerignore:"true"`
	KeyName  string `json:"key_name" validate:"required,max=255"`
	Value    string `json:"value" validate:"required"`
}

func (r *ReqAddSecret) Validate() error {
	return validator.New().Struct(r)
}

func (r *ReqAddSecret) Transform() *domain.Secret {
	return &domain.Secret{
		ID:       ulid.Make().String(),
		FileID:   r.FileID,
		KeyName:  r.KeyName,
		ValueEnc: r.Value,
	}
}

type ReqUpdateSecret struct {
	ID    string `uri:"id" validate:"required" swaggerignore:"true"`
	Value string `json:"value" validate:"required"`
}

func (r *ReqUpdateSecret) Validate() error {
	return validator.New().Struct(r)
}

func (r *ReqUpdateSecret) Transform() *domain.Secret {
	return &domain.Secret{
		ID:       r.ID,
		ValueEnc: r.Value,
	}
}

type ReqRevealSecret struct {
	ID string `uri:"id" validate:"required"`
}

func (r *ReqRevealSecret) Validate() error {
	return validator.New().Struct(r)
}

type ReqDeleteSecret struct {
	ID string `uri:"id" validate:"required"`
}

func (r *ReqDeleteSecret) Validate() error {
	return validator.New().Struct(r)
}

type ResSecret struct {
	ID        string `json:"id"`
	FileID    string `json:"file_id"`
	KeyName   string `json:"key_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (r *ResSecret) Transform(s *domain.Secret) {
	r.ID = s.ID
	r.FileID = s.FileID
	r.KeyName = s.KeyName
	r.CreatedAt = s.CreatedAt
	r.UpdatedAt = s.UpdatedAt
}

type ResSecretList []ResSecret

func (r *ResSecretList) Transform(secrets []*domain.Secret) {
	for _, s := range secrets {
		var item ResSecret
		item.Transform(s)
		*r = append(*r, item)
	}
}

type ResRevealSecret struct {
	Value string `json:"value"`
}
