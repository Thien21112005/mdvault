import { writable, derived, get } from 'svelte/store'
import type { Folder, Note, WorkspaceData } from '../types/domain'
import { WailsAdapter } from '../wails/wailsAdapter'
import { toast } from './toastStore'


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

export async function createNewNote(title: string = 'Untitled Note'): Promise<boolean> {
  const currentWs = get(rootPath)
  const targetFolder = get(selectedFolderId) || currentWs || 'root'
  if (!currentWs && get(folders).length === 0) {
    toast.warning('Chưa mở Workspace', 'Vui lòng chọn thư mục lưu trữ để bắt đầu tạo ghi chú!', {
      action: {
        label: 'Open Vault',
        onClick: () => {
          openWorkspaceDialog()
        }
      }
    })
    return false
  }
  const note = await WailsAdapter.createNote(targetFolder, title)
  if (note) {
    await loadInitialWorkspace()
    selectedNoteId.set(note.id)
    return true
  }
  return false
}

export async function createNewFolder(name: string): Promise<boolean> {
  const currentWs = get(rootPath)
  const targetParent = get(selectedFolderId) || currentWs || 'root'
  if (!currentWs && get(folders).length === 0) {
    toast.warning('Chưa mở Workspace', 'Vui lòng chọn thư mục lưu trữ để bắt đầu tạo thư mục!', {
      action: {
        label: 'Open Vault',
        onClick: () => {
          openWorkspaceDialog()
        }
      }
    })
    return false
  }
  const success = await WailsAdapter.createFolder(targetParent, name)
  if (success) {
    await loadInitialWorkspace()
    return true
  }
  return false
}


export async function saveActiveNote(id: string, content: string, title: string, tags: string[]) {
  const success = await WailsAdapter.saveNoteContent(id, content, title, tags)
  if (success) {
    await loadInitialWorkspace()
  }
}

export async function deleteNote(id: string): Promise<boolean> {
  const currentNotes = get(notes)
  const currentIndex = currentNotes.findIndex(n => n.id === id)
  const success = await WailsAdapter.deleteNote(id)
  if (success) {
    await loadInitialWorkspace()
    const updatedNotes = get(notes)
    if (get(selectedNoteId) === id) {
      if (updatedNotes.length > 0) {
        const nextIndex = Math.min(Math.max(0, currentIndex), updatedNotes.length - 1)
        selectedNoteId.set(updatedNotes[nextIndex]?.id || updatedNotes[0].id)
      } else {
        selectedNoteId.set(null)
      }
    }
    return true
  }
  return false
}

export async function deleteFolder(path: string): Promise<boolean> {
  const currentFolder = get(selectedFolderId)
  const success = await WailsAdapter.deleteFolder(path)
  if (success) {
    if (currentFolder === path || (currentFolder && currentFolder.startsWith(path))) {
      selectedFolderId.set(null)
    }
    await loadInitialWorkspace()
    return true
  }
  return false
}

export async function renameNote(oldPath: string, newTitle: string): Promise<Note | null> {
  const updatedNote = await WailsAdapter.renameNote(oldPath, newTitle)
  if (updatedNote) {
    await loadInitialWorkspace()
    if (get(selectedNoteId) === oldPath) {
      selectedNoteId.set(updatedNote.id)
    }
    return updatedNote
  }
  return null
}

export async function renameFolder(oldPath: string, newName: string): Promise<boolean> {
  const success = await WailsAdapter.renameFolder(oldPath, newName)
  if (success) {
    if (get(selectedFolderId) === oldPath) {
      selectedFolderId.set(null)
    }
    await loadInitialWorkspace()
    return true
  }
  return false
}

let unsubscribeWatcher: (() => void) | null = null
export function initWatcher() {
  if (unsubscribeWatcher) unsubscribeWatcher()
  unsubscribeWatcher = WailsAdapter.onFileChange(() => {
    loadInitialWorkspace()
  })
}
