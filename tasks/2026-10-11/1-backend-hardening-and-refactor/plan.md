# Kế Hoạch Cải Tạo & Gia Cố Toàn Diện Backend Go (Backend Hardening & Refactor Master Plan)

> **Điều phối bởi:** `mdvault-orchestrator` (Skill Mẹ)  
> **Các kỹ năng chuyên môn tham gia:**  
> - `planning-and-task-breakdown` (Lập lộ trình, chia nhỏ lát cắt mỏng, tiêu chí nghiệm thu)  
> - `golang-design-patterns` (Clean Architecture, In-Memory Caching, Interface Segregation, Contract First)  
> - `golang-concurrency` (Quản lý Goroutines, sync.RWMutex, fsnotify lifecycle, debouncing 50ms, triệt tiêu race condition)  
> - `golang-performance` (Tối ưu RAM < 50MB, tra cứu In-memory < 5ms, giảm thiểu cấp phát bộ nhớ rác)  
> - `golang-security` (Chống Path Traversal `filepath.Clean`, Atomic Write chống hỏng file)  
> - `golang-testing` & `test-driven-development` (Bộ kiểm thử Table-driven, cờ kiểm tra đa luồng `-race`, 100% test pass)  
> - `incremental-implementation` & `code-review` (Triển khai từng lát cắt & thẩm định 2 trục)  
>  
> **Phiên bản:** 2.1.0-PRO Backend  
> **Mục tiêu:** Khắc phục toàn bộ các lỗi tồn đọng ở backend Go (mất dữ liệu Frontmatter, ghi tệp không nguyên tử, bão sự kiện Watcher, thiếu in-memory cache, lệch path Windows, thiếu CRUD và 0% test coverage) để hệ thống đạt chuẩn công nghiệp ổn định và an toàn tuyệt đối.

---

## 1. Báo Cáo Thẩm Định Hiện Trạng (Backend Codebase Audit)

Đối chiếu với các tài liệu chuẩn kỹ thuật ([SRS.md](../../../markdown/SRS.md), [HLD.md](../../../markdown/HLD.md), [LLD.md](../../../markdown/LLD.md)), Backend hiện tại có các lỗi nghiêm trọng sau:

| Mã Vấn Đề | Hạng Mục | Hiện Trạng Code | Hậu Quả & Rủi Ro | Giải Pháp Chuẩn Kiến Trúc |
| :---: | :--- | :--- | :--- | :--- |
| **ISSUE-01** | **Bảo toàn Frontmatter** | `SaveNote` chỉ lưu `Title` và `Tags`, bỏ qua các trường YAML khác. | Xóa sạch các metadata như `icon`, `aliases`, `date`, `author` có sẵn trong file `.md`. | Lưu trữ toàn vẹn Frontmatter (`map[string]interface{}` hoặc raw YAML injection) để bảo toàn 100% thuộc tính. |
| **ISSUE-02** | **Atomic Persistence** | `SaveNote` dùng `os.WriteFile` ghi đè trực tiếp. | Nếu mất điện hoặc crash giữa chừng, file bị rỗng (0 KB corruption). | Triển khai kỹ thuật Atomic Write: Ghi file `.tmp` $\rightarrow$ `Sync()` $\rightarrow$ `os.Rename`. |
| **ISSUE-03** | **Bão Sự Kiện Watcher** | `SaveNote`/`CreateNote` tự gọi `cb()` bằng tay, kết hợp với OS `fsnotify` bắn sự kiện thêm lần nữa. | Phát sinh 2-3 sự kiện `on-file-change` liên tiếp, frontend quét lại đĩa 3 lần, gây giật lag và nhảy con trỏ. | Loại bỏ gọi tay; bổ sung **Debounce 50ms** bằng timer channel để gom các sự kiện ghi dồn dập. |
| **ISSUE-04** | **Thư Mục Ngoại Vi** | `listenEvents` không gọi `watcher.Add` khi có thư mục mới được tạo từ ngoài (VS Code/Explorer). | Trên Windows, các file tạo trong thư mục con mới sẽ hoàn toàn không được phát hiện. | Tự động kiểm tra `IsDir()` khi nhận sự kiện `Create` và tự đăng ký `r.watcher.Add(newDir)`. |
| **ISSUE-05** | **Vòng Đời Watcher** | Khi gọi `StopWatcher()`, `watcher.Close()` làm goroutine kết thúc vĩnh viễn. | Người dùng mở sang Workspace khác (`StartWatcher` lần 2) thì Watcher ngừng hoạt động hoàn toàn. | Hỗ trợ tái khởi tạo watcher và quản lý goroutine bằng context / channel đóng an toàn. |
| **ISSUE-06** | **In-Memory Cache** | `VaultUseCase` không lưu `cachedData`, mỗi lần gọi `SearchNotes` hay `ScanFolder` lại đọc đĩa từ đầu. | Tốn CPU/Disk I/O không cần thiết; Frontend phải bỏ backend để tự lọc JS. | Lưu `currentData` trên RAM bảo vệ bằng `sync.RWMutex`; tìm kiếm < 5ms trực tiếp trên bộ nhớ. |
| **ISSUE-07** | **Chuẩn Hóa Đường Dẫn** | Trên Windows dùng lẫn `\` và `/`; `folderID` trong `CreateNote` và `GetNote` không khớp với `ScanWorkspace`. | Frontend lọc folder bị sai lệch, note tạo ở root không nhận diện đúng. | Chuẩn hóa tất cả bằng `filepath.ToSlash`; quy ước note ở root có `folderID = "root"`. |
| **ISSUE-08** | **Lỗ Hổng Bảo Mật** | Chưa sanitize đường dẫn từ IPC truyền vào (`filePath`, `folderPath`). | Rủi ro Path Traversal (`../../etc/passwd` hoặc ghi đè file hệ thống ngoài workspace). | Kiểm tra `filepath.Clean` và ràng buộc `strings.HasPrefix(targetPath, absWorkspaceRoot)`. |
| **ISSUE-09** | **Thiếu Core CRUD** | Chưa có các hàm `DeleteNote`, `DeleteFolder`, `RenameNote`, `RenameFolder`. | Người dùng không thể xóa hoặc đổi tên note/folder từ giao diện. | Bổ sung đầy đủ giao diện và implementation CRUD từ Domain $\rightarrow$ UseCase $\rightarrow$ Controller. |
| **ISSUE-10** | **Kiểm Thử (Testing)** | 0% Test Coverage (`[no test files]`). | Không có cơ chế bảo vệ chống hồi quy (regression) và không phát hiện được Data Race. | Xây dựng bộ Unit Test và chạy kiểm tra đa luồng bằng cờ `go test -race ./...`. |

---

## 2. Thiết Kế Kiến Trúc & Quyết Định Kỹ Thuật (Architecture Decisions)

### 2.1. Quản lý trạng thái Cache an toàn đồng thời (`sync.RWMutex`)
Trong `usecase/vault_usecase.go`:
```go
type vaultUseCaseImpl struct {
    repo          repository.VaultRepository
    currentWsPath string
    cachedData    *model.WorkspaceData
    mu            sync.RWMutex
    onChange      func(path string)
}
```
- Thao tác đọc (`GetCurrentWorkspace`, `SearchNotes`): Sử dụng `mu.RLock()` cho phép nhiều luồng đọc đồng thời mà không bị block.
- Thao tác ghi/làm mới cache (`OpenWorkspace`, cập nhật sau khi nhận file change): Sử dụng `mu.Lock()`.

### 2.2. Cơ chế Debounce 50ms cho Watcher
Sử dụng `sync.Mutex` kết hợp `*time.Timer`:
- Khi nhận sự kiện ghi file từ `fsnotify`, nếu sự kiện cùng một file xảy ra liên tiếp trong 50ms, bộ đếm timer được reset.
- Chỉ khi hết 50ms yên tĩnh, callback mới được phát lên Wails IPC (`on-file-change`).
- Loại bỏ hoàn toàn các lệnh gọi thủ công `cb(filePath)` ở tầng ghi đĩa.

### 2.3. Quy trình Ghi nguyên tử (Atomic Write Protocol)
1. Xác định đường dẫn file tạm: `tmpPath := filePath + ".tmp"`
2. Tạo và ghi dữ liệu ra file tạm với quyền `0644`.
3. Gọi `tmpFile.Sync()` để xả dữ liệu từ bộ đệm OS xuống phần cứng đĩa.
4. Đóng file tạm.
5. Trên Windows, nếu file đích đã tồn tại, dùng chiến lược hoán đổi/ghi đè an toàn thông qua `os.Rename(tmpPath, filePath)`. Nếu lỗi, dọn dẹp file tạm.

### 2.4. Phân tích & Bảo toàn Frontmatter Linh Hoạt
Bóc tách phần đầu file:
- Đọc raw frontmatter string giữa cặp `---`.
- Sử dụng `yaml.Node` hoặc `map[string]interface{}` để lưu trữ toàn bộ các cặp key-value tùy biến.
- Khi lưu lại: Cập nhật `title` và `tags` trong map, giữ nguyên tất cả các key còn lại (như `icon`, `date`, `aliases`), sau đó marshal lại YAML header.

---

## 3. Lộ Trình Triển Khai Theo Lát Cắt Mỏng (Incremental Slices)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ PHASE 1: WATCHER LIFECYCLE & DEBOUNCE (Gia Cố Bộ Lắng Nghe & Chống Bão)     │
│ Task 1: Sửa vòng đời Watcher & Auto-add thư mục mới                         │
│ Task 2: Debounce 50ms & Triệt tiêu callback trùng lặp                       │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
┌──────────────────────────────────────▼──────────────────────────────────────┐
│ PHASE 2: ATOMIC PERSISTENCE & FRONTMATTER (Bảo Toàn Tệp & Dữ Liệu)          │
│ Task 3: Kỹ thuật Atomic Write (.tmp -> os.Rename)                           │
│ Task 4: Giữ nguyên 100% metadata Frontmatter khi Save                       │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
┌──────────────────────────────────────▼──────────────────────────────────────┐
│ PHASE 3: IN-MEMORY CACHE & PATH NORMALIZATION (Hiệu Năng & Đồng Bộ)        │
│ Task 5: In-Memory Caching (sync.RWMutex) & Tra cứu < 5ms                   │
│ Task 6: Chuẩn hóa đường dẫn ToSlash & Quy tắc Folder "root"                 │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
┌──────────────────────────────────────▼──────────────────────────────────────┐
│ PHASE 4: SECURITY HARDENING & CORE CRUD (Bảo Mật & Tính Năng Còn Thiếu)     │
│ Task 7: Chặn đứng Path Traversal (Clean & Prefix check)                     │
│ Task 8: Bổ sung API Delete/Rename cho Note & Folder                         │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
┌──────────────────────────────────────▼──────────────────────────────────────┐
│ PHASE 5: COMPREHENSIVE TESTING & CODE REVIEW (Kiểm Thử Toàn Diện)           │
│ Task 9: Bộ Unit Test & Race Test (go test -race ./...)                      │
│ Task 10: Thẩm định 2 trục (Standards & Spec) theo chuẩn code-review         │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Quy Ước Cam Kết Mã Nguồn (Conventional Commits Protocol)

Theo chỉ dẫn của `incremental-implementation` và `mdvault-orchestrator`, mỗi task sau khi kiểm thử đạt yêu cầu sẽ được commit riêng biệt theo chuẩn Conventional Commits:

| Task | Loại Commit | Thông Điệp Commit Chuẩn (Conventional Commit Message) |
| :---: | :---: | :--- |
| **Task 1** | `fix` | `fix(watcher): prevent goroutine leak and auto-watch externally created subdirectories` |
| **Task 2** | `perf` | `perf(watcher): implement 50ms debouncing and eliminate duplicate manual file callbacks` |
| **Task 3** | `fix` | `fix(storage): implement atomic write with temp file and fsync before rename` |
| **Task 4** | `fix` | `fix(parser): preserve all custom frontmatter metadata fields during note persistence` |
| **Task 5** | `perf` | `perf(cache): implement in-memory workspace caching with sync.RWMutex for sub-5ms search` |
| **Task 6** | `refactor` | `refactor(path): normalize slash formatting on windows and standardize root folder ID` |
| **Task 7** | `sec` | `sec(storage): enforce path traversal boundary validation against workspace root` |
| **Task 8** | `feat` | `feat(crud): add delete and rename capabilities for notes and folders` |
| **Task 9** | `test` | `test(backend): add table-driven test suites and race condition verification` |
| **Task 10** | `chore` | `chore(review): complete two-axis engineering audit and finalize backend hardening` |

