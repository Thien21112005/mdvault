# Danh Sách Công Việc: Hoàn Thiện Frontend, Bổ Sung CRUD UI & Đồng Bộ Đặc Tả Kỹ Thuật (Todo Tasks)

> **Dự án:** MDVault Desktop Application  
> **Kế hoạch đợt việc:** `2-frontend-and-spec-alignment` (Ngày 2026-10-11)  
> **Bộ tiêu chuẩn:** IEEE 830 / PRB / SRS / Clean Architecture / Apple Design  
> **Điều phối bởi:** `mdvault-orchestrator`  

---

## Phase 1: An Toàn Nhập Liệu & Cầu Nối CRUD (Data Safety & Bridge)

### Task 1: isDirty Guard Chống Mất Buffer Trong NoteEditor & Dọn Dẹp Code Rác
**Description:** Khắc phục rủi ro mất dữ liệu khi đang gõ bằng cách bổ sung cờ `isDirty` và theo dõi `lastLoadedNoteId` trong `frontend/src/components/note/NoteEditor.svelte`. Khi watcher phát tín hiệu cập nhật tệp, nếu người dùng đang có thay đổi chưa lưu, hệ thống không tự ý ghi đè lên textarea mà hiển thị thông báo. Đồng thời xóa bỏ file rác `frontend/src/lib/Counter.svelte`.

**Acceptance criteria:**
- [x] Gõ nội dung trong `NoteEditor` chuyển trạng thái sang `isDirty = true`.
- [x] Khi có sự kiện `on-file-change` từ bên ngoài đối với cùng một file, nội dung đang gõ dở không bị đè mất.
- [x] Sau khi bấm Lưu (Ctrl+S), trạng thái trở về `isDirty = false`.
- [x] Chuyển sang note khác vẫn nạp bình thường và reset `isDirty`.
- [x] File thừa `frontend/src/lib/Counter.svelte` được xóa sạch khỏi cây thư mục.

**Verification & Testing:**
- [x] Type check: `npm run check` (0 errors, 0 warnings).
- [x] Manual test: Mở note, gõ vài từ chưa bấm Lưu, cờ `isDirty` kích hoạt và bảo vệ buffer không bị ghi đè.

**Git Commit:**
```bash
git rm frontend/src/lib/Counter.svelte
git add frontend/src/components/note/NoteEditor.svelte frontend/src/lib/stores/viewStore.ts
git commit -m "fix(editor): implement isDirty buffer guard against watcher reload and clean dead code"
```
*(Trạng thái: ✅ Đã commit `2a98867`)*

**Dependencies:** None  
**Files likely touched:**
- `frontend/src/components/note/NoteEditor.svelte`
- `frontend/src/lib/stores/viewStore.ts`
- `frontend/src/lib/Counter.svelte` (Xóa)

**Estimated scope:** Small (2-3 files)

---

### Task 2: Wrap 4 API CRUD Trong wailsAdapter.ts & vaultStore.ts
**Description:** Mở rộng tầng Anti-Corruption Layer `frontend/src/lib/wails/wailsAdapter.ts` để gọi 4 API CRUD mới của Backend (`DeleteNote`, `DeleteFolder`, `RenameNote`, `RenameFolder`). Bổ sung các store actions trong `frontend/src/lib/stores/vaultStore.ts` (`deleteNote`, `deleteFolder`, `renameNote`, `renameFolder`) kèm logic tự động làm tươi dữ liệu Workspace sau khi thực hiện.

**Acceptance criteria:**
- [ ] `WailsAdapter` có 4 phương thức tĩnh: `deleteNote`, `deleteFolder`, `renameNote`, `renameFolder` với cơ chế fallback gọi an toàn qua `window.go`.
- [ ] `vaultStore.ts` có 4 action functions: `deleteNote(id)`, `deleteFolder(path)`, `renameNote(id, newTitle)`, `renameFolder(path, newName)`.
- [ ] Xóa note thành công tự động chọn note kế tiếp hoặc reset `selectedNoteId` về `null`.
- [ ] Đổi tên note thành công tự động cập nhật `selectedNoteId` theo ID (đường dẫn) mới.

**Verification & Testing:**
- [ ] Type check: `npm run check` (0 errors, 0 warnings).
- [ ] Go test: `go test -race ./...` (vẫn đảm bảo 17 test suites PASS 100%).

**Git Commit:**
```bash
git add frontend/src/lib/wails/wailsAdapter.ts frontend/src/lib/stores/vaultStore.ts
git commit -m "feat(store): expose full CRUD actions in wailsAdapter and vaultStore"
```
*(Trạng thái: ⏳ Chưa thực hiện)*

**Dependencies:** Task 1  
**Files likely touched:**
- `frontend/src/lib/wails/wailsAdapter.ts`
- `frontend/src/lib/stores/vaultStore.ts`

**Estimated scope:** Small (2 files)

---

## 🛑 Checkpoint: Sau Phase 1 (Data Safety & Store Bridge)
- [ ] Buffer soạn thảo được bảo vệ an toàn 100% trước các sự kiện cập nhật ngoài.
- [ ] Tầng Store có đầy đủ vũ khí CRUD sẵn sàng kết nối vào UI.
- [ ] `npm run check` không có lỗi cú pháp hoặc type.

---

## Phase 2: Giao Diện Tương Tác CRUD Cho Note & Folder (UI Interaction)

### Task 3: Bổ Sung Action Menu & Xác Nhận Xóa/Đổi Tên Cho Notes
**Description:** Bổ sung giao diện thao tác cho Ghi chú:
1. Thêm nút hành động (Options / More menu) trên Header của `NoteEditor.svelte` và trên từng thẻ ghi chú của `NotesList.svelte` (hỗ trợ Rename và Delete).
2. Xây dựng Modal hoặc Popover xác nhận xóa kính mờ phong cách Apple, cảnh báo người dùng trước khi xóa file vĩnh viễn trên đĩa.
3. Bổ sung các lệnh "Delete Note" và "Rename Note" vào `CommandPalette.svelte`.

**Acceptance criteria:**
- [ ] Người dùng có thể đổi tên ghi chú nhanh thông qua popover/dialog hoặc đổi tên trực tiếp.
- [ ] Thao tác xóa ghi chú kích hoạt popover xác nhận tinh tế (không dùng `window.confirm`).
- [ ] Sau khi xác nhận xóa, file biến mất trên danh sách và hiển thị Toast thông báo thành công.
- [ ] `CommandPalette` có các hành động thao tác với note hiện tại.

**Verification & Testing:**
- [ ] Type check: `npm run check` (0 errors, 0 warnings).
- [ ] Thử nghiệm tạo note tạm thời, đổi tên và xóa thử nghiệm, kiểm tra file trên đĩa biến mất chuẩn xác.

**Git Commit:**
```bash
git add frontend/src/components/note/NotesList.svelte frontend/src/components/note/NoteEditor.svelte frontend/src/components/layout/CommandPalette.svelte
git commit -m "feat(notes): add delete and rename actions with modal confirmation for notes"
```
*(Trạng thái: ⏳ Chưa thực hiện)*

**Dependencies:** Task 2  
**Files likely touched:**
- `frontend/src/components/note/NotesList.svelte`
- `frontend/src/components/note/NoteEditor.svelte`
- `frontend/src/components/layout/CommandPalette.svelte`

**Estimated scope:** Medium (3 files)

---

### Task 4: Bổ Sung Context Actions Xóa & Đổi Tên Thư Mục Trong Sidebar
**Description:** Nâng cấp `frontend/src/components/sidebar/FolderTreeItem.svelte` và `Sidebar.svelte`. Khi di chuột qua thư mục hoặc click vào biểu tượng 3 chấm / click chuột phải, hiển thị menu nhỏ gồm:
1. "New Note Inside" (Tạo note trong thư mục này)
2. "New Subfolder" (Tạo thư mục con)
3. "Rename Folder" (Đổi tên thư mục)
4. "Delete Folder" (Xóa thư mục kèm cảnh báo số lượng tệp bên trong)

**Acceptance criteria:**
- [ ] Mỗi thư mục trên cây phân cấp hiển thị nút menu hoặc hành động nhanh khi hover.
- [ ] Đổi tên thư mục inline hoặc qua modal mượt mà, tự động cập nhật lại cây thư mục.
- [ ] Xóa thư mục hiển thị cảnh báo nguy hiểm rõ ràng (cấm xóa nếu là root vault).
- [ ] Hỗ trợ tạo note/thư mục con đúng vị trí thư mục được chọn.

**Verification & Testing:**
- [ ] Type check: `npm run check` (0 errors, 0 warnings).
- [ ] Thử nghiệm tạo thư mục con, đổi tên và xóa thư mục con, xác minh cây thư mục cập nhật tức thì.

**Git Commit:**
```bash
git add frontend/src/components/sidebar/Sidebar.svelte frontend/src/components/sidebar/FolderTreeItem.svelte
git commit -m "feat(sidebar): add folder context actions for rename and safe deletion"
```
*(Trạng thái: ⏳ Chưa thực hiện)*

**Dependencies:** Task 3  
**Files likely touched:**
- `frontend/src/components/sidebar/Sidebar.svelte`
- `frontend/src/components/sidebar/FolderTreeItem.svelte`

**Estimated scope:** Medium (2 files)

---

## 🛑 Checkpoint: Sau Phase 2 (Full CRUD In UI)
- [ ] Người dùng có thể tạo, đọc, lưu, đổi tên và xóa cả Ghi Chú và Thư Mục hoàn toàn từ UI.
- [ ] Các thao tác xóa đều có lớp bảo vệ xác nhận kính mờ, không gây xóa nhầm dữ liệu.

---

## Phase 3: Polish Giao Diện Chuẩn Đặc Tả PRB/SRS (Layout & Visual Polish)

### Task 5: Xây Dựng Thanh Trạng Thái Chân Trang StatusBar.svelte
**Description:** Xây dựng component `frontend/src/components/layout/StatusBar.svelte` theo đúng thiết kế Wireframe trong tài liệu [PRB.md](../../../markdown/PRB.md#L171). Gắn thanh đáy vào `frontend/src/App.svelte` với chiều cao 28px (`h-7`), thiết kế kính mờ sang trọng chuẩn Apple Design.

**Acceptance criteria:**
- [ ] Component `StatusBar.svelte` hiển thị ở chân trang màn hình với hiệu ứng `backdrop-blur`.
- [ ] Hiển thị chấm trạng thái xanh lá nhấp nháy: `● fsnotify active` kèm đường dẫn workspace hiện tại.
- [ ] Hiển thị số từ (`words`), số ký tự (`chars`), và trạng thái lưu (`Saved` / `Unsaved`).
- [ ] Hiển thị thông số tổng số ghi chú trong vault.
- [ ] Layout tổng thể của ứng dụng cân đối, không bị che khuất nội dung cuộn bên trong Editor hay NotesList.

**Verification & Testing:**
- [ ] Type check: `npm run check` (0 errors, 0 warnings).
- [ ] Kiểm tra giao diện trực quan: Thanh đáy cố định đẹp mắt, mượt mà ở cả Dark và Light mode.

**Git Commit:**
```bash
git add frontend/src/components/layout/StatusBar.svelte frontend/src/App.svelte
git commit -m "feat(layout): implement glassmorphic status bar with realtime watcher status"
```
*(Trạng thái: ⏳ Chưa thực hiện)*

**Dependencies:** Phase 1  
**Files likely touched:**
- `frontend/src/components/layout/StatusBar.svelte` (Tạo mới)
- `frontend/src/App.svelte`

**Estimated scope:** Small-Medium (2 files)

---

### Task 6: Hiển Thị Icon Tùy Biến Frontmatter & Highlight Từ Khóa Tìm Kiếm
**Description:**
1. Hỗ trợ hiển thị icon tùy biến khai báo trong YAML Frontmatter (`note.metadata.icon` như `server`, `cpu`, `briefcase`, `file-text` hoặc emoji) trên thẻ danh sách của `NotesList.svelte` và tiêu đề của `NoteEditor.svelte`.
2. Triển khai bộ lọc highlight từ khóa tìm kiếm trong `NotesList.svelte`: Khi người dùng nhập `searchQuery`, các đoạn ký tự trùng khớp trong `Title` và `Excerpt` được bọc thẻ `<mark>` được tạo kiểu màu sắc tinh tế, an toàn XSS.

**Acceptance criteria:**
- [ ] File có `icon` trong YAML Frontmatter (ví dụ file `BRD.md` có `icon: "briefcase"`, `SRS.md` có `icon: "cpu"`) hiển thị đúng icon tương ứng thay vì icon tệp đơn điệu.
- [ ] Khi gõ từ khóa vào ô tìm kiếm, các ký tự khớp trong `NotesList` được tô sáng rõ ràng.
- [ ] Không phát sinh lỗi XSS hay lỗi render HTML.

**Verification & Testing:**
- [ ] Type check: `npm run check` (0 errors, 0 warnings).
- [ ] Mở file `BRD.md` và `PRB.md`, kiểm tra icon hiển thị; gõ từ khóa "clean" hoặc "kiến trúc" kiểm tra highlight.

**Git Commit:**
```bash
git add frontend/src/components/note/NotesList.svelte frontend/src/components/note/NoteEditor.svelte frontend/src/lib/utils/markdown.ts
git commit -m "feat(ui): support dynamic frontmatter icons and search keyword highlighting"
```
*(Trạng thái: ⏳ Chưa thực hiện)*

**Dependencies:** Task 5  
**Files likely touched:**
- `frontend/src/components/note/NotesList.svelte`
- `frontend/src/components/note/NoteEditor.svelte`
- `frontend/src/lib/utils/markdown.ts`

**Estimated scope:** Medium (3 files)

---

## 🛑 Checkpoint: Sau Phase 3 (Layout & UX Excellence)
- [ ] StatusBar chân trang hoạt động mượt mà, phản ánh thông số thời gian thực.
- [ ] Icon Frontmatter hiển thị sống động trên từng ghi chú.
- [ ] Trải nghiệm tìm kiếm trực quan với từ khóa được highlight rõ nét.

---

## Phase 4: Đồng Bộ Đặc Tả Kỹ Thuật & Nghiệm Thu (Docs & Review)

### Task 7: Cập Nhật Tài Liệu Kỹ Thuật LLD.md Đồng Bộ 100% Với Codebase
**Description:** Cập nhật tệp [markdown/LLD.md](../../../markdown/LLD.md) để phản ánh chính xác các nâng cấp kiến trúc:
1. Bổ sung trường `Metadata` vào struct `model.Note` và trường `Extra` vào `model.Frontmatter`.
2. Cập nhật Class Diagram và danh sách API của `VaultController` thêm 4 phương thức CRUD (`DeleteNote`, `DeleteFolder`, `RenameNote`, `RenameFolder`).
3. Cập nhật mô tả thuật toán bố cục đồ thị tri thức sang mô hình Galaxy Phyllotaxis Spiral và phân bổ màu sắc vũ trụ.

**Acceptance criteria:**
- [ ] `markdown/LLD.md` có đầy đủ đặc tả 12 API Wails IPC (thay vì 8 API cũ).
- [ ] Các struct Go trong tài liệu trùng khớp 100% với `domain/model/models.go`.
- [ ] Thuật toán trong tài liệu mô tả chính xác logic thực tế của `graphBuilder.ts`.

**Verification & Testing:**
- [ ] Kiểm tra tính toàn vẹn markdown, các liên kết nội bộ và code blocks syntax.

**Git Commit:**
```bash
git add markdown/LLD.md
git commit -m "docs(spec): synchronize LLD specification with current models and CRUD APIs"
```
*(Trạng thái: ⏳ Chưa thực hiện)*

**Dependencies:** All previous tasks  
**Files likely touched:**
- `markdown/LLD.md`

**Estimated scope:** Small (1 file)

---

### Task 8: Thẩm Định Nghiệm Thu 2 Trục (Standards & Spec Review) & Build Tổng Thể
**Description:** Tiến hành thẩm định chất lượng toàn diện theo quy trình 2 trục của skill `code-review`:
1. **Trục Standards:** Kiểm tra Clean Architecture, an toàn TypeScript/Svelte 5 Runes, tính thẩm mỹ Apple Glassmorphism, 0 lint warnings, 0 type errors.
2. **Trục Spec:** Đối chiếu lại với 5 tài liệu trong `markdown/` để xác nhận toàn bộ các tính năng cam kết đã được đáp ứng 100%.
3. Biên dịch kiểm chứng: `npm run check` $\rightarrow$ `go test -race ./...` $\rightarrow$ `go build -v .`.

**Acceptance criteria:**
- [ ] `npm run check` đạt 0 errors, 0 warnings.
- [ ] `go test -race ./...` vượt qua 100% test suites.
- [ ] `go build -v .` biên dịch thành công nhị phân `mdvault.exe`.
- [ ] Báo cáo nghiệm thu 2 trục đạt trạng thái PASS.
- [ ] Cập nhật bảng đăng ký đợt việc tại `tasks/README.md`.

**Verification & Testing:**
- [ ] `npm run check` (exit code 0).
- [ ] `go test -race ./...` (exit code 0).
- [ ] `go build -v .` (exit code 0).

**Git Commit:**
```bash
git add tasks/2026-10-11/2-frontend-and-spec-alignment/ tasks/README.md
git commit -m "chore(review): complete two-axis engineering audit and finalize frontend alignment"
```
*(Trạng thái: ⏳ Chưa thực hiện)*

**Dependencies:** Task 7  
**Files likely touched:**
- `tasks/2026-10-11/2-frontend-and-spec-alignment/todo.md`
- `tasks/README.md`

**Estimated scope:** Review & Documentation

---

## 🏁 Tổng Kết Tiến Độ Thực Hiện (Progress Tracker)

| Phase | Tên Giai Đoạn | Số Task | Trạng Thái |
| :---: | :--- | :---: | :---: |
| **Phase 1** | An Toàn Nhập Liệu & Cầu Nối CRUD | 2 Tasks (T1-T2) | ⏳ Chưa thực hiện (0/2) |
| **Phase 2** | Giao Diện Tương Tác CRUD Cho Note & Folder | 2 Tasks (T3-T4) | ⏳ Chưa thực hiện (0/2) |
| **Phase 3** | Polish Giao Diện Chuẩn Đặc Tả PRB/SRS | 2 Tasks (T5-T6) | ⏳ Chưa thực hiện (0/2) |
| **Phase 4** | Đồng Bộ Đặc Tả Kỹ Thuật & Nghiệm Thu | 2 Tasks (T7-T8) | ⏳ Chưa thực hiện (0/2) |
