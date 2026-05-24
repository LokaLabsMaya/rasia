package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"github.com/lokalabsmaya/rasia/internal/core/domain"
)

type ReqCreateNamespace struct {
	ParentID *string `json:"parent_id"`
	Name     string  `json:"name" validate:"required,max=255"`
}

func (r *ReqCreateNamespace) Validate() error {
	return validator.New().Struct(r)
}

func (r *ReqCreateNamespace) Transform() *domain.Namespace {
	return &domain.Namespace{
		ID:       ulid.Make().String(),
		ParentID: r.ParentID,
		Name:     r.Name,
	}
}

type ReqDeleteNamespace struct {
	ID string `uri:"id" validate:"required"`
}

func (r *ReqDeleteNamespace) Validate() error {
	return validator.New().Struct(r)
}

type ResNamespace struct {
	ID        string         `json:"id"`
	ParentID  *string        `json:"parent_id"`
	Name      string         `json:"name"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	Children  []ResNamespace `json:"children,omitempty"`
}

func (r *ResNamespace) Transform(ns *domain.Namespace) {
	r.ID = ns.ID
	r.ParentID = ns.ParentID
	r.Name = ns.Name
	r.CreatedAt = ns.CreatedAt
	r.UpdatedAt = ns.UpdatedAt
	for _, child := range ns.Children {
		var c ResNamespace
		c.Transform(child)
		r.Children = append(r.Children, c)
	}
}

type ResNamespaceTree []ResNamespace

func (r *ResNamespaceTree) Transform(roots []*domain.Namespace) {
	for _, ns := range roots {
		var item ResNamespace
		item.Transform(ns)
		*r = append(*r, item)
	}
}
