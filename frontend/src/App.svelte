<script lang="ts">
  import { onMount } from 'svelte'
  import Navbar from './components/layout/Navbar.svelte'
  import Sidebar from './components/sidebar/Sidebar.svelte'
  import NotesList from './components/note/NotesList.svelte'
  import NoteEditor from './components/note/NoteEditor.svelte'
  import GraphView from './components/graph/GraphView.svelte'
  import { loadInitialWorkspace, initWatcher } from '$lib/stores/vaultStore'
  import { viewMode } from '$lib/stores/viewStore'
  import { isDark } from '$lib/stores/themeStore'

  onMount(async () => {
    isDark.set(true)
    document.documentElement.classList.add('dark')
    
    await loadInitialWorkspace()
    initWatcher()
  })
</script>

<main class="w-screen h-screen bg-zinc-950 text-zinc-100 flex flex-col overflow-hidden font-sans">
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
