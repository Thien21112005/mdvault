// [CHANGE] Đã loại bỏ directive 'use client' và kết nối hàm tạo folder xuống đĩa thật qua Wails Go
import { useState } from 'react'
import { X } from 'lucide-react'

const AVAILABLE_ICONS = ['📁', '📚', '🏗️', '💻', '🎨', '📝', '🔐', '🌟', '📊', '🚀', '⚙️', '🔧']
const FOLDER_TYPES = [
  { value: 'markdown', label: 'Markdown Notes' },
  { value: 'file', label: 'File Storage' },
]

interface CreateFolderDialogProps {
  isOpen: boolean
  onClose: () => void
  onCreateFolder: (folderName: string) => void
  parentId?: string
}

export function CreateFolderDialog({
  isOpen,
  onClose,
  onCreateFolder,
  parentId,
}: CreateFolderDialogProps) {
  const [folderName, setFolderName] = useState('')
  const [selectedIcon, setSelectedIcon] = useState('📁')
  const [selectedType, setSelectedType] = useState('markdown')

  const handleCreate = () => {
    if (!folderName.trim()) return

    // [CHANGE] Gọi hàm tạo folder thực tế với tên được nhập
    onCreateFolder(folderName.trim())
    
    setFolderName('')
    setSelectedIcon('📁')
    setSelectedType('markdown')
    onClose()
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-96 max-w-full mx-4">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Tạo thư mục mới</h2>
          <button
            onClick={onClose}
            className="p-1 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
          >
            <X size={20} className="text-gray-500" />
          </button>
        </div>

        {/* Body */}
        <div className="p-6 space-y-4">
          {/* Folder Name */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Tên thư mục
            </label>
            <input
              type="text"
              value={folderName}
              onChange={e => setFolderName(e.target.value)}
              placeholder="Nhập tên thư mục..."
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              autoFocus
            />
          </div>

          {/* Folder Type */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Loại thư mục
            </label>
            <div className="grid grid-cols-2 gap-2">
              {FOLDER_TYPES.map(type => (
                <button
                  key={type.value}
                  onClick={() => setSelectedType(type.value)}
                  className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                    selectedType === type.value
                      ? 'bg-blue-500 text-white'
                      : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
                  }`}
                >
                  {type.label}
                </button>
              ))}
            </div>
          </div>

          {/* Icon Selection */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Chọn icon
            </label>
            <div className="grid grid-cols-6 gap-2">
              {AVAILABLE_ICONS.map(icon => (
                <button
                  key={icon}
                  onClick={() => setSelectedIcon(icon)}
                  className={`p-2 rounded-lg text-2xl transition-all ${
                    selectedIcon === icon
                      ? 'bg-blue-500 ring-2 ring-blue-300 dark:ring-blue-700 scale-110'
                      : 'bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600'
                  }`}
                >
                  {icon}
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="flex gap-2 p-6 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900">
          <button
            onClick={onClose}
            className="flex-1 px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 font-medium hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
          >
            Hủy
          </button>
          <button
            onClick={handleCreate}
            disabled={!folderName.trim()}
            className="flex-1 px-4 py-2 rounded-lg bg-blue-500 text-white font-medium hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Tạo
          </button>
        </div>
      </div>
    </div>
  )
}
