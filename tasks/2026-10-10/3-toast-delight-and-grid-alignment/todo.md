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

### Task 3: Tinh Chỉnh Spring Physics & Khắc Phục Lỗi Unmount Container Đột Ngột
**Description:** Khắc phục triệt để hiện tượng xuất hiện cái rụp (thiếu entrance nảy lên từ từ) và mất hiệu ứng lướt ra khi chỉ có 1 thông báo duy nhất:
- [x] Gỡ bỏ block `{#if $toast.length > 0}` bọc ngoài, giữ container cố định trong DOM (`pointer-events-none`) để Svelte không hủy parent đột ngột làm đứt `out:` transition khi toast cuối cùng biến mất.
- [x] Xây dựng custom `toastSpringIn` với `backOut` physics (`translateY(32px) scale(0.92)` -> `translateY(0) scale(1)`) tạo độ nảy đàn hồi (spring overshoot) nổi lên từ từ cực kỳ tự nhiên.
- [x] Xây dựng `toastSlideOut` với `cubicOut` (`translateX(64px) scale(0.92)`) lướt mượt mà sang phải khi đóng.
- [x] Áp dụng `animate:flip={{ duration: 250, easing: cubicOut }}` từ `svelte/animate` giúp các thông báo còn lại tự động trượt êm ái vào vị trí mới khi một thông báo biến mất thay vì giật cục.
- [x] Giới hạn tối đa 4 toast trong `toastStore.ts` để tránh tràn màn hình khi xuất hiện hàng loạt.

**Verification & Testing:**
- [x] Type check: `npm run check` (0 errors)
- [x] Build test: `npm run build` (Exit code 0)
- [x] Kiểm thử thông báo đơn lẻ: nảy lên từ từ nhẹ nhàng (`spring overshoot`), khi hết hạn hoặc click `X` lướt êm ái sang phải trước khi biến mất hoàn toàn.
- [x] Kiểm thử nhiều thông báo xếp chồng: các thẻ trượt mượt mà bằng FLIP layout animation, không đè lấn giật cục.

**Git Commit:**
```bash
git add frontend/src/lib/stores/toastStore.ts frontend/src/components/ui/Toast.svelte
git commit -m "fix(toast): ensure persistent container for reliable exits and add spring entrance with flip reordering"
*(Commit `6e31fe8`)*

---

### Task 4: Loại Bỏ Hiện Tượng Thụt Xuống Dưới Khi Lướt Ra (Pure Horizontal Slide-out & Fix False Success Toasts)
**Description:** Khắc phục triệt để hiện tượng thông báo bị thụt xuống dưới trước khi lướt ra, và ngăn chặn việc kích hoạt đồng thời cả thông báo cảnh báo lẫn thành công:
- [x] Gỡ bỏ `animate:flip` trên phần tử đang thoát trong container `fixed bottom-6` (nguyên nhân gốc rễ gây ra độ dời thẳng đứng `dy = +80px` thụt xuống dưới khi container đáy co lại).
- [x] Cố định trục tung trong `toastSlideOut`: chỉ lướt theo phương ngang `transform: translateX(${translateX}px)` với `cubicOut`, triệt tiêu hoàn toàn bất kỳ độ dịch chuyển theo trục Y nào.
- [x] Xóa bỏ class `transition-transform` trên thẻ card để trình duyệt không tự suy diễn chuyển đổi xung đột với keyframes của Svelte.
- [x] Sửa `createNewNote` và `createNewFolder` trong `vaultStore.ts` trả về `Promise<boolean>`, đồng thời cập nhật `Sidebar.svelte` và `CommandPalette.svelte` chỉ bắn `toast.success` khi thao tác thực sự thành công, chấm dứt tình trạng 2 thông báo mâu thuẫn xuất hiện cùng 1 lúc.

**Verification & Testing:**
- [x] Type check: `npm run check` (0 errors)
- [x] Build test: `npm run build` (Exit code 0)
- [x] Kiểm thử trực quan trong trình duyệt (Browser Subagent): Toast xuất hiện nảy lên với spring physics, khi tắt bằng nút `X` hoặc hết giờ thì lướt thẳng tắp sang bên phải trên đúng cao độ hiện tại, hoàn toàn không bị thụt xuống dưới.

**Git Commit:**
```bash
git add frontend/src/components/ui/Toast.svelte frontend/src/components/sidebar/Sidebar.svelte frontend/src/lib/stores/vaultStore.ts frontend/src/components/layout/CommandPalette.svelte
git commit -m "fix(toast): eliminate vertical drop during exit with pure horizontal slide and fix phantom success toasts"
```

---

## 🏁 Checkpoint: Kiểm Thử Toàn Diện & Nghiệm Thu
- [x] Đường viền ngang giữa Sidebar, NotesList và NoteEditor thẳng tắp 100% (cùng `h-12`).
- [x] Toast đơn lẻ hay nhiều toast đều nổi lên nảy nhẹ đàn hồi (`toastSpringIn`).
- [x] Toast khi đóng lướt thẳng sang phải tại đúng cao độ hiện tại, triệt tiêu 100% hiện tượng thụt xuống dưới.
- [x] Chỉ hiển thị thông báo thành công khi thao tác lưu/tạo ghi chú thành công thật sự.
- [x] Cả 2 lệnh `npm run check` và `npm run build` chạy thành công với mã thoát 0.



