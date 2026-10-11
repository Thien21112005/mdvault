<script lang="ts">
  import { X, AlertTriangle, Edit3, Trash2 } from '@lucide/svelte'

  export interface ConfirmModalProps {
    open: boolean
    title: string
    description?: string
    confirmText?: string
    cancelText?: string
    destructive?: boolean
    mode?: 'confirm' | 'prompt'
    initialValue?: string
    placeholder?: string
    onConfirm: (val?: string) => unknown
    onCancel: () => unknown
  }

  let {
    open,
    title,
    description = '',
    confirmText = 'Confirm',
    cancelText = 'Cancel',
    destructive = false,
    mode = 'confirm',
    initialValue = '',
    placeholder = '',
    onConfirm,
    onCancel
  }: ConfirmModalProps = $props()

  let inputValue = $state('')
  let inputEl = $state<HTMLInputElement | null>(null)

  $effect(() => {
    if (open) {
      inputValue = initialValue
      setTimeout(() => {
        if (inputEl) {
          inputEl.focus()
          inputEl.select()
        }
      }, 50)
    }
  })

  function handleKeydown(e: KeyboardEvent) {
    if (!open) return
    if (e.key === 'Escape') {
      e.preventDefault()
      onCancel()
    } else if (e.key === 'Enter') {
      e.preventDefault()
      handleConfirm()
    }
  }

  function handleConfirm() {
    if (mode === 'prompt') {
      if (!inputValue.trim()) return
      onConfirm(inputValue.trim())
    } else {
      onConfirm()
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if open}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <div
    role="dialog"
    aria-modal="true"
    tabindex="-1"
    class="fixed inset-0 z-50 bg-black/60 backdrop-blur-md flex items-center justify-center p-4 select-none animate-in fade-in duration-150"
    onclick={(e) => {
      if (e.target === e.currentTarget) onCancel()
    }}
  >
    <div
      class="w-full max-w-sm bg-[var(--bg-surface-elevated)]/95 backdrop-blur-2xl border border-[var(--border-subtle)] rounded-2xl shadow-[var(--shadow-glass)] p-5 animate-in zoom-in-95 duration-150 relative space-y-4"
    >
      <div class="flex items-start justify-between">
        <div class="flex items-center space-x-2.5">
          <div class="p-2 rounded-xl {destructive ? 'bg-rose-500/15 text-rose-500' : 'bg-[var(--accent-primary)]/15 text-[var(--accent-primary)]'}">
            {#if destructive}
              <Trash2 size={16} strokeWidth={2.2} />
            {:else if mode === 'prompt'}
              <Edit3 size={16} strokeWidth={2.2} />
            {:else}
              <AlertTriangle size={16} strokeWidth={2.2} />
            {/if}
          </div>
          <div>
            <h3 class="font-bold text-[14px] text-[var(--text-primary)] tracking-tight">
              {title}
            </h3>
            {#if description}
              <p class="text-[11.5px] text-[var(--text-muted)] mt-0.5 leading-relaxed">
                {description}
              </p>
            {/if}
          </div>
        </div>

        <button
          onclick={onCancel}
          class="p-1 rounded-md text-[var(--text-muted)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-colors"
          title="Close (ESC)"
        >
          <X size={14} />
        </button>
      </div>

      {#if mode === 'prompt'}
        <div class="pt-1">
          <input
            bind:this={inputEl}
            type="text"
            bind:value={inputValue}
            {placeholder}
            class="w-full bg-[var(--bg-surface-active)] border border-[var(--border-subtle)] rounded-xl px-3 py-2 text-[12.5px] text-[var(--text-primary)] placeholder-[var(--text-muted)] focus:outline-none focus:border-[var(--accent-primary)] focus:ring-1 focus:ring-[var(--accent-primary)] transition-all font-medium"
            onkeydown={(e) => {
              if (e.key === 'Enter') handleConfirm()
            }}
          />
        </div>
      {/if}

      <div class="flex items-center justify-end space-x-2 pt-2">
        <button
          onclick={onCancel}
          class="px-3.5 py-1.5 rounded-xl text-[12px] font-semibold text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] border border-transparent hover:border-[var(--border-subtle)] transition-all pressable"
        >
          {cancelText}
        </button>
        <button
          onclick={handleConfirm}
          disabled={mode === 'prompt' && !inputValue.trim()}
          class="px-4 py-1.5 rounded-xl text-[12px] font-semibold transition-all pressable disabled:opacity-40 {destructive ? 'bg-rose-600 hover:bg-rose-700 text-white shadow-[0_0_12px_rgba(225,29,72,0.35)]' : 'bg-gradient-to-r from-indigo-500 to-indigo-600 hover:from-indigo-600 hover:to-indigo-700 text-white shadow-[0_0_12px_rgba(99,102,241,0.35)]'}"
        >
          {confirmText}
        </button>
      </div>
    </div>
  </div>
{/if}
