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

<div class="flex-1 bg-transparent h-[calc(100vh-3rem)] relative w-full overflow-hidden">
  <div class="absolute top-6 left-6 z-10 bg-white/[0.01] backdrop-blur-xl border border-white/[0.06] px-5 py-4 rounded-xl shadow-2xl max-w-sm">
    <h3 class="font-bold text-[13px] tracking-tighter text-zinc-100 flex items-center space-x-2">
      <span>Knowledge Graph</span>
    </h3>
    <p class="text-[11px] text-zinc-400 mt-1 leading-relaxed">Visual representation of your notes and [[wiki-link]] connections.</p>
    <div class="mt-3 pt-3 border-t border-white/[0.06] flex items-start space-x-2">
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
      <Background gap={24} size={1} color="#3f3f46" />
      <Controls class="!bg-white/[0.02] !border-white/[0.06] !text-zinc-400 backdrop-blur-md shadow-lg" />
      <MiniMap class="!bg-black/50 !border-white/[0.06] backdrop-blur-md rounded-lg shadow-lg" maskColor="#00000080" nodeColor="#52525b" />
    </SvelteFlow>
  {/if}
</div>
