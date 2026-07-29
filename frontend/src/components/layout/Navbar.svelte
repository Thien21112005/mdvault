<script lang="ts">
  import { Folder, Search, Network, FileText, Moon, Sun } from '@lucide/svelte'
  import { viewMode } from '$lib/stores/viewStore'
  import { searchQuery, openWorkspaceDialog, rootPath } from '$lib/stores/vaultStore'
  import { isDark, toggleTheme } from '$lib/stores/themeStore'

  function switchView(mode: 'list' | 'graph') {
    viewMode.set(mode)
  }
</script>

<header class="h-12 bg-zinc-950 border-b border-zinc-800/80 text-zinc-300 flex items-center justify-between px-4 select-none">
  <!-- Brand -->
  <div class="flex items-center space-x-3 w-64">
    <div class="flex items-center justify-center w-6 h-6 rounded bg-zinc-800 text-zinc-200 shadow-sm border border-zinc-700/50">
      <Folder size={14} strokeWidth={2.5} />
    </div>
    <span class="font-bold text-sm tracking-tight text-zinc-100">
      MDVault
    </span>
    {#if $rootPath}
      <span class="text-[10px] uppercase tracking-widest text-zinc-500 font-semibold px-2 py-0.5 rounded border border-zinc-800 bg-zinc-900/50 max-w-[120px] truncate" title={$rootPath}>
        {$rootPath.split('/').pop() || $rootPath}
      </span>
    {/if}
  </div>

  <!-- Search Bar -->
  <div class="flex-1 max-w-md flex justify-center">
    <div class="relative w-full max-w-sm group">
      <Search size={14} class="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-500 group-focus-within:text-zinc-300 transition-colors" />
      <input
        type="text"
        bind:value={$searchQuery}
        placeholder="Search notes, tags..."
        class="w-full bg-zinc-900/50 border border-zinc-800/80 rounded-md pl-9 pr-3 py-1.5 text-xs text-zinc-200 placeholder-zinc-500 focus:outline-none focus:border-zinc-700 focus:bg-zinc-800/50 transition-all shadow-inner"
      />
    </div>
  </div>

  <!-- Actions -->
  <div class="flex items-center justify-end space-x-2 w-64">
    <!-- View Switcher (Segmented Control) -->
    <div class="bg-zinc-900/80 p-0.5 rounded-md border border-zinc-800/80 flex items-center shadow-sm">
      <button
        onclick={() => switchView('list')}
        class="flex items-center justify-center px-2.5 py-1 rounded-[4px] text-[11px] font-medium transition-all {$viewMode === 'list' ? 'bg-zinc-700/50 text-zinc-100 shadow-sm' : 'text-zinc-400 hover:text-zinc-200 hover:bg-zinc-800/50'}"
      >
        <FileText size={14} class="mr-1.5" />
        List
      </button>
      <button
        onclick={() => switchView('graph')}
        class="flex items-center justify-center px-2.5 py-1 rounded-[4px] text-[11px] font-medium transition-all {$viewMode === 'graph' ? 'bg-zinc-700/50 text-zinc-100 shadow-sm' : 'text-zinc-400 hover:text-zinc-200 hover:bg-zinc-800/50'}"
      >
        <Network size={14} class="mr-1.5" />
        Graph
      </button>
    </div>

    <!-- Workspace button -->
    <button
      onclick={openWorkspaceDialog}
      class="bg-zinc-100 hover:bg-white text-zinc-950 px-3 py-1.5 rounded-md text-[11px] font-semibold transition-all active:scale-[0.97] flex items-center space-x-1.5 shadow-sm"
    >
      <span>Open</span>
    </button>

    <!-- Theme toggler -->
    <button
      onclick={toggleTheme}
      class="p-1.5 bg-transparent hover:bg-zinc-800/80 rounded-md text-zinc-400 hover:text-zinc-200 transition-all active:scale-[0.97]"
      title="Toggle Theme"
    >
      {#if $isDark}
        <Sun size={15} />
      {:else}
        <Moon size={15} />
      {/if}
    </button>
  </div>
</header>
