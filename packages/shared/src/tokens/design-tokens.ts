/**
 * TigerSoft CI Design Tokens
 * Source: docs/Tigersoft_CI_Branding.md
 *
 * RULES:
 *  - All developers MUST import from this file.
 *  - NEVER use raw hex, raw font names, or raw spacing values in components.
 *  - NEVER use #000000 (pure black) — use oxford_blue instead.
 *  - Tokens marked ⚠️ FLAG were not in Tigersoft_CI_Branding.md and require
 *    Orchestrator sign-off before they are used in production.
 */

// ---------------------------------------------------------------------------
// COLOR TOKENS
// ---------------------------------------------------------------------------

export const COLOR = {
  /** Primary: vivid red — CTAs, primary action buttons, active nav indicators */
  vivid_red: '#F4001A',

  /** Primary: white — page backgrounds, card surfaces, inverted text */
  white: '#FFFFFF',

  /** Primary: Oxford blue — body text, headings, dark surfaces. NEVER use #000000. */
  oxford_blue: '#0B1F3A',

  /** Secondary: medium grey — secondary text, labels, placeholders */
  medium_grey: '#A3A3A3',

  /** Secondary: light silver — dividers, input borders, table row separators */
  light_silver: '#DBE1E1',

  /** Secondary: success green — success states, "Hired" status badges */
  success_green: '#34D186',

  // -------------------------------------------------------------------------
  // Semantic aliases (derived from branding doc — no new hex values)
  // -------------------------------------------------------------------------
  semantic: {
    /** Page & card backgrounds */
    background: '#FFFFFF',
    /** Primary body / heading text */
    text_primary: '#0B1F3A',
    /** Secondary / muted text, labels */
    text_secondary: '#A3A3A3',
    /** Borders, dividers, input outlines */
    border: '#DBE1E1',
    /** Destructive actions, errors */
    danger: '#F4001A',
    /** Success / positive states */
    success: '#34D186',
    /** Primary CTA background */
    cta_bg: '#F4001A',
    /** Primary CTA text (on vivid_red background) */
    cta_text: '#FFFFFF',
    /** Focus ring color — ⚠️ FLAG: not explicitly specified in branding doc; using vivid_red as safe default */
    focus_ring: '#F4001A',
  },
} as const;

// ---------------------------------------------------------------------------
// TYPOGRAPHY TOKENS
// ---------------------------------------------------------------------------

export const TYPOGRAPHY = {
  font: {
    /** English text — all latin-script content */
    en: 'Plus Jakarta Sans',
    /** Thai text — all Thai-script content */
    th: 'FC Vision',
    /** Fallback stack if primary fonts are unavailable */
    fallback: 'sans-serif',
  },

  /**
   * Font weight scale
   * ⚠️ FLAG: specific weight values not listed in branding doc — using common web-safe
   * weights consistent with Plus Jakarta Sans variable font. Needs Orchestrator sign-off.
   */
  weight: {
    regular: 400,
    medium: 500,
    semibold: 600,
    bold: 700,
  },

  /**
   * Type scale — base unit 4px grid
   * ⚠️ FLAG: exact type scale sizes not in branding doc. Values below are derived
   * from the 4px base unit and standard web practice. Needs Orchestrator sign-off.
   */
  size: {
    xs: '12px',   // 3 × 4px — captions, tags
    sm: '14px',   // 3.5 × 4px — labels, helper text
    base: '16px', // 4 × 4px — body default
    md: '18px',   // 4.5 × 4px — sub-headings
    lg: '20px',   // 5 × 4px — section titles
    xl: '24px',   // 6 × 4px — page titles
    '2xl': '32px', // 8 × 4px — hero headings
    '3xl': '40px', // 10 × 4px — display headings
  },

  /**
   * Line height
   * ⚠️ FLAG: not explicitly specified in branding doc. Needs Orchestrator sign-off.
   */
  lineHeight: {
    tight: '1.2',
    normal: '1.5',
    relaxed: '1.75',
  },
} as const;

// ---------------------------------------------------------------------------
// SPACING TOKENS (4px base grid)
// ---------------------------------------------------------------------------

/**
 * All spacing values are multiples of 4px per branding doc base unit.
 * ⚠️ FLAG: exact named spacing scale not listed in branding doc — derived from
 * 4px base unit rule. Needs Orchestrator sign-off.
 */
export const SPACING = {
  '0': '0px',
  '1': '4px',
  '2': '8px',
  '3': '12px',
  '4': '16px',
  '5': '20px',
  '6': '24px',
  '8': '32px',
  '10': '40px',
  '12': '48px',
  '16': '64px',
  '20': '80px',
  '24': '96px',
} as const;

// ---------------------------------------------------------------------------
// BORDER RADIUS TOKENS
// ---------------------------------------------------------------------------

/**
 * Soft rounded / pill style as per TigerSoft branding.
 * ⚠️ FLAG: exact pixel values for each radius tier not in branding doc — values
 * below are consistent with "soft rounded, pill" direction. Needs Orchestrator sign-off.
 */
export const RADIUS = {
  none: '0px',
  sm: '4px',
  md: '8px',
  lg: '12px',
  xl: '16px',
  '2xl': '24px',
  pill: '9999px', // Full pill shape — for tags, badges, status chips
  full: '50%',    // Perfect circle — for avatars
} as const;

// ---------------------------------------------------------------------------
// SHADOW TOKENS
// ---------------------------------------------------------------------------

/**
 * ⚠️ FLAG: shadow values not specified in branding doc.
 * Needs Orchestrator sign-off before use in production components.
 */
export const SHADOW = {
  none: 'none',
  sm: '0 1px 2px 0 rgba(11, 31, 58, 0.05)',
  md: '0 4px 6px -1px rgba(11, 31, 58, 0.08)',
  lg: '0 10px 15px -3px rgba(11, 31, 58, 0.10)',
  xl: '0 20px 25px -5px rgba(11, 31, 58, 0.12)',
} as const;

// ---------------------------------------------------------------------------
// Z-INDEX TOKENS
// ---------------------------------------------------------------------------

/**
 * ⚠️ FLAG: z-index scale not in branding doc — standard web convention.
 * Needs Orchestrator sign-off.
 */
export const Z_INDEX = {
  base: 0,
  raised: 10,
  dropdown: 100,
  sticky: 200,
  overlay: 300,
  modal: 400,
  toast: 500,
  tooltip: 600,
} as const;

// ---------------------------------------------------------------------------
// BREAKPOINT TOKENS
// ---------------------------------------------------------------------------

/**
 * ⚠️ FLAG: responsive breakpoints not in branding doc.
 * Needs Orchestrator sign-off.
 */
export const BREAKPOINT = {
  sm: '640px',
  md: '768px',
  lg: '1024px',
  xl: '1280px',
  '2xl': '1536px',
} as const;

// ---------------------------------------------------------------------------
// STATUS COLOR MAP
// Semantic colors for ATS pipeline statuses
// Derived from branding palette — no new hex values introduced.
// ---------------------------------------------------------------------------

export const STATUS_COLOR = {
  /** Applicant submitted, pending review */
  pending:      { bg: '#DBE1E1', text: '#0B1F3A' },
  /** Under review */
  in_review:    { bg: '#DBE1E1', text: '#0B1F3A' },
  /** Invited to test */
  testing:      { bg: '#DBE1E1', text: '#0B1F3A' },
  /** Invited to interview */
  interviewing: { bg: '#DBE1E1', text: '#0B1F3A' },
  /** Applicant hired */
  hired:        { bg: '#34D186', text: '#FFFFFF' },
  /** Application dropped */
  dropped:      { bg: '#F4001A', text: '#FFFFFF' },
  /** Transferred to another position */
  transferred:  { bg: '#A3A3A3', text: '#FFFFFF' },
  /** Job posting is open/active */
  open:         { bg: '#34D186', text: '#FFFFFF' },
  /** Job posting is closed */
  closed:       { bg: '#A3A3A3', text: '#FFFFFF' },
  /** Blacklisted applicant */
  blacklisted:  { bg: '#F4001A', text: '#FFFFFF' },
} as const;

// ---------------------------------------------------------------------------
// MASTER BRAND EXPORT
// Convenience re-export for components that need a single import.
// ---------------------------------------------------------------------------

export const BRAND = {
  color:      COLOR,
  typography: TYPOGRAPHY,
  spacing:    SPACING,
  radius:     RADIUS,
  shadow:     SHADOW,
  zIndex:     Z_INDEX,
  breakpoint: BREAKPOINT,
  status:     STATUS_COLOR,
} as const;

export default BRAND;
