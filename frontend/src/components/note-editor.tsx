// [CHANGE] Đã loại bỏ directive 'use client' của Next.js
import { ChevronRight } from 'lucide-react'
import { Note } from '@/lib/types'
import { folders } from '@/lib/mock-data'

interface NoteEditorProps {
  note: Note | null
}

export function NoteEditor({ note }: NoteEditorProps) {
  if (!note) {
    return (
      <div className="flex-1 bg-white dark:bg-gray-800 flex flex-col items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">📄</div>
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
            Chọn thư mục để bắt đầu
          </h2>
          <p className="text-gray-600 dark:text-gray-400 max-w-md">
            Chọn một thư mục từ bên trái để xem ghi chú của bạn, hoặc tạo một ghi chú mới.
          </p>
        </div>
      </div>
    )
  }

  const getFolderPath = (folderId: string): string[] => {
    const path: string[] = []

    const findFolder = (id: string, depth = 0): boolean => {
      for (const folder of folders) {
        if (folder.id === id) {
          path.unshift(folder.name)
          return true
        }
        if (folder.children) {
          for (const child of folder.children) {
            if (child.id === id) {
              path.unshift(child.name)
              path.unshift(folder.name)
              return true
            }
          }
        }
      }
      return false
    }

    findFolder(folderId)
    return path
  }

  const pathItems = getFolderPath(note.folder)
  const folderName =
    folders.find(f => f.id === note.folder)?.name ||
    folders.flatMap(f => f.children || []).find(f => f.id === note.folder)?.name

  return (
    <div className="flex-1 bg-white dark:bg-gray-800 flex flex-col">
      {/* Header */}
      <div className="border-b border-gray-200 dark:border-gray-700 p-6">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-3">
          {pathItems.map((item, idx) => (
            <div key={idx} className="flex items-center gap-2">
              {idx > 0 && <ChevronRight size={16} />}
              <span>{item}</span>
            </div>
          ))}
        </div>

        {/* Title */}
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">{note.title}</h1>

        {/* Tags and metadata */}
        <div className="flex flex-wrap gap-2 items-center">
          {note.tags.map(tag => (
            <span
              key={tag}
              className="px-3 py-1 rounded-full bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 text-sm font-medium"
            >
              #{tag}
            </span>
          ))}
          <span className="text-sm text-gray-500 dark:text-gray-400 ml-auto">
            Modified {new Date(note.lastModified).toLocaleDateString('vi-VN')}
          </span>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-6">
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
      </div>
    </div>
  )
}
