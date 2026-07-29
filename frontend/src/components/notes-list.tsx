// [CHANGE] Đã loại bỏ directive 'use client' của Next.js
import { LayoutGrid, List } from 'lucide-react'
import { Note, ViewMode } from '@/lib/types'

interface NotesListProps {
  notes: Note[]
  selectedNoteId: string | null
  onNoteSelect: (noteId: string) => void
  viewMode: ViewMode
  onViewModeChange: (mode: ViewMode) => void
}

export function NotesList({
  notes,
  selectedNoteId,
  onNoteSelect,
  viewMode,
  onViewModeChange,
}: NotesListProps) {
  const formatDate = (date: Date | string) => {
    return new Date(date).toLocaleDateString('vi-VN', {
      month: 'short',
      day: 'numeric',
    })
  }

  if (notes.length === 0) {
    return (
      <div className="flex-1 border-r border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 flex flex-col items-center justify-center">
        <div className="text-center">
          <div className="text-4xl mb-4">📭</div>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
            Không có ghi chú
          </h3>
          <p className="text-gray-600 dark:text-gray-400">
            Chọn một thư mục khác hoặc hãy tạo ghi chú mới
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="flex-1 border-r border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 flex flex-col">
      {/* Header */}
      <div className="border-b border-gray-200 dark:border-gray-700 p-4 flex items-center justify-between">
        <h2 className="font-semibold text-gray-900 dark:text-white">
          {notes.length} ghi chú
        </h2>
        <div className="flex gap-2">
          <button
            onClick={() => onViewModeChange('list')}
            className={`p-2 rounded-lg transition-colors ${
              viewMode === 'list'
                ? 'bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-300'
                : 'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-400'
            }`}
          >
            <List size={18} />
          </button>
          <button
            onClick={() => onViewModeChange('grid')}
            className={`p-2 rounded-lg transition-colors ${
              viewMode === 'grid'
                ? 'bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-300'
                : 'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-400'
            }`}
          >
            <LayoutGrid size={18} />
          </button>
        </div>
      </div>

      {/* Notes Container */}
      <div className="flex-1 overflow-y-auto">
        {viewMode === 'list' ? (
          <div className="divide-y divide-gray-200 dark:divide-gray-700">
            {notes.map(note => (
              <button
                key={note.id}
                onClick={() => onNoteSelect(note.id)}
                className={`w-full p-4 text-left transition-colors border-l-4 ${
                  selectedNoteId === note.id
                    ? 'bg-blue-50 dark:bg-blue-900/20 border-l-blue-500'
                    : 'hover:bg-gray-50 dark:hover:bg-gray-700 border-l-transparent'
                }`}
              >
                <h3 className="font-semibold text-gray-900 dark:text-white mb-1 truncate">
                  {note.title}
                </h3>
                <p className="text-sm text-gray-600 dark:text-gray-400 line-clamp-2 mb-2">
                  {note.excerpt}
                </p>
                <div className="flex items-center justify-between">
                  <div className="flex gap-1 flex-wrap">
                    {note.tags.slice(0, 2).map(tag => (
                      <span
                        key={tag}
                        className="text-xs bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 px-2 py-1 rounded"
                      >
                        {tag}
                      </span>
                    ))}
                    {note.tags.length > 2 && (
                      <span className="text-xs text-gray-500 dark:text-gray-400">
                        +{note.tags.length - 2}
                      </span>
                    )}
                  </div>
                  <span className="text-xs text-gray-500 dark:text-gray-400 ml-auto">
                    {formatDate(note.lastModified)}
                  </span>
                </div>
              </button>
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-2 gap-4 p-4">
            {notes.map(note => (
              <button
                key={note.id}
                onClick={() => onNoteSelect(note.id)}
                className={`p-4 rounded-lg border-2 text-left transition-all ${
                  selectedNoteId === note.id
                    ? 'bg-blue-50 dark:bg-blue-900/20 border-blue-500'
                    : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600 bg-white dark:bg-gray-800'
                }`}
              >
                <h3 className="font-semibold text-gray-900 dark:text-white mb-2 line-clamp-2">
                  {note.title}
                </h3>
                <p className="text-sm text-gray-600 dark:text-gray-400 line-clamp-3 mb-3">
                  {note.excerpt}
                </p>
                <div className="flex flex-wrap gap-1 mb-2">
                  {note.tags.slice(0, 2).map(tag => (
                    <span
                      key={tag}
                      className="text-xs bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 px-2 py-1 rounded"
                    >
                      {tag}
                    </span>
                  ))}
                </div>
                <div className="text-xs text-gray-500 dark:text-gray-400">
                  {formatDate(note.lastModified)}
                </div>
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
