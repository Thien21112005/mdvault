import { writable } from 'svelte/store'

export type ToastType = 'success' | 'info' | 'warning' | 'error'

export interface ToastItem {
  id: string
  type: ToastType
  title: string
  description?: string
  duration?: number
}

function createToastStore() {
  const { subscribe, update } = writable<ToastItem[]>([])

  function add(item: Omit<ToastItem, 'id'>) {
    const id = Math.random().toString(36).substring(2, 9)
    const duration = item.duration ?? 3200
    const newToast: ToastItem = { ...item, id, duration }

    update(list => [...list, newToast])

    if (duration > 0) {
      setTimeout(() => {
        remove(id)
      }, duration)
    }
    return id
  }

  function remove(id: string) {
    update(list => list.filter(t => t.id !== id))
  }

  return {
    subscribe,
    add,
    remove,
    success: (title: string, description?: string) => add({ type: 'success', title, description }),
    info: (title: string, description?: string) => add({ type: 'info', title, description }),
    warning: (title: string, description?: string) => add({ type: 'warning', title, description }),
    error: (title: string, description?: string) => add({ type: 'error', title, description })
  }
}

export const toast = createToastStore()
