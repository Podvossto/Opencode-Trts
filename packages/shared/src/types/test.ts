// Purpose: Shared types for Test feature (M07)
// Ref: DB tables: test_library, test_questions, test_rounds, test_assignments,
//      test_sessions, test_submissions, test_grades
// Mirrors Go API response contracts

export type TestScope = 'general' | 'technical';
export type QuestionType = 'mcq' | 'essay' | 'short_answer' | 'true_false';
export type TestAssignmentStatus = 'invited' | 'in_progress' | 'submitted' | 'graded';

export interface TestLibraryItem {
  id: string;
  title: string;
  scope: TestScope;
  department_id: string | null;
  time_limit_min: number;
  question_count: number;
  created_by: string | null;
  created_at: string;
}

export interface McqOption {
  id: string;
  text: string;
  is_correct: boolean;
}

export interface TestQuestion {
  id: string;
  test_id: string;
  question: string;
  type: QuestionType;
  options: McqOption[] | null;
  score: number;
  order_index: number;
  created_at: string;
}

export interface TestRound {
  id: string;
  pipeline_action_id: string;
  test_id: string;
  test_title: string;
  shuffle_seed: number | null;
  phase1_cutoff: number | null;
  total_score: number | null;
  created_at: string;
}

export interface TestAssignment {
  id: string;
  test_round_id: string;
  application_id: string;
  candidate_name: string;
  status: TestAssignmentStatus;
  total_score: number | null;
  invited_at: string;
}

export interface TestSession {
  id: string;
  test_assignment_id: string;
  started_at: string | null;
  submitted_at: string | null;
  tab_switch_count: number;
  copy_paste_count: number;
  last_heartbeat_at: string | null;
}

export interface TestSubmission {
  id: string;
  test_assignment_id: string;
  answer_data: AnswerEntry[];
  submitted_at: string;
}

export interface AnswerEntry {
  question_id: string;
  answer: string;
}

export interface TestGrade {
  id: string;
  submission_id: string;
  question_id: string;
  score: number;
  is_correct: boolean | null;
  grader_note: string | null;
  graded_by: string | null;
  graded_at: string | null;
}

export interface LeaderboardEntry {
  rank: number;
  application_id: string;
  candidate_name: string;
  total_score: number;
  max_score: number;
  status: TestAssignmentStatus;
}
