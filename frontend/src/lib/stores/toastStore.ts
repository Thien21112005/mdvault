import { writable } from 'svelte/store'

export type ToastType = 'success' | 'info' | 'warning' | 'error'

export interface ToastAction {
  label: string
  onClick: () => void
}

export interface ToastOptions {
  description?: string
  duration?: number
  action?: ToastAction
}

export interface ToastItem {
  id: string
  type: ToastType
  title: string
  description?: string
  duration: number
  createdAt: number
  action?: ToastAction
}

function createToastStore() {
  const { subscribe, update } = writable<ToastItem[]>([])
  const timers = new Map<string, ReturnType<typeof setTimeout>>()

  function add(item: Omit<ToastItem, 'id' | 'createdAt'>) {
    const id = Math.random().toString(36).substring(2, 9)
    const duration = item.duration ?? 4200
    const newToast: ToastItem = { ...item, id, duration, createdAt: Date.now() }

    update(list => {
      const next = [...list, newToast]
      if (next.length > 4) {
        const excess = next.slice(0, next.length - 4)
        excess.forEach(t => {
          const tmr = timers.get(t.id)
          if (tmr) {
            clearTimeout(tmr)
            timers.delete(t.id)
          }
        })
        return next.slice(-4)
      }
      return next
    })

    if (duration > 0) {
      const timer = setTimeout(() => {
        remove(id)
      }, duration)
      timers.set(id, timer)
    }
    return id
  }

  function remove(id: string) {
    const timer = timers.get(id)
    if (timer) {
      clearTimeout(timer)
      timers.delete(id)
    }
    update(list => list.filter(t => t.id !== id))
  }

  function pause(id: string) {
    const timer = timers.get(id)
    if (timer) {
      clearTimeout(timer)
      timers.delete(id)
    }
  }

  function resume(id: string, remainingMs: number = 2500) {
    if (timers.has(id)) return
    const timer = setTimeout(() => {
      remove(id)
    }, remainingMs)
    timers.set(id, timer)
  }

  function dispatch(
    type: ToastType,
    title: string,
    descOrOptions?: string | ToastOptions,
    extraOptions?: ToastOptions
  ) {
    let description: string | undefined
    let options: ToastOptions | undefined

    if (typeof descOrOptions === 'string') {
      description = descOrOptions
      options = extraOptions
    } else if (typeof descOrOptions === 'object' && descOrOptions !== null) {
      options = descOrOptions
      description = options.description
    }

    return add({
      type,
      title,
      description,
      duration: options?.duration ?? 4200,
      action: options?.action
    })
  }

  return {
    subscribe,
    add,
    remove,
    pause,
    resume,
    success: (title: string, descOrOptions?: string | ToastOptions, options?: ToastOptions) =>
      dispatch('success', title, descOrOptions, options),
    info: (title: string, descOrOptions?: string | ToastOptions, options?: ToastOptions) =>
      dispatch('info', title, descOrOptions, options),
    warning: (title: string, descOrOptions?: string | ToastOptions, options?: ToastOptions) =>
      dispatch('warning', title, descOrOptions, options),
    error: (title: string, descOrOptions?: string | ToastOptions, options?: ToastOptions) =>
      dispatch('error', title, descOrOptions, options)
  }
}

export const toast = createToastStore()
