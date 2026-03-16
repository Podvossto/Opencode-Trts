// Purpose: HTTP handlers for Resume Review — applicant list, action buttons (Interest/Reject/Star),
//
//	position marks, weight popup, and score retrieval.
//
// Ref: API Contract — GET /api/v1/jobs/{jobId}/resume-review, POST /api/v1/applications/{id}/action
// DB tables: applications, candidates, application_documents, candidate_similarity_scores, position_marks
// Dev must implement: GetResumeReviewList, ActionInterest, ActionReject, ActionStar, RetryOCR, GetScores
package resume_review
