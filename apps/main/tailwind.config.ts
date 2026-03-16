import type { Config } from 'tailwindcss'

const config: Config = {
  darkMode: ['class'],
  content: [
    './src/**/*.{ts,tsx}',
    './src/app/**/*.{ts,tsx}',
    './src/components/**/*.{ts,tsx}',
    './src/features/**/*.{ts,tsx}',
  ],
  theme: {
    container: {
      center: true,
      padding: '2rem',
      screens: {
        '2xl': '1280px',
      },
    },
    extend: {
      colors: {
        // TigerSoft Brand — Primary palette (85%)
        'vivid-red': {
          DEFAULT: '#F4001A',
          hover: '#C8001A',
          50: '#FFF1F2',
          100: '#FFE0E3',
          200: '#FFC7CC',
          300: '#FF9AA3',
          400: '#FF5C6B',
          500: '#F4001A',
          600: '#C8001A',
          700: '#9A0013',
          800: '#6E000E',
          900: '#42000A',
        },
        'oxford-blue': {
          DEFAULT: '#0B1F3A',
          50: '#E8ECF1',
          100: '#D1D9E2',
          200: '#A3B3C6',
          300: '#758DA9',
          400: '#47678D',
          500: '#1A4170',
          600: '#0B1F3A',
          700: '#081830',
          800: '#061125',
          900: '#030A1B',
        },
        'quick-silver': '#A3A3A3',
        serene: '#DBE1E1',
        'ufo-green': '#34D186',

        // shadcn/ui semantic tokens mapped to TigerSoft brand
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
        },
        secondary: {
          DEFAULT: 'hsl(var(--secondary))',
          foreground: 'hsl(var(--secondary-foreground))',
        },
        destructive: {
          DEFAULT: 'hsl(var(--destructive))',
          foreground: 'hsl(var(--destructive-foreground))',
        },
        muted: {
          DEFAULT: 'hsl(var(--muted))',
          foreground: 'hsl(var(--muted-foreground))',
        },
        accent: {
          DEFAULT: 'hsl(var(--accent))',
          foreground: 'hsl(var(--accent-foreground))',
        },
        popover: {
          DEFAULT: 'hsl(var(--popover))',
          foreground: 'hsl(var(--popover-foreground))',
        },
        card: {
          DEFAULT: 'hsl(var(--card))',
          foreground: 'hsl(var(--card-foreground))',
        },

        // ATS Feature-specific semantic tokens
        'score-high': '#34D186',
        'score-mid': '#F4001A',
        'score-low': '#A3A3A3',
      },
      fontFamily: {
        sans: ['var(--font-plus-jakarta)', 'Plus Jakarta Sans', 'system-ui', 'sans-serif'],
        thai: ['var(--font-fc-vision)', 'FC Vision', 'Sarabun', 'sans-serif'],
      },
      fontSize: {
        xs: ['0.75rem', { lineHeight: '1rem' }],
        sm: ['0.875rem', { lineHeight: '1.25rem' }],
        base: ['1rem', { lineHeight: '1.5rem' }],
        md: ['1.125rem', { lineHeight: '1.75rem' }],
        lg: ['1.25rem', { lineHeight: '1.75rem' }],
        xl: ['1.5rem', { lineHeight: '2rem' }],
        '2xl': ['1.875rem', { lineHeight: '2.25rem' }],
        '3xl': ['2.25rem', { lineHeight: '2.5rem' }],
        '4xl': ['3rem', { lineHeight: '1.1' }],
        '5xl': ['3.75rem', { lineHeight: '1.1' }],
        '6xl': ['4.5rem', { lineHeight: '1.05' }],
      },
      borderRadius: {
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 4px)',
        pill: '9999px',
        card: '16px',
        input: '8px',
        modal: '24px',
        image: '12px',
        badge: '9999px',
      },
      boxShadow: {
        sm: '0 1px 2px rgba(11, 31, 58, 0.08)',
        md: '0 4px 8px rgba(11, 31, 58, 0.12)',
        lg: '0 8px 16px rgba(11, 31, 58, 0.16)',
        cta: '0 4px 12px rgba(244, 0, 26, 0.30)',
      },
      spacing: {
        'page-x': '2.5rem',
        'section-y': '5rem',
        'card-p': '1.5rem',
        'component': '1rem',
        'tight': '0.5rem',
      },
      keyframes: {
        'accordion-down': {
          from: { height: '0' },
          to: { height: 'var(--radix-accordion-content-height)' },
        },
        'accordion-up': {
          from: { height: 'var(--radix-accordion-content-height)' },
          to: { height: '0' },
        },
      },
      animation: {
        'accordion-down': 'accordion-down 0.2s ease-out',
        'accordion-up': 'accordion-up 0.2s ease-out',
      },
    },
  },
  plugins: [require('tailwindcss-animate')],
}

export default config
