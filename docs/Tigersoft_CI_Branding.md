# TigerSoft CI — Web App Design System
> CI New Vision 2025 | theory & design tokens reference สำหรับ Web Application

---

## Colors

### Primary (85% of design)
| Token | Value | Role |
|---|---|---|
| `vivid-red` | `#F4001A` | CTA, accent, brand highlight — 20% |
| `white` | `#FFFFFF` | Background, card, surface — 45% |
| `oxford-blue` | `#0B1F3A` | Text, heading, dark bg — 20% |

### Secondary (15% of design)
| Token | Value | Role |
|---|---|---|
| `quick-silver` | `#A3A3A3` | Border, divider, disabled — 5% |
| `serene` | `#DBE1E1` | Subtle bg, light card — 5% |
| `ufo-green` | `#34D186` | Success, active, positive — 5% |

### Semantic Mapping
| Role | Value |
|---|---|
| `bg-primary` | white |
| `bg-secondary` | serene |
| `bg-dark` | oxford-blue |
| `text-primary` | oxford-blue (**ห้ามใช้ #000 เด็ดขาด**) |
| `text-inverse` | white |
| `text-secondary` | quick-silver |
| `action-primary` | vivid-red |
| `action-hover` | `#C8001A` (darken 15%) |
| `action-success` | ufo-green |
| `border-default` | serene |
| `border-focus` | vivid-red |

### Gradients (allowed combinations only)
White → Vivid Red → Oxford Blue → Quick Silver → Serene → White (ทิศทางเดียวกันเท่านั้น)

---

## Typography

### Fonts
- **English:** Plus Jakarta Sans (Google Fonts)
- **Thai:** FC Vision (custom — `guide/CI Toolkit/Font/TH/FC Vision/`)

### Weights
| Weight | Value | Usage |
|---|---|---|
| Light | 300 | Body, paragraph, description |
| Medium | 500 | Heading, subheading, button, label |

### Size Scale (base 16px)
`12 / 14 / 16 / 18 / 20 / 24 / 30 / 36 / 48 / 60 / 72px`

### Roles
| Role | Weight | Line Height | Letter Spacing |
|---|---|---|---|
| Heading | 500 | 1.2 | −0.02em |
| Subheading | 500 | 1.4 | 0 |
| Body | 300 | 1.6 | 0 |
| Label / Badge | 500 | — | +0.05em |

---

## Spacing

Base unit: **4px** — scale: `4 / 8 / 12 / 16 / 20 / 24 / 32 / 40 / 48 / 64 / 80 / 96 / 128px`

| Semantic Token | Value | Usage |
|---|---|---|
| `page-x` | 40px | Horizontal page margin |
| `section-y` | 80px | Section vertical padding |
| `card` | 24px | Card internal padding |
| `component` | 16px | Between inner elements |
| `tight` | 8px | Dense / compact content |

---

## Border Radius (Soft Edges)

`4 / 8 / 12 / 16 / 24 / 32 / 9999px`

| Token | Value | Usage |
|---|---|---|
| `button` | 9999px | Full pill — all buttons |
| `card` | 16px | Cards, panels |
| `input` | 8px | Input fields |
| `modal` | 24px | Modals, drawers |
| `image` | 12px | Clipping mask containers |
| `badge` | 9999px | Tags, chips, avatars |

---

## Shadows & Elevation

Shadows ใช้โทน **Oxford Blue** (ไม่ใช่ black) เพื่อความนุ่มนวล
ระดับ: `xs → sm → md → lg → xl → 2xl` — เพิ่ม blur และ spread ตามลำดับ

- **Red shadow** — ใช้เฉพาะ CTA button เท่านั้น (Vivid Red tint)
- **Glassmorphism** — `blur(12px)` + white overlay 12% + white border 20% — ใช้เฉพาะ **hero / featured section** เท่านั้น

---

## Grid & Layout

| Property | Value |
|---|---|
| Columns | 4 (mobile) / 8 (tablet) / 12 (desktop) |
| Gutter | 24px |
| Page margin | 16px / 32px / 40px |
| Max content width | 1280px |

---

## Motion

| Token | Value | Usage |
|---|---|---|
| `fast` | 150ms | Hover, color change |
| `normal` | 250ms | Standard transition |
| `slow` | 400ms | Page element entrance |
| `ease-brand` | `cubic-bezier(0.22, 1, 0.36, 1)` | Default easing |
| `ease-spring` | `cubic-bezier(0.34, 1.56, 0.64, 1)` | Bounce / delight |

---

## Component Principles

| Component | Key Rules |
|---|---|
| **Button** | Full pill · Vivid Red primary · Oxford Blue outline secondary · Medium 500 · Red shadow on primary |
| **Card** | 16px radius · 24px padding · serene border · elevate on hover |
| **Input** | 8px radius · serene border default · vivid-red border on focus · Light 300 |
| **Badge** | Full pill · 12px font · Medium 500 · tinted background per status |
| **Icon** | Outline style · 1.5px stroke · Vivid Red accent on key icons · white container + shadow |
| **Navbar** | White bg · serene border-bottom · 64px height · 40px horizontal padding |
| **Section** | 80px vertical · 40px horizontal · max-width 1280px centered |

---

## Accessibility

| Pair | Contrast | Grade |
|---|---|---|
| Oxford Blue on White | 14.8:1 | AAA ✅ |
| Vivid Red on White | 4.6:1 | AA ✅ |
| White on Oxford Blue | 14.8:1 | AAA ✅ |
| UFO Green on Oxford Blue | 6.2:1 | AA ✅ |

- Focus ring: Vivid Red ring `rgba(244, 0, 26, 0.30)` ขนาด 3px
- Minimum touch target: **44px**

---

## Quick Rules

| ✅ Do | ❌ Don't |
|---|---|
| ใช้ Oxford Blue แทน #000 เสมอ | ใช้ pure black (#000) |
| White space ≥ 45% ของ layout | ยัด element จนแน่น |
| Soft edge ทุก component | ใช้ sharp corner โดยไม่จำเป็น |
| Vivid Red เฉพาะ CTA และ accent | ใช้ red เป็น body หรือ bg หลัก |
| Glassmorphism เฉพาะ hero | ใช้ glass effect ทั่วไป |
| Primary 85% / Secondary 15% | ผสมสีเกินสัดส่วน |
| Gradient ตาม chain ที่กำหนด | ใช้ gradient combination ที่ไม่อนุญาต |

---