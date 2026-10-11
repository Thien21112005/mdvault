# 📋 MDVault Tasks & Roadmap Archive

Quy ước lưu trữ kế hoạch (Plan) và danh sách công việc (Todo) của dự án **MDVault**:
- Mỗi ngày có một thư mục riêng biệt theo ngày: `tasks/YYYY-MM-DD/`.
- Trong một ngày, các đợt lập kế hoạch được đánh số thứ tự tăng dần: `1-<tên-nhiệm-vụ>/`, `2-.../`.
- Mỗi thư mục đợt việc bao gồm:
  - `plan.md`: Bản thiết kế kiến trúc, phân tích kỹ thuật và lộ trình tổng quan.
  - `todo.md`: Danh sách các lát cắt nhiệm vụ mỏng (Task Slices), tiêu chí nghiệm thu và lịch sử git commit.

---

## 🗂️ Danh Mục Kế Hoạch (Task Registry)

### 📅 Ngày 2026-10-10

| STT | Nhiệm Vụ / Kế Hoạch | Trạng Thái | Đường Dẫn Plan & Todo | Nội Dung Chính |
| :---: | :--- | :---: | :--- | :--- |
| **01** | **Đại Tu Giao Diện Toàn Bộ Website & App** | ✅ Hoàn thành (8/8) | [`1-ui-revamp/`](file:///d:/Tu-hoc/project/MDVault/tasks/2026-10-10/1-ui-revamp) | Glassmorphism Tokens, Header Chrome, Sidebar Hierarchy, Multi-layer Notes Cards, XSS Sanitizer, NoteEditor Split-View, Galaxy Graph, Command Palette (Ctrl+K), Sonner Toasts |
| **02** | **Xóa Bỏ Native Alert & Custom Glassmorphic Controls** | ✅ Hoàn thành (3/3) | [`2-custom-ui-controls/`](file:///d:/Tu-hoc/project/MDVault/tasks/2026-10-10/2-custom-ui-controls) | Loại bỏ blocking alert() bằng Sonner Toast, xây dựng CustomDropdown kính mờ thay thế HTML `<select>` mặc định |
| **03** | **Toast Delight Motion & Viền Ngang Đồng Bộ** | ✅ Hoàn thành (4/4) | [`3-toast-delight-and-grid-alignment/`](file:///d:/Tu-hoc/project/MDVault/tasks/2026-10-10/3-toast-delight-and-grid-alignment) | Đồng bộ chiều cao viền ngang Header (h-12), nâng cấp Toast Spring Motion nảy nhẹ vào/lướt thẳng sang phải khi ra, triệt tiêu độ thụt dọc, Ambient Glow, thanh đếm ngược và Action Button tương tác tức thì |

---

### 📅 Ngày 2026-10-11

| STT | Nhiệm Vụ / Kế Hoạch | Trạng Thái | Đường Dẫn Plan & Todo | Nội Dung Chính |
| :---: | :--- | :--- | :--- | :--- |
| **01** | **Gia Cố & Cải Tạo Toàn Diện Backend Go** | ⏳ Đang thực hiện (6/10) | [`tasks/2026-10-11/1-backend-hardening-and-refactor/`](file:///d:/Tu-hoc/project/MDVault/tasks/2026-10-11/1-backend-hardening-and-refactor/) | Sửa bão sự kiện Watcher & Debounce 50ms, Triển khai Atomic Write (.tmp -> Rename), Bảo toàn 100% Frontmatter YAML, In-memory Caching (sync.RWMutex), Chuẩn hóa đường dẫn ToSlash, Chống Path Traversal, Bổ sung CRUD (Delete/Rename) và Bộ Unit Test Go đầy đủ (-race) |



