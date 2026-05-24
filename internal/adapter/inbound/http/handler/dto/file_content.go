package dto

import "github.com/go-playground/validator/v10"

type ReqGetFileContent struct {
	FileID string `uri:"file_id" validate:"required"`
}

func (r *ReqGetFileContent) Validate() error {
	return validator.New().Struct(r)
}

type ReqSaveFileContent struct {
	FileID  string `uri:"file_id" validate:"required" swaggerignore:"true"`
	Content string `json:"content" validate:"required"`
}

func (r *ReqSaveFileContent) Validate() error {
	return validator.New().Struct(r)
}

type ResFileContent struct {
	Content string `json:"content"`
}

type ReqExport struct {
	Path string `query:"path" validate:"required"`
}

func (r *ReqExport) Validate() error {
	return validator.New().Struct(r)
}
