---
name: mdvault-orchestrator
description: Master orchestrator and central brain that coordinates all 18 specialized skills in MDVault. Automatically analyzes user requests, identifies the execution phase (Planning, Architecture, Frontend UI/UX, Go Concurrency, Security, Performance, Testing, Review), and routes work through the proper child skills in disciplined order. Use when starting any task, designing features, writing Go backend or Svelte frontend code, auditing security, or ensuring engineering excellence.
---

# MDVault Brain: Master Engineering Orchestrator

## 1. Tổng Quan & Triết Lý Bộ Não (Executive Overview)

`mdvault-orchestrator` đóng vai trò là **Bộ Não Trung Tâm & Tổng Công Trình Sư (Master Orchestrator)** của hệ thống AI Agent trong dự án MDVault. 

Thay vì để Agent tự do suy đoán hoặc nhảy bổ vào gõ code bừa bãi, skill này thiết lập một **hệ thần kinh điều phối trung tâm**: nó lắng nghe mọi yêu cầu từ người dùng, đối chiếu với bối cảnh dự án (Clean Architecture Go + Svelte 5 + Wails v2 + Local-First), sau đó **tự động kích hoạt đúng tổ hợp kỹ năng chuyên môn** theo quy trình kỹ thuật phần mềm nghiêm ngặt.

```
                      ┌──────────────────────────────────────────────┐
                      │    USER REQUEST / NEW FEATURE / BUG / AUDIT   │
                      └──────────────────────┬───────────────────────┘
                                             │
                                             ▼
                      ┌──────────────────────────────────────────────┐
                      │       [BRAIN] mdvault-orchestrator           │
                      │    Phase Detection & Routing Matrix          │
                      └──────┬───────────────┬────────────────┬──────┘
                             │               │                │
            ┌────────────────┴───┐           │         ┌──────┴─────────────┐
            ▼                    ▼           │         ▼                    ▼
     [GIAI ĐOẠN 1]        [GIAI ĐOẠN 2]      │  [GIAI ĐOẠN 4]        [GIAI ĐOẠN 5]
    LẬP KẾ HOẠCH          KIẾN TRÚC          │   BẢO MẬT & TEST       CODE REVIEW
  planning-and-task-   golang-design-        │  golang-security       code-review
      breakdown            patterns          │  security-and-hardening   (2 Trục)
                       improve-codebase-     │  golang-testing
                          architecture       │  test-driven-development
                                             │
                                             ▼
                                      [GIAI ĐOẠN 3]
                                   THỰC THI LÁT CẮT MỎNG
                                incremental-implementation
                                     ┌───────┴───────┐
                                     ▼               ▼
                              [FRONTEND UI]    [BACKEND GO]
                               impeccable      golang-concurrency
                              apple-design     golang-performance
                                animate
                             emil-design-eng
```

---

## 2. Bản Đồ 18 Kỹ Năng Con Trong Hệ Thần Kinh (Skills Registry)

Bộ Não nắm giữ toàn bộ danh mục 18 skill được cài đặt trong thư mục `.agents/skills/`:

### Nhóm A: Quy Trình & Kế Hoạch (Workflow & Planning)
1. **`planning-and-task-breakdown`** (`.agents/skills/planning-and-task-breakdown`):  
   *Chức năng:* Phân rã yêu cầu lớn thành danh sách các task mỏng có thứ tự (Ordered Tasks), vẽ dependency graph và tiêu chí nghiệm thu (Acceptance Criteria).
2. **`incremental-implementation`** (`.agents/skills/incremental-implementation`):  
   *Chức năng:* Triển khai từng lát cắt mỏng (thin vertical slices). Quy trình: Code $\rightarrow$ Test $\rightarrow$ Verify $\rightarrow$ Commit $\rightarrow$ Next Slice. Cấm viết dồn dập >100 dòng không kiểm chứng.
3. **`code-review`** (`.agents/skills/code-review`):  
   *Chức năng:* Rà soát nghiệm thu theo 2 trục độc lập: **Trục Standards** (chuẩn Clean Arch, smell baseline) và **Trục Spec** (đối chiếu yêu cầu ban đầu).

### Nhóm B: Frontend Polish, Chống AI-Slop & Chuyển Động (UI/UX & Motion)
4. **`impeccable`** (`.agents/skills/impeccable`):  
   *Chức năng:* Tinh chỉnh thẩm mỹ cao cấp, chống giao diện AI-slop, xử lý visual hierarchy, whitespace, typography, layout critique.
5. **`apple-design`** (`.agents/skills/apple-design`):  
   *Chức năng:* Thiết kế giao diện phong cách Apple: vật liệu kính mờ (glassmorphism/translucent), độ sâu (depth), đàn hồi vật lý (spring physics).
6. **`animate`** (`.agents/skills/animate`):  
   *Chức năng:* Xây dựng animation từ số 0, chuẩn hóa easing curves, spring parameters, duration, interruptible state.
7. **`emil-design-eng`** (`.agents/skills/emil-design-eng`):  
   *Chức năng:* Tiêu chuẩn Design Engineering của Emil Kowalski, tinh chỉnh các chi tiết vi mô vô hình giúp app có cảm giác đắt tiền.
8. **`frontend-design`** (`.agents/skills/frontend-design`):  
   *Chức năng:* Chuẩn thiết kế hiện đại từ Anthropic, phối màu tương phản, bố cục không gian sạch sẽ.
9. **`design-taste-frontend`** (`.agents/skills/design-taste-frontend`):  
   *Chức năng:* Anti-slop cho web apps, ngăn chặn sử dụng các template giao diện rập khuôn.

### Nhóm C: Lõi Backend Go & Kiến Trúc (Go Backend & Clean Architecture)
10. **`golang-design-patterns`** (`.agents/skills/golang-design-patterns`):  
    *Chức năng:* Triển khai Clean Architecture (Domain $\rightarrow$ UseCase $\rightarrow$ Infrastructure $\rightarrow$ Controller), Dependency Inversion, Functional Options, Interface Segregation.
11. **`golang-concurrency`** (`.agents/skills/golang-concurrency`):  
    *Chức năng:* Quản trị goroutines, `sync.Mutex`/`sync.RWMutex`, kênh channels, `fsnotify` daemon watcher, triệt tiêu race condition và deadlock.
12. **`golang-performance`** (`.agents/skills/golang-performance`):  
    *Chức năng:* Tối ưu hóa bộ nhớ RAM < 50MB, giảm thiểu allocation rác, tối ưu thuật toán tìm kiếm toàn văn và duyệt đệ quy `WalkDir`.
13. **`improve-codebase-architecture`** (`.agents/skills/improve-codebase-architecture`):  
    *Chức năng:* Quét ranh giới kiến trúc module, triệt tiêu circular dependencies, bảo vệ cấu trúc Clean Architecture.

### Nhóm D: Bảo Mật & Kiểm Thử (Security & Testing)
14. **`golang-security`** (`.agents/skills/golang-security`):  
    *Chức năng:* Chống Path Traversal khi chọn file/thư mục, audit an toàn file I/O, kiểm tra race condition với `-race`, phân tích tĩnh với `gosec`.
15. **`security-and-hardening`** (`.agents/skills/security-and-hardening`):  
    *Chức năng:* Bảo vệ Sandbox WebView2, chống XSS khi render markdown HTML, cô lập untrusted markdown input.
16. **`golang-testing`** (`.agents/skills/golang-testing`):  
    *Chức năng:* Viết Table-driven tests, benchmark profiling, phát hiện goroutine leak với `uber-go/goleak`.
17. **`test-driven-development`** (`.agents/skills/test-driven-development`):  
    *Chức năng:* Chu trình Red-Green-Refactor cho các hàm logic nghiệp vụ nhạy cảm.
18. **`find-skills`** (`.agents/skills/find-skills`):  
    *Chức năng:* Khám phá và cài đặt bổ sung kỹ năng mới từ hệ sinh thái mở `skills.sh`.

---

## 3. Ma Trận Điều Phối Tự Động (Automatic Routing Matrix)

Khi nhận yêu cầu, Bộ Não quét từ khóa và ngữ cảnh để kích hoạt đúng pipeline:

| Loại Yêu Cầu / Tình Huống | Pipeline Kỹ Năng Kích Hoạt Theo Thứ Tự | Kết Quả Đầu Ra Bắt Buộc |
| :--- | :--- | :--- |
| **Tính năng mới toàn diện (Full Feature)** | `planning-and-task-breakdown` $\rightarrow$ `incremental-implementation` $\rightarrow$ Backend/Frontend Skills $\rightarrow$ `golang-testing` $\rightarrow$ `code-review` | Bản kế hoạch `tasks/plan.md`, code từng lát cắt, test chạy xanh, review 2 trục PASS |
| **Làm giao diện / Sửa UI / Thêm Component** | `impeccable` $\rightarrow$ `apple-design` $\rightarrow$ `animate` $\rightarrow$ `emil-design-eng` | UI đạt chuẩn chống AI-slop, có glassmorphism/depth, micro-animations mượt mà |
| **Sửa file watcher / Đa luồng / IPC backend** | `golang-concurrency` $\rightarrow$ `golang-security` $\rightarrow$ `golang-performance` $\rightarrow$ `golang-testing` | Bọc `sync.Mutex` an toàn, test `-race` không cảnh báo, không memory leak |
| **Kiểm tra an toàn / Rà soát lỗ hổng** | `golang-security` $\rightarrow$ `security-and-hardening` | Báo cáo kiểm toán Path Traversal, WebView Sandbox, File I/O |
| **Tối ưu tốc độ / Giảm ngốn RAM** | `golang-performance` $\rightarrow$ `golang-testing` (Benchmark) | Báo cáo RAM < 50MB, tốc độ tìm kiếm < 20ms |
| **Refactor / Dọn dẹp cấu trúc dự án** | `improve-codebase-architecture` $\rightarrow$ `golang-design-patterns` $\rightarrow$ `code-review` | Cấu trúc Clean Architecture chuẩn, không circular import |

---

## 4. Năm Quy Trình Tác Chiến Tiêu Chuẩn (5 Standard Playbooks)

### Playbook 1: Phát Triển Tính Năng Mới (End-to-End Feature Playbook)
1. **Bước 1 (Planning):** Kích hoạt `planning-and-task-breakdown`. Đọc [SRS.md](markdown/SRS.md) và [HLD.md](markdown/HLD.md). Phân rã thành các Task 1, 2, 3... Ghi nhận tiêu chí kiểm thử cho từng task. **Tuyệt đối không viết code trong bước này.**
2. **Bước 2 (Contract First):** Kích hoạt `golang-design-patterns`. Khai báo structs trong `domain/model/models.go` và interface trong `domain/repository/`. Đồng bộ types sang `frontend/src/lib/types/domain.ts`.
3. **Bước 3 (Slicing):** Kích hoạt `incremental-implementation`. Thực hiện từng task:
   - Viết test trước với `test-driven-development`.
   - Triển khai Go backend với `golang-concurrency` (nếu có I/O hoặc Goroutine).
   - Triển khai UI Svelte 5 với `impeccable` + `apple-design`.
   - Chạy test xác minh lát cắt hoạt động.
4. **Bước 4 (Review):** Kích hoạt `code-review`. Đối chiếu Standards và Spec.

### Playbook 2: Tinh Chỉnh Giao Diện Chống AI-Slop (UI Polish Playbook)
1. **Bước 1 (Audit & Critique):** Dùng `impeccable` soi xét layout hiện tại: kiểm tra độ tương phản, khoảng thở giữa các card, padding, font size hierarchy.
2. **Bước 2 (Material & Depth):** Dùng `apple-design` bổ sung hiệu ứng kính mờ (backdrop-blur), viền mỏng tinh tế (`border-white/10`), đổ bóng đa tầng (layered shadow).
3. **Bước 3 (Motion & Feedback):** Dùng `animate` và `emil-design-eng` gắn hiệu ứng hover nhẹ nhàng (scale 1.01, duration 150ms, spring curve), active tactile feedback khi click nút hoặc chọn note.

### Playbook 3: Xử Lý Đa Luồng & Watcher Backend (Concurrency Playbook)
1. **Bước 1 (Thread Safety):** Kích hoạt `golang-concurrency`. Kiểm tra mọi biến dùng chung (`currentWsPath`, `currentData`) đã được bảo vệ bằng `sync.Mutex` chưa.
2. **Bước 2 (Path Validation):** Kích hoạt `golang-security`. Áp dụng `filepath.Clean` và kiểm tra `strings.HasPrefix` để triệt tiêu lỗ hổng Path Traversal (`../`).
3. **Bước 3 (Race Test):** Kích hoạt `golang-testing`. Chạy lệnh `go test -race ./...` để bảo đảm không xảy ra xung đột đọc/ghi đồng thời.

---

## 5. Luật Thép Bắt Buộc Của Bộ Não (Non-Negotiable Rules)

1. **Không có Kế hoạch thì Không có Code:** Mọi yêu cầu chạm vào từ 2 file trở lên bắt buộc phải gọi `planning-and-task-breakdown` để chốt roadmap trước khi gõ code.
2. **Tuyệt đối không sinh ra AI-Slop UI:** Không dùng màu sắc default thô ráp, không dùng layout rập khuôn; mọi component UI đều phải tuân theo chuẩn của `impeccable` và `apple-design`.
3. **Bảo toàn Mutex Guard:** Mọi tác vụ I/O đĩa hoặc watcher trong Go phải được bọc trong `sync.Mutex` theo chỉ dẫn của `golang-concurrency`.
4. **Không tự xưng "xong" khi chưa Review:** Sau khi hoàn thành code, Agent bắt buộc phải thực hiện tự rà soát theo tiêu chuẩn 2 trục của `code-review`.
