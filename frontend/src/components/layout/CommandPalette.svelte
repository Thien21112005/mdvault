<script lang="ts">
  import {
    Search,
    FileText,
    FilePlus,
    Network,
    Sun,
    Moon,
    PanelLeft,
    FolderOpen,
    Sparkles,
    CornerDownLeft,
    X
  } from '@lucide/svelte'
  import { notes, selectedNoteId, createNewNote, openWorkspaceDialog } from '$lib/stores/vaultStore'
  import {
    viewMode,
    switchView,
    toggleSidebar,
    isCommandPaletteOpen,
    closeCommandPalette,
    setEditorMode
  } from '$lib/stores/viewStore'
  import { isDark, toggleTheme } from '$lib/stores/themeStore'
  import { toast } from '$lib/stores/toastStore'

  let query = $state('')
  let selectedIndex = $state(0)
  let inputEl = $state<HTMLInputElement | null>(null)
  let listContainerEl = $state<HTMLDivElement | null>(null)

  // Focus input automatically when palette opens
  $effect(() => {
    if ($isCommandPaletteOpen) {
      query = ''
      selectedIndex = 0
      requestAnimationFrame(() => {
        if (inputEl) inputEl.focus()
      })
    }
  })

  // Quick Action Commands
  const quickActions = [
    {
      id: 'cmd-new-note',
      type: 'action' as const,
      title: 'Create New Note',
      subtitle: 'Create a blank markdown note',
      icon: FilePlus,
      run: async () => {
        closeCommandPalette()
        await createNewNote('Untitled Note')
        setEditorMode('edit')
        toast.success('Note Created', 'Ready for editing')
      }
    },
    {
      id: 'cmd-view-graph',
      type: 'action' as const,
      title: 'Open Knowledge Graph',
      subtitle: 'Explore notes and links in galaxy map',
      icon: Network,
      run: () => {
        closeCommandPalette()
        switchView('graph')
        toast.info('Switched to Graph View')
      }
    },
    {
      id: 'cmd-view-list',
      type: 'action' as const,
      title: 'Open Notes List',
      subtitle: 'Browse and edit your documents',
      icon: FileText,
      run: () => {
        closeCommandPalette()
        switchView('list')
      }
    },
    {
      id: 'cmd-toggle-theme',
      type: 'action' as const,
      title: 'Toggle Color Theme',
      subtitle: 'Switch between Dark and Light mode',
      icon: Sun,
      run: () => {
        toggleTheme()
        toast.info('Theme Toggled')
      }
    },
    {
      id: 'cmd-toggle-sidebar',
      type: 'action' as const,
      title: 'Toggle Sidebar',
      subtitle: 'Show or hide folder tree',
      icon: PanelLeft,
      run: () => {
        toggleSidebar()
      }
    },
    {
      id: 'cmd-open-vault',
      type: 'action' as const,
      title: 'Open Vault Folder',
      subtitle: 'Select workspace root directory',
      icon: FolderOpen,
      run: async () => {
        closeCommandPalette()
        await openWorkspaceDialog()
      }
    }
  ]

  // Flattened items based on search query
  const matchingNotes = $derived.by(() => {
    const q = query.trim().toLowerCase()
    if (!q) {
      // Show first 6 recent notes
      return $notes.slice(0, 6)
    }
    return $notes.filter(n => {
      const titleMatch = (n.title || '').toLowerCase().includes(q)
      const tagsMatch = (n.tags || []).some(t => t.toLowerCase().includes(q))
      const excerptMatch = (n.excerpt || '').toLowerCase().includes(q)
      return titleMatch || tagsMatch || excerptMatch
    })
  })

  const matchingActions = $derived.by(() => {
    const q = query.trim().toLowerCase()
    if (!q) return quickActions
    return quickActions.filter(a => {
      return a.title.toLowerCase().includes(q) || a.subtitle.toLowerCase().includes(q)
    })
  })

  interface PaletteItem {
    id: string
    category: 'action' | 'note'
    title: string
    subtitle?: string
    tags?: string[]
    icon?: any
    run: () => void
  }

  const flattenedItems = $derived.by<PaletteItem[]>(() => {
    const items: PaletteItem[] = []

    // 1. Actions
    matchingActions.forEach(a => {
      items.push({
        id: a.id,
        category: 'action',
        title: a.title,
        subtitle: a.subtitle,
        icon: a.icon,
        run: a.run
      })
    })

    // 2. Notes
    matchingNotes.forEach(n => {
      items.push({
        id: `note-${n.id}`,
        category: 'note',
        title: n.title || 'Untitled Note',
        subtitle: n.excerpt ? n.excerpt.substring(0, 70) : undefined,
        tags: n.tags || [],
        icon: FileText,
        run: () => {
          closeCommandPalette()
          selectedNoteId.set(n.id)
          switchView('list')
        }
      })
    })

    return items
  })

  // Keep selected index within bounds
  $effect(() => {
    if (selectedIndex >= flattenedItems.length) {
      selectedIndex = Math.max(0, flattenedItems.length - 1)
    }
  })

  function handleKeydown(e: KeyboardEvent) {
    if (!$isCommandPaletteOpen) return

    if (e.key === 'Escape') {
      e.preventDefault()
      closeCommandPalette()
    } else if (e.key === 'ArrowDown') {
      e.preventDefault()
      if (flattenedItems.length > 0) {
        selectedIndex = (selectedIndex + 1) % flattenedItems.length
        scrollSelectedIntoView()
      }
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      if (flattenedItems.length > 0) {
        selectedIndex = (selectedIndex - 1 + flattenedItems.length) % flattenedItems.length
        scrollSelectedIntoView()
      }
    } else if (e.key === 'Enter') {
      e.preventDefault()
      const item = flattenedItems[selectedIndex]
      if (item) {
        item.run()
      }
    }
  }

  function scrollSelectedIntoView() {
    requestAnimationFrame(() => {
      if (listContainerEl) {
        const selectedEl = listContainerEl.querySelector(`[data-index="${selectedIndex}"]`) as HTMLElement | null
        if (selectedEl) {
          selectedEl.scrollIntoView({ block: 'nearest' })
        }
      }
    })
  }

  function handleBackdropClick(e: MouseEvent) {
    if (e.target === e.currentTarget) {
      closeCommandPalette()
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if $isCommandPaletteOpen}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <div
    role="dialog"
    aria-modal="true"
    tabindex="-1"
    onclick={handleBackdropClick}
    class="fixed inset-0 z-50 bg-black/60 backdrop-blur-md flex items-start justify-center pt-[12vh] px-4 select-none"
  >
    <div
      class="w-full max-w-xl bg-[var(--bg-surface-elevated)]/95 backdrop-blur-2xl border border-[var(--border-subtle)] rounded-2xl shadow-[var(--shadow-glass)] overflow-hidden flex flex-col transition-all animate-scale-in max-h-[75vh]"
    >
      <!-- Top Search Input Bar -->
      <div class="h-14 px-4 border-b border-[var(--border-subtle)] flex items-center space-x-3 shrink-0">
        <Search size={18} class="text-[var(--accent-primary)] shrink-0" strokeWidth={2.2} />
        <input
          bind:this={inputEl}
          type="text"
          bind:value={query}
          placeholder="Type a command or search notes..."
          class="flex-1 bg-transparent border-none text-[14px] text-[var(--text-primary)] placeholder-[var(--text-muted)] focus:outline-none focus:ring-0 tracking-tight"
        />
        {#if query}
          <button
            onclick={() => { query = ''; if (inputEl) inputEl.focus(); }}
            class="p-1 rounded-md text-[var(--text-muted)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-colors"
          >
            <X size={14} />
          </button>
        {/if}
        <kbd class="px-1.5 py-0.5 rounded text-[10px] font-mono font-semibold bg-[var(--bg-surface-active)] text-[var(--text-muted)] border border-[var(--border-subtle)]">
          ESC
        </kbd>
      </div>

      <!-- Results List Area -->
      <div
        bind:this={listContainerEl}
        class="flex-1 overflow-y-auto p-2 space-y-1 divide-y divide-[var(--border-subtle)]/30"
      >
        {#if flattenedItems.length === 0}
          <div class="py-12 flex flex-col items-center justify-center text-[var(--text-muted)] space-y-2">
            <Sparkles size={24} strokeWidth={1.5} class="text-[var(--text-muted)]" />
            <p class="text-[13px] font-medium">No results found for "{query}"</p>
          </div>
        {:else}
          {#each flattenedItems as item, idx}
            <button
              type="button"
              data-index={idx}
              onclick={() => item.run()}
              onmouseenter={() => (selectedIndex = idx)}
              class="w-full text-left px-3.5 py-2.5 rounded-xl flex items-center justify-between transition-all pressable {selectedIndex === idx ? 'bg-[var(--accent-soft)] text-[var(--text-primary)] border border-[var(--accent-primary)]/30' : 'text-[var(--text-secondary)] hover:bg-[var(--bg-surface-hover)] border border-transparent'}"
            >
              <div class="flex items-center space-x-3 truncate">
                <div class="p-1.5 rounded-lg {selectedIndex === idx ? 'bg-[var(--accent-primary)] text-white' : 'bg-[var(--bg-surface-active)] text-[var(--text-secondary)]'} transition-colors shrink-0">
                  <item.icon size={15} strokeWidth={2} />
                </div>
                <div class="truncate">
                  <div class="flex items-center space-x-2">
                    <span class="text-[13px] font-semibold tracking-tight text-[var(--text-primary)] truncate">
                      {item.title}
                    </span>
                    {#if item.category === 'action'}
                      <span class="px-1.5 py-0.2 rounded text-[9.5px] font-mono uppercase tracking-wider bg-[var(--bg-surface-active)] text-[var(--text-muted)] border border-[var(--border-subtle)]">
                        Action
                      </span>
                    {/if}
                  </div>
                  {#if item.subtitle}
                    <p class="text-[11px] text-[var(--text-muted)] truncate mt-0.5">
                      {item.subtitle}
                    </p>
                  {/if}
                </div>
              </div>

              {#if selectedIndex === idx}
                <div class="flex items-center space-x-1 text-[var(--accent-primary)] shrink-0 ml-2 animate-scale-in">
                  <span class="text-[10px] font-mono">Select</span>
                  <CornerDownLeft size={12} strokeWidth={2.5} />
                </div>
              {:else if item.tags && item.tags.length > 0}
                <div class="hidden sm:flex items-center space-x-1 shrink-0 ml-2">
                  {#each item.tags.slice(0, 2) as tag}
                    <span class="text-[10px] font-mono text-[var(--text-muted)] bg-[var(--bg-surface-active)] px-1.5 py-0.5 rounded border border-[var(--border-subtle)]">
                      #{tag}
                    </span>
                  {/each}
                </div>
              {/if}
            </button>
          {/each}
        {/if}
      </div>

      <!-- Bottom Keyboard Navigation Hints -->
      <div class="h-10 px-4 bg-[var(--bg-surface-active)]/80 border-t border-[var(--border-subtle)] flex items-center justify-between text-[11px] text-[var(--text-muted)] font-mono shrink-0">
        <div class="flex items-center space-x-3">
          <span class="flex items-center space-x-1">
            <kbd class="px-1 py-0.2 rounded bg-[var(--bg-surface-hover)] border border-[var(--border-subtle)] text-[10px]">↑</kbd>
            <kbd class="px-1 py-0.2 rounded bg-[var(--bg-surface-hover)] border border-[var(--border-subtle)] text-[10px]">↓</kbd>
            <span>Navigate</span>
          </span>
          <span class="flex items-center space-x-1">
            <kbd class="px-1.5 py-0.2 rounded bg-[var(--bg-surface-hover)] border border-[var(--border-subtle)] text-[10px]">↵</kbd>
            <span>Select</span>
          </span>
          <span class="flex items-center space-x-1">
            <kbd class="px-1.5 py-0.2 rounded bg-[var(--bg-surface-hover)] border border-[var(--border-subtle)] text-[10px]">ESC</kbd>
            <span>Close</span>
          </span>
        </div>
        <div>
          <span>{flattenedItems.length} items</span>
        </div>
      </div>
    </div>
  </div>
{/if}
