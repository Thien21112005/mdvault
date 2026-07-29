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

<aside class="w-64 bg-slate-900 border-r border-slate-800 text-slate-300 flex flex-col h-[calc(100vh-3.5rem)] select-none overflow-y-auto">
  <!-- Actions Header -->
  <div class="p-3 border-b border-slate-800 flex items-center justify-between">
    <span class="text-xs font-bold text-slate-400 uppercase tracking-wider">Workspace Tree</span>
    <div class="flex space-x-1">
      <button
        onclick={() => (isCreatingFolder = !isCreatingFolder)}
        class="p-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded transition"
        title="Tạo Thư mục mới"
      >
        <FolderPlus size={14} />
      </button>
      <button
        onclick={() => (isCreatingNote = !isCreatingNote)}
        class="p-1.5 bg-blue-600 hover:bg-blue-500 text-white rounded transition"
        title="Tạo Ghi chú mới"
      >
        <FilePlus size={14} />
      </button>
    </div>
  </div>

  <!-- Dialogs -->
  {#if isCreatingFolder}
    <div class="p-2 bg-slate-800/50 border-b border-slate-700 flex items-center space-x-2">
      <input
        type="text"
        bind:value={newFolderName}
        placeholder="Tên thư mục..."
        class="flex-1 bg-slate-900 border border-slate-700 rounded px-2 py-1 text-xs text-slate-200 focus:outline-none focus:border-blue-500"
        onkeydown={(e) => e.key === 'Enter' && handleCreateFolder()}
      />
      <button onclick={handleCreateFolder} class="text-xs bg-emerald-600 hover:bg-emerald-500 text-white px-2 py-1 rounded">Lưu</button>
    </div>
  {/if}

  {#if isCreatingNote}
    <div class="p-2 bg-slate-800/50 border-b border-slate-700 flex items-center space-x-2">
      <input
        type="text"
        bind:value={newNoteTitle}
        placeholder="Tiêu đề note..."
        class="flex-1 bg-slate-900 border border-slate-700 rounded px-2 py-1 text-xs text-slate-200 focus:outline-none focus:border-blue-500"
        onkeydown={(e) => e.key === 'Enter' && handleCreateNote()}
      />
      <button onclick={handleCreateNote} class="text-xs bg-blue-600 hover:bg-blue-500 text-white px-2 py-1 rounded">Tạo</button>
    </div>
  {/if}

  <!-- Folders -->
  <div class="flex-1 p-2 overflow-y-auto">
    <button
      onclick={() => selectedFolderId.set(null)}
      class="w-full text-left px-2.5 py-1.5 rounded-lg text-xs font-semibold flex items-center space-x-2 transition {$selectedFolderId === null ? 'bg-slate-800 text-blue-400 font-bold' : 'text-slate-400 hover:bg-slate-800/50 hover:text-slate-200'}"
    >
      <FolderIcon size={16} />
      <span>Tất cả ghi chú</span>
    </button>

    <div class="mt-2 space-y-0.5">
      {#each $folders as folder (folder.id)}
        <FolderTreeItem {folder} />
      {/each}
    </div>
  </div>

  <!-- Tags Section -->
  {#if $allTags.length > 0}
    <div class="p-3 border-t border-slate-800">
      <span class="text-xs font-bold text-slate-400 uppercase tracking-wider block mb-2">Tags ({$allTags.length})</span>
      <div class="flex flex-wrap gap-1.5 max-h-40 overflow-y-auto">
        {#each $allTags as tag}
          <button
            onclick={() => selectTag(tag)}
            class="px-2 py-0.5 rounded-full text-xs font-medium flex items-center space-x-1 transition {$selectedTag === tag ? 'bg-blue-600 text-white shadow' : 'bg-slate-800 text-slate-300 hover:bg-slate-700'}"
          >
            <Tag size={10} />
            <span>#{tag}</span>
          </button>
        {/each}
      </div>
    </div>
  {/if}
</aside>
