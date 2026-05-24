package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"github.com/lokalabsmaya/rasia/internal/core/domain"
)

type ReqListSecretFiles struct {
	NamespaceID string `uri:"namespace_id" validate:"required"`
}

func (r *ReqListSecretFiles) Validate() error {
	return validator.New().Struct(r)
}

type ReqCreateSecretFile struct {
	NamespaceID string                `uri:"namespace_id" validate:"required" swaggerignore:"true"`
	Name        string                `json:"name" validate:"required,max=255"`
	Ext         domain.SecretFileExt  `json:"ext" validate:"required,oneof=env yaml json txt"`
}

func (r *ReqCreateSecretFile) Validate() error {
	return validator.New().Struct(r)
}

func (r *ReqCreateSecretFile) Transform() *domain.SecretFile {
	return &domain.SecretFile{
		ID:          ulid.Make().String(),
		NamespaceID: r.NamespaceID,
		Name:        r.Name,
		Ext:         r.Ext,
	}
}

type ReqDeleteSecretFile struct {
	ID string `uri:"id" validate:"required"`
}

func (r *ReqDeleteSecretFile) Validate() error {
	return validator.New().Struct(r)
}

type ResSecretFile struct {
	ID          string               `json:"id"`
	NamespaceID string               `json:"namespace_id"`
	Name        string               `json:"name"`
	Ext         domain.SecretFileExt `json:"ext"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}

func (r *ResSecretFile) Transform(f *domain.SecretFile) {
	r.ID = f.ID
	r.NamespaceID = f.NamespaceID
	r.Name = f.Name
	r.Ext = f.Ext
	r.CreatedAt = f.CreatedAt
	r.UpdatedAt = f.UpdatedAt
}

type ResSecretFileList []ResSecretFile

func (r *ResSecretFileList) Transform(files []*domain.SecretFile) {
	for _, f := range files {
		var item ResSecretFile
		item.Transform(f)
		*r = append(*r, item)
	}
}
