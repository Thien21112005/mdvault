export interface Note {
  id: string
  title: string
  content: string
  tags: string[]
  folder: string
  lastModified: Date | string
  excerpt: string
}

export interface Folder {
  id: string
  name: string
  icon: string
  children?: Folder[]
  level: number
  type?: 'markdown' | 'file'
}

export type ViewMode = 'list' | 'grid'
