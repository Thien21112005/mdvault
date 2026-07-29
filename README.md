# MDVault - Desktop Markdown Note Taking App

**MDVault** là ứng dụng quản lý ghi chú và tri thức cá nhân nền tảng Markdown, được phát triển bằng **Wails v2** (sử dụng backend Go và giao diện Vite + React + Tailwind CSS v4).

## ✨ Tính năng nổi bật
- 📂 **Quản lý Workspace hệ thống tệp:** Chọn và theo dõi mọi thay đổi trong thư mục lưu trữ Ghi chú của bạn trên ổ đĩa theo thời gian thực (realtime filesystem watcher).
- 📝 **Soạn thảo và Quản lý YAML Frontmatter:** Parse, phân loại Tag, tạo và sửa tệp `.md` nhanh chóng với giao diện hiện đại.
- 🎨 **Giao diện Hài hoà & Hiệu suất Cao:** Tận dụng hiệu năng vi mô từ Go và webview tĩnh, tiết kiệm bộ nhớ gấp nhiều lần so với các ứng dụng Electron truyền thống.

## 🚀 Hướng dẫn cài đặt và chạy (Quick Start)

### Yêu cầu tiên quyết:
- [Go 1.22+](https://golang.org/dl/)
- [Node.js 18+ & npm](https://nodejs.org/)
- [Wails CLI v2](https://wails.io/):
  ```powershell
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

### Chạy chế độ phát triển (Live reload & HMR):
Tại thư mục gốc dự án, thực hiện lệnh:
```powershell
wails dev
```
Wails sẽ tự động tải các module Go, cài đặt gói NPM frontend và hiển thị giao diện ứng dụng.

### Biên dịch bản phát hành (Production Build):
```powershell
wails build -clean -s
```
Tệp thi hành ứng dụng thu gọn sẽ được sinh ra tại `build/bin/mdvault.exe`.
