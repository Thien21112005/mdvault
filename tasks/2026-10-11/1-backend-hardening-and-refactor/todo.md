# MDVault Backend Hardening & Refactor - Danh Sách Nhiệm Vụ (Todo Checklist)

> **Điều phối bởi:** `mdvault-orchestrator`  
> **Quy chuẩn thực thi:** `incremental-implementation` & `planning-and-task-breakdown`  
> **Quy trình chuẩn mỗi task:** **Code $\rightarrow$ Test & Verify $\rightarrow$ Review $\rightarrow$ Git Commit (Conventional Commits) $\rightarrow$ Next Task**

---

## 📌 Nguyên Tắc Kiểm Thử & Cam Kết Git (Testing & Git Protocol)

1. **Nguyên tắc Kiểm Thử Bắt Buộc:**
   - **Race Detector:** Chạy `go test -race ./...` để đảm bảo 0 cảnh báo Data Race trên mọi goroutine và mutex.
   - **Unit Test Pass:** Mọi test case cho Frontmatter, Watcher, Atomic Save và Cache phải chạy xanh 100%.
   - **Build Verification:** Chạy `go build -v .` để đảm bảo toàn bộ backend biên dịch không lỗi cú pháp hoặc gãy interface Clean Architecture.
2. **Quy tắc Git Commit:**
   - Mỗi lát cắt (Task) sau khi kiểm thử đạt yêu cầu sẽ được commit riêng biệt theo chuẩn **Conventional Commits** (`fix:`, `feat:`, `refactor:`, `test:`, `sec:`, `perf:`, `chore:`).
   - Tuyệt đối không gộp nhiều task vào 1 commit lớn.

---

## Phase 1: Watcher Lifecycle & Debounce (Gia Cố Bộ Lắng Nghe & Chống Bão)

### Task 1: Sửa Vòng Đời Watcher & Tự Động Theo Dõi Thư Mục Ngoại Vi Mới
**Description:** Sửa lỗi vòng đời của `fsnotify.Watcher` trong `infrastructure/filesystem/vault_repository_impl.go`. Cho phép tái khởi động watcher khi mở workspace mới hoặc sau khi `StopWatcher`, đồng thời tự động gọi `r.watcher.Add(newDir)` khi có thư mục mới được tạo từ bên ngoài (VS Code/Explorer). Không quét thư mục dự án khi `rootPath == ""`.

**Acceptance criteria:**
- [x] `StartWatcher` khởi tạo watcher mới an toàn nếu watcher cũ đã bị đóng hoặc chưa tồn tại.
- [x] Khi nhận sự kiện `fsnotify.Create` cho một thư mục, hệ thống tự động đăng ký theo dõi thư mục đó vào `watcher`.
- [x] Không tự động quét thư mục hiện tại (`.`) nếu đường dẫn `rootPath` là chuỗi rỗng.
- [x] Khi đóng watcher (`StopWatcher`), goroutine thoát sạch sẽ mà không để lại rò rỉ (leak).

**Verification & Testing:**
- [x] Tests pass: Kiểm tra hàm `StartWatcher` và `StopWatcher` không gây panic hoặc deadlock (`go test -race ./infrastructure/filesystem/...` PASS).
- [x] Build succeeds: `go build -v .` (exit code 0).
- [x] Manual check: Tạo thư mục mới và file `.md` bên trong từ ngoài, sự kiện được phát hiện chính xác.

**Git Commit:**
```bash
git add infrastructure/filesystem/vault_repository_impl.go infrastructure/filesystem/vault_repository_impl_test.go
git commit -m "fix(watcher): prevent goroutine leak and auto-watch externally created subdirectories"
```
*(Trạng thái: ✅ Đã commit `81d6626`)*

**Dependencies:** None  
**Files likely touched:**
- `infrastructure/filesystem/vault_repository_impl.go`

**Estimated scope:** Small (1-2 files)

---

### Task 2: Triển Khai Debounce 50ms & Triệt Tiêu Callback Trùng Lặp
**Description:** Bổ sung cơ chế Debounce 50ms cho các sự kiện file thay đổi trong `vaultRepositoryImpl` theo đúng chuẩn SRS FR-04. Gom các sự kiện ghi file liên tiếp và loại bỏ các lệnh gọi thủ công `cb(filePath)` thừa thãi trong `SaveNote`, `CreateNote`, `CreateFolder`.

**Acceptance criteria:**
- [x] Thao tác lưu hoặc sửa file chỉ phát đúng 1 sự kiện `on-file-change` sau 50ms kể từ lần ghi cuối cùng.
- [x] Không còn hiện tượng gọi tay `cb(filePath)` trùng lặp với sự kiện của `fsnotify`.
- [x] Lọc bỏ các file tạm (như `.tmp`, `.swp`, file ẩn bắt đầu bằng dấu chấm `.`).

**Verification & Testing:**
- [x] Tests pass: Viết test giả lập ghi file liên tiếp trong 20ms, đảm bảo callback chỉ được gọi duy nhất 1 lần sau 50ms (`go test -race ./infrastructure/filesystem/...` PASS).
- [x] Build succeeds: `go build -v .` (exit code 0).

**Git Commit:**
```bash
git add infrastructure/filesystem/vault_repository_impl.go infrastructure/filesystem/vault_repository_impl_test.go
git commit -m "perf(watcher): implement 50ms debouncing and eliminate duplicate manual file callbacks"
```
*(Trạng thái: ✅ Đã commit `ee68a54`)*

**Dependencies:** Task 1  
**Files likely touched:**
- `infrastructure/filesystem/vault_repository_impl.go`

**Estimated scope:** Small (1 file)

---

## 🛑 Checkpoint: Sau Phase 1 (Watcher & Debounce)
- [x] Watcher hoạt động ổn định, không bắn bão sự kiện.
- [x] Thư mục mới tạo từ bên ngoài được tự động theo dõi.
- [x] `go build -v .` thành công không cảnh báo.

---

## Phase 2: Atomic Persistence & Toàn Vẹn Frontmatter (Bảo Vệ Tệp & Dữ Liệu)

### Task 3: Kỹ Thuật Ghi Tệp Nguyên Tử (Atomic Write: .tmp $\rightarrow$ os.Rename)
**Description:** Thay thế việc gọi trực tiếp `os.WriteFile` bằng hàm `atomicWrite` an toàn trong `infrastructure/filesystem/vault_repository_impl.go`. Ghi nội dung vào `filePath + ".tmp"`, gọi `Sync()`, sau đó dùng `os.Rename` thay thế tệp đích (có cơ chế xử lý tương thích an toàn trên Windows).

**Acceptance criteria:**
- [x] Dữ liệu luôn được ghi vào file tạm `.tmp` và `Sync()` trước khi thay thế file gốc.
- [x] Nếu quá trình ghi thất bại, file gốc giữ nguyên 100% nội dung (không bị 0 KB).
- [x] Tệp tạm tự động được dọn dẹp nếu có lỗi phát sinh.

**Verification & Testing:**
- [x] Tests pass: Viết test ghi đè tệp bằng Atomic Write và kiểm tra nội dung sau khi lưu (`go test -race ./infrastructure/filesystem/...` PASS).
- [x] Build succeeds: `go build -v .` (exit code 0).

**Git Commit:**
```bash
git add infrastructure/filesystem/vault_repository_impl.go infrastructure/filesystem/vault_repository_impl_test.go
git commit -m "fix(storage): implement atomic write with temp file and fsync before rename"
```
*(Trạng thái: ✅ Đã commit `62afe55`)*

**Dependencies:** None  
**Files likely touched:**
- `infrastructure/filesystem/vault_repository_impl.go`

**Estimated scope:** Small (1 file)

---

### Task 4: Bảo Toàn 100% Siêu Dữ Liệu Frontmatter Khi Lưu
**Description:** Cải tiến logic trích xuất và serialize Frontmatter trong `ParseNote` và `SaveNote`. Thay vì chỉ lưu `Title` và `Tags`, hệ thống lưu giữ toàn bộ các key-value tùy biến (ví dụ: `icon`, `date`, `aliases`, `author`) bằng cách lưu trữ map hoặc YAML raw node, đảm bảo khi lưu lại không làm mất các trường có sẵn trong file.

**Acceptance criteria:**
- [x] Mọi trường metadata có trong khối YAML ban đầu (ngoài `title` và `tags`) được bảo toàn nguyên vẹn khi ghi lại.
- [x] Cập nhật chính xác `title` và `tags` mới mà không làm xáo trộn các thuộc tính khác.
- [x] File không có Frontmatter nhưng có H1 vẫn được xử lý chuẩn mực.

**Verification & Testing:**
- [x] Tests pass: Viết test phân tích và lưu file mẫu có `icon: "server"` và `aliases: ["clean-arch"]`, đảm bảo các trường này không bị mất sau khi `SaveNote` (`go test -race ./infrastructure/filesystem/...` PASS).
- [x] Build succeeds: `go build -v .` (exit code 0).
- [x] Type check frontend: `npm run check` (0 errors, 0 warnings).

**Git Commit:**
```bash
git add domain/model/models.go frontend/src/lib/types/domain.ts infrastructure/filesystem/vault_repository_impl.go infrastructure/filesystem/vault_repository_impl_test.go
git commit -m "fix(parser): preserve all custom frontmatter metadata fields during note persistence"
```
*(Trạng thái: ✅ Đã commit `509b4bf`)*

**Dependencies:** Task 3  
**Files likely touched:**
- `domain/model/models.go`
- `infrastructure/filesystem/vault_repository_impl.go`

**Estimated scope:** Small (2 files)

---

## 🛑 Checkpoint: Sau Phase 2 (Data Integrity)
- [x] File markdown khi lưu không còn bị mất các trường YAML mở rộng.
- [x] Kỹ thuật Atomic Write hoạt động trơn tru trên Windows.
- [x] Kiểm tra thực tế trên file `PRB.md` giữ nguyên được `icon: "file-text"`.

---

## Phase 3: In-Memory Caching & Chuẩn Hóa Đường Dẫn (Hiệu Năng & Đồng Bộ)

### Task 5: Triển Khai In-Memory Caching & Tra Cứu Siêu Tốc Trong VaultUseCase
**Description:** Cải tổ `usecase/vault_usecase.go` để duy trì bộ nhớ đệm `cachedData *model.WorkspaceData` trong RAM với khóa `sync.RWMutex`. Hàm `GetCurrentWorkspace()` trả về dữ liệu tức thì (< 1ms) từ RAM, và `SearchNotes()` tìm kiếm trực tiếp trên cache RAM với độ trễ < 5ms thay vì đọc đĩa từ đầu.

**Acceptance criteria:**
- [x] `cachedData` được nạp khi `OpenWorkspace` và được cập nhật khi có sự kiện thay đổi file.
- [x] `GetCurrentWorkspace()` không thực hiện quét lại ổ đĩa nếu cache đang hợp lệ.
- [x] `SearchNotes(query)` thực thi tìm kiếm tức thời trên cache RAM, hỗ trợ không phân biệt hoa thường trên `Title`, `Tags`, và `Content`.
- [x] Đọc/ghi cache được bảo vệ an toàn bằng `sync.RWMutex`.

**Verification & Testing:**
- [x] Tests pass: Benchmark tìm kiếm và đọc 1000 lần hoàn thành chỉ mất 734µs (< 1ms), tìm kiếm hoàn thành < 1.2ms (`go test -race ./usecase/...` PASS).
- [x] Race detection: `go test -race ./usecase/...` (0 cảnh báo với 22 goroutines đồng thời).
- [x] Build succeeds: `go build -v .` (exit code 0).

**Git Commit:**
```bash
git add usecase/vault_usecase.go usecase/vault_usecase_test.go
git commit -m "perf(cache): implement in-memory workspace caching with sync.RWMutex for sub-5ms search"
```
*(Trạng thái: ✅ Đã commit `7c97d64`)*

**Dependencies:** Phase 1, Phase 2  
**Files likely touched:**
- `usecase/vault_usecase.go`

**Estimated scope:** Medium (1-2 files)

---

### Task 6: Chuẩn Hóa Đường Dẫn (filepath.ToSlash) & Quy Ước Folder "root"
**Description:** Chuẩn hóa toàn bộ đường dẫn trên Windows thành định dạng Slash chuẩn `/` (sử dụng `filepath.ToSlash`). Đồng bộ quy ước: Note nằm tại thư mục gốc workspace luôn có thuộc tính `Folder: "root"`, áp dụng nhất quán xuyên suốt giữa `ScanWorkspace`, `CreateNote`, và `GetNote`.

**Acceptance criteria:**
- [x] Không còn dấu gạch chéo ngược Windows `\` trong các trường `ID` và `Folder` trả về cho Frontend.
- [x] Note mới tạo ở thư mục gốc luôn có `Folder = "root"` khớp với quy ước của `ScanWorkspace`.
- [x] Frontend lọc ghi chú theo thư mục (`matchesFolder`) hoạt động chính xác 100%.

**Verification & Testing:**
- [x] Tests pass: Test tạo note ở root và subfolder, kiểm tra trường `Folder` và `ID` (`TestVaultUseCase_PathNormalization_And_RootFolder` và `TestPathNormalization_WindowsNoBackslashes` PASS 100%).
- [x] Race detection: `go test -race ./...` (0 cảnh báo).
- [x] Build succeeds: `go build -v .` (exit code 0).

**Git Commit:**
```bash
git add infrastructure/filesystem/vault_repository_impl.go infrastructure/filesystem/vault_repository_impl_test.go usecase/vault_usecase.go usecase/vault_usecase_test.go
git commit -m "refactor(path): normalize slash formatting on windows and standardize root folder ID"
```
*(Trạng thái: ✅ Đã commit `d3476aa`)*

**Dependencies:** Task 5  
**Files likely touched:**
- `infrastructure/filesystem/vault_repository_impl.go`
- `infrastructure/filesystem/vault_repository_impl_test.go`
- `usecase/vault_usecase.go`
- `usecase/vault_usecase_test.go`

**Estimated scope:** Small (2-4 files)

---

## 🛑 Checkpoint: Sau Phase 3 (Performance & Path Consistency)
- [x] In-memory cache phản hồi tức thì < 1ms và tìm kiếm < 5ms.
- [x] Không còn dấu gạch chéo ngược Windows `\` trong bất kỳ ID hoặc Folder path nào.
- [x] Ghi chú tại gốc luôn mang Folder = "root", lọc ghi chú phía Frontend hoạt động hoàn hảo.

---

## Phase 4: Bảo Mật Path Traversal & Hoàn Thiện CRUD (Security & APIs)

### Task 7: Chặn Đứng Lỗ Hổng Path Traversal (golang-security)
**Description:** Bổ sung lớp kiểm tra an toàn đường dẫn trong `usecase/vault_usecase.go` và `controller/vault_controller.go`. Sử dụng `filepath.Clean` và kiểm tra ràng buộc `strings.HasPrefix(targetPath, absWorkspaceRoot)` để ngăn chặn việc đọc/ghi/tạo tệp ngoài phạm vi Workspace được chỉ định.

**Acceptance criteria:**
- [x] Mọi đường dẫn truyền vào từ IPC được làm sạch bằng `filepath.Clean`.
- [x] Bất kỳ yêu cầu truy cập nào chứa `../` thoát khỏi workspace root đều bị từ chối với mã lỗi an toàn ("truy cập bị từ chối: đường dẫn vượt ra ngoài ranh giới workspace").
- [x] Không cho phép tạo note/folder vượt ra ngoài ranh giới workspace.

**Verification & Testing:**
- [x] Tests pass: Test cố tình truyền `../../system_file.md`, `../../shadow.txt`, malicious folder/note titles với `..` đều bị từ chối (`TestVaultUseCase_PathTraversal_Rejected` PASS).
- [x] Race detection: `go test -race ./...` (0 cảnh báo Data Race).
- [x] Build succeeds: `go build -v .` (exit code 0).

**Git Commit:**
```bash
git add usecase/vault_usecase.go usecase/vault_usecase_test.go controller/vault_controller.go infrastructure/filesystem/vault_repository_impl.go
git commit -m "sec(storage): enforce path traversal boundary validation against workspace root"
```
*(Trạng thái: ✅ Đã commit `6f192b5`)*

**Dependencies:** Phase 3  
**Files likely touched:**
- `usecase/vault_usecase.go`
- `usecase/vault_usecase_test.go`
- `controller/vault_controller.go`
- `infrastructure/filesystem/vault_repository_impl.go`

**Estimated scope:** Small (2-4 files)

---

### Task 8: Bổ Sung Đầy Đủ API CRUD Cho Ghi Chú & Thư Mục (Delete & Rename)
**Description:** Mở rộng giao diện `domain/repository/vault_repository.go` và các phân lớp liên quan để bổ sung các chức năng còn thiếu: `DeleteNote`, `DeleteFolder`, `RenameNote`, `RenameFolder`. Kết nối API vào `controller/vault_controller.go` để sẵn sàng cho Wails IPC.

**Acceptance criteria:**
- [ ] Giao diện `VaultRepository` và `VaultUseCase` có đầy đủ các phương thức `DeleteNote`, `DeleteFolder`, `RenameNote`.
- [ ] Xóa note/folder di chuyển vào thùng rác hoặc xóa an toàn với quyền hạn kiểm soát trong workspace.
- [ ] Đổi tên note tự động cập nhật tên file và tiêu đề, đồng thời cập nhật cache RAM.

**Verification & Testing:**
- [ ] Tests pass: Test xóa note, đổi tên note và kiểm tra tính nhất quán của cache.
- [ ] Build succeeds: `go build -v .`

**Git Commit:**
```bash
git add domain/ domain/repository/ infrastructure/filesystem/ usecase/ controller/
git commit -m "feat(crud): add delete and rename capabilities for notes and folders"
```
*(Trạng thái: ⏳ Chờ commit)*

**Dependencies:** Task 7  
**Files likely touched:**
- `domain/repository/vault_repository.go`
- `infrastructure/filesystem/vault_repository_impl.go`
- `usecase/vault_usecase.go`
- `controller/vault_controller.go`

**Estimated scope:** Medium (4 files)

---

## Phase 5: Kiểm Thử Toàn Diện & Đóng Gói Nghiệm Thu (Testing & Review)

### Task 9: Xây Dựng Bộ Test Suite Tự Động & Kiểm Tra Đa Luồng (-race)
**Description:** Viết toàn bộ các bộ kiểm thử Table-driven test cho `infrastructure/filesystem` và `usecase`. Đạt tỷ lệ bao phủ các kịch bản quan trọng (Frontmatter, Atomic Write, In-memory Search, Path Security) và chạy kiểm tra đồng thời với `go test -race ./...`.

**Acceptance criteria:**
- [ ] Tạo các file test:
  - `infrastructure/filesystem/vault_repository_impl_test.go`
  - `usecase/vault_usecase_test.go`
- [ ] Lệnh `go test -race ./...` vượt qua 100% không có cảnh báo Data Race hoặc leak goroutine.

**Verification & Testing:**
- [ ] `go test -race -v ./...` exit code 0.

**Git Commit:**
```bash
git add infrastructure/filesystem/vault_repository_impl_test.go usecase/vault_usecase_test.go
git commit -m "test(backend): add table-driven test suites and race condition verification"
```
*(Trạng thái: ⏳ Chờ commit)*

**Dependencies:** All previous tasks  
**Files likely touched:**
- `infrastructure/filesystem/vault_repository_impl_test.go`
- `usecase/vault_usecase_test.go`

**Estimated scope:** Medium (2 files)

---

### Task 10: Thẩm Định Nghiệm Thu 2 Trục (Standards & Spec Review)
**Description:** Rà soát lại toàn bộ codebase sau khi hoàn thành theo tiêu chuẩn của skill `code-review`: Trục 1 (Standards - Clean Architecture, Go Concurrency Best Practices) và Trục 2 (Spec - Đáp ứng đầy đủ các yêu cầu trong SRS và PRB).

**Acceptance criteria:**
- [ ] Báo cáo nghiệm thu 2 trục đạt trạng thái PASS.
- [ ] Biên dịch `wails build` hoặc `go build` thành công hoàn toàn.

**Verification & Testing:**
- [ ] `go vet ./...` (0 warnings)
- [ ] `go build -v .` (exit code 0)

**Git Commit:**
```bash
git add tasks/2026-10-11/1-backend-hardening-and-refactor/
git commit -m "chore(review): complete two-axis engineering audit and finalize backend hardening"
```
*(Trạng thái: ⏳ Chờ commit)*

**Dependencies:** Task 9  
**Estimated scope:** Review & Documentation

---

## 🏁 Tổng Kết Tiến Độ Thực Hiện (Progress Tracker)

| Phase | Tên Giai Đoạn | Số Task | Trạng Thái |
| :---: | :--- | :---: | :---: |
| **Phase 1** | Watcher Lifecycle & Debounce | 2 Tasks (T1-T2) | ✅ Hoàn thành (2/2) |
| **Phase 2** | Atomic Persistence & Frontmatter | 2 Tasks (T3-T4) | ✅ Hoàn thành (2/2) |
| **Phase 3** | In-Memory Cache & Path Normalization | 2 Tasks (T5-T6) | ✅ Hoàn thành (2/2) |
| **Phase 4** | Security Hardening & Core CRUD | 2 Tasks (T7-T8) | ⏳ Sẵn sàng thực hiện |
| **Phase 5** | Comprehensive Testing & Code Review | 2 Tasks (T9-T10) | ⏳ Chờ Phase 4 |
