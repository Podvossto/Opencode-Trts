// Purpose: TypeScript types for the tests (online test) feature — HR side
// Ref: API Contract — GET/POST /api/v1/tests, DB table: test_assignments
// Dev must implement: use these types in test management components

export type TestAssignmentStatus = 'pending' | 'sent' | 'in_progress' | 'completed' | 'expired';
export type TestResult = 'pass' | 'fail' | null;

export interface TestAssignment {
  id: string;
  applicationId: string;
  candidateName: string;
  candidateEmail: string;
  jobTitle: string;
  formId: string;
  formTitle: string;
  token: string;
  status: TestAssignmentStatus;
  score: number | null;
  maxScore: number;
  result: TestResult;
  sentAt: string | null;
  startedAt: string | null;
  completedAt: string | null;
  expiresAt: string | null;
  createdAt: string;
}

export interface AssignTestRequest {
  applicationId: string;
  formId: string;
  expiresAt?: string;
}

export interface TestListResponse {
  data: TestAssignment[];
  total: number;
  page: number;
  limit: number;
}

export interface TestResultDetail {
  assignment: TestAssignment;
  answers: {
    questionId: string;
    questionText: string;
    answer: string;
    isCorrect: boolean | null;
    score: number;
  }[];
}
