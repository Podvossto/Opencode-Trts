// tokens/brand.ts — Applicant Portal — generated from docs/Tigersoft_CI_Branding.md
// Portal shares the same brand tokens. Import from this file for all portal UI.
// NEVER use raw hex values in components.

// Re-export from apps/main via tsconfig @/* alias is not available cross-app.
// Instead, we duplicate the essential subset needed for portal.
// Source of truth: apps/main/src/tokens/brand.ts

export const BRAND = {
  color: {
    vivid_red:    '#F4001A',
    white:        '#FFFFFF',
    oxford_blue:  '#0B1F3A',
    quick_silver: '#A3A3A3',
    serene:       '#DBE1E1',
    ufo_green:    '#34D186',

    bg_primary:    '#FFFFFF',
    bg_secondary:  '#DBE1E1',
    text_primary:  '#0B1F3A',
    text_inverse:  '#FFFFFF',
    text_secondary:'#A3A3A3',
    action_primary:'#F4001A',
    action_hover:  '#C8001A',
    action_success:'#34D186',
    border_default:'#DBE1E1',
    border_focus:  '#F4001A',
  },
  font: {
    en: 'Plus Jakarta Sans',
    th: 'FC Vision',
    weight_light:  300,
    weight_medium: 500,
  },
  radius: {
    button: 9999,
    card:   16,
    input:   8,
    modal:  24,
    image:  12,
    badge:  9999,
  },
} as const;

export type BrandToken = typeof BRAND;
