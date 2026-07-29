<script lang="ts">
  import { Edit3, Eye, Save, Tag, FileText, Calendar } from '@lucide/svelte'
  import { selectedNote, saveActiveNote } from '$lib/stores/vaultStore'
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
</script>

<div class="flex-1 bg-[#0a0a0a] flex flex-col h-[calc(100vh-3rem)] overflow-hidden relative">
  {#if !$selectedNote}
    <div class="flex-1 flex flex-col items-center justify-center text-zinc-600 space-y-4">
      <div class="p-4 rounded-2xl bg-zinc-900/30 border border-zinc-800/50">
        <FileText size={32} strokeWidth={1.5} class="text-zinc-500" />
      </div>
      <p class="text-[13px] font-medium tracking-wide">Select a note to view or edit</p>
    </div>
  {:else}
    <!-- Header bar -->
    <div class="h-12 border-b border-zinc-800/80 px-6 flex items-center justify-between select-none bg-[#0a0a0a]/80 backdrop-blur-sm sticky top-0 z-10">
      <div class="flex items-center space-x-3 text-[11px] font-medium text-zinc-500">
        <span class="flex items-center space-x-1.5">
          <Calendar size={12} />
          <span>Modified { $selectedNote.lastModified ? $selectedNote.lastModified.substring(0, 10) : '' }</span>
        </span>
      </div>

      <div class="flex items-center space-x-3">
        <!-- Segmented Control for Edit/Preview -->
        <div class="bg-zinc-900/80 p-0.5 rounded-md border border-zinc-800/80 flex items-center shadow-sm">
          <button
            onclick={() => isEditingNote.set(false)}
            class="flex items-center justify-center px-3 py-1 rounded-[4px] text-[11px] font-medium transition-all {!$isEditingNote ? 'bg-zinc-700/50 text-zinc-100 shadow-sm' : 'text-zinc-500 hover:text-zinc-300 hover:bg-zinc-800/50'}"
          >
            <Eye size={12} class="mr-1.5" />
            Preview
          </button>
          <button
            onclick={() => isEditingNote.set(true)}
            class="flex items-center justify-center px-3 py-1 rounded-[4px] text-[11px] font-medium transition-all {$isEditingNote ? 'bg-zinc-700/50 text-zinc-100 shadow-sm' : 'text-zinc-500 hover:text-zinc-300 hover:bg-zinc-800/50'}"
          >
            <Edit3 size={12} class="mr-1.5" />
            Edit
          </button>
        </div>

        {#if $isEditingNote}
          <div class="h-4 w-px bg-zinc-800"></div>
          <button
            onclick={handleSave}
            class="flex items-center space-x-1.5 bg-zinc-100 hover:bg-white text-zinc-950 px-3 py-1 rounded-md text-[11px] font-semibold transition-all active:scale-[0.97] shadow-sm"
          >
            <Save size={12} />
            <span>Save Note</span>
          </button>
        {/if}
      </div>
    </div>

    <!-- Main Workspace Area -->
    <div class="flex-1 overflow-y-auto w-full">
      <div class="max-w-[760px] mx-auto w-full px-8 py-12">
        {#if $isEditingNote}
          <div class="space-y-6">
            <div>
              <input
                type="text"
                bind:value={editTitle}
                placeholder="Note Title"
                class="w-full bg-transparent border-none px-0 py-2 text-3xl font-bold text-zinc-100 focus:outline-none focus:ring-0 placeholder-zinc-700 tracking-tight"
              />
            </div>
            
            <div class="flex items-center space-x-2 border-b border-zinc-800/50 pb-4">
              <Tag size={14} class="text-zinc-600" />
              <input
                type="text"
                bind:value={editTags}
                placeholder="Add tags separated by commas..."
                class="flex-1 bg-transparent border-none text-[13px] text-zinc-300 focus:outline-none focus:ring-0 placeholder-zinc-600"
              />
            </div>

            <div class="flex-1 pt-2">
              <textarea
                bind:value={editContent}
                placeholder="Start writing..."
                class="w-full min-h-[500px] bg-transparent border-none p-0 font-mono text-[13px] text-zinc-300 focus:outline-none focus:ring-0 leading-relaxed resize-y placeholder-zinc-700 selection:bg-zinc-800"
              ></textarea>
            </div>
          </div>
        {:else}
          <!-- Markdown Preview Mode using markdown-it -->
          <div class="pb-16">
            <h1 class="text-4xl font-extrabold text-zinc-100 mb-4 tracking-tight leading-tight">{$selectedNote.title}</h1>
            
            {#if $selectedNote.tags && $selectedNote.tags.length > 0}
              <div class="flex items-center flex-wrap gap-2 mb-10">
                {#each $selectedNote.tags as tag}
                  <span class="bg-zinc-900 border border-zinc-800 text-zinc-400 px-2.5 py-0.5 rounded-sm text-[11px] font-medium tracking-wide">
                    {tag}
                  </span>
                {/each}
              </div>
            {/if}
            
            <div class="prose prose-invert prose-zinc max-w-none leading-loose text-zinc-300 prose-headings:text-zinc-100 prose-headings:tracking-tight prose-a:text-zinc-300 prose-a:underline-offset-4 hover:prose-a:text-zinc-100 prose-strong:text-zinc-100 prose-code:text-zinc-200 prose-code:bg-zinc-900 prose-code:px-1 prose-code:py-0.5 prose-code:rounded-sm prose-code:before:content-none prose-code:after:content-none prose-pre:bg-[#0d0d0d] prose-pre:border prose-pre:border-zinc-800/80">
              {@html renderMarkdown($selectedNote.content)}
            </div>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>
