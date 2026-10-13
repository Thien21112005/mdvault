<script lang="ts">
  import {
    FileText,
    Calendar,
    Tag as TagIcon,
    ArrowDownUp,
    Hash,
    X,
    Plus,
    SearchX,
    Clock,
    Trash2
  } from '@lucide/svelte'
  import {
    filteredNotes,
    selectedNoteId,
    sortOption,
    selectedTag,
    searchQuery,
    createNewNote,
    deleteNote,
    type SortOption
  } from '$lib/stores/vaultStore'
  import type { Note } from '$lib/types/domain'
  import { viewMode } from '$lib/stores/viewStore'
  import CustomDropdown, { type DropdownOption } from '../ui/CustomDropdown.svelte'
  import ConfirmModal from '../ui/ConfirmModal.svelte'
  import NoteIcon from '../ui/NoteIcon.svelte'
  import { highlightMatch } from '$lib/utils/markdown'
  import { toast } from '$lib/stores/toastStore'

  const sortOptions: DropdownOption<SortOption>[] = [
    { value: 'newest', label: 'Newest', icon: Clock },
    { value: 'oldest', label: 'Oldest', icon: Calendar },
    { value: 'title-asc', label: 'Title (A-Z)', icon: FileText },
    { value: 'title-desc', label: 'Title (Z-A)', icon: FileText }
  ]


  function handleSelectNote(id: string) {
    selectedNoteId.set(id)
    viewMode.set('list')
  }

  function handleQuickNewNote() {
    createNewNote('Untitled Note')
  }

  let noteToDelete = $state<Note | null>(null)

  function clearFilters() {
    selectedTag.set(null)
    searchQuery.set('')
  }

  function formatDate(isoString: string): string {
    if (!isoString) return ''
    try {
      const d = new Date(isoString)
      if (isNaN(d.getTime())) return isoString.substring(0, 10)
      
      const now = new Date()
      const diffMs = now.getTime() - d.getTime()
      const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

      if (diffDays === 0) return 'Today'
      if (diffDays === 1) return 'Yesterday'
      if (diffDays < 7) return `${diffDays}d ago`
      
      return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
    } catch {
      return isoString.substring(0, 10)
    }
  }
</script>

<div class="w-80 bg-[var(--bg-surface)] backdrop-blur-xl border-r border-[var(--border-subtle)] flex flex-col h-full select-none transition-colors duration-200">
  <!-- Header Bar with Count & Sort Control -->
  <div class="h-12 px-3.5 border-b border-[var(--border-subtle)] flex items-center justify-between sticky top-0 z-10 bg-[var(--bg-surface)]/90 backdrop-blur-md shrink-0">
    <div class="flex items-center space-x-1.5">
      <span class="text-[10px] font-bold text-[var(--text-muted)] uppercase tracking-wider">Notes</span>
      <span class="px-1.5 py-0.2 rounded-full text-[9px] font-mono font-semibold bg-[var(--bg-surface-active)] text-[var(--text-muted)]">
        {$filteredNotes.length}
      </span>
    </div>

    <!-- Custom Glassmorphic Sort Dropdown -->
    <CustomDropdown
      options={sortOptions}
      bind:value={$sortOption}
      align="right"
    />
  </div>

  <!-- Active Filter Notification Pill -->
  {#if $selectedTag || $searchQuery}
    <div class="px-3 py-1.5 bg-[var(--accent-soft)] border-b border-[var(--border-subtle)] flex items-center justify-between text-[10.5px]">
      <div class="flex items-center space-x-1.5 truncate text-[var(--accent-primary)] font-medium">
        {#if $selectedTag}
          <span class="truncate">Tag: #{ $selectedTag }</span>
        {/if}
        {#if $searchQuery}
          <span class="truncate">Search: "{ $searchQuery }"</span>
        {/if}
      </div>
      <button
        onclick={clearFilters}
        class="p-0.5 rounded text-[var(--accent-primary)] hover:opacity-80 transition-opacity ml-1"
        title="Clear Filter"
      >
        <X size={11} />
      </button>
    </div>
  {/if}

  <!-- Scrollable Notes List -->
  <div class="flex-1 overflow-y-auto p-2 space-y-1.5">
    {#if $filteredNotes.length === 0}
      <!-- Empty State -->
      <div class="h-64 flex flex-col items-center justify-center text-center p-6 text-[var(--text-muted)] space-y-3">
        <div class="w-10 h-10 rounded-xl bg-[var(--bg-surface-active)] border border-[var(--border-subtle)] flex items-center justify-center">
          <SearchX size={20} strokeWidth={1.75} class="opacity-60" />
        </div>
        <div>
          <p class="text-[12px] font-semibold text-[var(--text-secondary)]">No notes found</p>
          <p class="text-[10.5px] text-[var(--text-muted)] mt-0.5 leading-relaxed">
            {#if $searchQuery || $selectedTag}
              No documents match your active filters.
            {:else}
              Your workspace has no notes yet.
            {/if}
          </p>
        </div>
        {#if $searchQuery || $selectedTag}
          <button
            onclick={clearFilters}
            class="text-[11px] font-medium text-[var(--accent-primary)] hover:underline"
          >
            Clear active filters
          </button>
        {:else}
          <button
            onclick={handleQuickNewNote}
            class="flex items-center space-x-1 px-3 py-1.5 rounded-lg bg-[var(--accent-primary)] text-white text-[11px] font-semibold transition-all pressable shadow-xs"
          >
            <Plus size={12} strokeWidth={2.5} />
            <span>Create First Note</span>
          </button>
        {/if}
      </div>
    {:else}
      {#each $filteredNotes as note (note.id)}
        {@const isSelected = $selectedNoteId === note.id}
        <div
          role="button"
          tabindex="0"
          onclick={() => handleSelectNote(note.id)}
          onkeydown={(e) => (e.key === 'Enter' || e.key === ' ') && handleSelectNote(note.id)}
          class="w-full text-left p-3 rounded-xl transition-all duration-150 ease-[var(--ease-apple)] border relative overflow-hidden group cursor-pointer pressable {isSelected ? 'bg-[var(--accent-soft)] border-[var(--accent-primary)]/30 shadow-[var(--shadow-glass)]' : 'bg-[var(--bg-surface-active)] hover:bg-[var(--bg-surface-hover)] border-transparent hover:border-[var(--border-subtle)] hover:-translate-y-[1px]'}"
        >
          <!-- Glowing Left Accent Bar when Selected -->
          {#if isSelected}
            <div class="absolute left-0 top-2 bottom-2 w-[3px] bg-[var(--accent-primary)] rounded-r-full shadow-[0_0_10px_var(--accent-glow)]"></div>
          {/if}

          <!-- Note Title with Frontmatter Icon & Hover Quick Delete Action -->
          <div class="flex items-start justify-between">
            <div class="flex items-start space-x-2 flex-1 min-w-0">
              <span class="mt-0.5 p-0.5 rounded bg-[var(--bg-surface-elevated)] text-[var(--accent-primary)] shrink-0 flex items-center justify-center border border-[var(--border-subtle)]/60 shadow-xs">
                <NoteIcon icon={note.metadata?.icon} size={12} strokeWidth={2.2} />
              </span>
              <span class="font-semibold text-[13px] tracking-tight truncate flex-1 {isSelected ? 'text-[var(--text-primary)] font-bold' : 'text-[var(--text-secondary)] group-hover:text-[var(--text-primary)]'}">
                {@html highlightMatch(note.title || 'Untitled Note', $searchQuery)}
              </span>
            </div>
            <button
              onclick={(e) => {
                e.stopPropagation()
                noteToDelete = note
              }}
              class="opacity-0 group-hover:opacity-100 p-1 rounded-md text-[var(--text-muted)] hover:text-rose-500 hover:bg-rose-500/10 transition-all pressable ml-1.5 shrink-0"
              title="Delete note"
              aria-label="Delete note"
            >
              <Trash2 size={12} strokeWidth={2.2} />
            </button>
          </div>

          <!-- Note Excerpt with Search Highlight -->
          <p class="text-[11.5px] mt-1.5 line-clamp-2 leading-relaxed {isSelected ? 'text-[var(--text-secondary)]' : 'text-[var(--text-muted)] group-hover:text-[var(--text-secondary)]'}">
            {@html highlightMatch(note.excerpt || 'No content preview available...', $searchQuery)}
          </p>
          
          <!-- Metadata Footer (Date & Tags) -->
          <div class="mt-2.5 flex items-center justify-between text-[10.5px] font-medium text-[var(--text-muted)]">
            <div class="flex items-center space-x-1">
              <Clock size={10.5} class="opacity-75" />
              <span class="font-mono text-[10px]">{formatDate(note.lastModified)}</span>
            </div>

            <!-- Tags Mini Chips -->
            {#if note.tags && note.tags.length > 0}
              <div class="flex items-center space-x-1 max-w-[130px] overflow-hidden">
                {#each note.tags.slice(0, 2) as tag}
                  <span class="inline-flex items-center space-x-0.5 px-1.5 py-0.2 rounded text-[9.5px] font-mono {isSelected ? 'bg-[var(--accent-primary)]/15 text-[var(--accent-primary)] font-semibold' : 'bg-[var(--bg-surface-hover)] text-[var(--text-muted)]'}">
                    <span>#{tag}</span>
                  </span>
                {/each}
                {#if note.tags.length > 2}
                  <span class="text-[9px] text-[var(--text-muted)] font-mono">+{note.tags.length - 2}</span>
                {/if}
              </div>
            {/if}
          </div>
        </div>
      {/each}
    {/if}
  </div>
</div>

<!-- Delete Note Confirmation Modal -->
<ConfirmModal
  open={Boolean(noteToDelete)}
  destructive={true}
  title="Delete Note"
  description={`Are you sure you want to permanently delete "${noteToDelete?.title || 'this note'}"? This action cannot be undone.`}
  confirmText="Delete Note"
  onConfirm={async () => {
    if (noteToDelete) {
      const title = noteToDelete.title
      const ok = await deleteNote(noteToDelete.id)
      noteToDelete = null
      if (ok) {
        toast.success('Note Deleted', `"${title}" has been deleted`)
      } else {
        toast.error('Delete Failed', 'Could not delete note from disk')
      }
    }
  }}
  onCancel={() => {
    noteToDelete = null
  }}
/>
