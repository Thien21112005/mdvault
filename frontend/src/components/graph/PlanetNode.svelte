<script lang="ts">
  import { Handle, Position, type NodeProps } from '@xyflow/svelte'

  let { data }: NodeProps = $props()

  // Fallback defaults
  const size = $derived((data?.size as number) || 40)
  const color = $derived((data?.color as string) || '#6366f1')
  const glowColor = $derived((data?.glowColor as string) || 'rgba(99, 102, 241, 0.4)')
  const connections = $derived((data?.connections as number) || 0)
  const label = $derived((data?.label as string) || 'Untitled Note')
</script>

<div class="planet-wrapper relative flex flex-col items-center justify-center cursor-pointer select-none group">
  <!-- Invisible Center Handles for Edges -->
  <Handle
    type="target"
    position={Position.Top}
    class="!top-1/2 !left-1/2 !-translate-x-1/2 !-translate-y-1/2 !w-1 !h-1 !opacity-0 !border-0 !bg-transparent pointer-events-none"
  />
  <Handle
    type="source"
    position={Position.Bottom}
    class="!top-1/2 !left-1/2 !-translate-x-1/2 !-translate-y-1/2 !w-1 !h-1 !opacity-0 !border-0 !bg-transparent pointer-events-none"
  />

  <!-- Planet Sphere Core -->
  <div
    class="planet-core rounded-full relative flex items-center justify-center transition-all duration-300 ease-out group-hover:scale-125"
    style="
      width: {size}px;
      height: {size}px;
      background: radial-gradient(circle at 35% 35%, #ffffff 0%, {color} 55%, rgba(0, 0, 0, 0.8) 100%);
      box-shadow: 0 0 {size * 0.45}px {glowColor}, inset 0 0 {size * 0.25}px rgba(255, 255, 255, 0.6);
    "
  >
    <!-- Soft pulsing aura for hub nodes -->
    {#if connections > 1}
      <div
        class="absolute -inset-1.5 rounded-full animate-ping opacity-25 pointer-events-none"
        style="background: {color}; animation-duration: 3s;"
      ></div>
    {/if}

    <!-- Specular highlight reflection -->
    <div
      class="absolute top-[18%] left-[22%] w-[28%] h-[28%] rounded-full bg-white/60 blur-[0.6px] pointer-events-none"
    ></div>

    <!-- Mini connection count badge on hover or if hub -->
    {#if connections > 0}
      <span
        class="absolute -top-1.5 -right-1.5 min-w-[17px] h-[17px] px-1 rounded-full text-[9px] font-mono font-bold bg-zinc-950/90 text-zinc-100 border border-white/20 flex items-center justify-center shadow-sm"
      >
        {connections}
      </span>
    {/if}
  </div>

  <!-- Planet Label Pill -->
  <div
    class="mt-2 px-2 py-0.5 rounded-md bg-[var(--bg-surface-elevated)]/90 backdrop-blur-md border border-[var(--border-subtle)] shadow-sm max-w-[130px] transition-all duration-200 group-hover:border-[var(--accent-primary)] group-hover:shadow-[var(--shadow-glass)] group-hover:scale-105 pointer-events-none"
  >
    <p class="text-[10.5px] font-semibold text-[var(--text-primary)] truncate text-center tracking-tight">
      {label}
    </p>
  </div>
</div>
