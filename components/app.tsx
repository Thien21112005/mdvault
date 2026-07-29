'use client'

import { useState, useMemo } from 'react'
import { notes, folders as initialFolders } from '@/lib/mock-data'
import { ViewMode, Folder } from '@/lib/types'
import { Sidebar } from './sidebar'
import { NotesList } from './notes-list'
import { NoteEditor } from './note-editor'
import { CreateFolderDialog } from './create-folder-dialog'

export function App() {
  const [folders, setFolders] = useState<Folder[]>(initialFolders)
  const [selectedFolder, setSelectedFolder] = useState<string | null>(null)
  const [selectedTag, setSelectedTag] = useState<string | null>(null)
  const [selectedNoteId, setSelectedNoteId] = useState<string | null>(null)
  const [viewMode, setViewMode] = useState<ViewMode>('list')
  const [isDark, setIsDark] = useState(false)
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)

  // Filter notes based on selected folder and tag
  const filteredNotes = useMemo(() => {
    return notes.filter(note => {
      const matchesFolder = selectedFolder === null || note.folder === selectedFolder
      const matchesTag = selectedTag === null || note.tags.includes(selectedTag)
      return matchesFolder && matchesTag
    })
  }, [selectedFolder, selectedTag])

  // Get selected note
  const selectedNote = filteredNotes.find(n => n.id === selectedNoteId) || null

  // Reset selected note if it's filtered out
  useMemo(() => {
    if (selectedNoteId && !filteredNotes.find(n => n.id === selectedNoteId)) {
      setSelectedNoteId(filteredNotes[0]?.id || null)
    }
  }, [filteredNotes, selectedNoteId])

  const handleThemeToggle = () => {
    setIsDark(!isDark)
    if (!isDark) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  const handleCreateFolder = (newFolder: Folder) => {
    setFolders([...folders, newFolder])
  }

  return (
    <div className={`flex h-screen ${isDark ? 'dark' : ''}`}>
      <Sidebar
        folders={folders}
        selectedFolder={selectedFolder}
        onFolderSelect={setSelectedFolder}
        selectedTag={selectedTag}
        onTagSelect={setSelectedTag}
        isDark={isDark}
        onThemeToggle={handleThemeToggle}
        onCreateFolder={() => setIsCreateDialogOpen(true)}
      />
      <NotesList
        notes={filteredNotes}
        selectedNoteId={selectedNoteId}
        onNoteSelect={setSelectedNoteId}
        viewMode={viewMode}
        onViewModeChange={setViewMode}
      />
      <NoteEditor note={selectedNote} />
      <CreateFolderDialog
        isOpen={isCreateDialogOpen}
        onClose={() => setIsCreateDialogOpen(false)}
        onCreateFolder={handleCreateFolder}
      />
    </div>
  )
}
