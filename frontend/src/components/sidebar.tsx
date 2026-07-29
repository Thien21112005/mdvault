// [CHANGE] Đã loại bỏ directive 'use client' của Next.js
import { useState } from 'react'
import { Search, ChevronDown, ChevronRight, FolderOpen, Moon, Sun } from 'lucide-react'
import { allTags } from '@/lib/mock-data'
import { Folder } from '@/lib/types'

interface SidebarProps {
  folders: Folder[]
  selectedFolder: string | null
  onFolderSelect: (folderId: string) => void
  selectedTag: string | null
  onTagSelect: (tag: string | null) => void
  isDark: boolean
  onThemeToggle: () => void
  onCreateFolder: () => void
}

export function Sidebar({
  folders,
  selectedFolder,
  onFolderSelect,
  selectedTag,
  onTagSelect,
  isDark,
  onThemeToggle,
  onCreateFolder,
}: SidebarProps) {
  const [expandedFolders, setExpandedFolders] = useState<Set<string>>(new Set(['1', '2', '3']))
  const [searchQuery, setSearchQuery] = useState('')

  const toggleFolder = (folderId: string) => {
    const newExpanded = new Set(expandedFolders)
    if (newExpanded.has(folderId)) {
      newExpanded.delete(folderId)
    } else {
      newExpanded.add(folderId)
    }
    setExpandedFolders(newExpanded)
  }

  const renderFolder = (folder: Folder, depth: number = 0) => {
    const isExpanded = expandedFolders.has(folder.id)
    const hasChildren = folder.children && folder.children.length > 0

    return (
      <div key={folder.id}>
        <button
          onClick={() => {
            onFolderSelect(folder.id)
            if (hasChildren) toggleFolder(folder.id)
          }}
          className={`w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm transition-colors ${
            selectedFolder === folder.id
              ? 'bg-blue-500/20 text-blue-600 dark:text-blue-400'
              : 'hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300'
          }`}
          style={{ paddingLeft: `${12 + depth * 16}px` }}
        >
          {hasChildren && (
            <span className="w-4 flex items-center justify-center">
              {isExpanded ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
            </span>
          )}
          {!hasChildren && <span className="w-4" />}
          <span className="text-base">{folder.icon}</span>
          <span className="flex-1 text-left truncate">{folder.name}</span>
        </button>

        {hasChildren && isExpanded && (
          <div>
            {folder.children?.map(child => renderFolder(child, depth + 1))}
          </div>
        )}
      </div>
    )
  }

  return (
    <div className="w-64 border-r border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900 flex flex-col h-screen">
      {/* Header with theme toggle */}
      <div className="p-4 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center justify-between mb-4">
          <h1 className="font-bold text-lg text-gray-900 dark:text-white">Notes</h1>
          <button
            onClick={onThemeToggle}
            className="p-2 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-lg transition-colors"
          >
            {isDark ? <Sun size={18} /> : <Moon size={18} />}
          </button>
        </div>

        {/* Search */}
        <div className="relative">
          <Search size={18} className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
          <input
            type="text"
            placeholder="Search notes..."
            value={searchQuery}
            onChange={e => setSearchQuery(e.target.value)}
            className="w-full pl-10 pr-3 py-2 rounded-lg bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 text-sm text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      {/* Folders */}
      <div className="flex-1 overflow-y-auto p-4">
        <div className="mb-6">
          <div className="flex items-center justify-between mb-3">
            <div className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase">
              Chọn thư mục
            </div>
            <button
              onClick={onCreateFolder}
              className="p-1 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-lg transition-colors text-gray-500 dark:text-gray-400 font-semibold text-lg"
              title="Tạo thư mục mới"
            >
              +
            </button>
          </div>
          {folders.map(folder => renderFolder(folder))}
        </div>

        {/* Tags */}
        <div>
          <div className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase mb-3">
            Tags
          </div>
          <div className="flex flex-wrap gap-2">
            {allTags.map(tag => (
              <button
                key={tag}
                onClick={() => onTagSelect(selectedTag === tag ? null : tag)}
                className={`px-3 py-1 rounded-full text-xs font-medium transition-colors ${
                  selectedTag === tag
                    ? 'bg-blue-500 text-white'
                    : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'
                }`}
              >
                {tag}
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
