import { writable, derived, get } from 'svelte/store'
import type { Folder, Note, WorkspaceData } from '../types/domain'
import { WailsAdapter } from '../wails/wailsAdapter'

export const folders = writable<Folder[]>([])
export const notes = writable<Note[]>([])
export const allTags = writable<string[]>([])
export const rootPath = writable<string>('')

export const selectedFolderId = writable<string | null>(null)
export const selectedTag = writable<string | null>(null)
export const selectedNoteId = writable<string | null>(null)
export const searchQuery = writable<string>('')

export type SortOption = 'newest' | 'oldest' | 'title-asc' | 'title-desc'
export const sortOption = writable<SortOption>('newest')

// Derived filtered & sorted notes
export const filteredNotes = derived(
  [notes, selectedFolderId, selectedTag, searchQuery, sortOption],
  ([$notes, $folderId, $tag, $query, $sort]) => {
    const list = $notes.filter(note => {
      const matchesFolder = !$folderId || note.folder === $folderId || (note.folder && note.folder.startsWith($folderId))
      const matchesTag = !$tag || (note.tags && note.tags.includes($tag))
      const q = $query.toLowerCase().trim()
      const matchesQuery = !q || (note.title && note.title.toLowerCase().includes(q)) || (note.content && note.content.toLowerCase().includes(q))
      return matchesFolder && matchesTag && matchesQuery
    })

    return list.sort((a, b) => {
      if ($sort === 'newest') {
        return (b.lastModified || '').localeCompare(a.lastModified || '')
      } else if ($sort === 'oldest') {
        return (a.lastModified || '').localeCompare(b.lastModified || '')
      } else if ($sort === 'title-asc') {
        return (a.title || '').localeCompare(b.title || '')
      } else if ($sort === 'title-desc') {
        return (b.title || '').localeCompare(a.title || '')
      }
      return 0
    })
  }
)

// Derived selected note object
export const selectedNote = derived(
  [notes, selectedNoteId],
  ([$notes, $id]) => $notes.find(n => n.id === $id) || null
)

// Store actions
export async function loadInitialWorkspace() {
  const data = await WailsAdapter.scanFolder()
  if (data) updateWorkspaceState(data)
}

export async function openWorkspaceDialog() {
  const data = await WailsAdapter.selectWorkspace()
  if (data) {
    updateWorkspaceState(data)
    selectedFolderId.set(null)
    selectedTag.set(null)
    if (data.notes && data.notes.length > 0) {
      selectedNoteId.set(data.notes[0].id)
    } else {
      selectedNoteId.set(null)
    }
  }
}

export function updateWorkspaceState(data: WorkspaceData) {
  folders.set(data.folders || [])
  notes.set(data.notes || [])
  allTags.set(data.allTags || [])
  rootPath.set(data.rootPath || '')

  const currentSelected = get(selectedNoteId)
  if (currentSelected && !data.notes?.some(n => n.id === currentSelected)) {
    selectedNoteId.set(data.notes?.[0]?.id || null)
  } else if (!currentSelected && data.notes?.length > 0) {
    selectedNoteId.set(data.notes[0].id)
  }
}

export async function createNewNote(title: string) {
  const currentWs = get(rootPath)
  const targetFolder = get(selectedFolderId) || currentWs || 'root'
  if (!currentWs && get(folders).length === 0) {
    alert('Vui lòng mở Workspace trước khi tạo ghi chú!')
    return
  }
  const note = await WailsAdapter.createNote(targetFolder, title)
  if (note) {
    await loadInitialWorkspace()
    selectedNoteId.set(note.id)
  }
}

export async function createNewFolder(name: string) {
  const targetParent = get(selectedFolderId) || get(rootPath) || 'root'
  const success = await WailsAdapter.createFolder(targetParent, name)
  if (success) {
    await loadInitialWorkspace()
  }
}

export async function saveActiveNote(id: string, content: string, title: string, tags: string[]) {
  const success = await WailsAdapter.saveNoteContent(id, content, title, tags)
  if (success) {
    await loadInitialWorkspace()
  }
}

let unsubscribeWatcher: (() => void) | null = null
export function initWatcher() {
  if (unsubscribeWatcher) unsubscribeWatcher()
  unsubscribeWatcher = WailsAdapter.onFileChange(() => {
    loadInitialWorkspace()
  })
}
