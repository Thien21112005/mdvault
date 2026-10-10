import { writable } from 'svelte/store'
import type { ViewMode } from '../types/domain'

export const viewMode = writable<ViewMode>('list')
export const isEditingNote = writable<boolean>(false)
export const isSidebarOpen = writable<boolean>(true)

export function toggleSidebar() {
  isSidebarOpen.update(v => !v)
}

export function switchView(mode: ViewMode) {
  viewMode.set(mode)
}
