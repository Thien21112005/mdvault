<script lang="ts">
  import { Folder as FolderIcon, Plus, Tag, FolderPlus, FilePlus } from '@lucide/svelte'
  import { folders, allTags, selectedFolderId, selectedTag, createNewNote, createNewFolder } from '$lib/stores/vaultStore'
  import FolderTreeItem from './FolderTreeItem.svelte'

  let isCreatingFolder = $state(false)
  let newFolderName = $state('')
  let isCreatingNote = $state(false)
  let newNoteTitle = $state('')

  function handleCreateFolder() {
    if (!newFolderName.trim()) return
    createNewFolder(newFolderName.trim())
    newFolderName = ''
    isCreatingFolder = false
  }

  function handleCreateNote() {
    if (!newNoteTitle.trim()) return
    createNewNote(newNoteTitle.trim())
    newNoteTitle = ''
    isCreatingNote = false
  }

  function selectTag(t: string) {
    if ($selectedTag === t) selectedTag.set(null)
    else selectedTag.set(t)
  }
</script>

<aside class="w-64 bg-white/[0.01] border-r border-white/[0.06] text-zinc-400 flex flex-col h-[calc(100vh-3rem)] select-none overflow-y-auto">
  <!-- Actions Header -->
  <div class="px-4 py-3 flex items-center justify-between">
    <span class="text-[10px] font-semibold text-zinc-500 uppercase tracking-widest">Workspace</span>
    <div class="flex space-x-0.5">
      <button
        onclick={() => (isCreatingFolder = !isCreatingFolder)}
        class="p-1 text-zinc-500 hover:text-zinc-200 hover:bg-white/[0.05] rounded transition-all active:scale-95"
        title="New Folder"
      >
        <FolderPlus size={14} />
      </button>
      <button
        onclick={() => (isCreatingNote = !isCreatingNote)}
        class="p-1 text-zinc-500 hover:text-zinc-200 hover:bg-white/[0.05] rounded transition-all active:scale-95"
        title="New Note"
      >
        <FilePlus size={14} />
      </button>
    </div>
  </div>

  <!-- Dialogs -->
  {#if isCreatingFolder}
    <div class="px-4 py-2 bg-white/[0.02] border-y border-white/[0.04] flex items-center space-x-2">
      <input
        type="text"
        bind:value={newFolderName}
        placeholder="Folder name..."
        class="flex-1 bg-black/50 border border-white/[0.08] rounded px-2 py-1 text-[11px] text-zinc-200 focus:outline-none focus:border-white/[0.15] transition-all shadow-inner"
        onkeydown={(e) => e.key === 'Enter' && handleCreateFolder()}
      />
      <button onclick={handleCreateFolder} class="text-[11px] font-medium bg-zinc-200 hover:bg-white text-black px-2.5 py-1 rounded transition-colors shadow-sm">Save</button>
    </div>
  {/if}

  {#if isCreatingNote}
    <div class="px-4 py-2 bg-white/[0.02] border-y border-white/[0.04] flex items-center space-x-2">
      <input
        type="text"
        bind:value={newNoteTitle}
        placeholder="Note title..."
        class="flex-1 bg-black/50 border border-white/[0.08] rounded px-2 py-1 text-[11px] text-zinc-200 focus:outline-none focus:border-white/[0.15] transition-all shadow-inner"
        onkeydown={(e) => e.key === 'Enter' && handleCreateNote()}
      />
      <button onclick={handleCreateNote} class="text-[11px] font-medium bg-zinc-200 hover:bg-white text-black px-2.5 py-1 rounded transition-colors shadow-sm">Create</button>
    </div>
  {/if}

  <!-- Folders -->
  <div class="flex-1 px-3 py-1 overflow-y-auto">
    <button
      onclick={() => selectedFolderId.set(null)}
      class="w-full text-left px-2 py-1.5 rounded-md text-[11px] font-medium flex items-center space-x-2 transition-all active:scale-[0.98] {$selectedFolderId === null ? 'bg-white/[0.06] text-zinc-100 shadow-[inset_0_1px_1px_rgba(255,255,255,0.05)]' : 'text-zinc-400 hover:bg-white/[0.03] hover:text-zinc-200'}"
    >
      <FolderIcon size={14} class={$selectedFolderId === null ? "text-zinc-300" : "text-zinc-500"} />
      <span>All Notes</span>
    </button>

    <div class="mt-1 space-y-0.5">
      {#each $folders as folder (folder.id)}
        <FolderTreeItem {folder} />
      {/each}
    </div>
  </div>

  <!-- Tags Section -->
  {#if $allTags.length > 0}
    <div class="p-4 border-t border-white/[0.04]">
      <span class="text-[10px] font-semibold text-zinc-500 uppercase tracking-widest block mb-2.5">Tags ({$allTags.length})</span>
      <div class="flex flex-wrap gap-1.5 max-h-40 overflow-y-auto">
        {#each $allTags as tag}
          <button
            onclick={() => selectTag(tag)}
            class="px-2 py-0.5 rounded-[4px] text-[10px] font-medium flex items-center space-x-1 transition-all border {$selectedTag === tag ? 'bg-zinc-200 text-black border-zinc-200 shadow-sm' : 'bg-transparent text-zinc-400 border-white/[0.08] hover:bg-white/[0.04] hover:text-zinc-200'}"
          >
            <Tag size={10} />
            <span>{tag}</span>
          </button>
        {/each}
      </div>
    </div>
  {/if}
</aside>
