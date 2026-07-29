<script lang="ts">
  import { FileText, Calendar, Tag } from '@lucide/svelte'
  import { filteredNotes, selectedNoteId } from '$lib/stores/vaultStore'
  import { viewMode } from '$lib/stores/viewStore'
</script>

<div class="w-80 bg-slate-900/90 border-r border-slate-800 flex flex-col h-[calc(100vh-3.5rem)] select-none">
  <div class="p-3 border-b border-slate-800 flex items-center justify-between bg-slate-900/50">
    <span class="text-xs font-bold text-slate-300">GHI CHÚ ({$filteredNotes.length})</span>
    <span class="text-xs text-slate-500 font-medium">Sắp xếp: Mới nhất</span>
  </div>

  <div class="flex-1 overflow-y-auto p-2 space-y-1.5">
    {#if $filteredNotes.length === 0}
      <div class="p-8 text-center text-slate-500 text-sm italic">
        Không có ghi chú nào. Bấm nút "+ Note" để tạo!
      </div>
    {:else}
      {#each $filteredNotes as note (note.id)}
        <button
          onclick={() => { selectedNoteId.set(note.id); viewMode.set('list') }}
          class="w-full text-left p-3 rounded-xl transition border {$selectedNoteId === note.id ? 'bg-blue-600/10 border-blue-500/50 text-white shadow-md' : 'bg-slate-800/40 border-slate-800/80 text-slate-300 hover:bg-slate-800/80 hover:border-slate-700'}"
        >
          <div class="flex items-start justify-between">
            <span class="font-bold text-sm truncate w-full text-slate-200">{note.title || 'Untitled Note'}</span>
          </div>
          <p class="text-xs text-slate-400 mt-1.5 line-clamp-2 font-normal leading-relaxed">{note.excerpt}</p>
          
          <div class="mt-2 flex items-center justify-between text-[11px] text-slate-500">
            <div class="flex items-center space-x-1">
              <Calendar size={11} />
              <span>{note.lastModified ? note.lastModified.substring(0, 10) : ''}</span>
            </div>
            {#if note.tags && note.tags.length > 0}
              <div class="flex items-center space-x-1 bg-slate-800 px-1.5 py-0.5 rounded text-blue-400">
                <Tag size={10} />
                <span>{note.tags.length}</span>
              </div>
            {/if}
          </div>
        </button>
      {/each}
    {/if}
  </div>
</div>
