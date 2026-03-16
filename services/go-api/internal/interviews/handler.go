// Purpose: HTTP handlers for Interview pipeline — round init, scoring criteria, interviewer assignment,
//
//	invite link generation/revocation, Phase 1 grade submission, Phase 2 leaderboard decisions.
//
// Ref: API Contract — /api/v1/jobs/{jobId}/interview-rounds/*
// DB tables: interview_rounds, interview_criteria, interview_assignments, temp_interviewers,
//
//	invite_links, interview_scores
//
// Dev must implement: InitRound, SetCriteria, AssignInterviewer, GenerateInviteLink, RevokeLink,
//
//	SubmitScores, OpenPhase2, MakeDecision
package interviews
