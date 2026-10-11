export interface Note {
  id: string
  title: string
  content: string
  tags: string[]
  folder: string
  lastModified: string
  excerpt: string
  metadata?: Record<string, any>
}

export interface Folder {
  id: string
  name: string
  icon: string
  children?: Folder[]
  level: number
  type?: string
}

export interface WorkspaceData {
  rootPath: string
  folders: Folder[]
  notes: Note[]
  allTags: string[]
}

export type ViewMode = 'list' | 'graph'
