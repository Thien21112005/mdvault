# MDVault - Professional Desktop Markdown Vault & Knowledge Graph

[![Clean Architecture](https://img.shields.io/badge/Backend-Clean%20Architecture-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Svelte 5](https://img.shields.io/badge/Frontend-Svelte%205%20%2B%20Vite-FF3E00?style=flat-square&logo=svelte)](https://svelte.dev/)
[![Tailwind CSS v4](https://img.shields.io/badge/Styling-Tailwind%20CSS%20v4-06B6D4?style=flat-square&logo=tailwindcss)](https://tailwindcss.com/)
[![Wails v2](https://img.shields.io/badge/Framework-Wails%20v2-E34F26?style=flat-square)](https://wails.io/)

**MDVault** là giải pháp phần mềm Quản lý Ghi chú và Hệ thống Cấu trúc Tri thức (Personal Knowledge Management - PKM) cho hệ điều hành Desktop. Ứng dụng theo đuổi triết lý **Local-First** và hiệu năng tối cao, xóa bỏ nỗi ám ảnh ngốn RAM của các app Electron truyền thống thông qua việc kết hợp con tim **Go (Wails v2)** và khối óc phản ứng siêu thanh **Svelte 5 + Vite**.

---

## 🏛️ Đặc Tả Kiến Trúc & Triết Lý Phát Triển

Dự án được xây dựng và duy trì dưới quy mô tiêu chuẩn kỹ thuật phần mềm cao nhất. Xin vui lòng tham khảo chi tiết tài liệu kiến trúc chuyên sâu của hệ thống tại:
👉 **[Tài Liệu Kỹ Thuật Chuyên Sâu & Đặc Tả Hệ Thống (system-architecture-spec.md)](file:///d:/Tu-hoc/project/MDVault/markdown/system-architecture-spec.md)**

### 🔥 Những Điểm Nhấn Kỹ Thuật Cốt Lõi:
1. **Backend Go theo chuẩn Clean Architecture (4 Phân Lớp):**
   - Sự phân tách nhiệm vụ rõ ràng: `domain/` (Cốt lõi Mô hình), `usecase/` (Logic nghiệp vụ), `infrastructure/` (Bộ điều khiển ổ đĩa & Watcher) và `controller/` (Cầu nối Wails IPC).
   - Tự động lắng nghe thời gian thực biến động file system bằng **`fsnotify`** và giải mã siêu dữ liệu Frontmatter bằng **`goccy/go-yaml`**.
2. **Frontend Svelte theo phong cách Feature-Driven Design (FDD):**
   - Xử lý trạng thái thông minh qua cụm Reactive Svelte Stores (`lib/stores/`), cùng Lớp chống phân ly (`WailsAdapter`) cô lập UI khỏi rào cản RPC gốc.
   - **Render Markdown:** Sức mạnh trình bày đỉnh cao, hỗ trợ Typography rực rỡ với thư viện **`markdown-it`**.
   - **Mạng Lưới Tri Thức (Knowledge Graph View):** Tự động truy hồi cú pháp **`[[wiki-link]]`** trên toàn vùng lưu trữ để dựng lên biểu đồ tương tác 2D trực quan bằng sức mạnh của **`@xyflow/svelte`**.

---

## 🚀 Hướng Dẫn Sử Dụng Và Dịch Mạng (Quick Start)

### 1. Yêu Cầu Tiền Trình:
- **Go 1.22+** & **Node.js 18+ / 20+**
- **Cài đặt Wails CLI v2:**
  ```powershell
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

### 2. Kích Hoạt Môi Trường Sáng Tạo (Development Mode):
Cung cấp chế độ thay thế mô-đun siêu tốc (Hot Module Replacement) cho Svelte và gỡ lỗi Trìu tượng cho Go:
```powershell
wails dev
```
*(Ngay sau khi lệnh xuất phát, cửa sổ MDVault sẽ hiển thị; bất cứ edit nào trên code Svelte sẽ tức khắc hiện lên mà không mất trạng thái App).*

### 3. Biên Dịch Xuất Xưởng Tệp Thi Hành (.EXE):
Xuất bản gói nhị phân độc lập thu gọn tối đa và tối ưu hiệu suất bộ nhớ:
```powershell
wails build -clean -s -o mdvault.exe
```
📦 **Sản phẩm thu thập được sẽ ngụ tại:** `build/bin/mdvault.exe`. Nhấp đúp chuột và bắt đầu kiến tạo thế giới ghi chú tri thức của riêng bạn!
