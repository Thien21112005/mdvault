// [CHANGE] Đã loại bỏ directive 'use client' và thêm nút Mở Workspace + Tạo Note thật từ Wails
import { useState } from 'react'
import { Search, ChevronDown, ChevronRight, FolderOpen, Moon, Sun, FilePlus, FolderPlus } from 'lucide-react'
import { Folder } from '@/lib/types'

interface SidebarProps {
  folders: Folder[]
  selectedFolder: string | null
  onFolderSelect: (folderId: string | null) => void
  selectedTag: string | null
  onTagSelect: (tag: string | null) => void
  allTags: string[]
  isDark: boolean
  onThemeToggle: () => void
  onCreateFolder: () => void
  onSelectWorkspace: () => void
  onCreateNote: () => void
}

export function Sidebar({
  folders,
  selectedFolder,
  onFolderSelect,
  selectedTag,
  onTagSelect,
  allTags,
  isDark,
  onThemeToggle,
  onCreateFolder,
  onSelectWorkspace,
  onCreateNote,
}: SidebarProps) {
  const [expandedFolders, setExpandedFolders] = useState<Set<string>>(new Set(['root', '1', '2', '3']))
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
            onFolderSelect(selectedFolder === folder.id ? null : folder.id)
            if (hasChildren) toggleFolder(folder.id)
          }}
          className={`w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm transition-colors ${
            selectedFolder === folder.id
              ? 'bg-blue-500/20 text-blue-600 dark:text-blue-400 font-semibold'
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
          <span className="text-base">{folder.icon || '📁'}</span>
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
      {/* Header with theme toggle & workspace action */}
      <div className="p-4 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center justify-between mb-3">
          <h1 className="font-bold text-lg text-gray-900 dark:text-white">MDVault</h1>
          <div className="flex items-center gap-1">
            <button
              onClick={onThemeToggle}
              className="p-2 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-lg transition-colors text-gray-600 dark:text-gray-300"
              title="Chuyển chế độ tối/sings"
            >
              {isDark ? <Sun size={18} /> : <Moon size={18} />}
            </button>
          </div>
        </div>

        {/* [CHANGE] Nút Mở Workspace thật và tạo Ghi chú qua Wails */}
        <div className="flex gap-2 mb-3">
          <button
            onClick={onSelectWorkspace}
            className="flex-1 flex items-center justify-center gap-2 py-2 px-3 rounded-lg bg-blue-600 text-white text-xs font-medium hover:bg-blue-700 transition-all shadow-sm"
            title="Mở thư mục Workspace trên máy"
          >
            <FolderOpen size={15} />
            <span>Mở Workspace</span>
          </button>
          <button
            onClick={onCreateNote}
            className="flex items-center justify-center p-2 rounded-lg bg-emerald-600 text-white hover:bg-emerald-700 transition-all shadow-sm"
            title="Tạo ghi chú markdown mới"
          >
            <FilePlus size={16} />
          </button>
        </div>

        {/* Search */}
        <div className="relative">
          <Search size={16} className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
          <input
            type="text"
            placeholder="Search notes..."
            value={searchQuery}
            onChange={e => setSearchQuery(e.target.value)}
            className="w-full pl-9 pr-3 py-1.5 rounded-lg bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 text-xs text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      {/* Folders */}
      <div className="flex-1 overflow-y-auto p-4">
        <div className="mb-6">
          <div className="flex items-center justify-between mb-2">
            <div className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase">
              Thư mục
            </div>
            <button
              onClick={onCreateFolder}
              className="p-1 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-lg transition-colors text-gray-500 dark:text-gray-400 hover:text-blue-500"
              title="Tạo thư mục mới trong Workspace"
            >
              <FolderPlus size={16} />
            </button>
          </div>
          {folders.length === 0 ? (
            <div className="text-xs text-gray-400 italic py-2 text-center">
              Chưa tải thư mục nào
            </div>
          ) : (
            folders.map(folder => renderFolder(folder))
          )}
        </div>

        {/* Tags */}
        {allTags && allTags.length > 0 && (
          <div>
            <div className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase mb-2">
              Tags
            </div>
            <div className="flex flex-wrap gap-1.5">
              {allTags.map(tag => (
                <button
                  key={tag}
                  onClick={() => onTagSelect(selectedTag === tag ? null : tag)}
                  className={`px-2.5 py-1 rounded-full text-xs font-medium transition-colors ${
                    selectedTag === tag
                      ? 'bg-blue-500 text-white shadow-sm'
                      : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'
                  }`}
                >
                  #{tag}
                </button>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
