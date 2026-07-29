# Prompt cho Agent: Xây dựng App Quản Lý Ghi Chú Markdown

## Bối cảnh & Mục tiêu

Xây dựng một ứng dụng desktop (không phải web) giúp người dùng đọc, quét và tổ chức các file ghi chú `.md` (markdown) theo chủ đề/nhóm kiến thức. Ứng dụng đọc trực tiếp các file trên máy tính người dùng (một thư mục gốc chứa nhiều file/folder con `.md`), phân loại chúng theo tag hoặc thư mục, cho phép xem trước nội dung, tìm kiếm nhanh, và (tùy chọn) hiển thị sơ đồ liên kết giữa các ghi chú.

## Yêu cầu chức năng

1. **Chọn & quét thư mục gốc**: người dùng chọn 1 folder trên máy → app quét đệ quy tất cả file `.md` bên trong.
2. **Đọc frontmatter**: mỗi file `.md` có thể có YAML frontmatter (title, tags, category) ở đầu file — app cần parse phần này riêng với phần nội dung.
3. **Tổ chức theo nhóm kiến thức**: hiển thị cây thư mục ở sidebar, có thể nhóm thêm theo `tags`/`category` lấy từ frontmatter (không chỉ theo cấu trúc folder vật lý).
4. **Xem trước nội dung**: click vào 1 file → hiển thị nội dung markdown đã render thành HTML ở panel chính.
5. **Tìm kiếm**: tìm theo tên file, tag, hoặc full-text trong nội dung.
6. **Theo dõi thay đổi real-time**: nếu file trong folder bị thêm/sửa/xóa bên ngoài app (VD sửa bằng VS Code), app tự cập nhật lại danh sách mà không cần quét lại thủ công.
7. **(Tùy chọn) Graph view**: hiển thị sơ đồ liên kết giữa các note nếu trong nội dung có cú pháp `[[wiki-link]]` trỏ đến note khác.

## Yêu cầu phi chức năng

- Chạy offline hoàn toàn, không cần server hay internet.
- App nhẹ, khởi động nhanh, không ngốn RAM nhiều.
- Đóng gói được thành file thực thi độc lập cho Windows (ưu tiên), có thể mở rộng sang macOS/Linux sau.

## Tech stack chỉ định

- **Framework**: Wails (`wailsapp/wails`) — Go ở backend, webview hệ điều hành để hiển thị UI.
- **Backend (Go)**:
  - `yuin/goldmark` — parse markdown sang HTML
  - `goccy/go-yaml` — parse YAML frontmatter
  - `fsnotify/fsnotify` — theo dõi thay đổi file/folder real-time
  - `modernc.org/sqlite` — lưu index (path, title, tags, excerpt) để search nhanh, không cần cgo
- **Frontend**: Svelte + Vite, Tailwind CSS cho style, `markdown-it` để render preview (nếu không render sẵn từ backend), `xyflow/xyflow` nếu làm graph view.

## Đầu ra mong muốn từ agent

1. Cấu trúc thư mục project rõ ràng (backend Go tách biệt frontend Svelte theo chuẩn Wails).
2. Code khởi tạo cho: scanner (quét folder), parser (đọc frontmatter + nội dung), watcher (fsnotify), và API expose ra frontend qua Wails binding.
3. Giao diện frontend cơ bản: sidebar cây thư mục + panel preview nội dung + ô search.
4. Hướng dẫn build & chạy thử (`wails dev`, `wails build`).

## Ràng buộc

- Ưu tiên code Go idiomatic, có xử lý lỗi rõ ràng (không bỏ qua error).
- Giữ frontend đơn giản trước, graph view làm sau nếu còn thời gian.
- Không cần authentication/multi-user — đây là app cá nhân, single-user, local-only.
