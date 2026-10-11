<script lang="ts">
  import {
    Folder as FolderIcon,
    FolderOpen,
    ChevronRight,
    MoreHorizontal,
    FilePlus,
    FolderPlus,
    Edit3,
    Trash2
  } from '@lucide/svelte'
  import {
    selectedFolderId,
    notes,
    deleteFolder,
    renameFolder,
    createNewNote,
    createNewFolder
  } from '$lib/stores/vaultStore'
  import { switchView, setEditorMode } from '$lib/stores/viewStore'
  import { toast } from '$lib/stores/toastStore'
  import type { Folder } from '$lib/types/domain'
  import ConfirmModal from '../ui/ConfirmModal.svelte'
  import FolderTreeItem from './FolderTreeItem.svelte'

  let { folder }: { folder: Folder } = $props()
  let isExpanded = $state(false)
  let isMenuOpen = $state(false)
  let menuEl = $state<HTMLDivElement | null>(null)

  // Modals state
  let isDeleteModalOpen = $state(false)
  let isRenameModalOpen = $state(false)
  let isNewNoteModalOpen = $state(false)
  let isNewSubfolderModalOpen = $state(false)

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

  function toggleMenu(e: MouseEvent) {
    e.stopPropagation()
    isMenuOpen = !isMenuOpen
  }

  function handleContextMenu(e: MouseEvent) {
    e.preventDefault()
    e.stopPropagation()
    selectedFolderId.set(folder.id)
    isMenuOpen = true
  }

  // Action handlers
  async function handleConfirmDelete() {
    isDeleteModalOpen = false
    const name = folder.name
    const success = await deleteFolder(folder.id)
    if (success) {
      toast.success('Folder Deleted', `Deleted "${name}" and its content`)
    } else {
      toast.error('Delete Failed', `Could not delete folder "${name}"`)
    }
  }

  async function handleConfirmRename(newName?: string) {
    if (!newName || !newName.trim()) return
    const cleanName = newName.trim()
    isRenameModalOpen = false
    const success = await renameFolder(folder.id, cleanName)
    if (success) {
      toast.success('Folder Renamed', `Renamed to "${cleanName}"`)
    } else {
      toast.error('Rename Failed', `Could not rename folder "${folder.name}"`)
    }
  }

  async function handleConfirmNewNote(title?: string) {
    const cleanTitle = (title && title.trim()) ? title.trim() : 'Untitled Note'
    isNewNoteModalOpen = false
    const success = await createNewNote(cleanTitle, folder.id)
    if (success) {
      isExpanded = true
      switchView('list')
      setEditorMode('edit')
      toast.success('Note Created', `Created "${cleanTitle}" in "${folder.name}"`)
    } else {
      toast.error('Create Failed', `Could not create note in "${folder.name}"`)
    }
  }

  async function handleConfirmNewSubfolder(name?: string) {
    if (!name || !name.trim()) return
    const cleanName = name.trim()
    isNewSubfolderModalOpen = false
    const success = await createNewFolder(cleanName, folder.id)
    if (success) {
      isExpanded = true
      toast.success('Folder Created', `Created "${cleanName}" inside "${folder.name}"`)
    } else {
      toast.error('Create Failed', `Could not create subfolder in "${folder.name}"`)
    }
  }
</script>

<svelte:window
  onclick={(e) => {
    if (isMenuOpen && menuEl && !menuEl.contains(e.target as Node)) {
      isMenuOpen = false
    }
  }}
  onkeydown={(e) => {
    if (e.key === 'Escape' && isMenuOpen) {
      isMenuOpen = false
    }
  }}
/>

<div class="select-none relative">
  <div
    role="button"
    tabindex="0"
    onclick={handleSelectFolder}
    oncontextmenu={handleContextMenu}
    onkeydown={(e) => (e.key === 'Enter' || e.key === ' ') && handleSelectFolder(e as unknown as MouseEvent)}
    class="w-full text-left px-2 py-1.5 rounded-lg text-[12px] font-medium flex items-center justify-between group transition-all pressable cursor-pointer relative {isSelected ? 'bg-[var(--accent-soft)] text-[var(--accent-primary)] border border-[var(--accent-primary)]/20 shadow-xs font-semibold' : 'text-[var(--text-secondary)] hover:bg-[var(--bg-surface-hover)] hover:text-[var(--text-primary)] border border-transparent'}"
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

    <!-- Right Controls: Note count pill and hover action button -->
    <div class="flex items-center space-x-1">
      {#if noteCount > 0}
        <span
          class="px-1.5 py-0.2 rounded-full text-[9.5px] font-mono {isSelected ? 'bg-[var(--accent-primary)]/20 text-[var(--accent-primary)] font-bold' : 'bg-[var(--bg-surface-active)] text-[var(--text-muted)] group-hover:text-[var(--text-secondary)]'}"
        >
          {noteCount}
        </span>
      {/if}

      <!-- More Actions Menu Trigger Button -->
      <button
        onclick={toggleMenu}
        class="opacity-0 group-hover:opacity-100 p-1 rounded-md text-[var(--text-muted)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-all pressable {isMenuOpen ? '!opacity-100 bg-[var(--bg-surface-hover)] text-[var(--text-primary)]' : ''}"
        title="Folder Options"
        aria-label="Folder Options"
      >
        <MoreHorizontal size={13} strokeWidth={2.2} />
      </button>
    </div>
  </div>

  <!-- Context Menu Popover -->
  {#if isMenuOpen}
    <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
    <div
      bind:this={menuEl}
      role="menu"
      tabindex="-1"
      class="absolute right-2 top-full mt-1 z-30 w-44 bg-[var(--bg-surface-elevated)]/95 backdrop-blur-2xl border border-[var(--border-subtle)] rounded-xl shadow-[var(--shadow-glass)] p-1 text-[11.5px] font-medium space-y-0.5 animate-in fade-in zoom-in-95 duration-100"
      onclick={(e) => e.stopPropagation()}
    >
      <button
        onclick={() => {
          isMenuOpen = false
          isNewNoteModalOpen = true
        }}
        class="w-full flex items-center space-x-2 px-2.5 py-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-colors pressable"
      >
        <FilePlus size={13} class="text-[var(--accent-primary)]" />
        <span>New Note Inside</span>
      </button>

      <button
        onclick={() => {
          isMenuOpen = false
          isNewSubfolderModalOpen = true
        }}
        class="w-full flex items-center space-x-2 px-2.5 py-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-colors pressable"
      >
        <FolderPlus size={13} class="text-[var(--accent-primary)]" />
        <span>New Subfolder</span>
      </button>

      <button
        onclick={() => {
          isMenuOpen = false
          isRenameModalOpen = true
        }}
        class="w-full flex items-center space-x-2 px-2.5 py-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-colors pressable"
      >
        <Edit3 size={13} />
        <span>Rename Folder</span>
      </button>

      <div class="h-px bg-[var(--border-subtle)] my-0.5"></div>

      <button
        onclick={() => {
          isMenuOpen = false
          isDeleteModalOpen = true
        }}
        class="w-full flex items-center space-x-2 px-2.5 py-1.5 rounded-lg text-rose-500 hover:bg-rose-500/10 transition-colors pressable"
      >
        <Trash2 size={13} />
        <span>Delete Folder</span>
      </button>
    </div>
  {/if}

  <!-- Modals -->
  <ConfirmModal
    open={isDeleteModalOpen}
    destructive={true}
    title="Delete Folder"
    description={`Are you sure you want to permanently delete folder "${folder.name}"${noteCount > 0 ? ` and its ${noteCount} note(s)` : ''}? All content inside will be deleted.`}
    confirmText="Delete Folder"
    onConfirm={handleConfirmDelete}
    onCancel={() => {
      isDeleteModalOpen = false
    }}
  />

  <ConfirmModal
    open={isRenameModalOpen}
    mode="prompt"
    title="Rename Folder"
    initialValue={folder.name}
    placeholder="Enter new folder name..."
    confirmText="Rename"
    onConfirm={handleConfirmRename}
    onCancel={() => {
      isRenameModalOpen = false
    }}
  />

  <ConfirmModal
    open={isNewNoteModalOpen}
    mode="prompt"
    title={`New Note in "${folder.name}"`}
    initialValue="Untitled Note"
    placeholder="Enter note title..."
    confirmText="Create Note"
    onConfirm={handleConfirmNewNote}
    onCancel={() => {
      isNewNoteModalOpen = false
    }}
  />

  <ConfirmModal
    open={isNewSubfolderModalOpen}
    mode="prompt"
    title={`New Subfolder in "${folder.name}"`}
    initialValue=""
    placeholder="Enter subfolder name..."
    confirmText="Create Folder"
    onConfirm={handleConfirmNewSubfolder}
    onCancel={() => {
      isNewSubfolderModalOpen = false
    }}
  />

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
