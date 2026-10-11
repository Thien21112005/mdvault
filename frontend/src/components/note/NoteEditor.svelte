<script lang="ts">
  import {
    Eye,
    Edit3,
    Columns2,
    Save,
    Check,
    Tag,
    FileText,
    Calendar,
    X,
    Bold,
    Italic,
    Code,
    Heading,
    List,
    ListTodo,
    Quote,
    Link,
    Minus,
    AlertCircle,
    Trash2
  } from '@lucide/svelte'
  import { selectedNote, saveActiveNote, notes, selectedNoteId, deleteNote, renameNote } from '$lib/stores/vaultStore'
  import { editorMode, setEditorMode, isEditorDirty } from '$lib/stores/viewStore'
  import { renderMarkdown } from '$lib/utils/markdown'
  import { toast } from '$lib/stores/toastStore'
  import ConfirmModal from '../ui/ConfirmModal.svelte'

  let editTitle = $state('')
  let editTagsList = $state<string[]>([])
  let newTagInput = $state('')
  let editContent = $state('')
  let isSavedToast = $state(false)
  let textareaEl = $state<HTMLTextAreaElement | null>(null)
  let previewContainerEl = $state<HTMLDivElement | null>(null)
  let splitPreviewContainerEl = $state<HTMLDivElement | null>(null)

  // Buffer protection & dirty tracking state
  let lastLoadedNoteId = $state<string | null>(null)
  let loadedTitle = $state('')
  let loadedContent = $state('')
  let loadedTagsJson = $state('[]')
  let isExternalModified = $state(false)

  // Derived dirty flag comparing active working copy against loaded baseline
  const isDirty = $derived(
    Boolean(
      $selectedNote && (
        editContent !== loadedContent ||
        editTitle !== loadedTitle ||
        JSON.stringify(editTagsList) !== loadedTagsJson
      )
    )
  )

  // Broadcast dirty state to viewStore for global UI awareness (e.g. StatusBar)
  $effect(() => {
    isEditorDirty.set(isDirty)
  })

  // Sync state whenever selected note changes, protecting unsaved edits from being wiped
  $effect(() => {
    const current = $selectedNote
    if (!current) {
      lastLoadedNoteId = null
      isExternalModified = false
      return
    }

    // Case 1: Switching to a different note entirely
    if (current.id !== lastLoadedNoteId) {
      editTitle = current.title || ''
      editTagsList = [...(current.tags || [])]
      editContent = current.content || ''
      loadedTitle = editTitle
      loadedContent = editContent
      loadedTagsJson = JSON.stringify(editTagsList)
      lastLoadedNoteId = current.id
      isExternalModified = false
      return
    }

    // Case 2: Same note refreshed by file watcher or background scan
    if (!isDirty) {
      // User has not typed anything unsaved, so we can seamlessly accept disk changes
      if (current.content !== editContent || current.title !== editTitle) {
        editTitle = current.title || ''
        editTagsList = [...(current.tags || [])]
        editContent = current.content || ''
        loadedTitle = editTitle
        loadedContent = editContent
        loadedTagsJson = JSON.stringify(editTagsList)
      }
      isExternalModified = false
    } else {
      // User HAS unsaved changes! Protect editContent and editTitle from being clobbered!
      if (current.content !== loadedContent) {
        isExternalModified = true
      }
    }
  })

  // Word count calculation
  const wordsCount = $derived(
    (editContent || '').trim() ? (editContent || '').trim().split(/\s+/).length : 0
  )

  function handleSave() {
    if (!$selectedNote) return
    // Auto-commit any non-empty tag still typed in the input
    if (newTagInput.trim()) {
      const clean = newTagInput.trim().replace(/^#+/, '')
      if (clean && !editTagsList.some(t => t.toLowerCase() === clean.toLowerCase())) {
        editTagsList = [...editTagsList, clean]
        newTagInput = ''
      }
    }

    saveActiveNote($selectedNote.id, editContent, editTitle, editTagsList)
    loadedTitle = editTitle
    loadedContent = editContent
    loadedTagsJson = JSON.stringify(editTagsList)
    isExternalModified = false
    toast.success('Note Saved', 'All changes stored safely')
    isSavedToast = true
    setTimeout(() => {
      isSavedToast = false
    }, 1800)
  }

  function handleKeydown(e: KeyboardEvent) {
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's') {
      e.preventDefault()
      handleSave()
    }
  }

  // Note CRUD modals state & handlers
  let isDeleteModalOpen = $state(false)
  let isRenameModalOpen = $state(false)

  function openDeleteModal() {
    if (!$selectedNote) return
    isDeleteModalOpen = true
  }

  async function handleConfirmDelete() {
    if (!$selectedNote) return
    const title = $selectedNote.title
    const ok = await deleteNote($selectedNote.id)
    isDeleteModalOpen = false
    if (ok) {
      toast.success('Note Deleted', `"${title}" has been deleted`)
    } else {
      toast.error('Delete Failed', 'Could not delete note from disk')
    }
  }

  function openRenameModal() {
    if (!$selectedNote) return
    isRenameModalOpen = true
  }

  async function handleConfirmRename(newTitle?: string) {
    if (!$selectedNote || !newTitle || !newTitle.trim()) return
    const cleanTitle = newTitle.trim()
    const updated = await renameNote($selectedNote.id, cleanTitle)
    isRenameModalOpen = false
    if (updated) {
      editTitle = updated.title
      toast.success('Note Renamed', `Renamed to "${updated.title}"`)
    } else {
      toast.error('Rename Failed', 'Could not rename note on disk')
    }
  }

  // Interactive tag operations
  function addTag() {
    const clean = newTagInput.trim().replace(/^#+/, '')
    if (clean && !editTagsList.some(t => t.toLowerCase() === clean.toLowerCase())) {
      editTagsList = [...editTagsList, clean]
    }
    newTagInput = ''
  }

  function removeTag(idx: number) {
    editTagsList = editTagsList.filter((_, i) => i !== idx)
  }

  function handleTagKeyDown(e: KeyboardEvent) {
    if (e.key === 'Enter' || e.key === ',') {
      e.preventDefault()
      addTag()
    } else if (e.key === 'Backspace' && newTagInput === '' && editTagsList.length > 0) {
      removeTag(editTagsList.length - 1)
    }
  }

  // Formatting helpers with cursor selection preservation
  function wrapSelection(prefix: string, suffix: string, defaultText: string) {
    if (!textareaEl) return
    const start = textareaEl.selectionStart
    const end = textareaEl.selectionEnd
    const text = textareaEl.value

    if (start !== end) {
      const selected = text.substring(start, end)
      const replacement = `${prefix}${selected}${suffix}`
      editContent = text.substring(0, start) + replacement + text.substring(end)
      requestAnimationFrame(() => {
        if (textareaEl) {
          textareaEl.focus()
          textareaEl.setSelectionRange(start + prefix.length, end + prefix.length)
        }
      })
    } else {
      const replacement = `${prefix}${defaultText}${suffix}`
      editContent = text.substring(0, start) + replacement + text.substring(end)
      requestAnimationFrame(() => {
        if (textareaEl) {
          textareaEl.focus()
          textareaEl.setSelectionRange(start + prefix.length, start + prefix.length + defaultText.length)
        }
      })
    }
  }

  function insertLinePrefix(prefix: string) {
    if (!textareaEl) return
    const start = textareaEl.selectionStart
    const text = textareaEl.value
    const lineStart = text.lastIndexOf('\n', start - 1) + 1
    editContent = text.substring(0, lineStart) + prefix + text.substring(lineStart)
    requestAnimationFrame(() => {
      if (textareaEl) {
        textareaEl.focus()
        const newPos = start + prefix.length
        textareaEl.setSelectionRange(newPos, newPos)
      }
    })
  }

  function insertDivider(textToInsert: string) {
    if (!textareaEl) return
    const start = textareaEl.selectionStart
    const text = textareaEl.value
    editContent = text.substring(0, start) + textToInsert + text.substring(start)
    requestAnimationFrame(() => {
      if (textareaEl) {
        textareaEl.focus()
        const newPos = start + textToInsert.length
        textareaEl.setSelectionRange(newPos, newPos)
      }
    })
  }

  // Event delegation for both Preview and Split views (code copy & wiki-link)
  function handlePreviewClick(e: MouseEvent) {
    const target = e.target as HTMLElement

    // 1. Copy code button
    const copyBtn = target.closest('.copy-code-btn') as HTMLButtonElement | null
    if (copyBtn) {
      const rawCode = decodeURIComponent(copyBtn.dataset.code || '')
      if (rawCode) {
        navigator.clipboard.writeText(rawCode)
        copyBtn.innerText = 'Copied!'
        copyBtn.classList.add('text-emerald-500', 'font-bold')
        setTimeout(() => {
          copyBtn.innerText = 'Copy'
          copyBtn.classList.remove('text-emerald-500', 'font-bold')
        }, 1500)
      }
      return
    }

    // 2. Wiki-link click navigation
    const wikiLink = target.closest('.wiki-link') as HTMLAnchorElement | null
    if (wikiLink) {
      e.preventDefault()
      const targetTitle = wikiLink.dataset.wiki || ''
      if (targetTitle) {
        const match = $notes.find(n => n.title.toLowerCase() === targetTitle.toLowerCase())
        if (match) {
          selectedNoteId.set(match.id)
        }
      }
    }
  }

  $effect(() => {
    const p1 = previewContainerEl
    const p2 = splitPreviewContainerEl
    if (p1) p1.addEventListener('click', handlePreviewClick)
    if (p2) p2.addEventListener('click', handlePreviewClick)
    return () => {
      if (p1) p1.removeEventListener('click', handlePreviewClick)
      if (p2) p2.removeEventListener('click', handlePreviewClick)
    }
  })
</script>

<svelte:window onkeydown={handleKeydown} />

<!-- Reusable Tag Chips Manager Snippet -->
{#snippet tagChipsManager()}
  <div class="flex items-center flex-wrap gap-1.5 py-1">
    <Tag size={13} class="text-[var(--text-muted)] mr-1 shrink-0" />
    {#each editTagsList as tag, i}
      <span class="inline-flex items-center space-x-1 bg-[var(--bg-surface-elevated)] border border-[var(--border-subtle)] text-[var(--text-secondary)] px-2 py-0.5 rounded-md text-[11px] font-mono group transition-colors hover:border-[var(--border-strong)]">
        <span class="text-[var(--accent-primary)] font-medium">#</span>
        <span>{tag}</span>
        <button
          type="button"
          onclick={() => removeTag(i)}
          class="text-[var(--text-muted)] hover:text-rose-400 p-0.5 rounded transition-colors pressable"
          title="Remove tag"
          aria-label="Remove tag {tag}"
        >
          <X size={10} />
        </button>
      </span>
    {/each}
    <div class="inline-flex items-center">
      <input
        type="text"
        bind:value={newTagInput}
        onkeydown={handleTagKeyDown}
        placeholder={editTagsList.length === 0 ? "Add tags (press Enter)..." : "+ tag..."}
        class="bg-transparent border-none text-[11px] font-mono text-[var(--text-secondary)] focus:outline-none focus:ring-0 placeholder-[var(--text-muted)] min-w-[80px] px-1 py-0.5"
      />
    </div>
  </div>
{/snippet}

<!-- Reusable Markdown Format Toolbar Snippet -->
{#snippet formatToolbar()}
  <div class="flex items-center space-x-0.5 p-1 bg-[var(--bg-surface-elevated)]/80 backdrop-blur-md border border-[var(--border-subtle)] rounded-xl shadow-[var(--shadow-glass)] overflow-x-auto select-none">
    <button
      type="button"
      onclick={() => wrapSelection('**', '**', 'bold text')}
      class="p-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-all pressable"
      title="Bold (**text**)"
    >
      <Bold size={13} strokeWidth={2.2} />
    </button>
    <button
      type="button"
      onclick={() => wrapSelection('*', '*', 'italic text')}
      class="p-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-all pressable"
      title="Italic (*text*)"
    >
      <Italic size={13} strokeWidth={2.2} />
    </button>
    <button
      type="button"
      onclick={() => wrapSelection('`', '`', 'code')}
      class="p-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-all pressable"
      title="Inline Code (`code`)"
    >
      <Code size={13} strokeWidth={2.2} />
    </button>

    <div class="h-3.5 w-px bg-[var(--border-subtle)] mx-1"></div>

    <button
      type="button"
      onclick={() => insertLinePrefix('### ')}
      class="p-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-all pressable"
      title="Heading (### title)"
    >
      <Heading size={13} strokeWidth={2.2} />
    </button>
    <button
      type="button"
      onclick={() => insertLinePrefix('- ')}
      class="p-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-all pressable"
      title="Bullet List (- item)"
    >
      <List size={13} strokeWidth={2.2} />
    </button>
    <button
      type="button"
      onclick={() => insertLinePrefix('- [ ] ')}
      class="p-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-all pressable"
      title="Task List (- [ ] task)"
    >
      <ListTodo size={13} strokeWidth={2.2} />
    </button>
    <button
      type="button"
      onclick={() => insertLinePrefix('> ')}
      class="p-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-all pressable"
      title="Quote (> quote)"
    >
      <Quote size={13} strokeWidth={2.2} />
    </button>

    <div class="h-3.5 w-px bg-[var(--border-subtle)] mx-1"></div>

    <button
      type="button"
      onclick={() => wrapSelection('[[', ']]', 'Page Title')}
      class="p-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-all pressable"
      title="Wiki Link ([[page]])"
    >
      <Link size={13} strokeWidth={2.2} />
    </button>
    <button
      type="button"
      onclick={() => insertDivider('\n\n---\n\n')}
      class="p-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-all pressable"
      title="Horizontal Divider (---)"
    >
      <Minus size={13} strokeWidth={2.2} />
    </button>
  </div>
{/snippet}

<div class="flex-1 bg-transparent flex flex-col h-[calc(100vh-3.25rem)] overflow-hidden relative">
  {#if !$selectedNote}
    <div class="flex-1 flex flex-col items-center justify-center text-[var(--text-muted)] space-y-4">
      <div class="p-5 rounded-3xl bg-[var(--bg-surface-elevated)] border border-[var(--border-subtle)] shadow-[var(--shadow-glass)]">
        <FileText size={32} strokeWidth={1.5} class="text-[var(--text-muted)]" />
      </div>
      <p class="text-[13px] font-medium tracking-wide">Select a note to view or edit</p>
    </div>
  {:else}
    <!-- Header bar -->
    <div class="h-12 border-b border-[var(--border-subtle)] px-6 flex items-center justify-between select-none bg-[var(--bg-surface)]/80 backdrop-blur-md sticky top-0 z-10 shrink-0">
      <div class="flex items-center space-x-3 text-[11px] font-medium text-[var(--text-muted)]">
        <span class="flex items-center space-x-1.5 font-mono">
          <Calendar size={12} />
          <span>Modified { $selectedNote.lastModified ? $selectedNote.lastModified.substring(0, 10) : '' }</span>
        </span>
        <span class="text-[var(--border-subtle)]">•</span>
        <span class="font-mono text-[10.5px]">
          {wordsCount} words
        </span>
        {#if isDirty}
          <span class="text-[var(--border-subtle)]">•</span>
          <span class="text-amber-500 font-mono text-[10px] flex items-center space-x-1 font-semibold animate-in fade-in duration-150">
            <span class="w-1.5 h-1.5 rounded-full bg-amber-400 animate-pulse"></span>
            <span>Unsaved</span>
          </span>
        {/if}
      </div>

      <div class="flex items-center space-x-3">
        <!-- Segmented Control for Preview / Split / Edit -->
        <div class="bg-[var(--bg-surface-active)] p-0.5 rounded-lg border border-[var(--border-subtle)] flex items-center shadow-inner">
          <button
            onclick={() => setEditorMode('preview')}
            class="flex items-center justify-center px-2.5 py-1 rounded-md text-[11px] font-medium transition-all pressable {$editorMode === 'preview' ? 'bg-[var(--bg-surface-elevated)] text-[var(--text-primary)] shadow-sm border border-[var(--border-sheen)]' : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] border border-transparent'}"
            title="Read-Only Preview Mode"
          >
            <Eye size={12} class="mr-1.5" />
            Preview
          </button>
          <button
            onclick={() => setEditorMode('split')}
            class="flex items-center justify-center px-2.5 py-1 rounded-md text-[11px] font-medium transition-all pressable {$editorMode === 'split' ? 'bg-[var(--bg-surface-elevated)] text-[var(--text-primary)] shadow-sm border border-[var(--border-sheen)]' : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] border border-transparent'}"
            title="Side-by-Side Split View"
          >
            <Columns2 size={12} class="mr-1.5" />
            Split
          </button>
          <button
            onclick={() => setEditorMode('edit')}
            class="flex items-center justify-center px-2.5 py-1 rounded-md text-[11px] font-medium transition-all pressable {$editorMode === 'edit' ? 'bg-[var(--bg-surface-elevated)] text-[var(--text-primary)] shadow-sm border border-[var(--border-sheen)]' : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] border border-transparent'}"
            title="Full Editor Focus Mode"
          >
            <Edit3 size={12} class="mr-1.5" />
            Edit
          </button>
        </div>

        {#if $editorMode !== 'preview'}
          <div class="h-4 w-px bg-[var(--border-subtle)]"></div>
          <button
            onclick={handleSave}
            class="flex items-center space-x-1.5 px-3 py-1 rounded-lg text-[11px] font-semibold transition-all pressable {isSavedToast ? 'bg-emerald-600 text-white shadow-[0_0_12px_rgba(16,185,129,0.35)]' : 'bg-gradient-to-r from-indigo-500 to-indigo-600 hover:from-indigo-600 hover:to-indigo-700 text-white shadow-[0_0_12px_rgba(99,102,241,0.35)]'}"
            title="Save Note (Ctrl+S)"
          >
            {#if isSavedToast}
              <Check size={12} strokeWidth={2.5} class="animate-scale-in" />
              <span>Saved!</span>
            {:else}
              <Save size={12} />
              <span>Save</span>
              {#if isDirty}
                <span class="w-1.5 h-1.5 rounded-full bg-amber-400 shadow-[0_0_6px_rgba(251,191,36,0.8)] animate-pulse" title="Unsaved changes"></span>
              {/if}
              <kbd class="text-[9px] font-mono opacity-70 ml-1 px-1 py-0.2 bg-black/20 rounded">Ctrl+S</kbd>
            {/if}
          </button>
        {/if}

        <div class="h-4 w-px bg-[var(--border-subtle)]"></div>

        <!-- Note Actions (Rename & Delete) -->
        <button
          onclick={openRenameModal}
          class="p-1.5 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] border border-transparent hover:border-[var(--border-subtle)] transition-all pressable"
          title="Rename Note"
          aria-label="Rename Note"
        >
          <Edit3 size={13} strokeWidth={2.2} />
        </button>

        <button
          onclick={openDeleteModal}
          class="p-1.5 rounded-lg text-[var(--text-muted)] hover:text-rose-500 hover:bg-rose-500/10 border border-transparent hover:border-rose-500/20 transition-all pressable"
          title="Delete Note"
          aria-label="Delete Note"
        >
          <Trash2 size={13} strokeWidth={2.2} />
        </button>
      </div>
    </div>

    <!-- External Modification Warning Banner -->
    {#if isExternalModified}
      <div class="mx-6 mt-3 px-3.5 py-2 rounded-xl bg-amber-500/10 border border-amber-500/25 flex items-center justify-between text-[11px] text-amber-500 animate-in fade-in slide-in-from-top-1 duration-200">
        <div class="flex items-center space-x-2 truncate">
          <AlertCircle size={14} class="shrink-0 text-amber-500" />
          <span class="truncate font-medium">File on disk was modified externally while editing. Press Save (Ctrl+S) to overwrite with your version.</span>
        </div>
        <button
          onclick={() => {
            if ($selectedNote) {
              editTitle = $selectedNote.title
              editContent = $selectedNote.content
              editTagsList = [...($selectedNote.tags || [])]
              loadedTitle = editTitle
              loadedContent = editContent
              loadedTagsJson = JSON.stringify(editTagsList)
              isExternalModified = false
              toast.info('Reloaded from disk')
            }
          }}
          class="ml-3 px-2 py-0.5 rounded bg-amber-500/20 hover:bg-amber-500/30 text-[10.5px] font-semibold text-amber-500 transition-colors shrink-0 pressable"
        >
          Reload from disk
        </button>
      </div>
    {/if}

    <!-- Main Workspace Area based on $editorMode -->
    {#if $editorMode === 'split'}
      <!-- SPLIT VIEW MODE: 50/50 side by side -->
      <div class="flex-1 overflow-hidden w-full">
        <div class="grid grid-cols-1 lg:grid-cols-2 h-full divide-y lg:divide-y-0 lg:divide-x divide-[var(--border-subtle)]">
          <!-- Left: Editor Pane -->
          <div class="h-full flex flex-col overflow-hidden bg-[var(--bg-surface)]/20">
            <div class="p-6 pb-3 space-y-3 border-b border-[var(--border-subtle)]/50 shrink-0">
              <input
                type="text"
                bind:value={editTitle}
                placeholder="Note Title"
                class="w-full bg-transparent border-none px-0 py-1 text-2xl font-bold text-[var(--text-primary)] focus:outline-none focus:ring-0 placeholder-[var(--text-muted)] tracking-tight heading-display"
              />
              {@render tagChipsManager()}
              {@render formatToolbar()}
            </div>
            <div class="flex-1 p-6 pt-3 overflow-y-auto">
              <textarea
                bind:this={textareaEl}
                bind:value={editContent}
                placeholder="Start writing in markdown..."
                class="w-full h-full min-h-[350px] bg-transparent border-none p-0 font-mono text-[13px] text-[var(--text-primary)] focus:outline-none focus:ring-0 leading-relaxed resize-none placeholder-[var(--text-muted)] selection:bg-[var(--accent-glow)] selection:text-white"
              ></textarea>
            </div>
          </div>

          <!-- Right: Live Preview Pane -->
          <div class="h-full overflow-y-auto p-6 lg:p-8 bg-[var(--bg-surface)]/10" bind:this={splitPreviewContainerEl}>
            <div class="max-w-[650px] mx-auto pb-16">
              <div class="text-[10px] font-mono font-semibold uppercase tracking-wider text-[var(--text-muted)] mb-2 flex items-center space-x-1.5">
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                <span>Live Preview</span>
              </div>
              <h1 class="text-2xl sm:text-3xl font-bold text-[var(--text-primary)] mb-3 tracking-tight leading-tight heading-display">
                {editTitle || 'Untitled Note'}
              </h1>
              {#if editTagsList.length > 0}
                <div class="flex items-center flex-wrap gap-1.5 mb-6">
                  {#each editTagsList as tag}
                    <span class="inline-flex items-center space-x-0.5 bg-[var(--bg-surface-active)] border border-[var(--border-subtle)] text-[var(--text-secondary)] px-2 py-0.5 rounded-md text-[10.5px] font-mono font-medium tracking-wide">
                      <span>#{tag}</span>
                    </span>
                  {/each}
                </div>
              {/if}
              <div class="prose prose-zinc dark:prose-invert max-w-none text-[var(--text-primary)] prose-headings:text-[var(--text-primary)] prose-headings:tracking-tight prose-a:text-[var(--accent-primary)] prose-a:underline-offset-4 prose-code:text-[var(--text-primary)] prose-code:bg-[var(--bg-surface-active)] prose-code:px-1.5 prose-code:py-0.5 prose-code:rounded prose-code:font-mono prose-code:before:content-none prose-code:after:content-none leading-relaxed">
                {@html renderMarkdown(editContent)}
              </div>
            </div>
          </div>
        </div>
      </div>

    {:else if $editorMode === 'edit'}
      <!-- EDIT ONLY FOCUS MODE: Centered distraction-free editor -->
      <div class="flex-1 overflow-y-auto w-full">
        <div class="max-w-[760px] mx-auto w-full px-8 py-8 space-y-4 pb-20">
          <input
            type="text"
            bind:value={editTitle}
            placeholder="Note Title"
            class="w-full bg-transparent border-none px-0 py-2 text-3xl sm:text-4xl font-bold text-[var(--text-primary)] focus:outline-none focus:ring-0 placeholder-[var(--text-muted)] tracking-tight heading-display"
          />

          <div class="border-y border-[var(--border-subtle)] py-2">
            {@render tagChipsManager()}
          </div>

          <div class="pt-1">
            {@render formatToolbar()}
          </div>

          <div class="pt-2">
            <textarea
              bind:this={textareaEl}
              bind:value={editContent}
              placeholder="Start writing in markdown..."
              class="w-full min-h-[550px] bg-transparent border-none p-0 font-mono text-[13.5px] text-[var(--text-primary)] focus:outline-none focus:ring-0 leading-relaxed resize-y placeholder-[var(--text-muted)] selection:bg-[var(--accent-glow)] selection:text-white"
            ></textarea>
          </div>
        </div>
      </div>

    {:else}
      <!-- PREVIEW ONLY MODE: Polished document reader view -->
      <div class="flex-1 overflow-y-auto w-full">
        <div bind:this={previewContainerEl} class="max-w-[760px] mx-auto w-full px-8 py-10 pb-20">
          <h1 class="text-3xl sm:text-4xl font-bold text-[var(--text-primary)] mb-4 tracking-tight leading-tight heading-display">
            {$selectedNote.title}
          </h1>

          {#if $selectedNote.tags && $selectedNote.tags.length > 0}
            <div class="flex items-center flex-wrap gap-1.5 mb-8">
              {#each $selectedNote.tags as tag}
                <span class="inline-flex items-center space-x-0.5 bg-[var(--bg-surface-active)] border border-[var(--border-subtle)] text-[var(--text-secondary)] px-2 py-0.5 rounded-md text-[10.5px] font-mono font-medium tracking-wide">
                  <span>#{tag}</span>
                </span>
              {/each}
            </div>
          {/if}

          <div class="prose prose-zinc dark:prose-invert max-w-none text-[var(--text-primary)] prose-headings:text-[var(--text-primary)] prose-headings:tracking-tight prose-a:text-[var(--accent-primary)] prose-a:underline-offset-4 prose-code:text-[var(--text-primary)] prose-code:bg-[var(--bg-surface-active)] prose-code:px-1.5 prose-code:py-0.5 prose-code:rounded prose-code:font-mono prose-code:before:content-none prose-code:after:content-none leading-relaxed">
            {@html renderMarkdown($selectedNote.content)}
          </div>
        </div>
      </div>
    {/if}
  {/if}
</div>

<!-- Reusable Note Action Modals -->
{#if $selectedNote}
  <ConfirmModal
    open={isDeleteModalOpen}
    destructive={true}
    title="Delete Note"
    description={`Are you sure you want to permanently delete "${$selectedNote.title || 'this note'}"? This action cannot be undone.`}
    confirmText="Delete Note"
    onConfirm={handleConfirmDelete}
    onCancel={() => {
      isDeleteModalOpen = false
    }}
  />

  <ConfirmModal
    open={isRenameModalOpen}
    mode="prompt"
    title="Rename Note"
    initialValue={$selectedNote.title || ''}
    placeholder="Enter new note title..."
    confirmText="Rename"
    onConfirm={handleConfirmRename}
    onCancel={() => {
      isRenameModalOpen = false
    }}
  />
{/if}
