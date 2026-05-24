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

// SecretFileHandler handles HTTP requests for secret file operations.
type SecretFileHandler struct {
	cfg *configs.Config
	log logger.Logger
	svc inbound.SecretFile
}

// NewSecretFileHandler creates a new SecretFileHandler.
func NewSecretFileHandler(cfg *configs.Config, log logger.Logger, svc inbound.SecretFile) *SecretFileHandler {
	return &SecretFileHandler{cfg: cfg, log: log, svc: svc}
}

// RegisterRoutes registers secret file HTTP routes.
func (h *SecretFileHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/api/namespaces/:namespace_id/files", h.ListSecretFiles)
	app.Post("/api/namespaces/:namespace_id/files", h.CreateSecretFile)
	app.Delete("/api/files/:id", h.DeleteSecretFile)
}

// ListSecretFiles godoc
// @Summary      List Secret Files
// @Description  Returns all files in a namespace
// @Tags         SecretFiles
// @Param        namespace_id  path  string  true  "Namespace ID"
// @Produce      json
// @Success      200  {object}  response.ResponseSuccess{data=dto.ResSecretFileList}
// @Router       /api/namespaces/{namespace_id}/files [get]
func (h *SecretFileHandler) ListSecretFiles(c fiber.Ctx) error {
	var (
		req dto.ReqListSecretFiles
		ctx = c.Context()
	)

	if err := c.Bind().URI(&req); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := req.Validate(); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	files, err := h.svc.ListSecretFiles(ctx, req.NamespaceID)
	if err != nil {
		return err
	}

	var res dto.ResSecretFileList
	res.Transform(files)

	return response.SuccessOK(c, res, "Secret files retrieved successfully")
}

// CreateSecretFile godoc
// @Summary      Create Secret File
// @Description  Creates a new file in a namespace
// @Tags         SecretFiles
// @Accept       json
// @Produce      json
// @Param        namespace_id  path  string                 true  "Namespace ID"
// @Param        file          body  dto.ReqCreateSecretFile true  "File data"
// @Success      201  {object}  response.ResponseSuccess{}
// @Router       /api/namespaces/{namespace_id}/files [post]
func (h *SecretFileHandler) CreateSecretFile(c fiber.Ctx) error {
	var (
		req dto.ReqCreateSecretFile
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

	file := req.Transform()

	if err := h.svc.CreateSecretFile(ctx, file); err != nil {
		return err
	}

	return response.SuccessCreated(c, fiber.Map{"id": file.ID}, "Secret file created successfully")
}

// DeleteSecretFile godoc
// @Summary      Delete Secret File
// @Description  Soft-deletes a file and its secrets/content
// @Tags         SecretFiles
// @Param        id  path  string  true  "File ID"
// @Produce      json
// @Success      200  {object}  response.ResponseSuccess{}
// @Router       /api/files/{id} [delete]
func (h *SecretFileHandler) DeleteSecretFile(c fiber.Ctx) error {
	var (
		req dto.ReqDeleteSecretFile
		ctx = c.Context()
	)

	if err := c.Bind().URI(&req); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := req.Validate(); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := h.svc.DeleteSecretFile(ctx, req.ID); err != nil {
		return err
	}

	return response.SuccessOK(c, nil, "Secret file deleted successfully")
}
