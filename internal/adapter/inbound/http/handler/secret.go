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

// SecretHandler handles HTTP requests for secret operations.
type SecretHandler struct {
	cfg *configs.Config
	log logger.Logger
	svc inbound.Secret
}

// NewSecretHandler creates a new SecretHandler.
func NewSecretHandler(cfg *configs.Config, log logger.Logger, svc inbound.Secret) *SecretHandler {
	return &SecretHandler{cfg: cfg, log: log, svc: svc}
}

// RegisterRoutes registers secret HTTP routes.
func (h *SecretHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/api/files/:file_id/secrets", h.ListSecrets)
	app.Post("/api/files/:file_id/secrets", h.AddSecret)
	app.Put("/api/secrets/:id", h.UpdateSecret)
	app.Get("/api/secrets/:id/reveal", h.RevealSecret)
	app.Delete("/api/secrets/:id", h.DeleteSecret)
}

// ListSecrets godoc
// @Summary      List Secrets
// @Description  Returns all secrets in a file (values masked)
// @Tags         Secrets
// @Param        file_id  path  string  true  "File ID"
// @Produce      json
// @Success      200  {object}  response.ResponseSuccess{data=dto.ResSecretList}
// @Router       /api/files/{file_id}/secrets [get]
func (h *SecretHandler) ListSecrets(c fiber.Ctx) error {
	var (
		req dto.ReqListSecrets
		ctx = c.Context()
	)

	if err := c.Bind().URI(&req); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := req.Validate(); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	secrets, err := h.svc.ListSecrets(ctx, req.FileID)
	if err != nil {
		return err
	}

	var res dto.ResSecretList
	res.Transform(secrets)

	return response.SuccessOK(c, res, "Secrets retrieved successfully")
}

// AddSecret godoc
// @Summary      Add Secret
// @Description  Encrypts and stores a new key-value pair
// @Tags         Secrets
// @Accept       json
// @Produce      json
// @Param        file_id  path  string           true  "File ID"
// @Param        secret   body  dto.ReqAddSecret true  "Secret data"
// @Success      201  {object}  response.ResponseSuccess{}
// @Router       /api/files/{file_id}/secrets [post]
func (h *SecretHandler) AddSecret(c fiber.Ctx) error {
	var (
		req dto.ReqAddSecret
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

	secret := req.Transform()

	if err := h.svc.AddSecret(ctx, secret); err != nil {
		return err
	}

	return response.SuccessCreated(c, fiber.Map{"id": secret.ID}, "Secret added successfully")
}

// UpdateSecret godoc
// @Summary      Update Secret
// @Description  Encrypts and updates an existing key-value pair
// @Tags         Secrets
// @Accept       json
// @Produce      json
// @Param        id      path  string              true  "Secret ID"
// @Param        secret  body  dto.ReqUpdateSecret true  "Secret data"
// @Success      200  {object}  response.ResponseSuccess{}
// @Router       /api/secrets/{id} [put]
func (h *SecretHandler) UpdateSecret(c fiber.Ctx) error {
	var (
		req dto.ReqUpdateSecret
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

	secret := req.Transform()

	if err := h.svc.UpdateSecret(ctx, secret); err != nil {
		return err
	}

	return response.SuccessOK(c, nil, "Secret updated successfully")
}

// RevealSecret godoc
// @Summary      Reveal Secret
// @Description  Decrypts and returns the plaintext value
// @Tags         Secrets
// @Param        id  path  string  true  "Secret ID"
// @Produce      json
// @Success      200  {object}  response.ResponseSuccess{data=dto.ResRevealSecret}
// @Router       /api/secrets/{id}/reveal [get]
func (h *SecretHandler) RevealSecret(c fiber.Ctx) error {
	var (
		req dto.ReqRevealSecret
		ctx = c.Context()
	)

	if err := c.Bind().URI(&req); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := req.Validate(); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	value, err := h.svc.RevealSecret(ctx, req.ID)
	if err != nil {
		return err
	}

	return response.SuccessOK(c, dto.ResRevealSecret{Value: value}, "Secret revealed successfully")
}

// DeleteSecret godoc
// @Summary      Delete Secret
// @Description  Soft-deletes a secret
// @Tags         Secrets
// @Param        id  path  string  true  "Secret ID"
// @Produce      json
// @Success      200  {object}  response.ResponseSuccess{}
// @Router       /api/secrets/{id} [delete]
func (h *SecretHandler) DeleteSecret(c fiber.Ctx) error {
	var (
		req dto.ReqDeleteSecret
		ctx = c.Context()
	)

	if err := c.Bind().URI(&req); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := req.Validate(); err != nil {
		return fail.Wrap(err).WithFailure(fail.ErrBadRequest)
	}

	if err := h.svc.DeleteSecret(ctx, req.ID); err != nil {
		return err
	}

	return response.SuccessOK(c, nil, "Secret deleted successfully")
}
