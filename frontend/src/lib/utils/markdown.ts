import MarkdownIt from 'markdown-it'
import DOMPurify from 'dompurify'

const md = new MarkdownIt({
  html: true,
  breaks: true,
  linkify: true,
  typographer: true,
})

// Custom code fence renderer with Top Bar, Language Badge, and Copy Button
md.renderer.rules.fence = (tokens, idx) => {
  const token = tokens[idx]
  const info = token.info ? token.info.trim() : ''
  const lang = info ? info.split(/\s+/)[0] : 'code'
  const rawCode = token.content
  const encoded = encodeURIComponent(rawCode)

  return `<div class="code-block-wrapper my-5 rounded-xl overflow-hidden border border-[var(--border-subtle)] bg-[var(--bg-surface-elevated)] shadow-[var(--shadow-glass)] select-text">
    <div class="code-block-header flex items-center justify-between px-3.5 py-1.5 bg-[var(--bg-surface-active)] border-b border-[var(--border-subtle)] text-[11px] font-mono text-[var(--text-muted)] select-none">
      <span class="font-semibold uppercase tracking-wider text-[10px] text-[var(--text-secondary)]">${lang}</span>
      <button class="copy-code-btn px-2.5 py-0.5 rounded-md text-[10.5px] font-sans font-medium text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] border border-transparent hover:border-[var(--border-subtle)] transition-all pressable cursor-pointer" data-code="${encoded}">
        Copy
      </button>
    </div>
    <pre class="p-4 m-0 overflow-x-auto text-[12.5px] leading-relaxed font-mono text-[var(--text-primary)] bg-transparent"><code>${md.utils.escapeHtml(rawCode)}</code></pre>
  </div>`
}

/**
 * Pre-process Markdown content to convert:
 * 1. [[Wiki-Links]] to clickable anchor tags
 * 2. Task list items (- [ ] and - [x]) to styled checkboxes
 */
function preprocessMarkdown(content: string): string {
  if (!content) return ''

  // 1. Convert [[Wiki-Link]] to markdown link with specific class
  let processed = content.replace(/\[\[(.*?)\]\]/g, (_, title) => {
    const cleanTitle = title.trim()
    return `<a href="#wiki:${encodeURIComponent(cleanTitle)}" class="wiki-link inline-flex items-center space-x-1 px-1.5 py-0.2 rounded font-medium bg-[var(--accent-soft)] text-[var(--accent-primary)] hover:underline border border-[var(--accent-primary)]/20" data-wiki="${cleanTitle}">[[${cleanTitle}]]</a>`
  })

  // 2. Convert task lists
  processed = processed.replace(/^([ \t]*)[-*][ \t]+\[ \][ \t]+(.*)$/gm, '$1- <input type="checkbox" disabled class="task-checkbox mr-2 align-middle accent-[var(--accent-primary)] rounded cursor-not-allowed"> $2')
  processed = processed.replace(/^([ \t]*)[-*][ \t]+\[[xX]\][ \t]+(.*)$/gm, '$1- <input type="checkbox" checked disabled class="task-checkbox mr-2 align-middle accent-[var(--accent-primary)] rounded cursor-not-allowed"> <del class="text-[var(--text-muted)]">$2</del>')

  return processed
}

/**
 * Render Markdown content securely:
 * 1. Preprocesses wiki-links and task lists
 * 2. Parses with MarkdownIt
 * 3. Sanitizes the resulting HTML with DOMPurify to eliminate all XSS vectors
 */
export function renderMarkdown(content: string): string {
  if (!content || !content.trim()) {
    return '<p class="text-[var(--text-muted)] italic text-[13px] py-4">Ghi chú trống. Nhấn Edit hoặc gõ nội dung để bắt đầu...</p>'
  }

  const preprocessed = preprocessMarkdown(content)
  const rawHtml = md.render(preprocessed)

  // Strict DOMPurify Sanitization against XSS in WebView2
  return DOMPurify.sanitize(rawHtml, {
    USE_PROFILES: { html: true },
    ADD_TAGS: ['input'],
    ADD_ATTR: ['target', 'data-code', 'data-wiki', 'checked', 'disabled', 'type'],
    FORBID_TAGS: ['script', 'iframe', 'object', 'embed', 'form', 'base'],
    FORBID_ATTR: ['onerror', 'onload', 'onclick', 'onmouseover', 'onfocus', 'onblur', 'style']
  })
}

/**
 * Safely highlight search keyword matches in text strings (Title & Excerpt).
 * Escapes raw text against XSS before wrapping matched substrings with an Apple-style <mark>.
 */
export function highlightMatch(text: string, query: string): string {
  if (!text) return ''
  if (!query || !query.trim()) {
    return md.utils.escapeHtml(text)
  }

  const escapedText = md.utils.escapeHtml(text)
  const escapedQuery = query.trim().replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const regex = new RegExp(`(${escapedQuery})`, 'gi')

  return escapedText.replace(
    regex,
    '<mark class="bg-amber-400/30 text-amber-700 dark:text-amber-300 rounded-[3px] px-0.5 font-semibold">$1</mark>'
  )
}

