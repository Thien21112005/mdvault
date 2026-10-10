<script lang="ts">
  import { onMount } from 'svelte'
  import Navbar from './components/layout/Navbar.svelte'
  import Sidebar from './components/sidebar/Sidebar.svelte'
  import NotesList from './components/note/NotesList.svelte'
  import NoteEditor from './components/note/NoteEditor.svelte'
  import GraphView from './components/graph/GraphView.svelte'
  import { loadInitialWorkspace, initWatcher } from '$lib/stores/vaultStore'
  import { viewMode } from '$lib/stores/viewStore'
  import { initTheme } from '$lib/stores/themeStore'

  onMount(async () => {
    initTheme()
    await loadInitialWorkspace()
    initWatcher()
  })
</script>

<main class="w-screen h-screen bg-[var(--bg-canvas)] text-[var(--text-primary)] flex flex-col overflow-hidden font-sans select-none transition-colors duration-200">
  <Navbar />
  <div class="flex-1 flex overflow-hidden">
    <Sidebar />
    {#if $viewMode === 'list'}
      <NotesList />
      <NoteEditor />
    {:else}
      <GraphView />
    {/if}
  </div>
</main>
