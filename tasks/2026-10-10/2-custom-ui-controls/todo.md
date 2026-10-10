# Danh Sách Công Việc (Todo List) - Đợt 2: Xóa Bỏ Native Alert & Custom Dropdown

> **Dự án:** MDVault  
> **Ngày lập:** 2026-10-10  
> **Người điều phối:** `mdvault-orchestrator` / `planning-and-task-breakdown`

---

## 📋 Danh Sách Nhiệm Vụ (Task Breakdown)

### Task 1: Loại Bỏ Native `alert()` & Thay Bằng Toast Cảnh Báo
**Description:** Loại bỏ triệt để lệnh `alert()` trong `vaultStore.ts` khi người dùng cố tạo note mà chưa mở Workspace. Chuyển sang dùng `toast.warning()` trang nhã kèm lời nhắc rõ ràng.
- [x] Xóa `alert('Vui lòng mở Workspace trước khi tạo ghi chú!')` trong `vaultStore.ts`.
- [x] Bổ sung `toast.warning('Chưa mở Workspace', 'Vui lòng nhấn "Open Vault" để chọn thư mục lưu trữ trước!')`.
- [x] Kiểm tra toàn bộ mã nguồn không còn bất kỳ lệnh `alert()`, `confirm()`, `prompt()` nào.

**Verification & Testing:**
- [x] Type check: `npm run check` (0 errors)
- [x] Build test: `npm run build` (Exit code 0)
- [x] Thử mở ứng dụng khi chưa chọn Vault và nhấn nút "Create Note" -> xuất hiện Toast cảnh báo màu vàng thay vì pop-up trình duyệt.

**Git Commit:**
```bash
git add frontend/src/lib/stores/vaultStore.ts
git commit -m "fix(vault): replace blocking browser alert with sonner toast warning"
```
*(Đã commit: `1225dea`)*


---

### Task 2: Xây Dựng Component `CustomDropdown.svelte` Kính Mờ Cao Cấp
**Description:** Tạo component tái sử dụng `frontend/src/components/ui/CustomDropdown.svelte`:
- [x] Trigger hiển thị nhãn hiện tại + Chevron xoay $180^\circ$ khi mở bằng `--ease-apple`.
- [x] Popover menu kính mờ `backdrop-blur-xl`, bo góc mềm mại, viền sáng sheen kim loại.
- [x] Mỗi option có icon checkmark `Check` khi active, hiệu ứng hover nhạy bén.
- [x] Tự động đóng menu khi click ra ngoài (`clickOutside` handler) hoặc nhấn phím `Escape`.
- [x] Hỗ trợ phím mũi tên `ArrowUp` / `ArrowDown` và `Enter` để chọn mục bằng bàn phím.

**Verification & Testing:**
- [x] Type check: `npm run check` (0 errors)
- [x] Build test: `npm run build` (Exit code 0)

**Git Commit:**
```bash
git add frontend/src/components/ui/CustomDropdown.svelte
git commit -m "feat(ui): create reusable glassmorphic custom dropdown component"
```
*(Đã commit: `7873e0e`)*


---

### Task 3: Tích Hợp `CustomDropdown` Vào Bộ Sắp Xếp `NotesList.svelte`
**Description:** Thay thế thẻ HTML `<select>` mặc định trong `NotesList.svelte` bằng `CustomDropdown.svelte`:
- [ ] Xóa bỏ `<select>` và các thẻ `<option>` nguyên thủy.
- [ ] Truyền danh sách tùy chọn sắp xếp (`Newest`, `Oldest`, `Title (A-Z)`, `Title (Z-A)`) vào `CustomDropdown`.
- [ ] Liên kết 2 chiều với store reactive `$sortOption`.
- [ ] Đảm bảo hiển thị hoàn hảo trên cả 2 theme Dark Mode và Light Mode.

**Verification & Testing:**
- [ ] Type check: `npm run check` (0 errors)
- [ ] Build test: `npm run build` (Exit code 0)
- [ ] Manual check: Nhấp vào menu sắp xếp -> menu kính mờ xuất hiện mượt mà, chọn đổi tùy chí sắp xếp -> danh sách ghi chú cập nhật ngay tức khắc.

**Git Commit:**
```bash
git add frontend/src/components/note/NotesList.svelte
git commit -m "feat(notes-list): replace native html select with glassmorphic custom dropdown"
```

---

## 🏁 Checkpoint: Kiểm Thử Toàn Diện & Hoàn Thành Đợt 2
- [ ] Không còn bất kỳ native `alert()` nào xuất hiện trong ứng dụng.
- [ ] Không còn native `<select>` HTML dropdown nào gây lệch phong cách thiết kế.
- [ ] Cả 3 lệnh `npm run check`, `npm run build` đều chạy thành công mỹ mãn với mã thoát 0.
- [ ] Đã hoàn tất 3 commits theo chuẩn Conventional Commits.
