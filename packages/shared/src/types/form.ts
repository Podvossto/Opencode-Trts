// Purpose: Shared TypeScript types for application form schema — consumed by apps/main (builder) and apps/portal (renderer)
// Ref: API Contract — POST /api/v1/forms, DB table: application_forms
// 14 field types per SRS M05

export type FormFieldType =
  | 'text'
  | 'textarea'
  | 'dropdown'
  | 'checkbox'
  | 'radio'
  | 'file_upload'
  | 'date_picker'
  | 'mcq'
  | 'short_answer'
  | 'essay'
  | 'true_false'
  | 'rating_scale'
  | 'image'
  | 'video'
  | 'link_url';

export interface FormField {
  field_id: string;        // unique within form, snake_case
  type: FormFieldType;
  label: string;
  required: boolean;
  options?: string[];      // required for: dropdown, checkbox, radio, mcq, true_false
  max_size_mb?: number;    // for file_upload, image, video
  accept?: string;         // MIME types for file_upload
  placeholder?: string;
}

export interface ApplicationForm {
  id: string;
  name: string;
  description?: string;
  form_schema: FormField[];
  created_at: string;
  updated_at?: string;
}
