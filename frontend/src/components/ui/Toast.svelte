<script lang="ts">
  import { CircleCheck, CircleAlert, TriangleAlert, Info, X, ArrowRight } from '@lucide/svelte'
  import { backOut, cubicOut } from 'svelte/easing'
  import { toast, type ToastType, type ToastItem } from '$lib/stores/toastStore'

  let hoveredToastId = $state<string | null>(null)

  function handleMouseEnter(item: ToastItem) {
    hoveredToastId = item.id
    toast.pause(item.id)
  }

  function handleMouseLeave(item: ToastItem) {
    if (hoveredToastId === item.id) {
      hoveredToastId = null
      toast.resume(item.id)
    }
  }

  // Apple-grade Spring Entrance: floats up from translateY(28px) with subtle overshoot bounce
  function toastSpringIn(_node: Element, { duration = 360 }: { duration?: number } = {}) {
    return {
      duration,
      css: (t: number) => {
        const eased = backOut(t)
        const translateY = (1 - eased) * 28
        const scale = 0.94 + eased * 0.06
        const opacity = Math.min(1, t * 2)
        return `
          transform: translateY(${translateY}px) scale(${scale});
          opacity: ${opacity};
        `
      }
    }
  }

  // Pure horizontal slide-out dismissal: slides directly to the right without any vertical dropping
  function toastSlideOut(_node: Element, { duration = 220 }: { duration?: number } = {}) {
    return {
      duration,
      css: (t: number) => {
        const eased = cubicOut(t)
        const translateX = (1 - eased) * 72
        return `
          transform: translateX(${translateX}px);
          opacity: ${eased};
        `
      }
    }
  }

  function getTypeStyles(type: ToastType) {
    switch (type) {
      case 'success':
        return {
          glow: 'shadow-[0_16px_36px_-6px_rgba(16,185,129,0.22),0_4px_16px_rgba(0,0,0,0.12)] border-emerald-500/35 dark:border-emerald-400/40',
          badge: 'bg-emerald-500/15 text-emerald-500 dark:text-emerald-400 border border-emerald-500/25',
          progress: 'bg-gradient-to-r from-emerald-500 to-emerald-400',
          actionBtn: 'bg-emerald-600 hover:bg-emerald-500 text-white shadow-[0_0_12px_rgba(16,185,129,0.35)]'
        }
      case 'error':
        return {
          glow: 'shadow-[0_16px_36px_-6px_rgba(244,63,94,0.22),0_4px_16px_rgba(0,0,0,0.12)] border-rose-500/35 dark:border-rose-400/40',
          badge: 'bg-rose-500/15 text-rose-500 dark:text-rose-400 border border-rose-500/25',
          progress: 'bg-gradient-to-r from-rose-500 to-rose-400',
          actionBtn: 'bg-rose-600 hover:bg-rose-500 text-white shadow-[0_0_12px_rgba(244,63,94,0.35)]'
        }
      case 'warning':
        return {
          glow: 'shadow-[0_16px_36px_-6px_rgba(245,158,11,0.25),0_4px_16px_rgba(0,0,0,0.12)] border-amber-500/35 dark:border-amber-400/40',
          badge: 'bg-amber-500/15 text-amber-500 dark:text-amber-400 border border-amber-500/25',
          progress: 'bg-gradient-to-r from-amber-500 to-amber-400',
          actionBtn: 'bg-amber-600 hover:bg-amber-500 text-white shadow-[0_0_12px_rgba(245,158,11,0.35)]'
        }
      case 'info':
      default:
        return {
          glow: 'shadow-[0_16px_36px_-6px_rgba(99,102,241,0.22),0_4px_16px_rgba(0,0,0,0.12)] border-indigo-500/35 dark:border-indigo-400/40',
          badge: 'bg-indigo-500/15 text-indigo-400 border border-indigo-500/25',
          progress: 'bg-gradient-to-r from-indigo-500 to-indigo-400',
          actionBtn: 'bg-indigo-600 hover:bg-indigo-500 text-white shadow-[0_0_12px_rgba(99,102,241,0.35)]'
        }
    }
  }
</script>

<!-- Outer container ALWAYS remains in DOM to ensure single-toast exits and entrances never get cut off by premature unmounting -->
<div
  class="fixed bottom-6 right-6 z-50 flex flex-col space-y-3 pointer-events-none max-w-sm w-full select-none"
  role="region"
  aria-label="Notifications"
>
  {#each $toast as item (item.id)}
    {@const styles = getTypeStyles(item.type)}
    {@const isHovered = hoveredToastId === item.id}
    <div
      in:toastSpringIn
      out:toastSlideOut
      onmouseenter={() => handleMouseEnter(item)}
      onmouseleave={() => handleMouseLeave(item)}
      role="alert"
      class="pointer-events-auto relative overflow-hidden flex flex-col rounded-2xl bg-[var(--bg-surface-elevated)]/95 backdrop-blur-2xl border {styles.glow} transition-shadow duration-200"
    >
      <!-- Subtle Top Gloss Sheen -->
      <div class="absolute inset-x-0 top-0 h-1/2 bg-gradient-to-b from-white/10 to-transparent pointer-events-none"></div>

      <!-- Main Content Row -->
      <div class="relative p-4 flex items-start justify-between space-x-3.5">
        <div class="flex items-start space-x-3 flex-1 min-w-0">
          <!-- Icon Pill Badge -->
          <div class="w-8 h-8 rounded-xl flex items-center justify-center shrink-0 {styles.badge} shadow-sm">
            {#if item.type === 'success'}
              <CircleCheck size={18} strokeWidth={2.2} />
            {:else if item.type === 'error'}
              <CircleAlert size={18} strokeWidth={2.2} />
            {:else if item.type === 'warning'}
              <TriangleAlert size={18} strokeWidth={2.2} />
            {:else}
              <Info size={18} strokeWidth={2.2} />
            {/if}
          </div>

          <!-- Text Content -->
          <div class="flex-1 min-w-0 pt-0.5">
            <h4 class="text-[13px] font-bold text-[var(--text-primary)] leading-tight tracking-tight">
              {item.title}
            </h4>
            {#if item.description}
              <p class="text-[11.5px] text-[var(--text-secondary)] mt-1 leading-relaxed font-normal">
                {item.description}
              </p>
            {/if}
          </div>
        </div>

        <!-- Close Dismiss Button -->
        <button
          onclick={() => toast.remove(item.id)}
          class="text-[var(--text-muted)] hover:text-[var(--text-primary)] p-1.5 rounded-lg hover:bg-[var(--bg-surface-hover)] border border-transparent hover:border-[var(--border-subtle)] transition-colors pressable shrink-0 -mr-1 -mt-1 cursor-pointer"
          title="Dismiss notification"
          aria-label="Dismiss notification"
        >
          <X size={14} strokeWidth={2} />
        </button>
      </div>

      <!-- Optional Action CTA Button -->
      {#if item.action}
        <div class="relative px-4 pb-3.5 pt-0 flex items-center justify-end">
          <button
            onclick={() => {
              item.action?.onClick()
              toast.remove(item.id)
            }}
            class="inline-flex items-center space-x-1.5 px-3.5 py-1.5 rounded-xl text-[11.5px] font-semibold {styles.actionBtn} hover:brightness-110 active:scale-[0.97] transition-all pressable cursor-pointer"
          >
            <span>{item.action.label}</span>
            <ArrowRight size={13} strokeWidth={2.2} />
          </button>
        </div>
      {/if}

      <!-- Countdown Progress Bar -->
      {#if item.duration > 0}
        <div class="absolute bottom-0 left-0 right-0 h-[2.5px] bg-black/5 dark:bg-white/5 overflow-hidden">
          <div
            class="h-full {styles.progress}"
            style="animation: toast-timer-progress {item.duration}ms linear forwards; animation-play-state: {isHovered ? 'paused' : 'running'};"
          ></div>
        </div>
      {/if}
    </div>
  {/each}
</div>

<style>
  @keyframes toast-timer-progress {
    from {
      width: 100%;
    }
    to {
      width: 0%;
    }
  }
</style>
