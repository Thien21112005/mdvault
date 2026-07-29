# MDVault - Desktop Markdown Note Taking & Knowledge App

**MDVault** là ứng dụng quản lý ghi chú và hệ thống tri thức cá nhân nền tảng Markdown cho máy tính Desktop. Ứng dụng được thiết kế theo các tiêu chuẩn kiến trúc hiện đại hàng đầu:
- **Backend:** **Go + Wails v2** (được mô-đun hoá theo **Clean Architecture**).
- **Frontend:** **Svelte 5 + TypeScript + Vite + Tailwind CSS v4**.
- **Hiển thị văn bản:** Render tốc độ cao với **`markdown-it`** và Typography.
- **Hiển thị đồ thị tri thức:** Mạng lưới liên kết Ghi chú dạng wiki-links `[[note]]` thông qua **`@xyflow/svelte`** (Svelte Flow).

---

## 🌟 Đặc điểm kỹ thuật & Tường lửa Kiến trúc

1. **Clean Architecture (Go Backend):**
   - Sự chia tách độc lập giữa các tầng: `domain/`, `usecase/`, `infrastructure/` và `controller/`.
   - Lắng nghe biến động ổ đĩa trực tiếp với `fsnotify`, tự động đồng bộ khi bạn chỉnh sửa file `.md` bằng editor ngoài (như VS Code, Neovim, Obsidian).
2. **Professional Modular Svelte App:**
   - Bộ quản lý trạng thái tập trung Svelte Stores (`lib/stores/`), cùng lớp Wails Adapter cô lập giao tiếp IPC.
   - Trình hiển thị Graph View tương tác cho phép bạn quan sát cấu trúc kết nối tri thức và chuyển hướng ngay bằng những nhấp chuột.

---

## 🚀 Hướng Dẫn Cài Đặt và Kích Hoạt

### Requirements:
- Go 1.22+
- Node.js 18+ & npm
- Wails CLI v2: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Chạy ứng dụng chế độ Dev (Hỗ trợ Live Reload):
Tại thư mục gốc của kho chứa, thực hiện lệnh sau:
```powershell
wails dev
```

### Đóng gói ứng dụng phát hành (.EXE cho Windows):
```powershell
wails build -clean -s
```
Tệp thực thi độc lập mang lại trải nghiệm tối ưu tài nguyên nhất sẽ ra mắt tại: `build/bin/mdvault.exe`.
