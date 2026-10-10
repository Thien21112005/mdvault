# MDVault UI/UX Redesign - Danh Sách Nhiệm Vụ Triển Khai (Todo Checklist)

> **Điều phối bởi:** `mdvault-orchestrator`  
> **Quy chuẩn thực thi:** `incremental-implementation` & `planning-and-task-breakdown`  
> **Quy trình chuẩn mỗi task:** **Code $\rightarrow$ Test & Verify $\rightarrow$ Review $\rightarrow$ Git Commit $\rightarrow$ Next Task**

---

## 📌 Nguyên Tắc Kiểm Thử & Cam Kết Git (Testing & Git Protocol)

1. **Có cần Test không? CỰC KỲ BẮT BUỘC!**
   - **Type & Compiler Check:** Chạy `npm run check` để đảm bảo hệ thống Svelte 5 Runes (`$state`, `$derived`, `$effect`), TypeScript types và bindings không phát sinh lỗi hoặc cảnh báo ngầm.
   - **Build Verification:** Chạy `npm run build` để xác minh Vite Rollup đóng gói asset, CSS tailwind và component sạch sẽ.
   - **Interactive / Visual Verification:** Kiểm tra trực quan trên màn hình (độ tương phản, hiệu ứng kính mờ, phản hồi xúc giác `:active scale(0.97)`, chuyển động mượt 60fps).
   - **Logic / Security Testing:** Với Task 5 (XSS Sanitizer), viết test kiểm thử khử trùng mã độc trong `markdown.ts` trước khi đưa vào preview.
2. **Quy tắc Git Commit:**
   - Mỗi task hoàn thành và kiểm thử PASS 100% sẽ được commit riêng biệt theo chuẩn **Conventional Commits** (`style:`, `feat:`, `refactor:`, `sec:`).
   - Tuyệt đối không dồn toàn bộ task vào 1 commit lớn; giúp dễ dàng revert nếu có sự cố.

---

## Phase 1: Nền Tảng Thiết Kế, Tokens & App Shell (Foundation)

### Task 1: Thiết Lập Design Tokens & Nâng Cấp Typography/Glassmorphism trong app.css
**Description:** Khởi tạo hệ thống biến CSS token cao cấp (Glassmorphism, Layered Elevation Shadows, Custom Apple Springs & Easing Curves, Surface Tints) và hoàn thiện hỗ trợ hai giao diện Dark & Light mode tương thích chuẩn Tailwind v4.

**Acceptance criteria:**
- [x] Định nghĩa các token biến CSS cho Glass surfaces, border glow, shadow đa tầng và easing curves (`--ease-apple`).
- [x] Hỗ trợ đầy đủ dark & light theme styling mượt mà, không bị chói hoặc gắt màu.
- [x] Typography Geist Sans & Geist Mono được cấu hình tracking và line-height theo chuẩn optical sizing.

**Verification & Testing:**
- [x] Type check: `cd frontend && npm run check` (0 errors)
- [x] Build succeeds: `cd frontend && npm run build` (Exit code 0)
- [x] Manual visual check: CSS tokens nạp đầy đủ trong DOM không sinh lỗi parse; áp dụng được cả dark và light mode.

**Git Commit:**
```bash
git add frontend/ ; git commit -m "style(design-system): implement glassmorphism tokens, optical typography and dual theme palette"
```
*(Đã commit: `32951ab`)*

**Dependencies:** None  
**Files likely touched:**
- `frontend/src/app.css`
- `frontend/src/App.svelte`

**Estimated scope:** Small (2 files)

---

### Task 2: Tái Cấu Trúc App Shell & Desktop Header Chrome (Navbar.svelte)
**Description:** Thiết kế lại thanh tiêu đề của ứng dụng thành Desktop Application Chrome: Tích hợp nút đóng/mở Sidebar, Breadcrumb thư mục đang mở, Command Palette trigger (`⌘K`/`Ctrl+K`), nút chuyển theme (Dark/Light), và hiển thị trạng thái kết nối vault trực quan.

**Acceptance criteria:**
- [x] Navbar có bố cục cân đối, mang dáng dấp ứng dụng Desktop bản địa (Native Desktop Feel) thay vì thanh điều hướng web thông thường.
- [x] Nút Toggle Sidebar có animation chuyển đổi trạng thái mượt mà.
- [x] Tích hợp nút chuyển theme (Sun/Moon icon) kết nối trực tiếp với `themeStore`.
- [x] Nút View Switcher (List vs Graph) sử dụng Segmented Pill Control có hiệu ứng trượt/nổi bật rõ rệt.

**Verification & Testing:**
- [x] Type check: `cd frontend && npm run check` (0 errors)
- [x] Build test: `cd frontend && npm run build` (Exit code 0)
- [x] Manual check: Thao tác bấm đổi theme Dark/Light, chuyển đổi view mode List/Graph và toggle sidebar hoạt động trơn tru.

**Git Commit:**
```bash
git add frontend/src/components/layout/Navbar.svelte frontend/src/App.svelte frontend/src/lib/stores/viewStore.ts
git commit -m "feat(layout): revamp desktop header chrome with theme switcher, breadcrumb and view controls"
```
*(Đã commit: `3ffef50`)*

**Dependencies:** Task 1  
**Files likely touched:**
- `frontend/src/components/layout/Navbar.svelte`
- `frontend/src/App.svelte`
- `frontend/src/lib/stores/viewStore.ts`

**Estimated scope:** Medium (3 files)

---

## Checkpoint 1: Nền Tảng Khung Giao Diện (Foundation Checkpoint)
- [x] `npm run check` và `npm run build` trong `frontend/` biên dịch sạch sẽ không có lỗi.
- [x] Khung Shell của ứng dụng hiển thị mượt mà với tông màu sang trọng, chuyển đổi theme trơn tru.
- [x] Đã hoàn thành 2 commits sạch cho Task 1 và Task 2 (`32951ab`, `3ffef50`).

---

## Phase 2: Cây Thư Mục & Không Gian Quản Trị Tệp (Sidebar Experience)

### Task 3: Tinh Chỉnh Cây Thư Mục & Popover Tạo Nhanh (Sidebar & FolderTreeItem)
**Description:** Nâng cấp `Sidebar.svelte` và `FolderTreeItem.svelte` với icon động, mũi tên đóng/mở thư mục dạng chevron xoay mượt mà, huy hiệu đếm số lượng tài liệu, và thay thế các ô nhập tạo mới dạng thô bằng Modal/Popover trang nhã.

**Acceptance criteria:**
- [x] Cây thư mục hiển thị rõ cấp bậc thụt lề, chevron mở/gập thư mục xoay với transition tự nhiên.
- [x] Hộp thoại tạo Folder/Note mới không còn chèn ép thanh điều hướng, thao tác trực quan với phím Enter và Esc để hủy.
- [x] Khu vực Thẻ (Tags) hiển thị dạng mây tag hiện đại với số lượng note trong mỗi tag và hiệu ứng active nổi bật.

**Verification & Testing:**
- [x] Type check: `cd frontend && npm run check` (0 errors)
- [x] Build test: `cd frontend && npm run build` (Exit code 0)
- [x] Manual check: Tạo thử folder mới, tạo note mới, mở/gập các thư mục con mượt mà; tag cloud lọc chính xác.

**Git Commit:**
```bash
git add frontend/src/components/sidebar/Sidebar.svelte frontend/src/components/sidebar/FolderTreeItem.svelte
git commit -m "feat(sidebar): enhance folder tree hierarchy, chevron rotation and smooth creation popovers"
```
*(Đã commit: `ed9f3a6`)*

**Dependencies:** Task 2  
**Files likely touched:**
- `frontend/src/components/sidebar/Sidebar.svelte`
- `frontend/src/components/sidebar/FolderTreeItem.svelte`

**Estimated scope:** Medium (2 files)

---

## Phase 3: Danh Sách Ghi Chú & Thẻ Bài Nghệ Thuật (Notes List & Search)

### Task 4: Nâng Cấp NotesList Với Note Card Đa Tầng & Bộ Lọc Nâng Cao
**Description:** Tái thiết kế toàn bộ `NotesList.svelte`: Note Card hiển thị tiêu đề, trích đoạn tóm tắt 2 dòng, ngày tháng sửa đổi dạng thân thiện, tag chips; bổ sung thanh công cụ sắp xếp (Mới nhất, Cũ nhất, Tên A-Z) và trạng thái rỗng (Empty State) bắt mắt.

**Acceptance criteria:**
- [ ] Mỗi thẻ Note Card có hiệu ứng hover nâng nhẹ và phản hồi nhấp chuột `:active scale(0.99)`.
- [ ] Note đang được chọn có dải sáng Accent Glow Indigo hoặc Amber phát quang mềm mại.
- [ ] Hỗ trợ chọn tiêu chí sắp xếp ghi chú.
- [ ] Khi không có ghi chú nào tìm thấy, hiển thị Empty State đồ họa tinh tế kèm lời gợi ý hữu ích.

**Verification & Testing:**
- [ ] Type check: `cd frontend && npm run check` (0 errors)
- [ ] Build test: `cd frontend && npm run build` (Exit code 0)
- [ ] Manual check: Lọc thử theo tag, gõ từ khóa tìm kiếm, kiểm tra đổi ghi chú trong danh sách mượt mà.

**Git Commit:**
```bash
git add frontend/src/components/note/NotesList.svelte frontend/src/lib/stores/vaultStore.ts
git commit -m "feat(notes-list): redesign note cards with layered depth, active glow and sorting options"
```

**Dependencies:** Task 3  
**Files likely touched:**
- `frontend/src/components/note/NotesList.svelte`
- `frontend/src/lib/stores/vaultStore.ts`

**Estimated scope:** Small (2 files)

---

## Checkpoint 2: Trải Nghiệm Điều Hướng & Quản Lý Danh Sách (Navigation Checkpoint)
- [ ] Hai cột bên trái (Sidebar và NotesList) hoạt động đồng bộ, bố cục chuẩn tỷ lệ vàng.
- [ ] Thao tác tìm kiếm và chuyển đổi giữa các note phản hồi tức thì dưới 16ms.
- [ ] Đã commit đầy đủ Task 3 và Task 4 vào git history.

---

## Phase 4: Trình Soạn Thảo, Preview & An Toàn Bảo Mật (Editor & Security)

### Task 5: Bảo Mật XSS Cho Markdown Renderer & Nâng Cấp Typography
**Description:** Khử khuẩn mã HTML đầu ra của `markdown-it` bằng thư viện hoặc bộ lọc an toàn để triệt tiêu nguy cơ Stored XSS trong Microsoft Edge WebView2 theo tiêu chuẩn `security-and-hardening`. Làm đẹp typography hiển thị (Code block có khung & nút copy, blockquote có viền gradient, danh sách công việc checkbox tương tác).

**Acceptance criteria:**
- [ ] Các thẻ HTML độc hại như `<script>` hoặc sự kiện `onerror` bị loại bỏ triệt để trước khi render lên DOM.
- [ ] Khối Code block trong markdown có nút "Sao chép" (Copy code) phản hồi trực quan.
- [ ] Kiểu chữ và khoảng cách dòng bài viết được tối ưu cho việc đọc dài lâu (Read mode cao cấp).

**Verification & Testing:**
- [ ] Security validation: Render thử chuỗi Markdown chứa payload `<script>alert(1)</script>` và `<img src=x onerror=alert(1)>`, xác minh script bị bóc tách an toàn.
- [ ] Type check: `cd frontend && npm run check` (0 errors)
- [ ] Build test: `cd frontend && npm run build` (Exit code 0)

**Git Commit:**
```bash
git add frontend/src/lib/utils/markdown.ts frontend/src/components/note/NoteEditor.svelte
git commit -m "sec(markdown): sanitize html rendering against xss and elevate reader typography"
```

**Dependencies:** Task 4  
**Files likely touched:**
- `frontend/src/lib/utils/markdown.ts`
- `frontend/src/components/note/NoteEditor.svelte`

**Estimated scope:** Small (2 files)

---

### Task 6: Tái Thiết Kế NoteEditor Với Chế Độ Split-View & Thanh Công Cụ Định Dạng
**Description:** Đại tu `NoteEditor.svelte`: Hỗ trợ 3 chế độ xem (Chỉ Soạn Thảo, Chỉ Xem Trước, và Xem Song Song Split-View); Bổ sung thanh công cụ định dạng văn bản nhanh (Format Toolbar); Trình quản lý Tags dạng chip badge tương tác xóa/thêm; Phím tắt lưu `Ctrl+S`.

**Acceptance criteria:**
- [ ] Người dùng có thể chuyển đổi mượt mà giữa chế độ Chỉnh sửa, Xem trước và Split View song song.
- [ ] Thanh công cụ định dạng (In đậm, In nghiêng, Tiêu đề, Danh sách, Liên kết) cho phép chèn cú pháp nhanh vào con trỏ.
- [ ] Quản lý Tags hiển thị dạng chip viên thuốc (badge pill) có thể nhấp để thêm/xóa riêng rẽ.
- [ ] Bấm `Ctrl+S` kích hoạt lưu ngay lập tức và phát phản hồi xúc giác.

**Verification & Testing:**
- [ ] Type check: `cd frontend && npm run check` (0 errors)
- [ ] Build test: `cd frontend && npm run build` (Exit code 0)
- [ ] Manual check: Thử gõ nội dung ở Split View, dùng Format bar, thêm bớt tag và nhấn `Ctrl+S` lưu file.

**Git Commit:**
```bash
git add frontend/src/components/note/NoteEditor.svelte frontend/src/lib/stores/viewStore.ts
git commit -m "feat(editor): implement live split-view, markdown formatting toolbar and interactive tag chips"
```

**Dependencies:** Task 5  
**Files likely touched:**
- `frontend/src/components/note/NoteEditor.svelte`
- `frontend/src/lib/stores/viewStore.ts`

**Estimated scope:** Medium (2 files)

---

## Checkpoint 3: Trình Soạn Thảo Đỉnh Cao (Editor Checkpoint)
- [ ] Trình soạn thảo và đọc ghi chú đạt chuẩn trải nghiệm PKM mượt mà, hỗ trợ Split View.
- [ ] Đã triệt tiêu hoàn toàn rủi ro XSS trong WebView2.
- [ ] Đã commit Task 5 và Task 6 vào git history.

---

## Phase 5: Đồ Thị Tri Thức Galaxy, Command Palette & Toast Notifications (Polish & Delight)

### Task 7: Nâng Cấp Đồ Thị Tri Thức Galaxy (GraphView)
**Description:** Nâng cấp `GraphView.svelte` từ các khối hộp cơ bản thành mạng lưới thiên hà tương tác sống động: Node hình tròn phát quang nhẹ với kích thước tỷ lệ theo số lượng liên kết, màu sắc phân biệt theo danh mục/thẻ, bảng điều khiển góc kính mờ sang trọng.

**Acceptance criteria:**
- [ ] Nodes trong đồ thị hiển thị dạng các điểm sáng hành tinh (Planet Nodes) phát quang nhẹ.
- [ ] Đường liên kết (Edges) có hiệu ứng mờ nhạt uyển chuyển, khi hover vào node sẽ làm sáng các kết nối liên quan.
- [ ] MiniMap và Controls được bọc kính mờ đồng bộ với toàn bộ theme.

**Verification & Testing:**
- [ ] Type check: `cd frontend && npm run check` (0 errors)
- [ ] Build test: `cd frontend && npm run build` (Exit code 0)
- [ ] Manual check: Chuyển sang Graph view, kéo thả di chuyển các node, phóng to/thu nhỏ và nhấp vào node để mở bài viết tương ứng.

**Git Commit:**
```bash
git add frontend/src/components/graph/GraphView.svelte frontend/src/lib/utils/graphBuilder.ts
git commit -m "feat(graph): style galaxy-inspired knowledge graph with glowing nodes and smooth physics"
```

**Dependencies:** Task 6  
**Files likely touched:**
- `frontend/src/components/graph/GraphView.svelte`
- `frontend/src/lib/utils/graphBuilder.ts`

**Estimated scope:** Medium (2 files)

---

### Task 8: Tích Hợp Command Palette (Ctrl+K) & Hệ Thống Toast Thông Báo (Sonner Style)
**Description:** Bổ sung Modal Command Palette (`Ctrl+K` / `⌘K`) cho phép tìm kiếm mọi note, mở nhanh view đồ thị, tạo note mới chỉ bằng bàn phím; Tích hợp hệ thống Toast thông báo mini góc màn hình cho các hành động (Đã lưu note, Mở workspace, Sao chép link).

**Acceptance criteria:**
- [ ] Nhấn `Ctrl+K` hoặc bấm vào thanh tìm kiếm trên Navbar sẽ mở Modal Command Palette với backdrop làm mờ màn hình.
- [ ] Có thể dùng phím mũi tên Lên/Xuống và Enter để chọn nhanh ghi chú từ kết quả tìm kiếm.
- [ ] Các hành động lưu tệp, tạo folder hiển thị Toast thông báo trang nhã, tự động biến mất sau 3 giây.

**Verification & Testing:**
- [ ] Type check: `cd frontend && npm run check` (0 errors)
- [ ] Build test: `cd frontend && npm run build` (Exit code 0)
- [ ] Manual check: Nhấn `Ctrl+K` thử tìm và nhảy tới note; lưu note kiểm tra Toast hiển thị mượt mà.

**Git Commit:**
```bash
git add frontend/src/components/layout/CommandPalette.svelte frontend/src/components/ui/Toast.svelte frontend/src/App.svelte
git commit -m "feat(ux): integrate ctrl+k command palette and sonner-style toast notifications"
```

**Dependencies:** Task 7  
**Files likely touched:**
- `frontend/src/components/layout/CommandPalette.svelte` (Tạo mới)
- `frontend/src/components/ui/Toast.svelte` (Tạo mới)
- `frontend/src/App.svelte`

**Estimated scope:** Medium (3 files)

---

## Checkpoint 4 & Final Review: Kiểm Thử Toàn Diện & Code Review 2 Trục
- [ ] Tất cả các bài kiểm tra type `npm run check` và lệnh build `npm run build` hoàn thành với mã thoát 0.
- [ ] Đã thực hiện đầy đủ 8 commits chuẩn hóa theo từng lát cắt mỏng.
- [ ] Tiến hành Code Review theo 2 trục của `mdvault-orchestrator`:
  - [ ] **Trục Standards:** Mã nguồn Svelte 5 sạch sẽ, tuân thủ Clean Architecture, không có smell hay duplicate code.
  - [ ] **Trục Spec:** Đáp ứng trọn vẹn tất cả yêu cầu trong SRS.md (FR-01 đến FR-07) và các tiêu chuẩn phi chức năng NFR.
