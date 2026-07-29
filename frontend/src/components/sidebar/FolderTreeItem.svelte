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
    class="w-full text-left px-2 py-1.5 rounded-md text-[11px] font-medium flex items-center justify-between transition-all active:scale-[0.98] {$selectedFolderId === folder.id ? 'bg-white/[0.06] text-zinc-100 shadow-[inset_0_1px_1px_rgba(255,255,255,0.05)]' : 'text-zinc-400 hover:bg-white/[0.03] hover:text-zinc-200'}"
    style="padding-left: {0.5 + folder.level * 0.75}rem;"
  >
    <div class="flex items-center space-x-2 truncate">
      <span class="text-[12px]">{folder.icon || '📁'}</span>
      <span class="truncate">{folder.name}</span>
    </div>
    {#if folder.children && folder.children.length > 0}
      {#if isExpanded}
        <ChevronDown size={12} class="text-zinc-500" />
      {:else}
        <ChevronRight size={12} class="text-zinc-500" />
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
