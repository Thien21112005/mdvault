<script lang="ts">
  import { FileText, Calendar, Tag } from '@lucide/svelte'
  import { filteredNotes, selectedNoteId } from '$lib/stores/vaultStore'
  import { viewMode } from '$lib/stores/viewStore'
</script>

<div class="w-80 bg-white/[0.01] border-r border-white/[0.06] flex flex-col h-[calc(100vh-3rem)] select-none shadow-[inset_-1px_0_0_rgba(0,0,0,0.2)]">
  <div class="px-4 py-3 border-b border-white/[0.04] flex items-center justify-between sticky top-0 z-10 backdrop-blur-md">
    <span class="text-[10px] font-semibold text-zinc-500 uppercase tracking-widest">Notes ({$filteredNotes.length})</span>
    <span class="text-[10px] text-zinc-500 font-medium tracking-wide">Sort: Newest</span>
  </div>

  <div class="flex-1 overflow-y-auto p-2 space-y-1">
    {#if $filteredNotes.length === 0}
      <div class="p-8 text-center text-zinc-500 text-[11px] mt-10">
        No notes found. Click "New Note" to create one.
      </div>
    {:else}
      {#each $filteredNotes as note (note.id)}
        <button
          onclick={() => { selectedNoteId.set(note.id); viewMode.set('list') }}
          class="w-full text-left p-3 rounded-lg transition-all active:scale-[0.99] border relative overflow-hidden group {$selectedNoteId === note.id ? 'bg-white/[0.06] border-white/[0.08] shadow-[inset_0_1px_1px_rgba(255,255,255,0.05)]' : 'bg-transparent border-transparent hover:bg-white/[0.03] hover:border-white/[0.04]'}"
        >
          {#if $selectedNoteId === note.id}
             <!-- glowing accent indicator line -->
             <div class="absolute left-0 top-1/2 -translate-y-1/2 h-1/2 w-[2px] bg-zinc-300 rounded-r-full shadow-[0_0_8px_rgba(255,255,255,0.4)]"></div>
          {/if}
          <div class="flex items-start justify-between">
            <span class="font-semibold text-[13px] tracking-tight truncate w-full {$selectedNoteId === note.id ? 'text-zinc-100' : 'text-zinc-300 group-hover:text-zinc-200'}">{note.title || 'Untitled Note'}</span>
          </div>
          <p class="text-[11px] mt-1 line-clamp-2 leading-relaxed {$selectedNoteId === note.id ? 'text-zinc-400' : 'text-zinc-500'}">{note.excerpt}</p>
          
          <div class="mt-2.5 flex items-center justify-between text-[10px] font-medium {$selectedNoteId === note.id ? 'text-zinc-500' : 'text-zinc-600'}">
            <div class="flex items-center space-x-1.5">
              <Calendar size={10} />
              <span>{note.lastModified ? note.lastModified.substring(0, 10) : ''}</span>
            </div>
            {#if note.tags && note.tags.length > 0}
              <div class="flex items-center space-x-1 bg-white/[0.05] border border-white/[0.08] px-1.5 py-0.5 rounded-sm {$selectedNoteId === note.id ? 'text-zinc-300' : 'text-zinc-400'}">
                <Tag size={9} />
                <span>{note.tags.length}</span>
              </div>
            {/if}
          </div>
        </button>
      {/each}
    {/if}
  </div>
</div>
