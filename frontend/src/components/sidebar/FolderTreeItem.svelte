<script lang="ts">
  import { Folder as FolderIcon, FolderOpen, ChevronRight } from '@lucide/svelte'
  import { selectedFolderId, notes } from '$lib/stores/vaultStore'
  import type { Folder } from '$lib/types/domain'
  import FolderTreeItem from './FolderTreeItem.svelte'

  let { folder }: { folder: Folder } = $props()
  let isExpanded = $state(false)

  // Count notes inside this folder recursively
  const noteCount = $derived(
    $notes.filter(n => n.folder === folder.id || (n.folder && n.folder.startsWith(folder.id))).length
  )

  const isSelected = $derived($selectedFolderId === folder.id)

  function handleSelectFolder(e: MouseEvent) {
    e.stopPropagation()
    selectedFolderId.set(folder.id)
    if (folder.children && folder.children.length > 0) {
      isExpanded = !isExpanded
    }
  }

  function handleToggleChevron(e: MouseEvent) {
    e.stopPropagation()
    isExpanded = !isExpanded
  }
</script>

<div class="select-none">
  <div
    role="button"
    tabindex="0"
    onclick={handleSelectFolder}
    onkeydown={(e) => (e.key === 'Enter' || e.key === ' ') && handleSelectFolder(e as unknown as MouseEvent)}
    class="w-full text-left px-2 py-1.5 rounded-lg text-[12px] font-medium flex items-center justify-between group transition-all pressable cursor-pointer {isSelected ? 'bg-[var(--accent-soft)] text-[var(--accent-primary)] border border-[var(--accent-primary)]/20 shadow-xs font-semibold' : 'text-[var(--text-secondary)] hover:bg-[var(--bg-surface-hover)] hover:text-[var(--text-primary)] border border-transparent'}"
    style="padding-left: {0.4 + folder.level * 0.75}rem;"
  >
    <div class="flex items-center space-x-1.5 truncate flex-1 min-w-0">
      <!-- Chevron expand toggle with smooth rotation -->
      {#if folder.children && folder.children.length > 0}
        <button
          onclick={handleToggleChevron}
          class="p-0.5 rounded text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-transform duration-200 ease-[var(--ease-apple)] flex items-center justify-center {isExpanded ? 'rotate-90' : 'rotate-0'}"
          title={isExpanded ? 'Collapse' : 'Expand'}
        >
          <ChevronRight size={12} strokeWidth={2.2} />
        </button>
      {:else}
        <span class="w-3.5 inline-block"></span>
      {/if}

      <!-- Folder Icon -->
      <span class="flex-shrink-0 flex items-center justify-center">
        {#if isExpanded}
          <FolderOpen size={13.5} strokeWidth={2} class="text-[var(--accent-primary)]" />
        {:else}
          <FolderIcon size={13.5} strokeWidth={2} class={isSelected ? "text-[var(--accent-primary)]" : "text-[var(--text-muted)] group-hover:text-[var(--text-secondary)]"} />
        {/if}
      </span>

      <!-- Folder Name -->
      <span class="truncate tracking-tight">{folder.name}</span>
    </div>

    <!-- Note count pill badge -->
    {#if noteCount > 0}
      <span
        class="ml-1 px-1.5 py-0.2 rounded-full text-[9.5px] font-mono {isSelected ? 'bg-[var(--accent-primary)]/20 text-[var(--accent-primary)] font-bold' : 'bg-[var(--bg-surface-active)] text-[var(--text-muted)] group-hover:text-[var(--text-secondary)]'}"
      >
        {noteCount}
      </span>
    {/if}
  </div>

  <!-- Recursive Child Folders with clean tree guide indentation line -->
  {#if isExpanded && folder.children && folder.children.length > 0}
    <div
      class="mt-0.5 space-y-0.5 ml-2.5 pl-2 border-l border-[var(--border-subtle)] transition-all duration-200 ease-[var(--ease-apple)]"
    >
      {#each folder.children as child (child.id)}
        <FolderTreeItem folder={child} />
      {/each}
    </div>
  {/if}
</div>
