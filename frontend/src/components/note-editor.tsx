// [CHANGE] Đã loại bỏ directive 'use client' và bổ sung chế độ Edit/Save ghi lại xuống ổ đĩa qua Wails
import { useState, useEffect } from 'react'
import { ChevronRight, Edit3, Save, Eye, FolderOpen } from 'lucide-react'
import { Note, Folder } from '@/lib/types'

interface NoteEditorProps {
  note: Note | null
  folders: Folder[]
  onSaveNote?: (id: string, content: string, title: string, tags: string[]) => void
  onOpenWorkspace?: () => void
}

export function NoteEditor({ note, folders, onSaveNote, onOpenWorkspace }: NoteEditorProps) {
  const [isEditing, setIsEditing] = useState(false)
  const [editContent, setEditContent] = useState('')
  const [editTitle, setEditTitle] = useState('')

  useEffect(() => {
    if (note) {
      setEditContent(note.content)
      setEditTitle(note.title)
      setIsEditing(false)
    }
  }, [note])

  if (!note) {
    return (
      <div className="flex-1 bg-white dark:bg-gray-800 flex flex-col items-center justify-center p-6">
        <div className="text-center max-w-md">
          <div className="text-6xl mb-4">📄</div>
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
            MDVault Desktop
          </h2>
          <p className="text-gray-600 dark:text-gray-400 mb-6">
            Chọn một thư mục Workspace hoặc một ghi chú từ bên trái để bắt đầu đọc và chỉnh sửa các tài liệu Markdown của bạn.
          </p>
          <button
            onClick={onOpenWorkspace}
            className="inline-flex items-center gap-2 px-6 py-3 rounded-xl bg-blue-600 text-white font-medium hover:bg-blue-700 transition-all shadow-md"
          >
            <FolderOpen size={18} />
            <span>Mở Workspace Markdown</span>
          </button>
        </div>
      </div>
    )
  }

  const getFolderPath = (folderId: string): string[] => {
    if (!folderId || folderId === 'root') return ['Workspace']
    const parts = folderId.split('/')
    return parts.slice(-2)
  }

  const handleSave = () => {
    if (onSaveNote && note) {
      onSaveNote(note.id, editContent, editTitle, note.tags || [])
      setIsEditing(false)
    }
  }

  const pathItems = getFolderPath(note.folder)

  return (
    <div className="flex-1 bg-white dark:bg-gray-800 flex flex-col h-screen overflow-hidden">
      {/* Header */}
      <div className="border-b border-gray-200 dark:border-gray-700 p-6 flex flex-col gap-4 shrink-0">
        <div className="flex items-center justify-between">
          {/* Breadcrumb */}
          <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
            {pathItems.map((item, idx) => (
              <div key={idx} className="flex items-center gap-2">
                {idx > 0 && <ChevronRight size={16} />}
                <span>{item}</span>
              </div>
            ))}
          </div>

          {/* [CHANGE] Nút chuyển chế độ Xem / Chỉnh sửa và Lưu file thật */}
          <div className="flex items-center gap-2">
            {isEditing ? (
              <>
                <button
                  onClick={() => setIsEditing(false)}
                  className="px-3 py-1.5 rounded-lg border border-gray-300 dark:border-gray-600 text-xs text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors inline-flex items-center gap-1.5"
                >
                  <Eye size={14} />
                  <span>Hủy</span>
                </button>
                <button
                  onClick={handleSave}
                  className="px-3 py-1.5 rounded-lg bg-blue-600 text-white text-xs font-medium hover:bg-blue-700 transition-all shadow-sm inline-flex items-center gap-1.5"
                >
                  <Save size={14} />
                  <span>Lưu file</span>
                </button>
              </>
            ) : (
              <button
                onClick={() => setIsEditing(true)}
                className="px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-xs font-medium text-gray-800 dark:text-gray-200 transition-colors inline-flex items-center gap-1.5"
              >
                <Edit3 size={14} />
                <span>Chỉnh sửa</span>
              </button>
            )}
          </div>
        </div>

        {/* Title */}
        {isEditing ? (
          <input
            type="text"
            value={editTitle}
            onChange={e => setEditTitle(e.target.value)}
            className="text-2xl font-bold text-gray-900 dark:text-white bg-transparent border-b border-gray-300 dark:border-gray-600 focus:outline-none focus:border-blue-500 pb-1"
          />
        ) : (
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white truncate">{note.title}</h1>
        )}

        {/* Tags and metadata */}
        <div className="flex flex-wrap gap-2 items-center">
          {note.tags && note.tags.map(tag => (
            <span
              key={tag}
              className="px-3 py-0.5 rounded-full bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 text-xs font-medium"
            >
              #{tag}
            </span>
          ))}
          <span className="text-xs text-gray-400 dark:text-gray-500 ml-auto font-mono">
            {note.id} • {note.lastModified ? new Date(note.lastModified).toLocaleString('vi-VN') : ''}
          </span>
        </div>
      </div>

      {/* Content Area */}
      <div className="flex-1 overflow-y-auto p-6">
        {isEditing ? (
          <textarea
            value={editContent}
            onChange={e => setEditContent(e.target.value)}
            className="w-full h-full p-4 bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100 font-mono text-sm rounded-xl border border-gray-200 dark:border-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 leading-relaxed resize-none"
            placeholder="Soạn thảo markdown..."
            autoFocus
          />
        ) : (
          <article className="prose dark:prose-invert max-w-none prose-sm">
            <div className="text-gray-800 dark:text-gray-200 leading-relaxed">
              {note.content.split('\n\n').map((paragraph, idx) => {
                if (paragraph.startsWith('#')) {
                  if (paragraph.startsWith('### ')) {
                    return (
                      <h3 key={idx} className="text-lg font-bold mt-4 mb-2 text-gray-900 dark:text-white">
                        {paragraph.replace(/^### /, '')}
                      </h3>
                    )
                  }
                  if (paragraph.startsWith('## ')) {
                    return (
                      <h2 key={idx} className="text-xl font-bold mt-6 mb-3 text-gray-900 dark:text-white">
                        {paragraph.replace(/^## /, '')}
                      </h2>
                    )
                  }
                  if (paragraph.startsWith('# ')) {
                    return (
                      <h1 key={idx} className="text-2xl font-bold mt-8 mb-4 text-gray-900 dark:text-white">
                        {paragraph.replace(/^# /, '')}
                      </h1>
                    )
                  }
                }

                if (paragraph.startsWith('```')) {
                  const [, lang, ...codeParts] = paragraph.split('\n')
                  const code = codeParts.join('\n').replace(/```$/, '')
                  return (
                    <pre
                      key={idx}
                      className="bg-gray-900 text-gray-100 p-4 rounded-lg my-4 overflow-x-auto text-sm font-mono"
                    >
                      <code>{code}</code>
                    </pre>
                  )
                }

                if (paragraph.startsWith('> ')) {
                  return (
                    <blockquote
                      key={idx}
                      className="border-l-4 border-blue-500 pl-4 italic text-gray-600 dark:text-gray-400 my-4"
                    >
                      {paragraph.replace(/^> /, '')}
                    </blockquote>
                  )
                }

                if (paragraph.startsWith('- ')) {
                  const items = paragraph.split('\n')
                  return (
                    <ul key={idx} className="list-disc pl-6 my-2">
                      {items.map((item, i) => (
                        <li key={i} className="text-gray-800 dark:text-gray-200 mb-1">
                          {item.replace(/^- /, '')}
                        </li>
                      ))}
                    </ul>
                  )
                }

                return (
                  <p key={idx} className="text-gray-800 dark:text-gray-200 mb-4 leading-relaxed">
                    {paragraph}
                  </p>
                )
              })}
            </div>
          </article>
        )}
      </div>
    </div>
  )
}
