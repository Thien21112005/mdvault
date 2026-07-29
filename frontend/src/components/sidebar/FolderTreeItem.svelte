<script lang="ts">
  import { Folder as FolderIcon, ChevronRight, ChevronDown } from '@lucide/svelte'
  import { selectedFolderId } from '$lib/stores/vaultStore'
  import type { Folder } from '$lib/types/domain'
  import FolderTreeItem from './FolderTreeItem.svelte'

  let { folder }: { folder: Folder } = $props()
  let isExpanded = $state(false)

  function handleClick(e: MouseEvent) {
    e.stopPropagation()
    selectedFolderId.set(folder.id)
    if (folder.children && folder.children.length > 0) {
      isExpanded = !isExpanded
    }
  }
</script>

<div>
  <button
    onclick={handleClick}
    class="w-full text-left px-2 py-1.5 rounded text-xs font-medium flex items-center justify-between transition {$selectedFolderId === folder.id ? 'bg-slate-800 text-blue-400 font-bold' : 'text-slate-400 hover:bg-slate-800/40 hover:text-slate-200'}"
    style="padding-left: {0.5 + folder.level * 0.75}rem;"
  >
    <div class="flex items-center space-x-2 truncate">
      <span>{folder.icon || '📁'}</span>
      <span class="truncate">{folder.name}</span>
    </div>
    {#if folder.children && folder.children.length > 0}
      {#if isExpanded}
        <ChevronDown size={14} class="text-slate-500" />
      {:else}
        <ChevronRight size={14} class="text-slate-500" />
      {/if}
    {/if}
  </button>

  {#if isExpanded && folder.children}
    <div class="mt-0.5 space-y-0.5">
      {#each folder.children as child (child.id)}
        <FolderTreeItem folder={child} />
      {/each}
    </div>
  {/if}
</div>
