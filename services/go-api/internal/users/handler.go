// Purpose: Users domain — handler layer (HTTP controllers)
// Ref: API Contract — GET/POST /api/v1/admin/users, GET/PUT/DELETE /api/v1/admin/users/:id
//
//	POST /api/v1/admin/users/:id/reset-password
//	GET /api/v1/users/me | PUT /api/v1/users/me
//	GET /api/v1/admin/departments | POST /api/v1/admin/departments
//
// ATS-008: List, ATS-009: Create, ATS-010: Edit/Delete, ATS-011: Admin reset password
package users

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tigersoft/ats-go-api/pkg/response"
)

// Handler holds user HTTP handlers
type Handler struct {
	svc *Service
}

// NewHandler creates a new users handler
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// ListUsers handles GET /api/v1/admin/users (ATS-008)
// Query params: page, limit, search, role, department_id, is_active
func (h *Handler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := ListFilter{
		Search: c.Query("search"),
		Role:   c.Query("role"),
		DeptID: c.Query("department_id"),
	}

	// Parse is_active if provided
	if activeStr := c.Query("is_active"); activeStr != "" {
		active := activeStr == "true"
		filter.IsActive = &active
	}

	result, err := h.svc.List(c.Request.Context(), page, limit, filter)
	if err != nil {
		log.Printf("list users error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, result)
}

// CreateUser handles POST /api/v1/admin/users (ATS-009)
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "body", Message: err.Error()}})
		return
	}

	user, tempPwd, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			response.BadRequest(c, []response.ValidationError{{Field: "email", Message: "email already exists"}})
			return
		}
		if errors.Is(err, ErrInvalidRole) {
			response.BadRequest(c, []response.ValidationError{{Field: "role", Message: "must be admin, hr, or supervisor"}})
			return
		}
		log.Printf("create user error: %v", err)
		response.InternalServerError(c)
		return
	}

	// Return user + temporary password (admin must communicate this to the user securely)
	c.JSON(http.StatusCreated, gin.H{
		"data":               user,
		"temporary_password": tempPwd,
	})
}

// GetUser handles GET /api/v1/admin/users/:id
func (h *Handler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "invalid UUID"}})
		return
	}

	user, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.NotFound(c)
			return
		}
		log.Printf("get user error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, user)
}

// UpdateUser handles PUT /api/v1/admin/users/:id (ATS-010)
func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "invalid UUID"}})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "body", Message: err.Error()}})
		return
	}

	user, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.NotFound(c)
			return
		}
		if errors.Is(err, ErrInvalidRole) {
			response.BadRequest(c, []response.ValidationError{{Field: "role", Message: "must be admin, hr, or supervisor"}})
			return
		}
		log.Printf("update user error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, user)
}

// DeleteUser handles DELETE /api/v1/admin/users/:id (ATS-010 — soft delete)
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "invalid UUID"}})
		return
	}

	callerID, _ := c.Get("user_id")
	callerStr, _ := callerID.(string)

	err = h.svc.Delete(c.Request.Context(), id, callerStr)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.NotFound(c)
			return
		}
		if errors.Is(err, ErrSelfDeactivate) {
			response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "cannot deactivate yourself"}})
			return
		}
		log.Printf("delete user error: %v", err)
		response.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deactivated"})
}

// AdminResetPassword handles POST /api/v1/admin/users/:id/reset-password (ATS-011)
func (h *Handler) AdminResetPassword(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "invalid UUID"}})
		return
	}

	var req AdminResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "body", Message: err.Error()}})
		return
	}

	err = h.svc.AdminResetPassword(c.Request.Context(), id, req.NewPassword)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.NotFound(c)
			return
		}
		log.Printf("admin reset password error: %v", err)
		response.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset, user must change on next login"})
}

// GetMe handles GET /api/v1/users/me
func (h *Handler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c)
		return
	}

	user, err := h.svc.GetMe(c.Request.Context(), userID.(string))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.NotFound(c)
			return
		}
		log.Printf("get me error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, user)
}

// UpdateMe handles PUT /api/v1/users/me
func (h *Handler) UpdateMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c)
		return
	}

	var req UpdateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "body", Message: err.Error()}})
		return
	}

	user, err := h.svc.UpdateMe(c.Request.Context(), userID.(string), req)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.NotFound(c)
			return
		}
		log.Printf("update me error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, user)
}

// ListDepartments handles GET /api/v1/admin/departments
func (h *Handler) ListDepartments(c *gin.Context) {
	depts, err := h.svc.ListDepartments(c.Request.Context())
	if err != nil {
		log.Printf("list departments error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, depts)
}

// CreateDepartment handles POST /api/v1/admin/departments
func (h *Handler) CreateDepartment(c *gin.Context) {
	var req CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "body", Message: err.Error()}})
		return
	}

	dept, err := h.svc.CreateDepartment(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrDeptExists) {
			response.BadRequest(c, []response.ValidationError{{Field: "name", Message: "department name already exists"}})
			return
		}
		log.Printf("create department error: %v", err)
		response.InternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse{Data: dept})
}
