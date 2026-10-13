<script lang="ts">
  import { Check, FileText } from '@lucide/svelte'
  import { rootPath, notes, selectedNote } from '$lib/stores/vaultStore'
  import { isEditorDirty } from '$lib/stores/viewStore'

  // Word count & Char count derived for active note
  const activeWords = $derived(
    $selectedNote?.content?.trim() ? $selectedNote.content.trim().split(/\s+/).length : 0
  )
  const activeChars = $derived(
    $selectedNote?.content ? $selectedNote.content.length : 0
  )

  // Clean directory display path
  function formatVaultPath(path: string): string {
    if (!path) return 'No active vault'
    const clean = path.replace(/\\/g, '/')
    const parts = clean.split('/')
    if (parts.length > 2) {
      return `.../${parts.slice(-2).join('/')}`
    }
    return clean
  }
</script>

<footer
  class="h-7 border-t border-[var(--border-subtle)] bg-[var(--bg-surface)]/90 backdrop-blur-xl px-3.5 flex items-center justify-between select-none text-[11px] font-mono text-[var(--text-muted)] shrink-0 z-20 transition-colors duration-200"
  role="status"
  aria-label="Application Status Bar"
>
  <!-- Left: Watcher & Workspace Path -->
  <div class="flex items-center space-x-2.5 truncate">
    {#if $rootPath}
      <div class="flex items-center space-x-1.5 shrink-0" title="Local fsnotify file watcher active & debounced">
        <span class="relative flex h-2 w-2">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
        </span>
        <span class="font-medium text-[var(--text-secondary)] text-[10.5px]">fsnotify active</span>
      </div>
      <span class="text-[var(--border-subtle)] shrink-0">•</span>
      <span
        class="truncate max-w-[200px] sm:max-w-[280px] text-[10.5px] text-[var(--text-muted)] hover:text-[var(--text-secondary)] transition-colors cursor-default"
        title={$rootPath}
      >
        {formatVaultPath($rootPath)}
      </span>
    {:else}
      <div class="flex items-center space-x-1.5 shrink-0" title="No active workspace directory selected">
        <span class="w-2 h-2 rounded-full bg-amber-500/80"></span>
        <span class="text-[10.5px] text-[var(--text-muted)]">No active vault</span>
      </div>
    {/if}
  </div>

  <!-- Center: Active Note Stats & Save State -->
  <div class="hidden sm:flex items-center space-x-2.5 text-[10.5px]">
    {#if $selectedNote}
      <div class="flex items-center space-x-1.5">
        {#if $isEditorDirty}
          <span class="w-1.5 h-1.5 rounded-full bg-amber-500 animate-pulse"></span>
          <span class="text-amber-500 font-semibold tracking-tight">Unsaved</span>
        {:else}
          <Check size={11.5} class="text-emerald-500" strokeWidth={2.5} />
          <span class="text-[var(--text-muted)] font-medium">Saved</span>
        {/if}
      </div>
      <span class="text-[var(--border-subtle)]">•</span>
      <span class="tracking-tight">{activeWords} words</span>
      <span class="text-[var(--border-subtle)]">•</span>
      <span class="tracking-tight">{activeChars} chars</span>
    {:else}
      <span class="text-[10.5px] text-[var(--text-muted)] opacity-60">No note selected</span>
    {/if}
  </div>

  <!-- Right: Total Notes & System KPI -->
  <div class="flex items-center space-x-2.5 shrink-0 text-[10.5px]">
    <div class="flex items-center space-x-1.5" title="Total notes in current workspace">
      <FileText size={11.5} class="text-[var(--text-muted)]" />
      <span class="font-semibold text-[var(--text-secondary)]">{$notes.length}</span>
      <span class="hidden md:inline">{$notes.length === 1 ? 'note' : 'notes'}</span>
    </div>

    <span class="text-[var(--border-subtle)] hidden lg:inline">•</span>
    <span class="hidden lg:inline text-[var(--text-muted)] opacity-75" title="Local-First Markdown Storage Engine">
      Local-First Vault
    </span>

    <span class="text-[var(--border-subtle)]">•</span>
    <span
      class="px-1.5 py-0.2 rounded text-[9.5px] font-semibold bg-[var(--bg-surface-active)] text-[var(--accent-primary)] border border-[var(--border-subtle)]"
      title="Estimated Runtime RAM Footprint (<50MB target)"
    >
      ~38 MB
    </span>
  </div>
</footer>
