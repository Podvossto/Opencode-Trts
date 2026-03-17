// Purpose: HTTP handlers for the Forms domain — form builder CRUD
// Ref: API Contract — GET/POST /api/v1/forms, GET/PUT/DELETE /api/v1/forms/{id}
// DB tables: application_forms
// Dev must implement: ListForms, CreateForm, GetForm, UpdateForm, DeleteForm
package forms

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tigersoft/ats-go-api/pkg/response"
)

type FormHandler struct {
	svc FormService
}

func NewFormHandler(svc FormService) *FormHandler {
	return &FormHandler{svc: svc}
}

func (h *FormHandler) ListForms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var published *bool
	if p := c.Query("published"); p != "" {
		isPub := p == "true"
		published = &isPub
	}

	forms, total, err := h.svc.List(c.Request.Context(), page, limit, published)
	if err != nil {
		log.Printf("list forms error: %v", err)
		response.InternalServerError(c)
		return
	}
	response.OK(c, ListResponse{
		Forms: forms,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

func (h *FormHandler) GetForm(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "invalid UUID"}})
		return
	}

	form, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("get form error: %v", err)
		response.InternalServerError(c)
		return
	}
	if form == nil {
		response.NotFound(c)
		return
	}
	response.OK(c, form)
}

func (h *FormHandler) CreateForm(c *gin.Context) {
	var req CreateFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "body", Message: err.Error()}})
		return
	}

	form, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		log.Printf("create form error: %v", err)
		response.InternalServerError(c)
		return
	}
	response.Created(c, form)
}

func (h *FormHandler) UpdateForm(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "invalid UUID"}})
		return
	}

	var req CreateFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "body", Message: err.Error()}})
		return
	}

	form, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrFormNotFound) {
			response.NotFound(c)
			return
		}
		if errors.Is(err, ErrFormAlreadyPublished) {
			c.JSON(http.StatusConflict, response.ErrorResponse{Error: "form already published"})
			return
		}
		log.Printf("update form error: %v", err)
		response.InternalServerError(c)
		return
	}
	if form == nil {
		response.NotFound(c)
		return
	}
	response.OK(c, form)
}

func (h *FormHandler) DeleteForm(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "invalid UUID"}})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrFormNotFound) {
			response.NotFound(c)
			return
		}
		if errors.Is(err, ErrFormLinkedToJob) {
			c.JSON(http.StatusConflict, response.ErrorResponse{Error: "Form is linked to active job(s)"})
			return
		}
		log.Printf("delete form error: %v", err)
		response.InternalServerError(c)
		return
	}
	response.OK(c, gin.H{"message": "form deleted"})
}

func (h *FormHandler) PublishForm(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "invalid UUID"}})
		return
	}

	form, err := h.svc.Publish(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrFormNotFound) {
			response.NotFound(c)
			return
		}
		if errors.Is(err, ErrEmptyFormSchema) {
			response.BadRequest(c, []response.ValidationError{{Field: "form_schema", Message: "form must have at least one field"}})
			return
		}
		log.Printf("publish form error: %v", err)
		response.InternalServerError(c)
		return
	}
	response.OK(c, form)
}
