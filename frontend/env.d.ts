/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly DEPLOY_BASE_URL: string
    
  }
  
  interface ImportMeta {
    readonly env: ImportMetaEnv
  }