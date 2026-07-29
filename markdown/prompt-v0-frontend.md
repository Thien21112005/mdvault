# Prompt cho v0: Giao diện App Quản Lý Ghi Chú Markdown

Thiết kế giao diện cho một ứng dụng desktop quản lý ghi chú markdown (giống Obsidian/Notion nhưng đơn giản hơn). Dùng React + Tailwind CSS + shadcn/ui. Đây chỉ là phần giao diện, chưa cần logic thật, có thể dùng data mẫu (mock data).

## Bố cục tổng thể

Layout 3 cột kiểu ứng dụng desktop:

1. **Sidebar trái** (cố định, có thể thu gọn):
   - Ô search ở trên cùng
   - Cây thư mục hiển thị các nhóm kiến thức/folder (có thể expand/collapse), mỗi node có icon folder hoặc file
   - Dưới cây thư mục là danh sách tags (dạng pill/badge), click vào để lọc theo tag
   - Nút "Chọn thư mục" ở đầu sidebar để giả lập việc chọn folder gốc

2. **Panel giữa** (danh sách file trong folder/tag đang chọn):
   - Danh sách các note dạng card hoặc list, mỗi item hiển thị: tiêu đề, vài dòng excerpt, tags, ngày sửa đổi
   - Có thể toggle giữa view dạng list và dạng grid card

3. **Panel phải** (nội dung chính, chiếm phần lớn màn hình):
   - Hiển thị nội dung markdown đã render thành HTML (heading, list, code block, blockquote đều có style rõ ràng)
   - Header của panel có: tiêu đề note, tags, đường dẫn file (breadcrumb)

## Yêu cầu về style

- Theme sáng/tối (dark mode toggle ở góc trên)
- Font chữ dễ đọc cho nội dung dài (ưu tiên font sans-serif rõ ràng cho UI, có thể dùng serif hoặc mono cho phần code block)
- Màu sắc trung tính, không quá sặc sỡ — cảm giác gọn gàng như app ghi chú chuyên nghiệp
- Responsive không cần thiết lắm (đây là app desktop, không cần tối ưu mobile)

## Trạng thái cần có

- Trạng thái rỗng: khi chưa chọn folder nào, hiển thị màn hình hướng dẫn "Chọn thư mục để bắt đầu"
- Trạng thái loading khi đang quét file
- Trạng thái không có kết quả tìm kiếm

## Data mẫu gợi ý

Dùng khoảng 8-10 note mẫu với các tag như "golang", "design-pattern", "note-app", chia vào 2-3 folder giả lập (VD: "Học Go", "Design Patterns", "Dự án cá nhân") để giao diện có nội dung demo đầy đủ.
