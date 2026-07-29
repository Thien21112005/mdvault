import { writable } from 'svelte/store'
import type { ViewMode } from '../types/domain'

export const viewMode = writable<ViewMode>('list')
export const isEditingNote = writable<boolean>(false)
