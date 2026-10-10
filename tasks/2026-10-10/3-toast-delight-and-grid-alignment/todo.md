# Danh Sách Công Việc (Todo List) - Đợt 3: Toast Delight Motion & Viền Ngang Đồng Bộ

> **Dự án:** MDVault  
> **Ngày lập:** 2026-10-10  
> **Người điều phối:** `mdvault-orchestrator`

---

## 📋 Danh Sách Nhiệm Vụ (Task Breakdown)

### Task 1: Cân Chỉnh Đường Viền Ngang Đồng Bộ (Pixel-Perfect Grid Baseline Alignment)
**Description:** Đồng bộ chiều cao Header giữa `Sidebar.svelte` và `NotesList.svelte` thành `h-12 shrink-0` (khớp hoàn hảo với `NoteEditor.svelte`), đảm bảo đường viền phân cách ngang `border-b` chạy thẳng tắp xuyên suốt từ trái sang phải.
- [x] Chỉnh header của `Sidebar.svelte` thành `h-12 px-3.5 border-b border-[var(--border-subtle)] flex items-center justify-between shrink-0`.
- [x] Chỉnh header của `NotesList.svelte` thành `h-12 px-3.5 border-b border-[var(--border-subtle)] flex items-center justify-between sticky top-0 z-10 bg-[var(--bg-surface)]/90 backdrop-blur-md shrink-0`.
- [x] Xác minh trực quan đường viền ngang của cả 3 cột nằm trên cùng một độ cao hoàn hảo.

**Verification & Testing:**
- [x] Type check: `npm run check` (0 errors)
- [x] Build test: `npm run build` (Exit code 0)

**Git Commit:**
```bash
git add frontend/src/components/sidebar/Sidebar.svelte frontend/src/components/note/NotesList.svelte
git commit -m "style(layout): align header baseline heights across sidebar and notes list"
```
*(Commit `e0d0d40`)*

---

### Task 2: Nâng Cấp Hệ Thống Toast "WOW" Effect & Smooth Motion
**Description:** Nâng cấp `toastStore.ts` và `Toast.svelte`:
- [x] Hỗ trợ Action Button (ví dụ `{ action: { label: 'Open Vault', onClick: () => openWorkspaceDialog() } }`).
- [x] Animation chuyển động lướt vào (`in:fly={{ y: 24, duration: 320 }}`) và lướt ra (`out:fly={{ x: 64, duration: 220 }}`) bằng `--ease-apple`.
- [x] Thêm thanh đếm ngược thời gian (Countdown Progress Timer Bar) mượt mà, tạm dừng khi hover.
- [x] Viền phát quang đa tầng (Ambient Halo Glow) theo từng trạng thái (Warning Amber, Success Emerald, Error Rose, Info Indigo).
- [x] Cập nhật gọi `toast.warning('Chưa mở Workspace', { description: '...', action: { label: 'Open Vault', onClick: ... } })` trong `vaultStore.ts`.

**Verification & Testing:**
- [x] Type check: `npm run check` (0 errors)
- [x] Build test: `npm run build` (Exit code 0)
- [x] Thử tạo note khi chưa mở vault -> Toast lướt vào mượt mà, có nút `Open Vault`, click vào mở ngay dialog chọn thư mục, lướt ra êm dịu khi đóng.

**Git Commit:**
```bash
git add frontend/src/lib/stores/toastStore.ts frontend/src/components/ui/Toast.svelte frontend/src/lib/stores/vaultStore.ts
git commit -m "feat(toast): elevate toast UX with smooth entrance/exit motion, ambient glow and action button"
```
*(Commit `ccc34e4`)*

---

## 🏁 Checkpoint: Kiểm Thử Toàn Diện & Nghiệm Thu
- [x] Đường viền ngang giữa Sidebar, NotesList và NoteEditor thẳng tắp 100% (cùng `h-12`).
- [x] Toast notification xuất hiện lướt vào/ra mượt mà, đẹp mắt, có nút bấm tương tác tức thì.
- [x] Cả 2 lệnh `npm run check` và `npm run build` chạy thành công với mã thoát 0.

