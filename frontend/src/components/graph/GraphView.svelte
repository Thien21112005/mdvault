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

<div class="flex-1 bg-slate-950 h-[calc(100vh-3.5rem)] relative w-full overflow-hidden">
  <div class="absolute top-4 left-4 z-10 bg-slate-900/90 border border-slate-800 px-4 py-2 rounded-xl shadow-lg">
    <h3 class="font-bold text-sm text-blue-400 flex items-center space-x-2">
      <span>🌐 Mạng Lưới Kiến Thức (Graph View)</span>
    </h3>
    <p class="text-xs text-slate-400 mt-0.5">Sơ đồ biểu diễn các Ghi chú và kết nối Wiki-link [[note]]</p>
    <p class="text-[11px] text-emerald-400 mt-1">💡 Mẹo: Nhấp vào nút bất kỳ để chuyển ngay tới ghi chú đó!</p>
  </div>

  {#if $notes.length === 0}
    <div class="w-full h-full flex items-center justify-center text-slate-600 text-sm">
      Chưa có ghi chú nào trong Workspace để hiển thị Đồ thị.
    </div>
  {:else}
    <SvelteFlow {nodes} {edges} fitView on:nodeclick={onNodeClick}>
      <Background gap={20} size={1.5} color="#334155" />
      <Controls class="bg-slate-800 text-white border-slate-700" />
      <MiniMap class="bg-slate-900 border border-slate-800 rounded-lg" />
    </SvelteFlow>
  {/if}
</div>
