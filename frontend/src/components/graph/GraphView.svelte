<script lang="ts">
  import { SvelteFlow, Background, Controls, MiniMap, type Node, type Edge } from '@xyflow/svelte'
  import '@xyflow/svelte/dist/style.css'
  import { Sparkles, Orbit, RotateCcw } from '@lucide/svelte'
  import { notes, selectedNoteId } from '$lib/stores/vaultStore'
  import { viewMode } from '$lib/stores/viewStore'
  import { buildVaultGraph, type GraphNodeData } from '$lib/utils/graphBuilder'
  import PlanetNode from './PlanetNode.svelte'

  const nodeTypes = {
    planet: PlanetNode
  }

  let nodes = $state<Node[]>([])
  let edges = $state<Edge[]>([])
  let hoveredNode = $state<GraphNodeData | null>(null)

  $effect(() => {
    const graph = buildVaultGraph($notes)
    nodes = graph.nodes as unknown as Node[]
    edges = graph.edges as unknown as Edge[]
  })

  function onNodeClick({ node }: { node: Node }) {
    const id = node.id
    if (id) {
      selectedNoteId.set(id)
      viewMode.set('list')
    }
  }

  function onNodePointerEnter({ node }: { node: Node }) {
    hoveredNode = (node.data as unknown as GraphNodeData) || null
    const targetId = node.id

    // Highlight connecting constellation paths
    edges = edges.map(edge => {
      const isConnected = edge.source === targetId || edge.target === targetId
      if (isConnected) {
        return {
          ...edge,
          animated: true,
          style: 'stroke: #818cf8; stroke-width: 2.5px; opacity: 1; filter: drop-shadow(0 0 6px #6366f1);'
        }
      }
      return {
        ...edge,
        animated: false,
        style: 'stroke: rgba(129, 140, 248, 0.1); stroke-width: 1px; opacity: 0.15;'
      }
    })
  }

  function onNodePointerLeave() {
    hoveredNode = null
    // Restore default starry edges
    edges = edges.map(edge => ({
      ...edge,
      animated: false,
      style: 'stroke: rgba(129, 140, 248, 0.4); stroke-width: 1.5px; opacity: 0.6;'
    }))
  }
</script>

<div class="flex-1 h-[calc(100vh-3.25rem)] relative w-full overflow-hidden bg-[radial-gradient(ellipse_at_center,_var(--bg-surface-elevated)_0%,_var(--bg-base)_100%)] select-none">
  <!-- Galaxy HUD Information Card (Top-Left) -->
  <div class="absolute top-5 left-5 z-10 bg-[var(--bg-surface-elevated)]/80 backdrop-blur-xl border border-[var(--border-subtle)] p-4 rounded-2xl shadow-[var(--shadow-glass)] max-w-sm transition-all duration-300">
    <div class="flex items-center justify-between space-x-2">
      <div class="flex items-center space-x-2">
        <div class="p-1.5 rounded-lg bg-[var(--accent-primary)]/15 text-[var(--accent-primary)]">
          <Orbit size={16} class="animate-spin" style="animation-duration: 20s;" />
        </div>
        <div>
          <h3 class="font-bold text-[13px] tracking-tight text-[var(--text-primary)] heading-display flex items-center space-x-1.5">
            <span>Galaxy Knowledge Graph</span>
          </h3>
          <p class="text-[10px] text-[var(--text-muted)] font-mono">
            {nodes.length} celestial notes • {edges.length} constellation links
          </p>
        </div>
      </div>
    </div>

    <!-- Hovered Planet Detail Pill -->
    {#if hoveredNode}
      <div class="mt-3 pt-3 border-t border-[var(--border-subtle)] space-y-1.5 animate-scale-in">
        <div class="flex items-center space-x-2">
          <span
            class="w-2.5 h-2.5 rounded-full shrink-0 shadow-sm"
            style="background: {hoveredNode.color}; box-shadow: 0 0 8px {hoveredNode.glowColor};"
          ></span>
          <p class="text-[12px] font-semibold text-[var(--text-primary)] truncate">
            {hoveredNode.label}
          </p>
        </div>
        <div class="flex items-center space-x-2 text-[10px] font-mono text-[var(--text-secondary)]">
          <span class="bg-[var(--bg-surface-active)] px-1.5 py-0.5 rounded border border-[var(--border-subtle)]">
            {hoveredNode.connections} {hoveredNode.connections === 1 ? 'connection' : 'connections'}
          </span>
          {#if hoveredNode.tags && hoveredNode.tags.length > 0}
            <span class="text-[var(--text-muted)] truncate">
              #{hoveredNode.tags.slice(0, 2).join(' #')}
            </span>
          {/if}
        </div>
      </div>
    {:else}
      <div class="mt-2.5 pt-2.5 border-t border-[var(--border-subtle)] flex items-start space-x-2">
        <Sparkles size={12} class="text-[var(--accent-primary)] mt-0.5 shrink-0" />
        <p class="text-[10.5px] text-[var(--text-muted)] leading-relaxed">
          Hover over planets to light up constellations. Click any planet to jump directly into its editor.
        </p>
      </div>
    {/if}
  </div>

  {#if $notes.length === 0}
    <div class="w-full h-full flex flex-col items-center justify-center text-[var(--text-muted)] space-y-3">
      <div class="p-4 rounded-2xl bg-[var(--bg-surface-elevated)] border border-[var(--border-subtle)]">
        <Orbit size={28} class="text-[var(--text-muted)]" />
      </div>
      <p class="text-[13px] font-medium tracking-wide">No notes in the vault to map the galaxy.</p>
    </div>
  {:else}
    <SvelteFlow
      {nodes}
      {edges}
      {nodeTypes}
      fitView
      onnodeclick={onNodeClick}
      onnodepointerenter={onNodePointerEnter}
      onnodepointerleave={onNodePointerLeave}
    >
      <Background gap={32} size={1.2} patternColor="rgba(129, 140, 248, 0.12)" />
      <Controls
        class="!bg-[var(--bg-surface-elevated)]/80 !border-[var(--border-subtle)] !text-[var(--text-secondary)] !shadow-[var(--shadow-glass)] backdrop-blur-xl rounded-2xl overflow-hidden [&_button]:!border-[var(--border-subtle)] [&_button]:!bg-transparent hover:[&_button]:!bg-[var(--bg-surface-hover)] [&_button_svg]:!fill-current"
      />
      <MiniMap
        class="!bg-[var(--bg-surface-elevated)]/80 !border-[var(--border-subtle)] backdrop-blur-xl rounded-2xl shadow-[var(--shadow-glass)] overflow-hidden"
        maskColor="rgba(0, 0, 0, 0.45)"
        nodeColor="#6366f1"
      />
    </SvelteFlow>
  {/if}
</div>
