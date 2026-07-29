<script lang="ts">
  import { Folder, Search, Network, FileText, Moon, Sun } from '@lucide/svelte'
  import { viewMode } from '$lib/stores/viewStore'
  import { searchQuery, openWorkspaceDialog, rootPath } from '$lib/stores/vaultStore'
  import { isDark, toggleTheme } from '$lib/stores/themeStore'

  function switchView(mode: 'list' | 'graph') {
    viewMode.set(mode)
  }
</script>

<header class="h-14 bg-slate-900 border-b border-slate-800 text-white flex items-center justify-between px-4 select-none">
  <div class="flex items-center space-x-3">
    <div class="bg-blue-600 p-2 rounded-lg text-white shadow-md shadow-blue-500/20">
      <Folder size={20} />
    </div>
    <span class="font-bold text-lg tracking-wide bg-gradient-to-r from-blue-400 to-teal-300 bg-clip-text text-transparent">
      MDVault
    </span>
    {#if $rootPath}
      <span class="text-xs text-slate-400 bg-slate-800/80 px-2.5 py-1 rounded-md max-w-xs truncate" title={$rootPath}>
        📁 {$rootPath.split('/').pop() || $rootPath}
      </span>
    {/if}
  </div>

  <!-- Search Bar -->
  <div class="flex-1 max-w-md mx-6">
    <div class="relative">
      <Search size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
      <input
        type="text"
        bind:value={$searchQuery}
        placeholder="Tìm kiếm Ghi chú, Nội dung hoặc Tag..."
        class="w-full bg-slate-800/80 border border-slate-700 rounded-full pl-9 pr-4 py-1.5 text-sm text-slate-200 placeholder-slate-500 focus:outline-none focus:border-blue-500 transition"
      />
    </div>
  </div>

  <div class="flex items-center space-x-2">
    <!-- View Switcher -->
    <div class="bg-slate-800/80 p-1 rounded-lg border border-slate-700 flex space-x-1">
      <button
        onclick={() => switchView('list')}
        class="flex items-center space-x-1 px-2.5 py-1 rounded text-xs font-semibold transition {$viewMode === 'list' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:text-slate-200'}"
      >
        <FileText size={14} />
        <span>Ghi chú</span>
      </button>
      <button
        onclick={() => switchView('graph')}
        class="flex items-center space-x-1 px-2.5 py-1 rounded text-xs font-semibold transition {$viewMode === 'graph' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:text-slate-200'}"
      >
        <Network size={14} />
        <span>Graph View</span>
      </button>
    </div>

    <!-- Workspace button -->
    <button
      onclick={openWorkspaceDialog}
      class="bg-blue-600 hover:bg-blue-500 text-white px-3 py-1.5 rounded-lg text-xs font-semibold shadow transition flex items-center space-x-1"
    >
      <span>Mở Workspace</span>
    </button>

    <!-- Theme toggler -->
    <button
      onclick={toggleTheme}
      class="p-2 bg-slate-800 hover:bg-slate-700 rounded-lg text-slate-300 transition"
      title="Đoi giao dien"
    >
      {#if $isDark}
        <Sun size={18} class="text-yellow-400" />
      {:else}
        <Moon size={18} class="text-slate-300" />
      {/if}
    </button>
  </div>
</header>
