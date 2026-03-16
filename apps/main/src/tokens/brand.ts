// tokens/brand.ts — generated from docs/Tigersoft_CI_Branding.md
// ATS Recruitment — Batch 2 Design Tokens
// ALL frontend components must import from this file.
// NEVER use raw hex values, raw font strings, or raw px values in components.

export const BRAND = {
  color: {
    // Primary palette (85% of design)
    vivid_red:    '#F4001A', // CTA, accent, brand highlight — buttons, focus rings
    white:        '#FFFFFF', // Background, card, surface
    oxford_blue:  '#0B1F3A', // Body text, headings, dark bg — NEVER use #000000

    // Secondary palette (15% of design)
    quick_silver: '#A3A3A3', // Border, divider, disabled state
    serene:       '#DBE1E1', // Subtle bg, light card, default border
    ufo_green:    '#34D186', // Success, active, positive status

    // Semantic tokens
    bg_primary:    '#FFFFFF',
    bg_secondary:  '#DBE1E1',
    bg_dark:       '#0B1F3A',
    text_primary:  '#0B1F3A', // ห้ามใช้ #000 เด็ดขาด
    text_inverse:  '#FFFFFF',
    text_secondary:'#A3A3A3',
    action_primary:'#F4001A',
    action_hover:  '#C8001A', // darken 15%
    action_success:'#34D186',
    border_default:'#DBE1E1',
    border_focus:  '#F4001A',

    // Feature-specific semantic tokens (Batch 2)
    score_high:    '#34D186', // Overall score ≥ 0.75
    score_mid:     '#F4001A', // Overall score 0.50–0.74
    score_low:     '#A3A3A3', // Overall score < 0.50
    star_active:   '#F4001A', // Starred candidate badge
    round_badge:   '#0B1F3A', // Round badge background (Test / Interview)
    pending:       '#A3A3A3', // Grader "Pending" indicator
    dropped:       '#A3A3A3', // Dropped applicant text
    hired:         '#34D186', // Hired applicant text
    transferred:   '#0B1F3A', // Transferred applicant text
    notification_unread: '#F4001A', // Unread badge
  },

  font: {
    en: 'Plus Jakarta Sans',  // Google Fonts — all English text
    th: 'FC Vision',          // Custom — Thai text (guide/CI Toolkit/Font/TH/)
    weight_light:  300,       // Body, paragraph, description
    weight_medium: 500,       // Heading, subheading, button, label
  },

  size: {
    // Type scale (base 16px)
    xs:   12,
    sm:   14,
    base: 16,
    md:   18,
    lg:   20,
    xl:   24,
    '2xl': 30,
    '3xl': 36,
    '4xl': 48,
    '5xl': 60,
    '6xl': 72,
  },

  spacing: {
    // 4px base unit scale
    1:  4,
    2:  8,
    3:  12,
    4:  16,
    5:  20,
    6:  24,
    8:  32,
    10: 40,
    12: 48,
    16: 64,
    20: 80,
    24: 96,
    32: 128,
    // Semantic spacing
    page_x:    40, // Horizontal page margin
    section_y: 80, // Section vertical padding
    card:      24, // Card internal padding
    component: 16, // Between inner elements
    tight:      8, // Dense / compact content
  },

  radius: {
    // Soft rounded edges — all values in px
    sm:     4,
    md:     8,
    lg:     12,
    xl:     16,
    '2xl':  24,
    '3xl':  32,
    full:   9999, // Full pill

    // Semantic radius tokens
    button: 9999, // Full pill — ALL buttons
    card:   16,   // Cards, panels, 3-panel layout
    input:   8,   // Input fields
    modal:  24,   // Modals, drawers
    image:  12,   // Clipping mask containers
    badge:  9999, // Tags, chips, round badges, avatars
  },

  motion: {
    fast:   '150ms',  // Hover, color change
    normal: '250ms',  // Standard transition
    slow:   '400ms',  // Page element entrance
    ease_brand:  'cubic-bezier(0.22, 1, 0.36, 1)',   // Default easing
    ease_spring: 'cubic-bezier(0.34, 1.56, 0.64, 1)', // Bounce / delight
  },

  grid: {
    columns_mobile:  4,
    columns_tablet:  8,
    columns_desktop: 12,
    gutter:          24,
    max_content_width: 1280,
  },

  shadow: {
    // All shadows use Oxford Blue tint (not black)
    // Red shadow ONLY on primary CTA button
    sm:  '0 1px 2px rgba(11, 31, 58, 0.08)',
    md:  '0 4px 8px rgba(11, 31, 58, 0.12)',
    lg:  '0 8px 16px rgba(11, 31, 58, 0.16)',
    cta: '0 4px 12px rgba(244, 0, 26, 0.30)', // Red shadow — CTA only
  },

  accessibility: {
    // Focus ring — Vivid Red 3px
    focus_ring: 'rgba(244, 0, 26, 0.30)',
    focus_ring_width: 3,
    min_touch_target: 44, // px — WCAG minimum
  },
} as const;

export type BrandToken = typeof BRAND;
