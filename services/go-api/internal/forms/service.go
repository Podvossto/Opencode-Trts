// Purpose: Business logic implementation for the Forms domain
// Ref: API Contract — FormService interface defined in routes.go
// DB tables: application_forms
// Dev must implement: formService struct satisfying FormService interface
package forms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrFormNotFound = errors.New("form not found")
var ErrFormLinkedToJob = errors.New("form is linked to active job(s)")
var ErrFormAlreadyPublished = errors.New("form already published")
var ErrEmptyFormSchema = errors.New("form schema cannot be empty")

type formRepository struct {
	db *pgxpool.Pool
}

func newFormRepository(db *pgxpool.Pool) *formRepository {
	return &formRepository{db: db}
}

type formService struct {
	repo *formRepository
}

func NewFormService(db *pgxpool.Pool) FormService {
	return &formService{repo: newFormRepository(db)}
}

func (s *formService) List(ctx context.Context, page, limit int, published *bool) ([]ApplicationForm, int, error) {
	return s.repo.List(ctx, page, limit, published)
}

func (s *formService) GetByID(ctx context.Context, id uuid.UUID) (*ApplicationForm, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *formService) Create(ctx context.Context, req CreateFormRequest) (*ApplicationForm, error) {
	schemaJSON, err := json.Marshal(req.FormSchema)
	if err != nil {
		return nil, fmt.Errorf("marshal form schema: %w", err)
	}
	return s.repo.Create(ctx, req.Name, schemaJSON)
}

func (s *formService) Update(ctx context.Context, id uuid.UUID, req CreateFormRequest) (*ApplicationForm, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrFormNotFound
	}
	if existing.IsPublished {
		return nil, ErrFormAlreadyPublished
	}

	schemaJSON, err := json.Marshal(req.FormSchema)
	if err != nil {
		return nil, fmt.Errorf("marshal form schema: %w", err)
	}
	return s.repo.Update(ctx, id, req.Name, schemaJSON)
}

func (s *formService) Delete(ctx context.Context, id uuid.UUID) error {
	form, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if form == nil {
		return ErrFormNotFound
	}

	linked, err := s.repo.CheckJobsUsingForm(ctx, id)
	if err != nil {
		return err
	}
	if linked {
		return ErrFormLinkedToJob
	}

	return s.repo.Delete(ctx, id)
}

func (s *formService) Publish(ctx context.Context, id uuid.UUID) (*ApplicationForm, error) {
	form, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if form == nil {
		return nil, ErrFormNotFound
	}

	var fields []FormField
	if err := json.Unmarshal(form.FormSchema, &fields); err != nil {
		return nil, fmt.Errorf("unmarshal form schema: %w", err)
	}
	if len(fields) == 0 {
		return nil, ErrEmptyFormSchema
	}

	return s.repo.Publish(ctx, id)
}

func (r *formRepository) List(ctx context.Context, page, limit int, published *bool) ([]ApplicationForm, int, error) {
	baseWhere := ""
	var args []interface{}
	argIdx := 1

	if published != nil {
		if *published {
			baseWhere = fmt.Sprintf("WHERE is_published = true")
		} else {
			baseWhere = fmt.Sprintf("WHERE is_published = false")
		}
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM application_forms %s", baseWhere)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count forms: %w", err)
	}

	offset := (page - 1) * limit
	dataQuery := fmt.Sprintf(`
		SELECT id, name, form_schema, is_published, created_at, updated_at
		FROM application_forms
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`,
		baseWhere, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list forms: %w", err)
	}
	defer rows.Close()

	var forms []ApplicationForm
	for rows.Next() {
		var f ApplicationForm
		if err := rows.Scan(&f.ID, &f.Name, &f.FormSchema, &f.IsPublished, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan form: %w", err)
		}
		forms = append(forms, f)
	}
	return forms, total, nil
}

func (r *formRepository) GetByID(ctx context.Context, id uuid.UUID) (*ApplicationForm, error) {
	var f ApplicationForm
	err := r.db.QueryRow(ctx, `
		SELECT id, name, form_schema, is_published, created_at, updated_at 
		FROM application_forms WHERE id = $1`, id).
		Scan(&f.ID, &f.Name, &f.FormSchema, &f.IsPublished, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("get form: %w", err)
	}
	return &f, nil
}

func (r *formRepository) Create(ctx context.Context, name string, schemaJSON json.RawMessage) (*ApplicationForm, error) {
	var f ApplicationForm
	err := r.db.QueryRow(ctx, `
		INSERT INTO application_forms (name, form_schema)
		VALUES ($1, $2)
		RETURNING id, name, form_schema, is_published, created_at, updated_at`, name, schemaJSON).
		Scan(&f.ID, &f.Name, &f.FormSchema, &f.IsPublished, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create form: %w", err)
	}
	return &f, nil
}

func (r *formRepository) Update(ctx context.Context, id uuid.UUID, name string, schemaJSON json.RawMessage) (*ApplicationForm, error) {
	var f ApplicationForm
	err := r.db.QueryRow(ctx, `
		UPDATE application_forms SET name = $1, form_schema = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, form_schema, is_published, created_at, updated_at`, name, schemaJSON, id).
		Scan(&f.ID, &f.Name, &f.FormSchema, &f.IsPublished, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("update form: %w", err)
	}
	return &f, nil
}

func (r *formRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM application_forms WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete form: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrFormNotFound
	}
	return nil
}

func (r *formRepository) Publish(ctx context.Context, id uuid.UUID) (*ApplicationForm, error) {
	var f ApplicationForm
	err := r.db.QueryRow(ctx, `
		UPDATE application_forms SET is_published = true, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, form_schema, is_published, created_at, updated_at`, id).
		Scan(&f.ID, &f.Name, &f.FormSchema, &f.IsPublished, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("publish form: %w", err)
	}
	return &f, nil
}

func (r *formRepository) CheckJobsUsingForm(ctx context.Context, formID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM jobs WHERE form_id = $1 AND status IN ('open', 'scheduled'))`, formID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check jobs using form: %w", err)
	}
	return exists, nil
}
