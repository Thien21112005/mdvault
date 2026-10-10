<script lang="ts">
  import { CircleCheck, CircleAlert, TriangleAlert, Info, X } from '@lucide/svelte'
  import { toast } from '$lib/stores/toastStore'
</script>

{#if $toast.length > 0}
  <div class="fixed bottom-5 right-5 z-50 flex flex-col space-y-2.5 pointer-events-none max-w-sm w-full select-none">
    {#each $toast as item (item.id)}
      <div
        class="pointer-events-auto flex items-start justify-between space-x-3 p-3.5 rounded-2xl bg-[var(--bg-surface-elevated)]/90 backdrop-blur-xl border border-[var(--border-subtle)] shadow-[var(--shadow-glass)] transition-all animate-scale-in"
      >
        <div class="flex items-start space-x-3">
          <div class="mt-0.5 shrink-0">
            {#if item.type === 'success'}
              <CircleCheck size={16} class="text-emerald-500" strokeWidth={2.2} />
            {:else if item.type === 'error'}
              <CircleAlert size={16} class="text-rose-500" strokeWidth={2.2} />
            {:else if item.type === 'warning'}
              <TriangleAlert size={16} class="text-amber-500" strokeWidth={2.2} />
            {:else}
              <Info size={16} class="text-indigo-400" strokeWidth={2.2} />
            {/if}
          </div>

          <div class="flex-1">
            <h4 class="text-[12px] font-semibold text-[var(--text-primary)] leading-tight tracking-tight">
              {item.title}
            </h4>
            {#if item.description}
              <p class="text-[11px] text-[var(--text-muted)] mt-1 leading-snug">
                {item.description}
              </p>
            {/if}
          </div>
        </div>

        <button
          onclick={() => toast.remove(item.id)}
          class="text-[var(--text-muted)] hover:text-[var(--text-primary)] p-1 rounded-lg hover:bg-[var(--bg-surface-hover)] transition-colors pressable shrink-0"
          title="Dismiss"
          aria-label="Dismiss notification"
        >
          <X size={13} />
        </button>
      </div>
    {/each}
  </div>
{/if}
