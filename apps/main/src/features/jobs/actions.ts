// Purpose: Jobs feature — server actions for job posting CRUD and reference data
// Ref: API Contract — GET/POST /api/v1/jobs | GET/PUT/DELETE /api/v1/jobs/{id}
//                     GET /api/v1/admin/departments | GET /api/v1/admin/employment-types
//                     GET /api/v1/admin/experience-levels | GET /api/v1/admin/hard-skills
// DB tables: jobs, job_hard_skills, departments, employment_types, experience_levels, hard_skills
// Role: HR/Admin
// Dev must implement:
//   listJobsAction(filters?): Promise<{ data: Job[]; total: number }>
//   getJobAction(id): Promise<Job>
//   createJobAction(body): Promise<Job>
//   updateJobAction(id, body): Promise<Job>
//   deleteJobAction(id): Promise<void>
//   listReferenceDataAction(): Promise<{ departments, employmentTypes, experienceLevels, hardSkills }>
