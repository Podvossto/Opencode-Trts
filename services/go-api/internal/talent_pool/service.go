// Purpose: Business logic for Talent Pool — enqueue similarity jobs on job publish,
//
//	remove from pool on blacklist, block re-add on unblacklist.
//
// Ref: US-M10-01, US-M10-02, API Contract — /api/v1/talent-pool
// Dev must implement: TalentPoolService interface — EnqueueForJob, RemoveOnBlacklist
// DB table: talent_pool
package talent_pool

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrEntryNotFound is returned when a talent pool entry ID does not exist
var ErrEntryNotFound = errors.New("talent pool entry not found")

// ErrAlreadyInPool is returned when the same application is added twice
var ErrAlreadyInPool = errors.New("application already exists in talent pool")

// talentPoolService is the concrete implementation of TalentPoolService
// Dev must implement: inject pgxpool.Pool
type talentPoolService struct {
	// db *pgxpool.Pool
}

// NewService constructs a production TalentPoolService
// Dev must implement: accept *pgxpool.Pool as parameter
func NewService() TalentPoolService {
	return &talentPoolService{}
}

// ListEntries retrieves talent pool entries with optional search and tag filters
// Dev must implement:
//
//	SELECT tp.*, c.full_name, c.email FROM talent_pool tp
//	JOIN candidates c ON c.id = tp.candidate_id
//	WHERE ($1='' OR c.full_name ILIKE '%'||$1||'%')
//	AND ($2::text[] IS NULL OR tp.tags && $2)
//	ORDER BY tp.created_at DESC LIMIT $3 OFFSET $4
func (s *talentPoolService) ListEntries(ctx context.Context, filter TalentPoolFilter) ([]TalentPoolEntry, int, error) {
	// TODO: build dynamic query from TalentPoolFilter fields
	// TODO: run separate COUNT(*) for pagination total
	return nil, 0, errors.New("not implemented")
}

// GetEntry retrieves a single talent pool entry by primary key
// Dev must implement: SELECT * FROM talent_pool WHERE id=$1
func (s *talentPoolService) GetEntry(ctx context.Context, id uuid.UUID) (*TalentPoolEntry, error) {
	// TODO: return ErrEntryNotFound if no rows
	return nil, errors.New("not implemented")
}

// AddEntry inserts a new talent pool entry
// Guard: check for duplicate application_id; derive candidate_id from applications JOIN
// Dev must implement:
//
//	SELECT COUNT(*) FROM talent_pool WHERE application_id=$1
//	If count>0: return ErrAlreadyInPool
//	INSERT INTO talent_pool (id, application_id, candidate_id, added_by_id, note, tags, expires_at)
//	SELECT gen_random_uuid(),$1,a.candidate_id,$2,$3,$4,$5 FROM applications a WHERE a.id=$1 RETURNING *
func (s *talentPoolService) AddEntry(ctx context.Context, addedByID uuid.UUID, req AddTalentPoolRequest) (*TalentPoolEntry, error) {
	return nil, errors.New("not implemented")
}

// RemoveEntry deletes a talent pool entry by primary key
// Called by: blacklist flow, manual removal
// Dev must implement: DELETE FROM talent_pool WHERE id=$1
func (s *talentPoolService) RemoveEntry(ctx context.Context, id uuid.UUID) error {
	// TODO: return ErrEntryNotFound if affected rows == 0
	return errors.New("not implemented")
}

// UpdateEntry patches note, tags, and expires_at on an existing entry
// Dev must implement:
//
//	UPDATE talent_pool SET note=$1, tags=$2, expires_at=$3 WHERE id=$4 RETURNING *
func (s *talentPoolService) UpdateEntry(ctx context.Context, id uuid.UUID, req UpdateTalentPoolRequest) (*TalentPoolEntry, error) {
	return nil, errors.New("not implemented")
}
