// Purpose: Forms feature — local types (extends packages/shared)
// Ref: API Contract — POST /api/v1/forms, DB table: application_forms
// Dev must implement: re-export and extend shared FormField type

export type { FormFieldType, FormField, ApplicationForm } from '@ats/shared';

// Field type display labels for the form builder UI
export const FIELD_TYPE_LABELS: Record<string, string> = {
  text:          'Short Text',
  textarea:      'Long Text',
  dropdown:      'Dropdown',
  checkbox:      'Checkbox',
  radio:         'Radio Button',
  file_upload:   'File Upload',
  date_picker:   'Date Picker',
  mcq:           'Multiple Choice (MCQ)',
  short_answer:  'Short Answer',
  essay:         'Essay',
  true_false:    'True / False',
  rating_scale:  'Rating Scale',
  image:         'Image',
  video:         'Video',
  link_url:      'Link / URL',
};

// Field types that require options array
export const OPTION_FIELD_TYPES = ['dropdown', 'checkbox', 'radio', 'mcq', 'true_false'] as const;

// Field types that are file-based (need max_size_mb + accept)
export const FILE_FIELD_TYPES = ['file_upload', 'image', 'video'] as const;
