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

  static onFileChange(callback: (changedPath: string) => void): () => void {
    EventsOn('on-file-change', callback)
    return () => EventsOff('on-file-change')
  }
}
