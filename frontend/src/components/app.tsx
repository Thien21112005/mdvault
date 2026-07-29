// [CHANGE] Đã loại bỏ directive 'use client' của Next.js và thay thế hoàn toàn mock-data bằng Wails Bindings & Events
import { useState, useEffect, useMemo } from 'react'
import { ViewMode, Folder, Note } from '@/lib/types'
import { Sidebar } from './sidebar'
import { NotesList } from './notes-list'
import { NoteEditor } from './note-editor'
import { CreateFolderDialog } from './create-folder-dialog'
import * as AppBindings from '../../wailsjs/go/main/App'
import { EventsOn, EventsOff } from '../../wailsjs/runtime'

export function App() {
  const [folders, setFolders] = useState<Folder[]>([])
  const [notes, setNotes] = useState<Note[]>([])
  const [allTags, setAllTags] = useState<string[]>([])
  const [rootPath, setRootPath] = useState<string>('')
  
  const [selectedFolder, setSelectedFolder] = useState<string | null>(null)
  const [selectedTag, setSelectedTag] = useState<string | null>(null)
  const [selectedNoteId, setSelectedNoteId] = useState<string | null>(null)
  const [viewMode, setViewMode] = useState<ViewMode>('list')
  const [isDark, setIsDark] = useState(false)
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)

  // [CHANGE] Hàm làm mới danh sách dữ liệu thực tế từ backend Wails
  const refreshWorkspace = async () => {
    try {
      const data = await AppBindings.ScanFolder()
      if (data) {
        setFolders((data.folders || []) as Folder[])
        setNotes((data.notes || []) as Note[])
        setAllTags(data.allTags || [])
        setRootPath(data.rootPath || '')
      }
    } catch (e) {
      console.log('Chưa có workspace mặc định:', e)
    }
  }

  // [CHANGE] Xử lý mở hộp thoại Wails chọn folder thực trên máy
  const handleSelectWorkspace = async () => {
    try {
      const data = await AppBindings.SelectWorkspace()
      if (data) {
        setFolders((data.folders || []) as Folder[])
        setNotes((data.notes || []) as Note[])
        setAllTags(data.allTags || [])
        setRootPath(data.rootPath || '')
        setSelectedFolder(null)
        setSelectedTag(null)
        setSelectedNoteId(data.notes?.[0]?.id || null)
      }
    } catch (err) {
      console.error('Lỗi khi mở workspace:', err)
    }
  }

  // [CHANGE] Lắng nghe sự kiện "on-file-change" từ backend Wails và tải dữ liệu ban đầu
  useEffect(() => {
    refreshWorkspace()

    EventsOn('on-file-change', () => {
      refreshWorkspace()
    })

    return () => {
      EventsOff('on-file-change')
    }
  }, [])

  // Filter notes based on selected folder and tag
  const filteredNotes = useMemo(() => {
    return notes.filter(note => {
      const matchesFolder = selectedFolder === null || note.folder === selectedFolder || note.folder.startsWith(selectedFolder)
      const matchesTag = selectedTag === null || (note.tags && note.tags.includes(selectedTag))
      return matchesFolder && matchesTag
    })
  }, [notes, selectedFolder, selectedTag])

  // Get selected note
  const selectedNote = filteredNotes.find(n => n.id === selectedNoteId) || null

  // Reset selected note if it's filtered out
  useEffect(() => {
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

  // [CHANGE] Tạo file markdown Note mới qua Wails Backend
  const handleCreateNote = async () => {
    if (!rootPath && folders.length === 0) {
      alert('Vui lòng mở một Workspace thư mục trước khi tạo ghi chú!')
      return
    }
    const title = prompt('Nhập tiêu đề cho ghi chú mới:', 'Ghi chú mới')
    if (!title) return

    try {
      const targetDir = selectedFolder || rootPath || 'root'
      const newNote = await AppBindings.CreateNote(targetDir, title)
      if (newNote) {
        await refreshWorkspace()
        setSelectedNoteId(newNote.id)
      }
    } catch (e) {
      alert(`Lỗi khi tạo ghi chú: ${e}`)
    }
  }

  // [CHANGE] Tạo folder trên filesystem thực
  const handleCreateFolder = async (folderName: string) => {
    try {
      const targetParent = selectedFolder || rootPath || 'root'
      await AppBindings.CreateFolder(targetParent, folderName)
      await refreshWorkspace()
      setIsCreateDialogOpen(false)
    } catch (e) {
      alert(`Lỗi khi tạo thư mục: ${e}`)
    }
  }

  // [CHANGE] Hàm lưu thay đổi nội dung file markdown
  const handleSaveNote = async (id: string, content: string, title: string, tags: string[]) => {
    try {
      await AppBindings.SaveNoteContent(id, content, title, tags)
      await refreshWorkspace()
    } catch (e) {
      alert(`Lỗi khi lưu ghi chú: ${e}`)
    }
  }

  return (
    <div className={`flex h-screen ${isDark ? 'dark' : ''}`}>
      <Sidebar
        folders={folders}
        selectedFolder={selectedFolder}
        onFolderSelect={setSelectedFolder}
        selectedTag={selectedTag}
        onTagSelect={setSelectedTag}
        allTags={allTags}
        isDark={isDark}
        onThemeToggle={handleThemeToggle}
        onCreateFolder={() => setIsCreateDialogOpen(true)}
        onSelectWorkspace={handleSelectWorkspace}
        onCreateNote={handleCreateNote}
      />
      <NotesList
        notes={filteredNotes}
        selectedNoteId={selectedNoteId}
        onNoteSelect={setSelectedNoteId}
        viewMode={viewMode}
        onViewModeChange={setViewMode}
      />
      <NoteEditor
        note={selectedNote}
        folders={folders}
        onSaveNote={handleSaveNote}
        onOpenWorkspace={handleSelectWorkspace}
      />
      <CreateFolderDialog
        isOpen={isCreateDialogOpen}
        onClose={() => setIsCreateDialogOpen(false)}
        onCreateFolder={handleCreateFolder}
      />
    </div>
  )
}
