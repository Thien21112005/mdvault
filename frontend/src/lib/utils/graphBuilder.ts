import type { Note } from '../types/domain'

export interface GraphNodeData {
  label: string
  id: string
  excerpt: string
  connections: number
  size: number
  color: string
  glowColor: string
  tags: string[]
}

export interface GraphNode {
  id: string
  type: string
  data: GraphNodeData
  position: { x: number; y: number }
}

export interface GraphEdge {
  id: string
  source: string
  target: string
  animated?: boolean
  style?: string
}

const GALAXY_PALETTES = [
  { color: '#6366f1', glow: 'rgba(99, 102, 241, 0.5)' }, // Indigo Pulsar
  { color: '#a855f7', glow: 'rgba(168, 85, 247, 0.5)' }, // Violet Nebula
  { color: '#38bdf8', glow: 'rgba(56, 189, 248, 0.5)' }, // Cyan Comet
  { color: '#10b981', glow: 'rgba(16, 185, 129, 0.5)' }, // Emerald Aurora
  { color: '#f59e0b', glow: 'rgba(245, 158, 11, 0.5)' }, // Solar Amber
  { color: '#f43f5e', glow: 'rgba(244, 63, 94, 0.5)' }, // Nova Rose
]

function getPlanetPalette(note: Note, index: number) {
  const tagsStr = (note.tags || []).join(' ').toLowerCase()
  if (tagsStr.includes('code') || tagsStr.includes('dev') || tagsStr.includes('tech') || tagsStr.includes('api')) {
    return GALAXY_PALETTES[0] // Indigo
  }
  if (tagsStr.includes('design') || tagsStr.includes('ui') || tagsStr.includes('art')) {
    return GALAXY_PALETTES[1] // Violet
  }
  if (tagsStr.includes('doc') || tagsStr.includes('guide') || tagsStr.includes('book')) {
    return GALAXY_PALETTES[2] // Cyan
  }
  if (tagsStr.includes('project') || tagsStr.includes('work') || tagsStr.includes('vault')) {
    return GALAXY_PALETTES[3] // Emerald
  }
  if (tagsStr.includes('daily') || tagsStr.includes('todo') || tagsStr.includes('task')) {
    return GALAXY_PALETTES[4] // Amber
  }
  if (tagsStr.includes('personal') || tagsStr.includes('life') || tagsStr.includes('idea')) {
    return GALAXY_PALETTES[5] // Rose
  }

  // Fallback to deterministic pseudo-random hash
  const hash = (note.title || '').split('').reduce((acc, char) => acc + char.charCodeAt(0), index)
  return GALAXY_PALETTES[Math.abs(hash) % GALAXY_PALETTES.length]
}

export function buildVaultGraph(notes: Note[]): { nodes: GraphNode[]; edges: GraphEdge[] } {
  if (notes.length === 0) {
    return { nodes: [], edges: [] }
  }

  const edges: GraphEdge[] = []
  const noteMap = new Map<string, string>()
  const connectionCounts = new Map<string, number>()

  notes.forEach(note => {
    if (note.title) {
      noteMap.set(note.title.toLowerCase().trim(), note.id)
    }
    connectionCounts.set(note.id, 0)
  })

  // Parse [[wiki-links]] to construct edges
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
            animated: false,
            style: 'stroke: rgba(129, 140, 248, 0.4); stroke-width: 1.5px; transition: all 0.3s ease;'
          })

          connectionCounts.set(sourceNote.id, (connectionCounts.get(sourceNote.id) || 0) + 1)
          connectionCounts.set(targetId, (connectionCounts.get(targetId) || 0) + 1)
        }
      }
    }
  })

  // Sort notes by connection count descending to put central hub planets near galaxy center
  const sortedNotes = [...notes].sort((a, b) => {
    return (connectionCounts.get(b.id) || 0) - (connectionCounts.get(a.id) || 0)
  })

  const nodes: GraphNode[] = []
  const centerX = 500
  const centerY = 400
  const goldenAngle = 2.399963 // ~137.5 degrees in radians
  const spacing = notes.length > 20 ? 80 : 110

  sortedNotes.forEach((note, i) => {
    const connections = connectionCounts.get(note.id) || 0
    const palette = getPlanetPalette(note, i)

    // Calculate planet diameter based on degree (34px - 62px)
    const size = Math.min(62, Math.max(34, 34 + connections * 7))

    // Position using golden-angle phyllotaxis spiral
    let x = centerX
    let y = centerY
    if (sortedNotes.length > 1) {
      const radius = 100 + Math.sqrt(i) * spacing
      const angle = i * goldenAngle
      x = Math.round(centerX + Math.cos(angle) * radius)
      y = Math.round(centerY + Math.sin(angle) * radius)
    }

    nodes.push({
      id: note.id,
      type: 'planet',
      data: {
        label: note.title || 'Untitled Note',
        id: note.id,
        excerpt: note.excerpt || '',
        connections,
        size,
        color: palette.color,
        glowColor: palette.glow,
        tags: note.tags || []
      },
      position: { x, y }
    })
  })

  return { nodes, edges }
}
