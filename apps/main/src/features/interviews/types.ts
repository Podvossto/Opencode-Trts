// Purpose: TypeScript types for the interviews feature
// Ref: API Contract — GET/POST/PUT /api/v1/interviews, DB table: interview_schedules
// Dev must implement: use these types in interview scheduling components

export type InterviewType = 'phone' | 'video' | 'onsite' | 'panel';
export type InterviewStatus = 'scheduled' | 'completed' | 'cancelled' | 'no_show';
export type InterviewResult = 'pass' | 'fail' | 'hold' | null;

export interface InterviewSchedule {
  id: string;
  applicationId: string;
  candidateName: string;
  candidateEmail: string;
  jobTitle: string;
  interviewerIds: string[];
  interviewerNames: string[];
  type: InterviewType;
  scheduledAt: string;
  durationMinutes: number;
  location: string | null;
  meetingLink: string | null;
  status: InterviewStatus;
  result: InterviewResult;
  feedback: string | null;
  round: number;
  createdAt: string;
}

export interface ScheduleInterviewRequest {
  applicationId: string;
  interviewerIds: string[];
  type: InterviewType;
  scheduledAt: string;
  durationMinutes: number;
  location?: string;
  meetingLink?: string;
  round: number;
}

export interface RecordInterviewResultRequest {
  result: Exclude<InterviewResult, null>;
  feedback?: string;
}

export interface InterviewListResponse {
  data: InterviewSchedule[];
  total: number;
  page: number;
  limit: number;
}
