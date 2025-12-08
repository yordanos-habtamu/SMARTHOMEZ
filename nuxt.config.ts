// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig(
  {
    compatibilityDate: "2025-07-15",
    devtools: { enabled: true },
    modules: ['@nuxtjs/google-fonts'],

    runtimeConfig: {
      public: {
        apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:4000/api/v1'
      }
    },

    googleFonts: {
      families: {
        Roboto: [100, 300, 400, 500, 700, 900],
        'Open Sans': [100, 300, 400, 500, 700],
        'Lato': [100, 300, 400, 500, 700, 900],
        'Montserrat': [100, 300, 400, 500, 700, 900],
        'Poppins': [100, 300, 400, 500, 700, 900],

      }
    },
    css: ['~/assets/css/main.css'],
    vite: { plugins: [tailwindcss(),], },
  }
);