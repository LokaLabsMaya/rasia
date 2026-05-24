package handler

import (
	"github.com/redhajuanda/komon/fail"
	"github.com/redhajuanda/komon/logger"
	"github.com/lokalabsmaya/rasia/configs"
	"github.com/lokalabsmaya/rasia/internal/adapter/inbound/http/handler/dto"
	"github.com/lokalabsmaya/rasia/internal/adapter/inbound/http/response"
	"github.com/lokalabsmaya/rasia/internal/core/port/inbound"

	"github.com/gofiber/fiber/v3"
)

// FileContentHandler handles HTTP requests for file content operations.
type FileContentHandler struct {
	cfg    *configs.Config
	log    logger.Logger
	svcFC  inbound.FileContent
	svcExp inbound.Export
}

// NewFileContentHandler creates a new FileContentHandler.
func NewFileContentHandler(cfg *configs.Config, log logger.Logger, svcFC inbound.FileContent, svcExp inbound.Export) *FileContentHandler {
	return &FileContentHandler{cfg: cfg, log: log, svcFC: svcFC, svcExp: svcExp}
}

// RegisterRoutes registers file content HTTP routes.
func (h *FileContentHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/api/files/:file_id/content", h.GetContent)
	app.Put("/api/files/:file_id/content", h.SaveContent)
	app.Get("/api/export", h.Export)
}

// GetContent godoc
// @Summary      Get File Content
// @Description  Returns decrypted file content
// @Tags         FileContent
// @Param        file_id  path  string  true  "File ID"
// @Produce      json
// @Success      200  {object}  response.ResponseSuccess{data=dto.ResFileContent}
// @Router       /api/files/{file_id}/content [get]
func (h *FileContentHandler) GetContent(c fiber.Ctx) error {
	var (
		req dto.ReqGetFileContent
		ctx = c.Context()
	)

	if err := c.Bind().URI(&req); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := req.Validate(); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	content, err := h.svcFC.GetContent(ctx, req.FileID)
	if err != nil {
		return err
	}

	return response.SuccessOK(c, dto.ResFileContent{Content: content}, "Content retrieved successfully")
}

// SaveContent godoc
// @Summary      Save File Content
// @Description  Encrypts and upserts file content
// @Tags         FileContent
// @Accept       json
// @Produce      json
// @Param        file_id  path  string                 true  "File ID"
// @Param        body     body  dto.ReqSaveFileContent true  "Content"
// @Success      200  {object}  response.ResponseSuccess{}
// @Router       /api/files/{file_id}/content [put]
func (h *FileContentHandler) SaveContent(c fiber.Ctx) error {
	var (
		req dto.ReqSaveFileContent
		ctx = c.Context()
	)

	if err := c.Bind().URI(&req); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := c.Bind().Body(&req); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := req.Validate(); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := h.svcFC.SaveContent(ctx, req.FileID, req.Content); err != nil {
		return err
	}

	return response.SuccessOK(c, nil, "Content saved successfully")
}

// Export godoc
// @Summary      Export Secrets
// @Description  Resolves a namespace path and returns all decrypted files
// @Tags         Export
// @Param        path  query  string  true  "Namespace path (e.g. production/application)"
// @Produce      json
// @Success      200  {object}  response.ResponseSuccess{}
// @Router       /api/export [get]
func (h *FileContentHandler) Export(c fiber.Ctx) error {
	var (
		req dto.ReqExport
		ctx = c.Context()
	)

	if err := c.Bind().Query(&req); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := req.Validate(); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	files, err := h.svcExp.ExportByPath(ctx, req.Path)
	if err != nil {
		return err
	}

	return response.SuccessOK(c, files, "Export successful")
}
