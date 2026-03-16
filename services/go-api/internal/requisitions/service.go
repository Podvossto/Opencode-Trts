// Purpose: Business logic for Requisition — Draft→Pending status transitions,
//
//	conversion pre-fill to jobs table, status chain to In-Progress/Published,
//	trigger FRNOT01 on submit, FRNOT02 on job publish.
//
// Ref: US-M12-01 through US-M12-03, API Contract — /api/v1/requisitions
// Dev must implement: RequisitionService interface — Create, Submit, Edit, Convert
// DB table: requisitions
package requisitions

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrRequisitionNotFound is returned when a requisition ID does not exist
var ErrRequisitionNotFound = errors.New("requisition not found")

// ErrInvalidStatusTransition is returned when a status change is not allowed
var ErrInvalidStatusTransition = errors.New("invalid status transition")

// ErrRequisitionInUse is returned when cancelling an approved/fulfilled requisition
var ErrRequisitionInUse = errors.New("cannot cancel approved or fulfilled requisition")

// requisitionService is the concrete implementation of RequisitionService
// Dev must implement: inject pgxpool.Pool and email.Sender
type requisitionService struct {
	// db    *pgxpool.Pool
	// email email.Sender
}

// NewService constructs a production RequisitionService
// Dev must implement: accept *pgxpool.Pool and email.Sender as parameters
func NewService() RequisitionService {
	return &requisitionService{}
}

// ListRequisitions retrieves all requisitions with optional filters
// Dev must implement:
//
//	SELECT * FROM requisitions WHERE (status = $1 OR $1 IS NULL) AND (department = $2 OR $2 IS NULL)
//	ORDER BY created_at DESC LIMIT $3 OFFSET $4
func (s *requisitionService) ListRequisitions(ctx context.Context) ([]Requisition, error) {
	return nil, errors.New("not implemented")
}

// GetRequisition retrieves a single requisition by primary key
// Dev must implement:
//
//	SELECT * FROM requisitions WHERE id = $1
func (s *requisitionService) GetRequisition(ctx context.Context, id uuid.UUID) (*Requisition, error) {
	// TODO: return ErrRequisitionNotFound if no rows
	return nil, errors.New("not implemented")
}

// CreateRequisition inserts a new requisition with status = "draft"
// Triggers FRNOT01: notify Admin/Supervisor of new headcount request
// Dev must implement:
//
//	INSERT INTO requisitions (id, requested_by_id, department, position, headcount, justification, status)
//	VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, 'draft') RETURNING *
func (s *requisitionService) CreateRequisition(ctx context.Context, requestedByID uuid.UUID, req CreateRequisitionRequest) (*Requisition, error) {
	// TODO: insert with status = draft
	// TODO: if auto-submit, transition to pending_approval and trigger FRNOT01 email
	return nil, errors.New("not implemented")
}

// ReviewRequisition approves or rejects a pending requisition
// Triggers FRNOT02 on approval: notify requester that job can now be published
// Dev must implement:
//
//	UPDATE requisitions SET status=$1, approved_by_id=$2, rejection_note=$3, updated_at=now()
//	WHERE id=$4 AND status='pending_approval' RETURNING *
func (s *requisitionService) ReviewRequisition(ctx context.Context, reviewerID uuid.UUID, id uuid.UUID, req ReviewRequisitionRequest) (*Requisition, error) {
	// TODO: fetch, validate status == pending_approval
	// TODO: update status to approved/rejected
	// TODO: send FRNOT01/02 email to requester
	return nil, errors.New("not implemented")
}

// CancelRequisition cancels a draft or pending requisition
// Guard: only HR (own) or Admin; cannot cancel approved or fulfilled
// Dev must implement:
//
//	UPDATE requisitions SET status='cancelled', updated_at=now() WHERE id=$1
func (s *requisitionService) CancelRequisition(ctx context.Context, requesterID uuid.UUID, id uuid.UUID) error {
	// TODO: validate ownership for HR role
	// TODO: validate status is draft or pending_approval, else return ErrRequisitionInUse
	return errors.New("not implemented")
}
