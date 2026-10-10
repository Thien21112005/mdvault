<script lang="ts">
  import {
    Folder,
    FolderOpen,
    Search,
    Network,
    FileText,
    PanelLeft,
    Sun,
    Moon,
    X,
    Sparkles
  } from '@lucide/svelte'
  import { viewMode, isSidebarOpen, toggleSidebar, switchView, openCommandPalette } from '$lib/stores/viewStore'
  import { searchQuery, openWorkspaceDialog, rootPath } from '$lib/stores/vaultStore'
  import { isDark, toggleTheme } from '$lib/stores/themeStore'

  function getVaultName(path: string): string {
    if (!path) return 'No Vault Open'
    const clean = path.replace(/\\/g, '/')
    const parts = clean.split('/').filter(Boolean)
    return parts.pop() || path
  }
</script>

<header class="h-13 bg-[var(--bg-surface)] backdrop-blur-2xl border-b border-[var(--border-subtle)] text-[var(--text-primary)] flex items-center justify-between px-3.5 select-none relative z-50 transition-colors duration-200 shadow-[var(--shadow-sm)]">
  <!-- Left Section: Navigation Controls & Vault Info -->
  <div class="flex items-center space-x-2.5 min-w-[240px]">
    <!-- Toggle Sidebar Button -->
    <button
      onclick={toggleSidebar}
      class="p-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] border border-transparent hover:border-[var(--border-subtle)] transition-all pressable flex items-center justify-center {$isSidebarOpen ? 'text-[var(--accent-primary)] bg-[var(--accent-soft)]' : ''}"
      title={$isSidebarOpen ? 'Hide Sidebar (Ctrl+B)' : 'Show Sidebar (Ctrl+B)'}
      aria-label="Toggle Sidebar"
    >
      <PanelLeft size={16} strokeWidth={2} />
    </button>

    <div class="h-4 w-px bg-[var(--border-subtle)]"></div>

    <!-- Brand Identity -->
    <div class="flex items-center space-x-2">
      <div class="flex items-center justify-center w-6 h-6 rounded-lg bg-gradient-to-br from-indigo-500/20 to-violet-600/30 text-[var(--accent-primary)] border border-indigo-500/30 shadow-[0_0_12px_rgba(99,102,241,0.2)]">
        <Sparkles size={13} strokeWidth={2.2} />
      </div>
      <span class="font-bold text-[13.5px] tracking-tight bg-gradient-to-r from-[var(--text-primary)] to-[var(--text-secondary)] bg-clip-text text-transparent">
        MDVault
      </span>
    </div>

    <!-- Vault Name & Live Status Badge -->
    {#if $rootPath}
      <div
        class="flex items-center space-x-1.5 px-2 py-0.5 rounded-full text-[11px] font-medium bg-[var(--bg-surface-active)] border border-[var(--border-subtle)] max-w-[160px] transition-all cursor-default"
        title={`Active Vault: ${$rootPath}`}
      >
        <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.7)] animate-pulse"></span>
        <span class="truncate font-mono text-[10.5px] text-[var(--text-secondary)]">{getVaultName($rootPath)}</span>
      </div>
    {:else}
      <div class="flex items-center space-x-1.5 px-2 py-0.5 rounded-full text-[10.5px] font-medium bg-[var(--bg-surface-hover)] border border-[var(--border-subtle)] text-[var(--text-muted)] cursor-default">
        <span class="w-1.5 h-1.5 rounded-full bg-amber-500/60"></span>
        <span>No Vault</span>
      </div>
    {/if}
  </div>

  <!-- Center Section: Quick Search & Command Hint -->
  <div class="flex-1 max-w-md mx-4">
    <button
      onclick={openCommandPalette}
      class="w-full bg-[var(--bg-surface-active)] hover:bg-[var(--bg-surface-hover)] border border-[var(--border-subtle)] hover:border-[var(--border-strong)] rounded-lg pl-9 pr-3 py-1.5 text-[12.5px] text-[var(--text-muted)] text-left transition-all shadow-inner flex items-center justify-between cursor-pointer group relative pressable"
      title="Search or execute commands (Ctrl+K)"
    >
      <Search
        size={13.5}
        class="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--text-muted)] group-hover:text-[var(--accent-primary)] transition-colors pointer-events-none"
      />
      
      <span class="truncate">Search notes, commands... (Ctrl+K)</span>

      <kbd class="hidden sm:inline-flex items-center px-1.5 py-0.5 text-[9.5px] font-mono font-semibold text-[var(--text-muted)] bg-[var(--bg-surface-hover)] border border-[var(--border-subtle)] rounded shadow-xs pointer-events-none group-hover:border-[var(--border-strong)]">
        Ctrl+K
      </kbd>
    </button>
  </div>

  <!-- Right Section: View Mode, Theme Toggle & Vault Button -->
  <div class="flex items-center justify-end space-x-2 min-w-[240px]">
    <!-- View Switcher Segmented Control -->
    <div class="bg-[var(--bg-surface-active)] p-0.5 rounded-lg border border-[var(--border-subtle)] flex items-center shadow-inner">
      <button
        onclick={() => switchView('list')}
        class="flex items-center justify-center px-2.5 py-1 rounded-md text-[11px] font-medium transition-all pressable {$viewMode === 'list' ? 'bg-[var(--bg-surface-elevated)] text-[var(--text-primary)] shadow-sm border border-[var(--border-sheen)]' : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] border border-transparent'}"
        title="Notes List View (Ctrl+1)"
      >
        <FileText size={12.5} class="mr-1.5 {$viewMode === 'list' ? 'text-[var(--accent-primary)]' : ''}" />
        <span>Notes</span>
      </button>

      <button
        onclick={() => switchView('graph')}
        class="flex items-center justify-center px-2.5 py-1 rounded-md text-[11px] font-medium transition-all pressable {$viewMode === 'graph' ? 'bg-[var(--bg-surface-elevated)] text-[var(--text-primary)] shadow-sm border border-[var(--border-sheen)]' : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] border border-transparent'}"
        title="Knowledge Graph View (Ctrl+2)"
      >
        <Network size={12.5} class="mr-1.5 {$viewMode === 'graph' ? 'text-[var(--accent-primary)]' : ''}" />
        <span>Graph</span>
      </button>
    </div>

    <!-- Theme Switcher Toggle -->
    <button
      onclick={toggleTheme}
      class="p-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] border border-transparent hover:border-[var(--border-subtle)] transition-all pressable flex items-center justify-center"
      title={$isDark ? 'Switch to Light Theme' : 'Switch to Dark Theme'}
      aria-label="Toggle Theme"
    >
      {#if $isDark}
        <Sun size={15} strokeWidth={2} class="text-amber-400" />
      {:else}
        <Moon size={15} strokeWidth={2} class="text-indigo-600" />
      {/if}
    </button>

    <div class="h-4 w-px bg-[var(--border-subtle)]"></div>

    <!-- Open Vault Button -->
    <button
      onclick={openWorkspaceDialog}
      class="bg-gradient-to-r from-indigo-500 to-indigo-600 hover:from-indigo-600 hover:to-indigo-700 text-white px-3 py-1.5 rounded-lg text-[11.5px] font-semibold transition-all pressable flex items-center space-x-1.5 shadow-[0_2px_10px_rgba(99,102,241,0.35)] border border-indigo-400/30"
      title="Open Workspace Directory"
    >
      <FolderOpen size={13} strokeWidth={2} />
      <span>Open Vault</span>
    </button>
  </div>
</header>
