// Purpose: Business logic for Resume Review — status transitions, position mark enforcement (max 3),
//
//	Star action (add to talent pool + drop + NO rejection email), pre-compute next applicant score.
//
// Ref: US-M06-01 through US-M06-05
// Dev must implement: ReviewService interface — ListApplicants, ApplyAction, MarkPositions, GetWeightedScore
package resume_review
