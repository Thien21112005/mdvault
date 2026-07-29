# Prompt cho Agent: Tích hợp giao diện v0 vào App Wails

## Bối cảnh

Tôi đã dùng v0 (Vercel) để thiết kế giao diện cho một app desktop quản lý ghi chú markdown. Code v0 xuất ra là React viết theo kiểu Next.js (có thể có file-based routing, server components, mock data). Tôi cần chuyển nó thành frontend chạy được trong **Wails** (framework desktop dùng Go ở backend + webview để hiển thị UI thuần React, không phải Next.js).

## Việc cần làm

1. **Bóc tách code khỏi cấu trúc Next.js**:
   - Loại bỏ file-based routing (`app/` hoặc `pages/`), gộp lại thành component React thuần trong 1 SPA
   - Loại bỏ mọi Server Component, Server Action nếu có — chuyển hết về Client Component vì Wails không có server runtime
   - Giữ nguyên các component UI (shadcn/ui, Tailwind classes) không cần đổi

2. **Thiết lập project Wails đúng chuẩn**:
   - Cấu trúc: `frontend/` (chứa React app) + backend Go ở thư mục gốc theo chuẩn `wailsapp/wails`
   - Build tool cho frontend dùng Vite (không dùng Next.js build)

3. **Thay mock data bằng dữ liệu thật qua Wails bindings**:
   - Thay các function fetch mock data bằng gọi hàm Go đã expose qua Wails (VD: `ScanFolder()`, `GetNoteContent(path string)`, `SearchNotes(query string)`)
   - Backend Go xử lý: quét folder (đệ quy tìm file `.md`), parse frontmatter (`goccy/go-yaml`), render markdown sang HTML (`yuin/goldmark`), theo dõi thay đổi file (`fsnotify`)
   - Frontend gọi các hàm này qua object binding tự sinh của Wails (thường nằm ở `wailsjs/go/main/App`)

4. **Xử lý state thật thay vì state giả lập**:
   - Cây thư mục và danh sách note load từ kết quả `ScanFolder()` thay vì mảng cứng
   - Cập nhật lại danh sách khi có sự kiện thay đổi file từ backend (Wails hỗ trợ event emit từ Go sang frontend qua `runtime.EventsEmit`)

## Yêu cầu khi code

- Giữ nguyên phần UI/UX đã thiết kế từ v0, chỉ đổi phần data layer
- Code Go cần xử lý lỗi rõ ràng, không bỏ qua error
- Comment lại những chỗ đã thay đổi so với code gốc từ v0 để dễ đối chiếu
- Nếu thiếu function Go nào cần thiết mà tôi chưa có, tự tạo thêm và giải thích ngắn gọn dùng để làm gì

## Input tôi sẽ cung cấp

Code React xuất ra từ v0 (component files, có thể kèm theo cả file cấu hình Tailwind/shadcn nếu có).
