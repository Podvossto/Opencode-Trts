// Purpose: Shared types for Interview feature (M08)
// Ref: DB tables: interview_rounds, scoring_criteria, interview_assignments,
//      interview_evaluations, invite_links
// Mirrors Go API response contracts

export type InterviewRoundStatus = 'open' | 'scoring' | 'closed';
export type InterviewAssignmentStatus = 'invited' | 'scored' | 'passed' | 'dropped';
export type InviteLinkStatus = 'active' | 'used' | 'revoked' | 'expired';

export interface InterviewRound {
  id: string;
  pipeline_action_id: string;
  phase1_cutoff: number | null;
  status: InterviewRoundStatus;
  criteria: ScoringCriterion[];
  created_at: string;
}

export interface ScoringCriterion {
  id: string;
  interview_round_id: string;
  name: string;
  description: string | null;
  max_score: number;
  order_index: number;
}

export interface InterviewAssignment {
  id: string;
  interview_round_id: string;
  application_id: string;
  candidate_name: string;
  status: InterviewAssignmentStatus;
  average_score: number | null;
  created_at: string;
}

export interface InterviewEvaluation {
  id: string;
  interview_assignment_id: string;
  criterion_id: string;
  criterion_name: string;
  interviewer_id: string;
  interviewer_name: string;
  score: number | null;
  note: string | null;
  submitted_at: string | null;
}

export interface InviteLink {
  id: string;
  interview_round_id: string;
  token: string;
  invited_email: string;
  invited_user_id: string | null;
  status: InviteLinkStatus;
  expires_at: string;
  revoked_by: string | null;
  revoked_at: string | null;
  created_by: string | null;
  created_at: string;
}

export interface InterviewLeaderboardEntry {
  rank: number;
  application_id: string;
  candidate_name: string;
  average_score: number;
  max_possible_score: number;
  evaluator_count: number;
  status: InterviewAssignmentStatus;
}
