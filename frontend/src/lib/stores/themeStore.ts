import { writable } from 'svelte/store'

export const isDark = writable<boolean>(true)

export function toggleTheme() {
  isDark.update(dark => {
    const next = !dark
    if (typeof document !== 'undefined') {
      if (next) {
        document.documentElement.classList.add('dark')
        localStorage.setItem('mdvault-theme', 'dark')
      } else {
        document.documentElement.classList.remove('dark')
        localStorage.setItem('mdvault-theme', 'light')
      }
    }
    return next
  })
}

export function initTheme() {
  if (typeof window === 'undefined') return
  const saved = localStorage.getItem('mdvault-theme')
  const shouldBeDark = saved !== null ? saved === 'dark' : true
  isDark.set(shouldBeDark)
  if (shouldBeDark) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}
