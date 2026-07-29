<script lang="ts">
  import { SvelteFlow, Background, Controls, MiniMap, type Node, type Edge } from '@xyflow/svelte'
  import '@xyflow/svelte/dist/style.css'
  import { notes, selectedNoteId } from '$lib/stores/vaultStore'
  import { viewMode } from '$lib/stores/viewStore'
  import { buildVaultGraph } from '$lib/utils/graphBuilder'

  let nodes = $state<Node[]>([])
  let edges = $state<Edge[]>([])

  $effect(() => {
    const graph = buildVaultGraph($notes)
    nodes = graph.nodes as unknown as Node[]
    edges = graph.edges as unknown as Edge[]
  })

  function onNodeClick(event: CustomEvent<{ event: MouseEvent; node: Node }>) {
    const id = event.detail.node.id
    if (id) {
      selectedNoteId.set(id)
      viewMode.set('list')
    }
  }
</script>

<div class="flex-1 bg-[#0a0a0a] h-[calc(100vh-3rem)] relative w-full overflow-hidden">
  <div class="absolute top-6 left-6 z-10 bg-[#0a0a0a]/90 backdrop-blur-md border border-zinc-800/80 px-5 py-4 rounded-xl shadow-lg max-w-sm">
    <h3 class="font-bold text-[13px] tracking-tight text-zinc-100 flex items-center space-x-2">
      <span>Knowledge Graph</span>
    </h3>
    <p class="text-[11px] text-zinc-400 mt-1 leading-relaxed">Visual representation of your notes and [[wiki-link]] connections.</p>
    <div class="mt-3 pt-3 border-t border-zinc-800/50 flex items-start space-x-2">
      <span class="text-[14px]">💡</span>
      <p class="text-[10px] text-zinc-500 font-medium leading-relaxed">Tip: Click any node to instantly jump to that note.</p>
    </div>
  </div>

  {#if $notes.length === 0}
    <div class="w-full h-full flex items-center justify-center text-zinc-600 text-[13px] font-medium tracking-wide">
      No notes in the workspace to display on the graph.
    </div>
  {:else}
    <SvelteFlow {nodes} {edges} fitView on:nodeclick={onNodeClick}>
      <Background gap={24} size={1} color="#27272a" />
      <Controls class="!bg-zinc-900 !border-zinc-800 !text-zinc-400" />
      <MiniMap class="!bg-[#0a0a0a] !border-zinc-800 rounded-lg" maskColor="#18181b80" nodeColor="#52525b" />
    </SvelteFlow>
  {/if}
</div>
