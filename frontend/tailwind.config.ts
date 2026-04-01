import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./lib/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          50: "#f6f8fb",
          100: "#e9eef5",
          600: "#1f3a5f",
          700: "#182f4d",
          900: "#0f1d30",
        },
        seat: {
          available: "#16a34a",
          held: "#ea580c",
          booked: "#6b7280",
        },
      },
    },
  },
  plugins: [],
};

export default config;
