/** @type {import('tailwindcss').Config} */
export default {
  content: ['packages/web/templates/**/*.templ'],
  theme: {
    extend: {
      colors: {
        ternary: '#F1F2F6',
      }
    },
  },
  plugins: [],
}

