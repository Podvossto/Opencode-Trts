// Purpose: HTTP handlers for Test pipeline — CRUD test library, round management,
//
//	applicant selection, MCQ auto-grade trigger, grading inbox, Phase 2 leaderboard.
//
// Ref: API Contract — /api/v1/tests/*, /api/v1/jobs/{jobId}/test-rounds/*
// DB tables: tests, test_questions, test_rounds, test_assignments, test_submissions, test_answers, test_scores
// Dev must implement: CreateTest, ListTests, GetRound, AssignApplicants, SubmitGrade, GetLeaderboard
package tests
