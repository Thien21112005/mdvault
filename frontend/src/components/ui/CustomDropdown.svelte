<script lang="ts">
  import { ChevronDown, Check } from '@lucide/svelte'

  export interface DropdownOption {
    value: string
    label: string
    icon?: any
    description?: string
  }

  interface Props {
    options: DropdownOption[]
    value?: string
    placeholder?: string
    align?: 'left' | 'right'
    class?: string
  }

  let {
    options = [],
    value = $bindable(''),
    placeholder = 'Select...',
    align = 'left',
    class: customClass = ''
  }: Props = $props()

  let isOpen = $state(false)
  let highlightedIndex = $state(0)
  let triggerEl = $state<HTMLButtonElement | null>(null)
  let menuEl = $state<HTMLDivElement | null>(null)

  const selectedOption = $derived(options.find(opt => opt.value === value))
  const displayLabel = $derived(selectedOption ? selectedOption.label : placeholder)

  function selectOption(val: string) {
    value = val
    isOpen = false
    triggerEl?.focus()
  }

  function handleClickOutside(e: MouseEvent) {
    if (!isOpen) return
    const target = e.target as Node
    if (triggerEl && triggerEl.contains(target)) return
    if (menuEl && menuEl.contains(target)) return
    isOpen = false
  }

  $effect(() => {
    if (isOpen) {
      const currIdx = options.findIndex(opt => opt.value === value)
      highlightedIndex = currIdx >= 0 ? currIdx : 0
      window.addEventListener('click', handleClickOutside)
      return () => window.removeEventListener('click', handleClickOutside)
    }
  })

  function handleKeydown(e: KeyboardEvent) {
    if (!isOpen) {
      if (e.key === 'ArrowDown' || e.key === 'Enter' || e.key === ' ') {
        e.preventDefault()
        isOpen = true
      }
      return
    }

    if (e.key === 'Escape') {
      e.preventDefault()
      isOpen = false
      triggerEl?.focus()
    } else if (e.key === 'ArrowDown') {
      e.preventDefault()
      if (options.length > 0) {
        highlightedIndex = (highlightedIndex + 1) % options.length
      }
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      if (options.length > 0) {
        highlightedIndex = (highlightedIndex - 1 + options.length) % options.length
      }
    } else if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault()
      const opt = options[highlightedIndex]
      if (opt) {
        selectOption(opt.value)
      }
    } else if (e.key === 'Tab') {
      isOpen = false
    }
  }
</script>

<div class="relative inline-block text-left {customClass}">
  <!-- Trigger Button -->
  <button
    bind:this={triggerEl}
    type="button"
    onclick={() => (isOpen = !isOpen)}
    onkeydown={handleKeydown}
    aria-haspopup="listbox"
    aria-expanded={isOpen}
    class="inline-flex items-center justify-between space-x-1.5 px-2.5 py-1 rounded-lg text-[11px] font-medium transition-all pressable bg-[var(--bg-surface-active)]/80 hover:bg-[var(--bg-surface-hover)] border border-[var(--border-subtle)] hover:border-[var(--border-strong)] text-[var(--text-primary)] shadow-xs select-none cursor-pointer"
  >
    <span class="truncate">{displayLabel}</span>
    <ChevronDown
      size={11.5}
      strokeWidth={2.2}
      class="text-[var(--text-muted)] transition-transform duration-200 ease-[var(--ease-apple)] {isOpen ? 'rotate-180 text-[var(--accent-primary)]' : ''}"
    />
  </button>

  <!-- Glassmorphic Dropdown Popover Menu -->
  {#if isOpen}
    <div
      bind:this={menuEl}
      role="listbox"
      tabindex="-1"
      class="absolute z-50 mt-1.5 min-w-[150px] {align === 'right' ? 'right-0' : 'left-0'} p-1 rounded-xl bg-[var(--bg-surface-elevated)]/95 backdrop-blur-2xl border border-[var(--border-subtle)] shadow-[var(--shadow-glass)] transition-all animate-scale-in select-none focus:outline-none"
    >
      {#each options as opt, idx (opt.value)}
        <button
          type="button"
          role="option"
          aria-selected={opt.value === value}
          onclick={() => selectOption(opt.value)}
          onmouseenter={() => (highlightedIndex = idx)}
          class="w-full text-left flex items-center justify-between px-2.5 py-1.5 rounded-lg text-[11px] font-medium transition-colors pressable cursor-pointer {opt.value === value ? 'bg-[var(--accent-soft)] text-[var(--accent-primary)] font-semibold' : highlightedIndex === idx ? 'bg-[var(--bg-surface-hover)] text-[var(--text-primary)]' : 'text-[var(--text-secondary)]'}"
        >
          <div class="flex items-center space-x-2 truncate">
            {#if opt.icon}
              <opt.icon size={12.5} class="shrink-0 text-[var(--text-muted)]" />
            {/if}
            <span class="truncate">{opt.label}</span>
          </div>
          {#if opt.value === value}
            <Check size={12} strokeWidth={2.5} class="text-[var(--accent-primary)] ml-2 shrink-0 animate-scale-in" />
          {/if}
        </button>
      {/each}
    </div>
  {/if}
</div>
