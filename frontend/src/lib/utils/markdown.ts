import MarkdownIt from 'markdown-it'

const md = new MarkdownIt({
  html: true,
  breaks: true,
  linkify: true,
  typographer: true,
})

export function renderMarkdown(content: string): string {
  if (!content) return '<p class="text-slate-400 italic">Ghi chú trống...</p>'
  return md.render(content)
}
