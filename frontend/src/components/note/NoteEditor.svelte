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

<div class="flex-1 bg-slate-950 flex flex-col h-[calc(100vh-3.5rem)] overflow-hidden">
  {#if !$selectedNote}
    <div class="flex-1 flex flex-col items-center justify-center text-slate-600 space-y-3">
      <FileText size={48} class="opacity-30" />
      <p class="text-sm font-medium">Chọn một ghi chú từ danh sách hoặc tạo ghi chú mới để bắt đầu</p>
    </div>
  {:else}
    <!-- Header bar -->
    <div class="h-12 bg-slate-900 border-b border-slate-800 px-6 flex items-center justify-between select-none">
      <div class="flex items-center space-x-3 text-xs text-slate-400">
        <span class="flex items-center space-x-1">
          <Calendar size={14} />
          <span>Sửa lần cuối: {$selectedNote.lastModified ? $selectedNote.lastModified.substring(0, 10) : ''}</span>
        </span>
      </div>

      <div class="flex items-center space-x-2">
        <button
          onclick={() => isEditingNote.set(!$isEditingNote)}
          class="flex items-center space-x-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition {$isEditingNote ? 'bg-slate-800 text-yellow-400 border border-yellow-500/30' : 'bg-slate-800 text-slate-300 hover:bg-slate-700'}"
        >
          {#if $isEditingNote}
            <Eye size={14} />
            <span>Xem trước</span>
          {:else}
            <Edit3 size={14} />
            <span>Chỉnh sửa</span>
          {/if}
        </button>

        {#if $isEditingNote}
          <button
            onclick={handleSave}
            class="flex items-center space-x-1.5 bg-emerald-600 hover:bg-emerald-500 text-white px-3.5 py-1.5 rounded-lg text-xs font-semibold shadow transition"
          >
            <Save size={14} />
            <span>Lưu file</span>
          </button>
        {/if}
      </div>
    </div>

    <!-- Main Workspace Area -->
    <div class="flex-1 overflow-y-auto p-8 max-w-4xl mx-auto w-full">
      {#if $isEditingNote}
        <div class="space-y-4">
          <div>
            <label class="text-xs font-bold text-slate-400 uppercase tracking-wider block mb-1">Tiêu đề Note</label>
            <input
              type="text"
              bind:value={editTitle}
              class="w-full bg-slate-900 border border-slate-700 rounded-xl px-4 py-2.5 text-lg font-bold text-slate-100 focus:outline-none focus:border-blue-500"
            />
          </div>
          <div>
            <label class="text-xs font-bold text-slate-400 uppercase tracking-wider block mb-1">Tags (phân cách bằng dấu phẩy)</label>
            <div class="relative">
              <Tag size={16} class="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-500" />
              <input
                type="text"
                bind:value={editTags}
                placeholder="tag1, tag2, kiến_thức_svelte..."
                class="w-full bg-slate-900 border border-slate-700 rounded-xl pl-10 pr-4 py-2 text-sm text-slate-200 focus:outline-none focus:border-blue-500"
              />
            </div>
          </div>
          <div class="flex-1">
            <label class="text-xs font-bold text-slate-400 uppercase tracking-wider block mb-1">Nội dung Markdown</label>
            <textarea
              bind:value={editContent}
              rows={20}
              class="w-full bg-slate-900 border border-slate-700 rounded-xl p-4 font-mono text-sm text-slate-200 focus:outline-none focus:border-blue-500 leading-relaxed resize-y"
            ></textarea>
          </div>
        </div>
      {:else}
        <!-- Markdown Preview Mode using markdown-it -->
        <div>
          <h1 class="text-3xl font-extrabold text-slate-100 mb-3 tracking-tight">{$selectedNote.title}</h1>
          {#if $selectedNote.tags && $selectedNote.tags.length > 0}
            <div class="flex items-center space-x-1.5 mb-6">
              {#each $selectedNote.tags as tag}
                <span class="bg-blue-600/20 text-blue-400 border border-blue-500/30 px-2.5 py-0.5 rounded-full text-xs font-semibold">
                  #{tag}
                </span>
              {/each}
            </div>
          {/if}
          <div class="border-t border-slate-800 pt-6">
            <div class="prose prose-invert prose-blue max-w-none leading-relaxed text-slate-300">
              {@html renderMarkdown($selectedNote.content)}
            </div>
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>
