<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import Navbar from './components/layout/Navbar.svelte'
  import Sidebar from './components/sidebar/Sidebar.svelte'
  import NotesList from './components/note/NotesList.svelte'
  import NoteEditor from './components/note/NoteEditor.svelte'
  import GraphView from './components/graph/GraphView.svelte'
  import CommandPalette from './components/layout/CommandPalette.svelte'
  import Toast from './components/ui/Toast.svelte'
  import { loadInitialWorkspace, initWatcher } from '$lib/stores/vaultStore'
  import { viewMode, isSidebarOpen, toggleSidebar, switchView, openCommandPalette } from '$lib/stores/viewStore'
  import { initTheme } from '$lib/stores/themeStore'

  function handleKeydown(e: KeyboardEvent) {
    const isMeta = e.ctrlKey || e.metaKey
    if (isMeta && e.key.toLowerCase() === 'b') {
      e.preventDefault()
      toggleSidebar()
    } else if (isMeta && e.key === '1') {
      e.preventDefault()
      switchView('list')
    } else if (isMeta && e.key === '2') {
      e.preventDefault()
      switchView('graph')
    } else if (isMeta && e.key.toLowerCase() === 'k') {
      e.preventDefault()
      openCommandPalette()
    }
  }

  onMount(async () => {
    initTheme()
    window.addEventListener('keydown', handleKeydown)
    await loadInitialWorkspace()
    initWatcher()
  })

  onDestroy(() => {
    if (typeof window !== 'undefined') {
      window.removeEventListener('keydown', handleKeydown)
    }
  })
</script>

<main class="w-screen h-screen bg-[var(--bg-canvas)] text-[var(--text-primary)] flex flex-col overflow-hidden font-sans select-none transition-colors duration-200">
  <Navbar />
  <div class="flex-1 flex overflow-hidden relative">
    <!-- Collapsible Sidebar with smooth Apple Spring transition -->
    <div
      class="transition-[width,opacity] duration-300 ease-[var(--ease-apple)] flex-shrink-0 overflow-hidden {$isSidebarOpen ? 'w-64 opacity-100' : 'w-0 opacity-0'}"
    >
      <Sidebar />
    </div>

    <!-- Main Content Workspace -->
    {#if $viewMode === 'list'}
      <NotesList />
      <NoteEditor />
    {:else}
      <GraphView />
    {/if}
  </div>

  <!-- Global Command Palette Modal & Sonner Toast Notifications -->
  <CommandPalette />
  <Toast />
</main>
