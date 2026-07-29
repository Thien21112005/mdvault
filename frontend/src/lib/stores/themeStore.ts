import { writable } from 'svelte/store'

export const isDark = writable<boolean>(false)

export function toggleTheme() {
  isDark.update(dark => {
    const next = !dark
    if (typeof document !== 'undefined') {
      if (next) {
        document.documentElement.classList.add('dark')
      } else {
        document.documentElement.classList.remove('dark')
      }
    }
    return next
  })
}
