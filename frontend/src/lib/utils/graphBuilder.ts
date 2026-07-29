import type { Note } from '../types/domain'

export interface GraphNode {
  id: string
  data: { label: string; id: string; excerpt: string }
  position: { x: number; y: number }
  class?: string
}

export interface GraphEdge {
  id: string
  source: string
  target: string
  animated?: boolean
  style?: string
}

export function buildVaultGraph(notes: Note[]): { nodes: GraphNode[]; edges: GraphEdge[] } {
  const nodes: GraphNode[] = []
  const edges: GraphEdge[] = []
  const noteMap = new Map<string, string>()

  const radius = Math.max(250, notes.length * 50)
  const angleStep = notes.length > 0 ? (2 * Math.PI) / notes.length : 0

  notes.forEach((note, i) => {
    const angle = i * angleStep
    const x = Math.round(Math.cos(angle) * radius) + 400
    const y = Math.round(Math.sin(angle) * radius) + 300

    nodes.push({
      id: note.id,
      data: { label: note.title || 'Untitled Note', id: note.id, excerpt: note.excerpt },
      position: { x, y },
      class: 'bg-white dark:bg-slate-800 border-2 border-blue-500 rounded-xl px-4 py-2.5 shadow-lg text-xs font-bold text-slate-800 dark:text-slate-100 cursor-pointer hover:border-emerald-500 transition-all min-w-[120px] text-center'
    })

    if (note.title) {
      noteMap.set(note.title.toLowerCase().trim(), note.id)
    }
  })

  // Parse [[wiki-links]] trong nội dung note để tạo kết nối Edge cho xyflow
  const wikiRegex = /\[\[(.*?)\]\]/g
  notes.forEach(sourceNote => {
    let match: RegExpExecArray | null
    const content = sourceNote.content || ''
    while ((match = wikiRegex.exec(content)) !== null) {
      const targetTitle = match[1].toLowerCase().trim()
      const targetId = noteMap.get(targetTitle)
      if (targetId && targetId !== sourceNote.id) {
        const edgeId = `e-${sourceNote.id}-${targetId}`
        if (!edges.some(e => e.id === edgeId)) {
          edges.push({
            id: edgeId,
            source: sourceNote.id,
            target: targetId,
            animated: true,
            style: 'stroke: #3b82f6; stroke-width: 2px;'
          })
        }
      }
    }
  })

  return { nodes, edges }
}
