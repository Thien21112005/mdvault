import React from 'react'
import ReactDOM from 'react-dom/client'
import { App } from './components/app'
import './index.css'

// [CHANGE] Khai báo điểm Render client cho SPA thay thế cho App Router của Next.js (app/page.tsx)
ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
