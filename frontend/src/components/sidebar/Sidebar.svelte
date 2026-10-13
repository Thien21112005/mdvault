<script lang="ts">
  import {
    Folder as FolderIcon,
    Tag,
    FolderPlus,
    FilePlus,
    Hash,
    Layers,
    X,
    FolderOpen,
    Check
  } from '@lucide/svelte'
  import {
    folders,
    notes,
    allTags,
    selectedFolderId,
    selectedTag,
    createNewNote,
    createNewFolder,
    rootPath
  } from '$lib/stores/vaultStore'
  import FolderTreeItem from './FolderTreeItem.svelte'
  import { toast } from '$lib/stores/toastStore'
  import { switchView, setEditorMode } from '$lib/stores/viewStore'


  let isCreatingFolder = $state(false)
  let newFolderName = $state('')
  let isCreatingNote = $state(false)
  let newNoteTitle = $state('')

  let folderInputEl = $state<HTMLInputElement | null>(null)
  let noteInputEl = $state<HTMLInputElement | null>(null)

  // Compute occurrences of each tag across notes
  const tagCountMap = $derived(
    $notes.reduce((acc, note) => {
      (note.tags || []).forEach(t => {
        acc[t] = (acc[t] || 0) + 1
      })
      return acc
    }, {} as Record<string, number>)
  )

  function openCreateFolder() {
    isCreatingNote = false
    isCreatingFolder = true
    newFolderName = ''
    setTimeout(() => folderInputEl?.focus(), 50)
  }

  function openCreateNote() {
    isCreatingFolder = false
    isCreatingNote = true
    newNoteTitle = ''
    setTimeout(() => noteInputEl?.focus(), 50)
  }

  async function handleCreateFolder() {
    if (!newFolderName.trim()) return
    const name = newFolderName.trim()
    const success = await createNewFolder(name)
    if (success) {
      toast.success('Folder Created', `Created folder "${name}"`)
    }
    newFolderName = ''
    isCreatingFolder = false
  }

  async function handleCreateNote() {
    if (!newNoteTitle.trim()) return
    const title = newNoteTitle.trim()
    const success = await createNewNote(title)
    if (success) {
      switchView('list')
      setEditorMode('edit')
      toast.success('Note Created', `Created note "${title}"`)
    }
    newNoteTitle = ''
    isCreatingNote = false
  }


  function selectTag(t: string) {
    if ($selectedTag === t) selectedTag.set(null)
    else selectedTag.set(t)
  }

  function getTargetFolderName(): string {
    if (!$selectedFolderId) return 'Root Vault'
    const clean = $selectedFolderId.replace(/\\/g, '/')
    return clean.split('/').pop() || 'Selected Folder'
  }
</script>

<aside class="w-64 bg-[var(--bg-surface)] backdrop-blur-xl border-r border-[var(--border-subtle)] text-[var(--text-secondary)] flex flex-col h-full select-none transition-colors duration-200">
  <!-- Explorer Header & Action Bar -->
  <div class="h-12 px-3.5 border-b border-[var(--border-subtle)] flex items-center justify-between shrink-0">
    <div class="flex items-center space-x-1.5">
      <span class="text-[10px] font-bold text-[var(--text-muted)] uppercase tracking-wider">Explorer</span>
      <span class="px-1.5 py-0.2 rounded-full text-[9px] font-mono font-semibold bg-[var(--bg-surface-active)] text-[var(--text-muted)]">
        {$notes.length}
      </span>
    </div>

    <!-- Quick Actions -->
    <div class="flex items-center space-x-1">
      <button
        onclick={openCreateFolder}
        class="p-1 rounded-md text-[var(--text-muted)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] border border-transparent hover:border-[var(--border-subtle)] transition-all pressable"
        title="New Folder"
        aria-label="New Folder"
      >
        <FolderPlus size={13.5} strokeWidth={2} />
      </button>

      <button
        onclick={openCreateNote}
        class="p-1 rounded-md text-[var(--text-muted)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] border border-transparent hover:border-[var(--border-subtle)] transition-all pressable"
        title="New Note"
        aria-label="New Note"
      >
        <FilePlus size={13.5} strokeWidth={2} />
      </button>
    </div>
  </div>

  <!-- Inline Origin-Aware Creation Card for Folder -->
  {#if isCreatingFolder}
    <div class="m-2 p-2.5 rounded-xl bg-[var(--bg-surface-elevated)] border border-[var(--border-sheen)] shadow-[var(--shadow-glass)] animate-in fade-in zoom-in-95 duration-150">
      <div class="flex items-center justify-between mb-1.5">
        <span class="text-[10px] font-semibold text-[var(--text-muted)] uppercase tracking-wider flex items-center space-x-1">
          <FolderPlus size={11} class="text-[var(--accent-primary)]" />
          <span>New Folder</span>
        </span>
        <button
          onclick={() => (isCreatingFolder = false)}
          class="p-0.5 rounded text-[var(--text-muted)] hover:text-[var(--text-primary)]"
        >
          <X size={11} />
        </button>
      </div>

      <p class="text-[10px] text-[var(--text-muted)] mb-2 truncate">
        Inside: <span class="font-mono text-[var(--text-secondary)]">{getTargetFolderName()}</span>
      </p>

      <div class="flex space-x-1.5">
        <input
          bind:this={folderInputEl}
          type="text"
          bind:value={newFolderName}
          placeholder="Folder name..."
          class="flex-1 bg-[var(--bg-surface-active)] border border-[var(--border-subtle)] rounded-md px-2 py-1 text-[11px] text-[var(--text-primary)] placeholder-[var(--text-muted)] focus:outline-none focus:border-[var(--accent-primary)] transition-all"
          onkeydown={(e) => {
            if (e.key === 'Enter') handleCreateFolder()
            if (e.key === 'Escape') isCreatingFolder = false
          }}
        />
        <button
          onclick={handleCreateFolder}
          disabled={!newFolderName.trim()}
          class="bg-[var(--accent-primary)] hover:opacity-90 disabled:opacity-40 text-white px-2.5 py-1 rounded-md text-[11px] font-semibold transition-all pressable flex items-center space-x-1"
        >
          <Check size={11} />
        </button>
      </div>
    </div>
  {/if}

  <!-- Inline Origin-Aware Creation Card for Note -->
  {#if isCreatingNote}
    <div class="m-2 p-2.5 rounded-xl bg-[var(--bg-surface-elevated)] border border-[var(--border-sheen)] shadow-[var(--shadow-glass)] animate-in fade-in zoom-in-95 duration-150">
      <div class="flex items-center justify-between mb-1.5">
        <span class="text-[10px] font-semibold text-[var(--text-muted)] uppercase tracking-wider flex items-center space-x-1">
          <FilePlus size={11} class="text-[var(--accent-primary)]" />
          <span>New Note</span>
        </span>
        <button
          onclick={() => (isCreatingNote = false)}
          class="p-0.5 rounded text-[var(--text-muted)] hover:text-[var(--text-primary)]"
        >
          <X size={11} />
        </button>
      </div>

      <p class="text-[10px] text-[var(--text-muted)] mb-2 truncate">
        In folder: <span class="font-mono text-[var(--text-secondary)]">{getTargetFolderName()}</span>
      </p>

      <div class="flex space-x-1.5">
        <input
          bind:this={noteInputEl}
          type="text"
          bind:value={newNoteTitle}
          placeholder="Note title..."
          class="flex-1 bg-[var(--bg-surface-active)] border border-[var(--border-subtle)] rounded-md px-2 py-1 text-[11px] text-[var(--text-primary)] placeholder-[var(--text-muted)] focus:outline-none focus:border-[var(--accent-primary)] transition-all"
          onkeydown={(e) => {
            if (e.key === 'Enter') handleCreateNote()
            if (e.key === 'Escape') isCreatingNote = false
          }}
        />
        <button
          onclick={handleCreateNote}
          disabled={!newNoteTitle.trim()}
          class="bg-[var(--accent-primary)] hover:opacity-90 disabled:opacity-40 text-white px-2.5 py-1 rounded-md text-[11px] font-semibold transition-all pressable flex items-center space-x-1"
        >
          <Check size={11} />
        </button>
      </div>
    </div>
  {/if}

  <!-- Tree Navigation View -->
  <div class="flex-1 px-2.5 py-1.5 overflow-y-auto space-y-0.5">
    <!-- All Notes Shortcut Button -->
    <button
      onclick={() => selectedFolderId.set(null)}
      class="w-full text-left px-2 py-1.5 rounded-lg text-[12px] font-medium flex items-center justify-between transition-all pressable {$selectedFolderId === null ? 'bg-[var(--accent-soft)] text-[var(--accent-primary)] border border-[var(--accent-primary)]/20 shadow-xs font-semibold' : 'text-[var(--text-secondary)] hover:bg-[var(--bg-surface-hover)] hover:text-[var(--text-primary)] border border-transparent'}"
    >
      <div class="flex items-center space-x-2 truncate">
        <Layers size={13.5} strokeWidth={2} class={$selectedFolderId === null ? 'text-[var(--accent-primary)]' : 'text-[var(--text-muted)]'} />
        <span>All Notes</span>
      </div>
      <span class="px-1.5 py-0.2 rounded-full text-[9.5px] font-mono {$selectedFolderId === null ? 'bg-[var(--accent-primary)]/20 text-[var(--accent-primary)] font-bold' : 'bg-[var(--bg-surface-active)] text-[var(--text-muted)]'}">
        {$notes.length}
      </span>
    </button>

    <!-- Folder Hierarchy -->
    <div class="mt-1 space-y-0.5">
      {#each $folders as folder (folder.id)}
        <FolderTreeItem {folder} />
      {/each}
    </div>

    <!-- Empty State when no folders exist -->
    {#if $folders.length === 0 && !$rootPath}
      <div class="py-12 px-3 text-center text-[var(--text-muted)] space-y-2">
        <FolderOpen size={24} strokeWidth={1.5} class="mx-auto opacity-40 text-[var(--text-muted)]" />
        <p class="text-[11px] font-medium leading-relaxed">No folder open.</p>
        <p class="text-[10px] opacity-75">Click "Open Vault" above to select a workspace.</p>
      </div>
    {/if}
  </div>

  <!-- Tags Section Cloud -->
  {#if $allTags.length > 0}
    <div class="p-3 border-t border-[var(--border-subtle)] bg-[var(--bg-surface)]">
      <div class="flex items-center justify-between mb-2">
        <div class="flex items-center space-x-1.5">
          <Tag size={11} strokeWidth={2.2} class="text-[var(--text-muted)]" />
          <span class="text-[10px] font-bold text-[var(--text-muted)] uppercase tracking-wider">Tags ({$allTags.length})</span>
        </div>

        {#if $selectedTag}
          <button
            onclick={() => selectedTag.set(null)}
            class="text-[9.5px] font-semibold text-[var(--accent-primary)] hover:underline flex items-center space-x-0.5"
            title="Clear tag filter"
          >
            <span>Clear</span>
            <X size={10} />
          </button>
        {/if}
      </div>

      <div class="flex flex-wrap gap-1 max-h-36 overflow-y-auto pr-0.5">
        {#each $allTags as tag}
          {@const count = tagCountMap[tag] || 0}
          <button
            onclick={() => selectTag(tag)}
            class="px-2 py-0.5 rounded-md text-[10px] font-medium flex items-center space-x-1 transition-all pressable border {$selectedTag === tag ? 'bg-amber-500/15 text-amber-500 border-amber-500/30 shadow-xs font-semibold' : 'bg-[var(--bg-surface-active)] text-[var(--text-secondary)] border-[var(--border-subtle)] hover:bg-[var(--bg-surface-hover)] hover:text-[var(--text-primary)]'}"
          >
            <Hash size={9.5} class={$selectedTag === tag ? 'text-amber-500' : 'text-[var(--text-muted)]'} />
            <span class="truncate max-w-[90px]">{tag}</span>
            {#if count > 0}
              <span class="text-[8.5px] opacity-60 font-mono">({count})</span>
            {/if}
          </button>
        {/each}
      </div>
    </div>
  {/if}
</aside>
