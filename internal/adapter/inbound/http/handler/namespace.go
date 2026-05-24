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

// NamespaceHandler handles HTTP requests for namespace operations.
type NamespaceHandler struct {
	cfg *configs.Config
	log logger.Logger
	svc inbound.Namespace
}

// NewNamespaceHandler creates a new NamespaceHandler.
func NewNamespaceHandler(cfg *configs.Config, log logger.Logger, svc inbound.Namespace) *NamespaceHandler {
	return &NamespaceHandler{cfg: cfg, log: log, svc: svc}
}

// RegisterRoutes registers namespace HTTP routes.
func (h *NamespaceHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/api/namespaces/tree", h.GetTree)
	app.Post("/api/namespaces", h.CreateNamespace)
	app.Delete("/api/namespaces/:id", h.DeleteNamespace)
}

// GetTree godoc
// @Summary      Get Namespace Tree
// @Description  Returns the full namespace tree
// @Tags         Namespaces
// @Produce      json
// @Success      200  {object}  response.ResponseSuccess{data=dto.ResNamespaceTree}
// @Router       /api/namespaces/tree [get]
func (h *NamespaceHandler) GetTree(c fiber.Ctx) error {
	ctx := c.Context()

	roots, err := h.svc.GetTree(ctx)
	if err != nil {
		return err
	}

	var res dto.ResNamespaceTree
	res.Transform(roots)

	return response.SuccessOK(c, res, "Namespace tree retrieved successfully")
}

// CreateNamespace godoc
// @Summary      Create Namespace
// @Description  Creates a new namespace node
// @Tags         Namespaces
// @Accept       json
// @Produce      json
// @Param        namespace  body      dto.ReqCreateNamespace  true  "Namespace data"
// @Success      201  {object}  response.ResponseSuccess{}
// @Router       /api/namespaces [post]
func (h *NamespaceHandler) CreateNamespace(c fiber.Ctx) error {
	var (
		req dto.ReqCreateNamespace
		ctx = c.Context()
	)

	if err := c.Bind().Body(&req); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := req.Validate(); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	ns := req.Transform()

	if err := h.svc.CreateNamespace(ctx, ns); err != nil {
		return err
	}

	return response.SuccessCreated(c, fiber.Map{"id": ns.ID}, "Namespace created successfully")
}

// DeleteNamespace godoc
// @Summary      Delete Namespace
// @Description  Soft-deletes a namespace and its descendants
// @Tags         Namespaces
// @Param        id   path  string  true  "Namespace ID"
// @Produce      json
// @Success      200  {object}  response.ResponseSuccess{}
// @Router       /api/namespaces/{id} [delete]
func (h *NamespaceHandler) DeleteNamespace(c fiber.Ctx) error {
	var (
		req dto.ReqDeleteNamespace
		ctx = c.Context()
	)

	if err := c.Bind().URI(&req); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := req.Validate(); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := h.svc.DeleteNamespace(ctx, req.ID); err != nil {
		return err
	}

	return response.SuccessOK(c, nil, "Namespace deleted successfully")
}
