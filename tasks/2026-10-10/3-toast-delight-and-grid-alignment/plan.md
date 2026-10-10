# Kế Hoạch Đợt 3: Toast Delight Motion & Cân Bằng Viền Ngang Pixel-Perfect

> **Mục tiêu:** Đồng bộ chiều cao Header giữa Sidebar, NotesList và NoteEditor để đường viền ngang thẳng hàng 100%; Đại tu Toast Notification với hiệu ứng lướt vào/ra mượt mà (`in:fly` / `out:fly`), Action Button tương tác tức thì và Ambient Halo Glow.

---

## 1. Phân Tích Nguyên Nhân & Vấn Đề

### 1.1. Viền ngang bị lệch giữa Sidebar và NotesList (Screenshot 2)
- **Nguyên nhân:**
  - Header của `Sidebar.svelte` dùng `py-3` không set `height` cố định, bên trong chỉ có icon nhỏ nên chiều cao tổng thể $\approx 44\text{px}$.
  - Header của `NotesList.svelte` dùng `py-3` nhưng bên trong có `CustomDropdown`, làm chiều cao tổng thể tăng lên $\approx 50\text{px}$.
  - Header của `NoteEditor.svelte` có `h-12` ($48\text{px}$).
- **Hệ quả:** Đường `border-b` của Sidebar bị cao hơn/lệch so với đường `border-b` của NotesList, tạo cảm giác đứt gãy thị giác.
- **Giải pháp:** Cố định chuẩn `h-12 shrink-0` ($48\text{px}$) cho cả 3 header. Đường viền `border-b` sẽ nằm trên cùng 1 trục $Y=48\text{px}$ xuyên suốt cả 3 cột.

### 1.2. Toast Notification thiếu tính "WOW" & Animation (Screenshot 1)
- **Nguyên nhân:**
  - Toast hiện tại chỉ là khung card tĩnh mờ nhạt, biến mất đột ngột khi hết giờ.
  - Khi báo "Chưa mở Workspace", người dùng phải tự di chuột lên thanh Navbar tìm nút "Open Vault".
- **Giải pháp:**
  - **Motion:** Tích hợp animation Svelte GPU-accelerated: lướt vào mượt mà từ góc dưới phải (`y: 20px`, `--ease-apple`) và lướt ra trơn tru (`x: 40px`, fade-out).
  - **Action Button:** Toast hỗ trợ nút hành động tức thì: Nút `Open Vault` ngay trên thông báo để người dùng click mở luôn dialog chọn thư mục!
  - **Aesthetics:** Thanh đếm ngược thời gian (Countdown Progress Bar), viền phát quang Ambient Glow theo màu trạng thái (Amber cho Warning, Emerald cho Success...).

---

## 2. Lộ Trình Triển Khai (Phased Tasks)

- **Task 1:** Đồng bộ chiều cao `h-12 shrink-0` cho Header của `Sidebar.svelte` và `NotesList.svelte` tạo đường viền ngang thẳng tắp.
- **Task 2:** Đại tu `toastStore.ts` & `Toast.svelte`: Animation lướt vào/ra (`fly`/`fade`), Action Button ngữ cảnh, Countdown Timer Bar và Glassmorphism Glow.
- **Task 3:** Kiểm thử toàn diện (`npm run check`, `npm run build`), nghiệm thu trực quan.
