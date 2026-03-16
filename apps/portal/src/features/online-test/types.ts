// Purpose: TypeScript types for applicant-facing online test delivery (M07)
// Ref: API Contract — GET/POST /api/v1/portal/test/:token
// DB tables: test_assignments, test_sessions, test_submissions, test_questions
// Dev must implement: use these types in portal online test components

export type QuestionType = 'mcq' | 'essay' | 'short_answer' | 'true_false';

export interface TestQuestion {
  id: string;
  question: string;
  type: QuestionType;
  options: McqOption[] | null; // null for non-MCQ
  orderIndex: number;
}

export interface McqOption {
  id: string;
  text: string;
  // Note: is_correct is NOT exposed to applicant
}

export interface OnlineTestSession {
  testTitle: string;
  timeLimitMin: number;
  questions: TestQuestion[];
  startedAt: string | null;
  remainingSeconds: number;
  savedAnswers: AnswerEntry[];
}

export interface AnswerEntry {
  questionId: string;
  answer: string; // MCQ: option ID, essay/short_answer: text, true_false: "true"|"false"
}

export interface SubmitAnswersRequest {
  answers: AnswerEntry[];
}

export interface AutoSaveRequest {
  answers: AnswerEntry[];
}
