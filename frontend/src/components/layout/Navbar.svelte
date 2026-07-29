<script lang="ts">
  import { Folder, Search, Network, FileText } from '@lucide/svelte'
  import { viewMode } from '$lib/stores/viewStore'
  import { searchQuery, openWorkspaceDialog, rootPath } from '$lib/stores/vaultStore'

  function switchView(mode: 'list' | 'graph') {
    viewMode.set(mode)
  }
</script>

<header class="h-12 bg-white/[0.01] backdrop-blur-md border-b border-white/[0.06] text-zinc-300 flex items-center justify-between px-4 select-none relative z-50">
  <!-- Brand -->
  <div class="flex items-center space-x-3 w-64">
    <div class="flex items-center justify-center w-5 h-5 rounded-md bg-white/[0.05] text-zinc-200 border border-white/[0.08] shadow-[inset_0_1px_1px_rgba(255,255,255,0.1)]">
      <Folder size={12} strokeWidth={2.5} />
    </div>
    <span class="font-bold text-[13px] tracking-tight text-zinc-100">
      MDVault
    </span>
    {#if $rootPath}
      <span class="text-[9px] uppercase tracking-widest text-zinc-500 font-semibold px-2 py-0.5 rounded border border-white/[0.06] bg-white/[0.02] max-w-[120px] truncate" title={$rootPath}>
        {$rootPath.split('/').pop() || $rootPath}
      </span>
    {/if}
  </div>

  <!-- Search Bar -->
  <div class="flex-1 max-w-md flex justify-center">
    <div class="relative w-full max-w-sm group">
      <Search size={13} class="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-500 group-focus-within:text-zinc-300 transition-colors" />
      <input
        type="text"
        bind:value={$searchQuery}
        placeholder="Search notes, tags..."
        class="w-full bg-white/[0.03] border border-white/[0.06] rounded-md pl-9 pr-3 py-1 text-[13px] text-zinc-200 placeholder-zinc-500 focus:outline-none focus:border-white/[0.12] focus:bg-white/[0.05] transition-all shadow-inner"
      />
    </div>
  </div>

  <!-- Actions -->
  <div class="flex items-center justify-end space-x-2 w-64">
    <!-- View Switcher -->
    <div class="bg-white/[0.03] p-0.5 rounded-md border border-white/[0.06] flex items-center shadow-inner">
      <button
        onclick={() => switchView('list')}
        class="flex items-center justify-center px-2.5 py-1 rounded-[4px] text-[11px] font-medium transition-all {$viewMode === 'list' ? 'bg-white/[0.08] text-zinc-100 shadow-[inset_0_1px_1px_rgba(255,255,255,0.1)]' : 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04]'}"
      >
        <FileText size={13} class="mr-1.5" />
        List
      </button>
      <button
        onclick={() => switchView('graph')}
        class="flex items-center justify-center px-2.5 py-1 rounded-[4px] text-[11px] font-medium transition-all {$viewMode === 'graph' ? 'bg-white/[0.08] text-zinc-100 shadow-[inset_0_1px_1px_rgba(255,255,255,0.1)]' : 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04]'}"
      >
        <Network size={13} class="mr-1.5" />
        Graph
      </button>
    </div>

    <!-- Workspace button -->
    <button
      onclick={openWorkspaceDialog}
      class="bg-zinc-100 hover:bg-white text-black px-3 py-1 rounded-md text-[11px] font-semibold transition-all active:scale-[0.97] flex items-center space-x-1.5 shadow-[0_0_10px_rgba(255,255,255,0.1)]"
    >
      <span>Open</span>
    </button>

  </div>
</header>
