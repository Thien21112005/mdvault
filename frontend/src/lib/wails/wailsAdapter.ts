// Tầng abstraction cho giao tiếp IPC giữa Svelte và Wails Go Backend
import * as VaultBindings from '../../../wailsjs/go/controller/VaultController'
import { EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime'
import type { WorkspaceData, Note } from '../types/domain'

export class WailsAdapter {
  static async selectWorkspace(): Promise<WorkspaceData | null> {
    try {
      return (await VaultBindings.SelectWorkspace()) as unknown as WorkspaceData
    } catch (e) {
      console.error('Lỗi khi mở workspace dialog:', e)
      return null
    }
  }

  static async scanFolder(): Promise<WorkspaceData | null> {
    try {
      return (await VaultBindings.ScanFolder()) as unknown as WorkspaceData
    } catch (e) {
      console.log('Chưa có workspace kích hoạt hoặc lỗi scan:', e)
      return null
    }
  }

  static async getNoteContent(path: string): Promise<Note | null> {
    try {
      return (await VaultBindings.GetNoteContent(path)) as unknown as Note
    } catch (e) {
      console.error('Lỗi khi tải chi tiết ghi chú:', e)
      return null
    }
  }

  static async saveNoteContent(id: string, content: string, title: string, tags: string[]): Promise<boolean> {
    try {
      await VaultBindings.SaveNoteContent(id, content, title, tags)
      return true
    } catch (e) {
      console.error('Lỗi khi lưu ghi chú:', e)
      return false
    }
  }

  static async createFolder(parentPath: string, folderName: string): Promise<boolean> {
    try {
      await VaultBindings.CreateFolder(parentPath, folderName)
      return true
    } catch (e) {
      console.error('Lỗi khi tạo thư mục:', e)
      return false
    }
  }

  static async createNote(folderPath: string, noteTitle: string): Promise<Note | null> {
    try {
      return (await VaultBindings.CreateNote(folderPath, noteTitle)) as unknown as Note
    } catch (e) {
      console.error('Lỗi khi tạo ghi chú:', e)
      return null
    }
  }

  static async searchNotes(query: string): Promise<Note[]> {
    try {
      return (await VaultBindings.SearchNotes(query)) as unknown as Note[]
    } catch (e) {
      console.error('Lỗi tìm kiếm:', e)
      return []
    }
  }

  static async deleteNote(filePath: string): Promise<boolean> {
    try {
      const bindings = VaultBindings as any
      if (typeof bindings.DeleteNote === 'function') {
        await bindings.DeleteNote(filePath)
      } else if (typeof (window as any).go?.controller?.VaultController?.DeleteNote === 'function') {
        await (window as any).go.controller.VaultController.DeleteNote(filePath)
      } else {
        console.warn('DeleteNote RPC method not found in bindings or window.go')
        return false
      }
      return true
    } catch (e) {
      console.error('Lỗi khi xóa ghi chú:', e)
      return false
    }
  }

  static async deleteFolder(folderPath: string): Promise<boolean> {
    try {
      const bindings = VaultBindings as any
      if (typeof bindings.DeleteFolder === 'function') {
        await bindings.DeleteFolder(folderPath)
      } else if (typeof (window as any).go?.controller?.VaultController?.DeleteFolder === 'function') {
        await (window as any).go.controller.VaultController.DeleteFolder(folderPath)
      } else {
        console.warn('DeleteFolder RPC method not found in bindings or window.go')
        return false
      }
      return true
    } catch (e) {
      console.error('Lỗi khi xóa thư mục:', e)
      return false
    }
  }

  static async renameNote(oldPath: string, newTitle: string): Promise<Note | null> {
    try {
      const bindings = VaultBindings as any
      if (typeof bindings.RenameNote === 'function') {
        return (await bindings.RenameNote(oldPath, newTitle)) as unknown as Note
      } else if (typeof (window as any).go?.controller?.VaultController?.RenameNote === 'function') {
        return (await (window as any).go.controller.VaultController.RenameNote(oldPath, newTitle)) as unknown as Note
      } else {
        console.warn('RenameNote RPC method not found in bindings or window.go')
        return null
      }
    } catch (e) {
      console.error('Lỗi khi đổi tên ghi chú:', e)
      return null
    }
  }

  static async renameFolder(oldPath: string, newName: string): Promise<boolean> {
    try {
      const bindings = VaultBindings as any
      if (typeof bindings.RenameFolder === 'function') {
        await bindings.RenameFolder(oldPath, newName)
      } else if (typeof (window as any).go?.controller?.VaultController?.RenameFolder === 'function') {
        await (window as any).go.controller.VaultController.RenameFolder(oldPath, newName)
      } else {
        console.warn('RenameFolder RPC method not found in bindings or window.go')
        return false
      }
      return true
    } catch (e) {
      console.error('Lỗi khi đổi tên thư mục:', e)
      return false
    }
  }

  static onFileChange(callback: (changedPath: string) => void): () => void {
    EventsOn('on-file-change', callback)
    return () => EventsOff('on-file-change')
  }
}
