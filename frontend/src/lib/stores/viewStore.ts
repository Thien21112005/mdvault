import { writable } from 'svelte/store'
import type { ViewMode } from '../types/domain'

export type EditorMode = 'preview' | 'split' | 'edit'

export const viewMode = writable<ViewMode>('list')
export const isEditingNote = writable<boolean>(false)
export const isSidebarOpen = writable<boolean>(true)
export const editorMode = writable<EditorMode>('preview')

export function toggleSidebar() {
  isSidebarOpen.update(v => !v)
}

export function switchView(mode: ViewMode) {
  viewMode.set(mode)
}

export function setEditorMode(mode: EditorMode) {
  editorMode.set(mode)
  isEditingNote.set(mode !== 'preview')
}

export const isCommandPaletteOpen = writable<boolean>(false)

export function openCommandPalette() {
  isCommandPaletteOpen.set(true)
}

export function closeCommandPalette() {
  isCommandPaletteOpen.set(false)
}

export function toggleCommandPalette() {
  isCommandPaletteOpen.update(v => !v)
}


