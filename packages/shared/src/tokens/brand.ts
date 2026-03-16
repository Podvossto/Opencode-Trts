/**
 * TigerSoft ATS — Design Tokens
 * Source: docs/Tigersoft_CI_Branding.md
 *
 * DEV RULES:
 * - NEVER use raw hex values in components. Always import from this file.
 * - NEVER use #000000 (pure black) — use BRAND.color.oxford_blue for dark text.
 * - CTA buttons MUST use vivid_red + pill radius + red shadow.
 * - White space must be ≥ 45% of any layout.
 * - Glassmorphism ONLY on hero sections.
 */

export const BRAND = {
  color: {
    // Primary palette
    vivid_red:    '#F4001A', // CTAs, primary action buttons, active states
    white:        '#FFFFFF', // Backgrounds, card surfaces
    oxford_blue:  '#0B1F3A', // Body text, headings — use instead of #000
    quick_silver: '#A3A3A3', // Placeholder text, disabled states, captions
    serene:       '#DBE1E1', // Borders, dividers, input outlines at rest
    ufo_green:    '#34D186', // Success states, hired badge, pass indicator

    // Semantic aliases (use these in components)
    text_primary:   '#0B1F3A', // oxford_blue
    text_secondary: '#A3A3A3', // quick_silver
    text_on_cta:    '#FFFFFF', // white — text on vivid_red buttons
    border_default: '#DBE1E1', // serene
    bg_surface:     '#FFFFFF', // white
    bg_page:        '#F8F9FA', // off-white page background (not in branding doc — flagged)
    status_success: '#34D186', // ufo_green
    status_error:   '#F4001A', // vivid_red
    status_warning: '#F59E0B', // amber — not in branding doc, flagged for Orchestrator review
  },

  font: {
    en: 'Plus Jakarta Sans', // All English text — imported via Google Fonts
    th: 'FC Vision',         // All Thai text — imported via brand font CDN
    fallback: 'system-ui, sans-serif',
  },

  fontSize: {
    // Scale in rem (base 16px)
    xs:   '0.75rem',  // 12px — labels, captions
    sm:   '0.875rem', // 14px — body small
    base: '1rem',     // 16px — body
    lg:   '1.125rem', // 18px — body large
    xl:   '1.25rem',  // 20px — subheading
    '2xl':'1.5rem',   // 24px — heading sm
    '3xl':'1.875rem', // 30px — heading md
    '4xl':'2.25rem',  // 36px — heading lg
    '5xl':'3rem',     // 48px — hero heading
  },

  fontWeight: {
    regular:   '400',
    medium:    '500',
    semibold:  '600',
    bold:      '700',
    extrabold: '800',
  },

  spacing: {
    // Base unit: 4px — all spacing is multiples of 4
    '0':  '0px',
    '1':  '4px',
    '2':  '8px',
    '3':  '12px',
    '4':  '16px',
    '5':  '20px',
    '6':  '24px',
    '8':  '32px',
    '10': '40px',
    '12': '48px',
    '16': '64px',
    '20': '80px',
    '24': '96px',
  },

  radius: {
    none:   '0px',
    sm:     '4px',   // small elements, tags
    input:  '8px',   // all form inputs
    card:   '16px',  // cards, modals, panels
    pill:   '9999px', // CTA buttons, badges — required for primary buttons
  },

  shadow: {
    // CTA button shadow — required on vivid_red buttons
    cta:    '0 4px 14px 0 rgba(244, 0, 26, 0.35)',
    card:   '0 2px 12px 0 rgba(11, 31, 58, 0.08)',
    modal:  '0 8px 40px 0 rgba(11, 31, 58, 0.16)',
    none:   'none',
  },

  motion: {
    // Duration
    fast:   '150ms',
    normal: '250ms',
    slow:   '400ms',
    // Easing
    ease_brand: 'cubic-bezier(0.22, 1, 0.36, 1)',
    ease_in:    'cubic-bezier(0.4, 0, 1, 1)',
    ease_out:   'cubic-bezier(0, 0, 0.2, 1)',
    ease_linear:'linear',
  },

  // Glassmorphism — ONLY for hero sections
  glass: {
    background: 'rgba(255, 255, 255, 0.12)',
    backdrop:   'blur(16px)',
    border:     '1px solid rgba(255, 255, 255, 0.20)',
  },

  // Z-index layers
  zIndex: {
    base:    '0',
    raised:  '10',
    dropdown:'100',
    sticky:  '200',
    overlay: '300',
    modal:   '400',
    toast:   '500',
  },

  // Layout
  layout: {
    max_width:      '1280px',
    content_width:  '1024px',
    sidebar_width:  '240px',
    header_height:  '64px',
    // White space rule: minimum 45% of layout area
    white_space_min: '45%',
  },
} as const;

// Type helper for consumers
export type BrandColor     = keyof typeof BRAND.color;
export type BrandFont      = keyof typeof BRAND.font;
export type BrandRadius    = keyof typeof BRAND.radius;
export type BrandShadow    = keyof typeof BRAND.shadow;
export type BrandMotion    = keyof typeof BRAND.motion;

// ⚠️  FLAGGED TOKENS — not found in docs/Tigersoft_CI_Branding.md
// Orchestrator must confirm before use in components:
//   - bg_page (#F8F9FA) — assumed off-white page bg
//   - status_warning (#F59E0B) — amber, not in branding doc
