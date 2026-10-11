# Kế Hoạch Hoàn Thiện Frontend, Bổ Sung CRUD UI & Đồng Bộ Đặc Tả Kỹ Thuật (Frontend CRUD, UI Polish & Spec Alignment Master Plan)

> **Điều phối bởi:** `mdvault-orchestrator` (Skill Mẹ)  
> **Các kỹ năng chuyên môn tham gia:**  
> - `planning-and-task-breakdown` (Lập lộ trình, phân rã 8 lát cắt mỏng, tiêu chí nghiệm thu)  
> - `impeccable` & `apple-design` (Thiết kế StatusBar glassmorphic, Modal/Menu CRUD tinh tế, phím tắt)  
> - `emil-design-eng` & `animate` (Micro-animations xác nhận xóa, transitions mượt mà)  
> - `security-and-hardening` (XSS sanitization khi highlight từ khóa tìm kiếm)  
> - `incremental-implementation` & `code-review` (Triển khai từng lát cắt & thẩm định 2 trục)  
>  
> **Phiên bản:** 2.2.0-PRO  
> **Mục tiêu:** Khắc phục toàn bộ 7 khoảng trống kỹ thuật phát hiện sau đợt kiểm toán hệ thống: bổ sung trọn vẹn UI tương tác CRUD (Xóa/Đổi tên Notes & Folders), xây dựng StatusBar chân trang theo PRB, hiển thị icon Frontmatter, highlight từ khóa tìm kiếm, chống race condition mất dữ liệu khi đang gõ (`isDirty` guard), dọn dẹp code thừa và đồng bộ hóa 100% tài liệu kỹ thuật trong `markdown/LLD.md`.

---

## 1. Phân Tích Thực Trạng & Mục Tiêu Kỹ Thuật (Executive Context)

Đợt kiểm toán toàn diện ngày 11/10/2026 đã chỉ ra các điểm nghẽn giữa tài liệu đặc tả ([PRB.md](../../../markdown/PRB.md), [SRS.md](../../../markdown/SRS.md), [LLD.md](../../../markdown/LLD.md)) và mã nguồn hiện tại:

| Mã Khoảng Trống | Hạng Mục | Thực Trạng Code Hiện Tại | Hậu Quả & Rủi Ro | Giải Pháp Kiến Trúc |
| :---: | :--- | :--- | :--- | :--- |
| **GAP-01** | **Bảo vệ Buffer Khi Gõ** | `$effect` trong `NoteEditor.svelte` tự động ghi đè `editContent` mỗi khi `loadInitialWorkspace` chạy lại do watcher. | Nếu file khác thay đổi hoặc watcher kích hoạt khi đang soạn thảo, nội dung chưa lưu sẽ bị xóa sạch. | Triển khai cờ phản ứng `isDirty` & `activeEditingId`; chỉ đồng bộ từ đĩa khi không có thay đổi chưa lưu, hoặc hiển thị toast cảnh báo. |
| **GAP-02** | **Bộ API CRUD Phía UI** | Backend Go đã có sẵn `DeleteNote`, `DeleteFolder`, `RenameNote`, `RenameFolder` nhưng `wailsAdapter.ts` và `vaultStore.ts` chưa wrap. | Người dùng hoàn toàn không có cách nào Xóa hoặc Đổi tên ghi chú/thư mục từ giao diện. | Mở rộng `WailsAdapter` và `vaultStore` với đầy đủ 4 actions CRUD; tích hợp vào `NotesList`, `NoteEditor`, `FolderTreeItem`, `Sidebar` và `CommandPalette`. |
| **GAP-03** | **Status Bar Chân Trang** | PRB.md mô tả wireframe có thanh đáy hiển thị trạng thái `fsnotify`, dung lượng, tiến trình, nhưng UI chưa có. | Thiếu phản hồi trực quan về tình trạng đồng bộ thời gian thực của hệ điều hành. | Xây dựng component `StatusBar.svelte` chuẩn Glassmorphism (h-7/28px) hiển thị trạng thái Watcher, số lượng từ/ký tự, đường dẫn workspace. |
| **GAP-04** | **Icon Frontmatter** | YAML Frontmatter hỗ trợ `icon` (đã trích xuất trong `Metadata`) nhưng UI chỉ dùng icon mặc định. | Mất đi trải nghiệm nhận diện ghi chú sinh động đã cam kết trong SRS FR-02 và PRB US-02. | Xây dựng dynamic icon renderer (hỗ trợ các biểu tượng Lucide phổ biến và emoji) trên `NotesList` và `NoteEditor`. |
| **GAP-05** | **Highlight Từ Khóa Search** | SRS FR-05 quy định làm nổi bật từ khóa tìm kiếm, nhưng UI chỉ lọc danh sách phẳng. | Người dùng khó nhận biết vị trí xuất hiện của từ khóa trong tiêu đề và đoạn trích. | Thêm bộ xử lý highlight chuỗi an toàn (tránh XSS) bằng thẻ `<mark>` được tạo kiểu tinh tế. |
| **GAP-06** | **Lỗi Thời Tài Liệu LLD.md** | `markdown/LLD.md` chưa ghi nhận trường `Metadata`/`Extra`, thiếu 4 API CRUD mới, và mô tả thuật toán Graph cũ. | Kỹ sư đọc tài liệu sẽ hiểu sai API contract thực tế của Backend. | Cập nhật `markdown/LLD.md` đồng bộ 100% với các struct và phương thức Go hiện có. |
| **GAP-07** | **Mã Nguồn Rác Tồn Dư** | Tệp `frontend/src/lib/Counter.svelte` là file mẫu thừa từ template Vite. | Tăng kích thước dự án, giảm tính chuyên nghiệp của codebase. | Xóa bỏ file `Counter.svelte`. |

---

## 2. Thiết Kế Kiến Trúc & Quyết Định Kỹ Thuật (Architecture Decisions)

### 2.1. Cơ chế `isDirty` Guard Chống Ghi Đè Buffer Trong Svelte 5
```typescript
// Trong NoteEditor.svelte
let isDirty = $state(false)
let lastLoadedNoteId = $state<string | null>(null)

$effect(() => {
  const currentNote = $selectedNote
  if (!currentNote) return

  // Nếu chuyển hẳn sang mở một note khác -> Luôn nạp mới
  if (currentNote.id !== lastLoadedNoteId) {
    editTitle = currentNote.title
    editTagsList = [...(currentNote.tags || [])]
    editContent = currentNote.content
    lastLoadedNoteId = currentNote.id
    isDirty = false
    return
  }

  // Nếu cùng một note nhưng nội dung đĩa thay đổi từ bên ngoài khi người dùng CHƯA gõ gì
  if (!isDirty && currentNote.content !== editContent) {
    editContent = currentNote.content
    editTitle = currentNote.title
  }
})
```
- Khi người dùng gõ vào textarea, đánh dấu `isDirty = true`.
- Khi bấm Save (Ctrl+S), sau khi lưu thành công chuyển `isDirty = false`.
- Loại bỏ hoàn toàn xung đột race condition giữa File Watcher và bàn phím người dùng.

### 2.2. Lớp Chuyển Đổi An Toàn (Dynamic RPC Binding Wrapper)
Do các file `wailsjs` sinh tự động chưa được cập nhật qua lệnh build, `WailsAdapter` cung cấp cơ chế gọi an toàn hai tầng:
```typescript
static async deleteNote(filePath: string): Promise<boolean> {
  try {
    const bindings = VaultBindings as any
    if (typeof bindings.DeleteNote === 'function') {
      await bindings.DeleteNote(filePath)
    } else if (typeof (window as any).go?.controller?.VaultController?.DeleteNote === 'function') {
      await (window as any).go.controller.VaultController.DeleteNote(filePath)
    }
    return true
  } catch (e) {
    console.error('Lỗi khi xóa ghi chú:', e)
    return false
  }
}
```

### 2.3. Hộp Thoại Xác Nhận Thao Tác Nhạy Cảm (Apple-style Dialog / Custom Toast Confirmation)
- Khi thực hiện thao tác xóa Note hoặc Folder (thao tác không thể hoàn tác trên đĩa):
  - Tuyệt đối không dùng `window.confirm()` thô ráp.
  - Sử dụng Modal/Popover kính mờ tối giản (`backdrop-blur-md`, viền `border-rose-500/30`, nút bấm xác nhận màu đỏ `pressable`) hoặc Action Toast có nút "Confirm Delete" ngay tức thì.

### 2.4. Thiết Kế StatusBar Chân Trang Chuẩn Glassmorphic (Apple Design)
- Chiều cao cố định `h-7` (28px), đặt cố định ở đáy viewport.
- Nền `bg-[var(--bg-surface)]/80 backdrop-blur-xl border-t border-[var(--border-subtle)]`.
- Hiển thị 3 cụm thông tin:
  - **Trái:** Tình trạng kết nối Watcher (`● fsnotify active` với pulse animation xanh lục) + Tên Vault.
  - **Giữa:** Tình trạng lưu (`All changes saved` / `● Unsaved changes`).
  - **Phải:** Thông số ghi chú đang mở (`{wordsCount} words`, `{charsCount} chars`, `{linesCount} lines`, UTF-8).

---

## 3. Lộ Trình Triển Khai 8 Lát Cắt Mỏng (Vertical Slices Roadmap)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ GIAI ĐOẠN 1: AN TOÀN NHẬP LIỆU & CẦU NỐI CRUD (DATA INTEGRITY & RPC)         │
│ Task 1: isDirty Guard chống mất buffer trong NoteEditor & Dọn dẹp rác       │
│ Task 2: Wrap 4 API CRUD trong wailsAdapter.ts & vaultStore.ts               │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
┌──────────────────────────────────────▼──────────────────────────────────────┐
│ GIAI ĐOẠN 2: GIAO DIỆN TƯƠNG TÁC CRUD CHO NOTE & FOLDER (UI INTERACTION)     │
│ Task 3: Bổ sung Action Menu & Xác nhận Xóa/Đổi tên cho Notes                │
│ Task 4: Bổ sung Context Menu / Actions Xóa & Đổi tên Thư mục trong Sidebar  │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
┌──────────────────────────────────────▼──────────────────────────────────────┐
│ GIAI ĐOẠN 3: POLISH GIAO DIỆN CHUẨN ĐẶC TẢ PRB/SRS (LAYOUT & VISUAL POLISH) │
│ Task 5: Xây dựng Thanh Trạng Thái Chân Trang StatusBar.svelte               │
│ Task 6: Hiển thị Icon tùy biến Frontmatter & Highlight từ khóa tìm kiếm     │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
┌──────────────────────────────────────▼──────────────────────────────────────┐
│ GIAI ĐOẠN 4: ĐỒNG BỘ ĐẶC TẢ KỸ THUẬT & NGHIỆM THU (DOCS & REVIEW)            │
│ Task 7: Cập nhật tài liệu kỹ thuật LLD.md đồng bộ 100% với Codebase         │
│ Task 8: Thẩm định nghiệm thu 2 trục (Standards & Spec Review) & Build       │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Quy Ước Cam Kết Mã Nguồn (Conventional Commits Protocol)

| Task | Commit Type | Conventional Commit Message |
| :---: | :---: | :--- |
| **Task 1** | `fix` | `fix(editor): implement isDirty buffer guard against watcher reload and clean dead code` |
| **Task 2** | `feat` | `feat(store): expose full CRUD actions in wailsAdapter and vaultStore` |
| **Task 3** | `feat` | `feat(notes): add delete and rename actions with modal confirmation for notes` |
| **Task 4** | `feat` | `feat(sidebar): add folder context actions for rename and safe deletion` |
| **Task 5** | `feat` | `feat(layout): implement glassmorphic status bar with realtime watcher status` |
| **Task 6** | `feat` | `feat(ui): support dynamic frontmatter icons and search keyword highlighting` |
| **Task 7** | `docs` | `docs(spec): synchronize LLD specification with current models and CRUD APIs` |
| **Task 8** | `chore` | `chore(review): complete two-axis engineering audit and finalize frontend alignment` |

---

## 5. Đánh Giá Rủi Ro & Chiến Lược Giảm Thiểu (Risks & Mitigations)

| Rủi Ro | Mức Độ | Biện Pháp Giảm Thiểu |
| :--- | :---: | :--- |
| **Lỗi Xóa Thư Mục Nhầm:** Người dùng vô tình bấm xóa thư mục cha chứa nhiều file. | **Cao** | Bắt buộc hiển thị modal xác nhận yêu cầu gõ lại tên thư mục hoặc bấm nút xác nhận nguy hiểm có cảnh báo số lượng file bên trong. |
| **XSS Khi Highlight Từ Khóa:** Dùng `{@html}` để tô đậm từ khóa tìm kiếm có thể lộ lỗ hổng chèn mã. | **Cao** | Escape HTML trước khi chèn thẻ `<mark>`, sau đó bọc qua DOMPurify trước khi render. |
| **Lỗi TypeScript Khi Gọi Binding Mới:** `wailsjs` chưa có type cho 4 hàm mới. | **Trung bình** | Sử dụng type assertion hoặc fallback window invocation an toàn, đảm bảo `npm run check` đạt 0 lỗi. |
