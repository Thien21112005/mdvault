<script lang="ts">
  import {
    FileText,
    Briefcase,
    Cpu,
    Server,
    BookOpen,
    Code,
    Database,
    Rocket,
    Sparkles,
    Folder,
    Layers,
    Globe,
    Tag,
    Shield,
    Zap,
    Heart,
    CheckSquare,
    AlertTriangle,
    Info,
    Settings,
    User,
    type Icon
  } from '@lucide/svelte'

  let {
    icon = '',
    size = 14,
    strokeWidth = 2,
    class: className = ''
  }: {
    icon?: string | null
    size?: number
    strokeWidth?: number
    class?: string
  } = $props()

  const iconMap: Record<string, typeof FileText> = {
    briefcase: Briefcase,
    cpu: Cpu,
    server: Server,
    'file-text': FileText,
    file: FileText,
    document: FileText,
    book: BookOpen,
    'book-open': BookOpen,
    code: Code,
    terminal: Code,
    database: Database,
    rocket: Rocket,
    star: Sparkles,
    sparkles: Sparkles,
    folder: Folder,
    layers: Layers,
    globe: Globe,
    network: Globe,
    tag: Tag,
    hash: Tag,
    shield: Shield,
    lock: Shield,
    zap: Zap,
    bolt: Zap,
    heart: Heart,
    check: CheckSquare,
    checklist: CheckSquare,
    alert: AlertTriangle,
    warning: AlertTriangle,
    info: Info,
    settings: Settings,
    user: User,
    users: User
  }

  // Check if string is an emoji character
  function isEmoji(str: string): boolean {
    if (!str) return false
    const emojiRegex = /\p{Extended_Pictographic}/u
    return emojiRegex.test(str.trim())
  }

  const cleanIcon = $derived((icon || '').trim().toLowerCase())
  const LucideComponent = $derived(iconMap[cleanIcon] || null)
  const emojiChar = $derived(icon && isEmoji(icon) ? icon.trim() : null)
</script>

{#if emojiChar}
  <span class="inline-flex items-center justify-center leading-none select-none {className}" style="font-size: {size}px;">
    {emojiChar}
  </span>
{:else if LucideComponent}
  {@const IconCmp = LucideComponent}
  <IconCmp {size} {strokeWidth} class={className} />
{:else}
  <FileText {size} {strokeWidth} class={className} />
{/if}
