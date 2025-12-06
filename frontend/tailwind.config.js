/** @type {import('tailwindcss').Config} */
export default {
  darkMode: ["class"],
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'ui-sans-serif', 'system-ui', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'sans-serif'],
        serif: ['Playfair Display', 'Georgia', 'Cambria', 'Times New Roman', 'Times', 'serif'],
      },
      colors: {
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
        // オークション専用カラー (拡張)
        auction: {
          gold: "#D4AF37",
          'gold-light': "#D4A84B", // シャンパンゴールド (HSL: 43 74% 49%)
          'gold-dark': "#C4982D",
          green: "#10B981",
          'green-racing': "#3A7D5C", // レーシンググリーン (HSL: 152 45% 35%)
          'green-light': "#4A9D6C",
          'green-dark': "#2A5D4C",
          red: "#EF4444",
          blue: "#3B82F6",
          burgundy: "#8B2942", // ボルドー (HSL: 345 65% 35%)
          'burgundy-light': "#A5374F",
          'burgundy-dark': "#6B1F32",
          platinum: "#B8BCC6", // プラチナ (HSL: 220 15% 75%)
          silver: "#C0C0C0",
          cream: "#FAF8F5", // クリーム (HSL: 40 33% 98%)
          ivory: "#FDFCFB", // ウォームホワイト (HSL: 40 20% 99%)
          'warm-gray': "#EDEBE8", // ウォームグレー (HSL: 40 15% 93%)
        },
        status: {
          pending: "hsl(43, 96%, 56%)",
          active: "hsl(152, 45%, 35%)",
          ended: "hsl(40, 15%, 93%)",
          cancelled: "hsl(0, 84%, 60%)",
        },
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
      boxShadow: {
        'luxury': '0 2px 8px rgba(212, 168, 75, 0.1), 0 1px 3px rgba(0, 0, 0, 0.05)',
        'luxury-lg': '0 4px 16px rgba(212, 168, 75, 0.15), 0 2px 6px rgba(0, 0, 0, 0.05)',
        'luxury-xl': '0 8px 32px rgba(212, 168, 75, 0.2), 0 4px 12px rgba(0, 0, 0, 0.1)',
        'gold-glow': '0 0 16px rgba(212, 168, 75, 0.4), 0 0 32px rgba(212, 168, 75, 0.2)',
        'inner-luxury': 'inset 0 2px 4px rgba(212, 168, 75, 0.1)',
      },
      keyframes: {
        "accordion-down": {
          from: { height: "0" },
          to: { height: "var(--radix-accordion-content-height)" },
        },
        "accordion-up": {
          from: { height: "var(--radix-accordion-content-height)" },
          to: { height: "0" },
        },
        "shimmer": {
          "0%": { backgroundPosition: "-1000px 0" },
          "100%": { backgroundPosition: "1000px 0" },
        },
        "fade-in": {
          from: { opacity: "0" },
          to: { opacity: "1" },
        },
        "slide-in-up": {
          from: { transform: "translateY(10px)", opacity: "0" },
          to: { transform: "translateY(0)", opacity: "1" },
        },
        "scale-in": {
          from: { transform: "scale(0.95)", opacity: "0" },
          to: { transform: "scale(1)", opacity: "1" },
        },
        "gold-shimmer": {
          "0%": { backgroundPosition: "-1000px 0" },
          "100%": { backgroundPosition: "1000px 0" },
        },
        "pulse-gold": {
          "0%, 100%": { opacity: "1", transform: "scale(1)" },
          "50%": { opacity: "0.9", transform: "scale(1.02)" },
        },
        "winning-glow": {
          "0%, 100%": { boxShadow: "0 0 8px rgba(212, 168, 75, 0.3)" },
          "50%": { boxShadow: "0 0 20px rgba(212, 168, 75, 0.6)" },
        },
      },
      animation: {
        "accordion-down": "accordion-down 0.2s ease-out",
        "accordion-up": "accordion-up 0.2s ease-out",
        "shimmer": "shimmer 2s infinite linear",
        "fade-in": "fade-in 0.3s ease-out",
        "slide-in-up": "slide-in-up 0.4s ease-out",
        "scale-in": "scale-in 0.3s ease-out",
        "gold-shimmer": "gold-shimmer 3s infinite linear",
        "pulse-gold": "pulse-gold 2s cubic-bezier(0.4, 0, 0.6, 1) infinite",
        "winning-glow": "winning-glow 2s ease-in-out infinite",
      },
      backgroundImage: {
        'gold-gradient': 'linear-gradient(135deg, #D4A84B 0%, #C4982D 100%)',
        'burgundy-gradient': 'linear-gradient(135deg, #8B2942 0%, #6B1F32 100%)',
        'luxury-gradient': 'linear-gradient(135deg, #FAF8F5 0%, #EDEBE8 100%)',
      },
    },
  },
  plugins: [],
}
