<script lang="ts">
  import { Edit3, Eye, Save, Tag, FileText, Calendar } from '@lucide/svelte'
  import { selectedNote, saveActiveNote, notes, selectedNoteId } from '$lib/stores/vaultStore'
  import { isEditingNote } from '$lib/stores/viewStore'
  import { renderMarkdown } from '$lib/utils/markdown'

  let editTitle = $state('')
  let editTags = $state('')
  let editContent = $state('')

  $effect(() => {
    if ($selectedNote) {
      editTitle = $selectedNote.title
      editTags = ($selectedNote.tags || []).join(', ')
      editContent = $selectedNote.content
    }
  })

  function handleSave() {
    if (!$selectedNote) return
    const tagsArr = editTags.split(',').map(t => t.trim()).filter(Boolean)
    saveActiveNote($selectedNote.id, editContent, editTitle, tagsArr)
    isEditingNote.set(false)
  }

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
  let previewContainerEl = $state<HTMLDivElement | null>(null)

  $effect(() => {
    const el = previewContainerEl
    if (el) {
      el.addEventListener('click', handlePreviewClick)
      return () => el.removeEventListener('click', handlePreviewClick)
    }
  })
</script>

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
    <div class="h-12 border-b border-[var(--border-subtle)] px-6 flex items-center justify-between select-none bg-[var(--bg-surface)]/80 backdrop-blur-md sticky top-0 z-10">
      <div class="flex items-center space-x-3 text-[11px] font-medium text-[var(--text-muted)]">
        <span class="flex items-center space-x-1.5 font-mono">
          <Calendar size={12} />
          <span>Modified { $selectedNote.lastModified ? $selectedNote.lastModified.substring(0, 10) : '' }</span>
        </span>
      </div>

      <div class="flex items-center space-x-3">
        <!-- Segmented Control for Edit/Preview -->
        <div class="bg-[var(--bg-surface-active)] p-0.5 rounded-lg border border-[var(--border-subtle)] flex items-center shadow-inner">
          <button
            onclick={() => isEditingNote.set(false)}
            class="flex items-center justify-center px-3 py-1 rounded-md text-[11px] font-medium transition-all pressable {!$isEditingNote ? 'bg-[var(--bg-surface-elevated)] text-[var(--text-primary)] shadow-sm border border-[var(--border-sheen)]' : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] border border-transparent'}"
          >
            <Eye size={12} class="mr-1.5" />
            Preview
          </button>
          <button
            onclick={() => isEditingNote.set(true)}
            class="flex items-center justify-center px-3 py-1 rounded-md text-[11px] font-medium transition-all pressable {$isEditingNote ? 'bg-[var(--bg-surface-elevated)] text-[var(--text-primary)] shadow-sm border border-[var(--border-sheen)]' : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] border border-transparent'}"
          >
            <Edit3 size={12} class="mr-1.5" />
            Edit
          </button>
        </div>

        {#if $isEditingNote}
          <div class="h-4 w-px bg-[var(--border-subtle)]"></div>
          <button
            onclick={handleSave}
            class="flex items-center space-x-1.5 bg-gradient-to-r from-indigo-500 to-indigo-600 hover:from-indigo-600 hover:to-indigo-700 text-white px-3 py-1 rounded-lg text-[11px] font-semibold transition-all pressable shadow-[0_0_12px_rgba(99,102,241,0.35)]"
          >
            <Save size={12} />
            <span>Save Note</span>
          </button>
        {/if}
      </div>
    </div>

    <!-- Main Workspace Area -->
    <div class="flex-1 overflow-y-auto w-full">
      <div class="max-w-[760px] mx-auto w-full px-8 py-10">
        {#if $isEditingNote}
          <div class="space-y-6">
            <div>
              <input
                type="text"
                bind:value={editTitle}
                placeholder="Note Title"
                class="w-full bg-transparent border-none px-0 py-2 text-3xl sm:text-4xl font-bold text-[var(--text-primary)] focus:outline-none focus:ring-0 placeholder-[var(--text-muted)] tracking-tight heading-display"
              />
            </div>
            
            <div class="flex items-center space-x-2 border-b border-[var(--border-subtle)] pb-3">
              <Tag size={13} class="text-[var(--text-muted)]" />
              <input
                type="text"
                bind:value={editTags}
                placeholder="Add tags separated by commas..."
                class="flex-1 bg-transparent border-none text-[12px] text-[var(--text-secondary)] focus:outline-none focus:ring-0 placeholder-[var(--text-muted)] tracking-wide"
              />
            </div>

            <div class="flex-1 pt-2">
              <textarea
                bind:value={editContent}
                placeholder="Start writing in markdown..."
                class="w-full min-h-[500px] bg-transparent border-none p-0 font-mono text-[13px] text-[var(--text-primary)] focus:outline-none focus:ring-0 leading-relaxed resize-y placeholder-[var(--text-muted)] selection:bg-[var(--accent-glow)] selection:text-white"
              ></textarea>
            </div>
          </div>
        {:else}
          <!-- Markdown Preview Mode with Security Sanitizer & Interactive Elements -->
          <div bind:this={previewContainerEl} class="pb-16">
            <h1 class="text-3xl sm:text-4xl font-bold text-[var(--text-primary)] mb-4 tracking-tight leading-tight heading-display">{$selectedNote.title}</h1>
            
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
        {/if}
      </div>
    </div>
  {/if}
</div>
