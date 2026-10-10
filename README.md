# MDVault

<p align="center">
  <strong>Lightweight Desktop Markdown Vault & Personal Knowledge Graph</strong>
</p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Backend-Go%201.22+-00ADD8?style=flat-square&logo=go" alt="Go Version"></a>
  <a href="https://wails.io/"><img src="https://img.shields.io/badge/Framework-Wails%20v2-DF0000?style=flat-square" alt="Wails Version"></a>
  <a href="https://svelte.dev/"><img src="https://img.shields.io/badge/Frontend-Svelte%205%20%2B%20Vite-FF3E00?style=flat-square&logo=svelte" alt="Svelte 5"></a>
  <a href="https://tailwindcss.com/"><img src="https://img.shields.io/badge/Styling-Tailwind%20CSS%20v4-06B6D4?style=flat-square&logo=tailwindcss" alt="Tailwind CSS"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="MIT License"></a>
</p>

<p align="center">
  <a href="#english">English</a> • <a href="#tiếng-việt">Tiếng Việt</a>
</p>

---

<a name="english"></a>
## English

### Overview

**MDVault** is a local-first, high-performance desktop application for Markdown note-taking and Personal Knowledge Management (PKM). Designed as a lightweight alternative to memory-heavy Electron solutions, MDVault combines a compiled **Go 1.22+** backend with the operating system's native WebView engine (Edge WebView2 on Windows, WebKit on macOS/Linux) and a reactive **Svelte 5** frontend.

### Core Architectural Principles

- **Local-First & Privacy by Design:** Notes remain on the local file system in standard `.md` format with YAML frontmatter. No cloud account required, zero network telemetry, and full compatibility with Git, Vim, and external text editors.
- **Resource Efficiency:** Idles under 50MB RAM and launches in under 300ms, using native WebView rendering instead of embedding an entire Chromium runtime.
- **Clean Architecture (Go Backend):** Strict separation of concerns across Domain, UseCase, Infrastructure, and Controller layers following Uncle Bob's Clean Architecture pattern.
- **Kernel-Level Realtime File Synchronization:** Uses `fsnotify` to listen for OS filesystem events directly, propagating external edits to the UI in less than 20ms without CPU-intensive polling.
- **Interactive Knowledge Graph:** Automatically parses bidirectional `[[Wiki-link]]` syntax and constructs an interactive 2D graph powered by `@xyflow/svelte`.
- **Atomic File Persistence:** Employs temporary file writes followed by atomic renames (`os.Rename`), preventing file corruption during unexpected power loss.

### Performance Benchmark

| Metric | MDVault (Go + Svelte 5) | Obsidian | Notion Desktop | Logseq |
| :--- | :---: | :---: | :---: | :---: |
| **Idle RAM Footprint** | **~35 - 50 MB** | 250 - 450 MB | 400 - 800 MB | 350 - 600 MB |
| **Cold Startup Time** | **< 300 ms** | 1.5 - 3.5 s | 3.0 - 6.0 s | 2.5 - 5.0 s |
| **Runtime Architecture** | **Wails v2 + Native WebView** | Electron | Electron | Electron |
| **Data Storage** | **100% Local (.md files)** | Local (.md files) | Cloud (Proprietary) | Local (.md/.org) |
| **Executable Size** | **~15 - 20 MB** | ~120 MB | ~150 MB | ~180 MB |

### Engineering Documentation Suite

Comprehensive technical specifications are available in the [`markdown/`](markdown/) directory:

| Document | Title | Description |
| :--- | :--- | :--- |
| [BRD.md](markdown/BRD.md) | Business Requirements Document | Problem statement, competitive analysis, scope, and business objectives |
| [PRB.md](markdown/PRB.md) | Product Requirements Brief | Personas, user stories, Gherkin acceptance criteria, and UI wireframe flows |
| [SRS.md](markdown/SRS.md) | Software Requirements Specification | IEEE 830 functional (FR-01 to FR-07) and non-functional requirements |
| [HLD.md](markdown/HLD.md) | High-Level Design | Clean Architecture blueprint, sequence diagrams, and cross-cutting concerns |
| [LLD.md](markdown/LLD.md) | Low-Level Design & API Specs | Go structs, Svelte stores, and complete Wails IPC API contracts |

### Technology Stack

- **Backend:** Go 1.22+ (`goccy/go-yaml`, `fsnotify/fsnotify`, `wailsapp/wails/v2`)
- **Frontend:** Svelte 5, TypeScript 5.x, Vite 5.x, Tailwind CSS v4
- **Markdown & Graph:** `markdown-it`, `@tailwindcss/typography`, `@xyflow/svelte`
- **Application Shell:** Wails v2 (Go IPC binding to OS WebView)

### Project Structure

```
MDVault/
├── .github/                 # Issue templates and PR templates
├── controller/              # Wails IPC controller bindings
├── domain/                  # Clean Architecture: entities, models, repository interfaces
│   ├── model/               # Note, Folder, WorkspaceData, Frontmatter
│   └── repository/          # Repository contracts
├── usecase/                 # Business logic, in-memory search, workspace orchestration
├── infrastructure/          # Physical implementations (filesystem scanner, fsnotify watcher)
├── frontend/                # Svelte 5 single page application
│   └── src/
│       ├── components/      # Feature modules (Navbar, Sidebar, NoteEditor, GraphView)
│       ├── lib/
│       │   ├── stores/      # Reactive Svelte stores (vaultStore, viewStore)
│       │   ├── utils/       # Markdown pipeline and Wiki-link graphBuilder
│       │   └── wails/       # Anti-Corruption Layer (WailsAdapter)
│       └── assets/          # Theme tokens and stylesheet architecture
├── markdown/                # Official technical specifications (BRD, PRB, SRS, HLD, LLD)
├── main.go                  # Application entrypoint & dependency injection wire-up
└── wails.json               # Wails build configuration
```

### Getting Started

#### Prerequisites
- **Go 1.22+**: [Download Go](https://go.dev/dl/)
- **Node.js 18+ or 20+ (LTS)**: [Download Node.js](https://nodejs.org/)
- **Wails CLI v2**:
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

#### Development Mode
Run the application with Hot Module Replacement (HMR) and Go live compile:
```bash
wails dev
```

#### Production Build
Build a standalone, stripped executable:
```bash
wails build -clean -s -o mdvault.exe
```
The compiled binary will be placed at `build/bin/mdvault.exe`.

---

<a name="tiếng-việt"></a>
## Tiếng Việt

### Tổng Quan

**MDVault** là ứng dụng máy tính cá nhân phục vụ quản lý ghi chú Markdown và xây dựng mạng lưới tri thức cá nhân (Personal Knowledge Management - PKM) theo triết lý **Local-First**. Được thiết kế nhằm thay thế các giải pháp cồng kềnh nền Electron, MDVault kết hợp lõi xử lý nhị phân bằng **Go 1.22+** với engine WebView gốc của hệ điều hành (Microsoft Edge WebView2 trên Windows, WebKit trên macOS/Linux) và giao diện phản ứng siêu tốc **Svelte 5**.

### Các Điểm Nhấn Kiến Trúc Cốt Lõi

- **Tự Chủ Dữ Liệu & Bảo Mật Tuyệt Đối:** Toàn bộ ghi chú được lưu trữ cục bộ dưới định dạng `.md` chuẩn kèm YAML Frontmatter. Không yêu cầu tài khoản, không gửi dữ liệu ra ngoài Internet, tương thích hoàn toàn với Git, Vim, VS Code.
- **Tối Ưu Tài Nguyên:** Mức tiêu thụ bộ nhớ tĩnh dưới 50MB RAM và thời gian khởi động dưới 300ms nhờ cơ chế tận dụng WebView dùng chung của hệ điều hành thay vì nhúng Chromium riêng.
- **Kiến Trúc Sạch (Clean Architecture Go):** Tách bạch rõ ràng giữa 4 phân lớp: Domain, UseCase, Infrastructure và Controller theo tiêu chuẩn kỹ thuật phần mềm công nghiệp.
- **Đồng Bộ Tập Tin Thời Gian Thực:** Sử dụng `fsnotify` kết nối trực tiếp với sự kiện nhân hệ điều hành, cập nhật thay đổi file từ bên ngoài lên giao diện trong dưới 20ms mà không cần quét tuần tự (polling).
- **Đồ Thị Tri Thức Tương Tác:** Tự động giải mã liên kết hai chiều `[[Wiki-link]]` để dựng đồ thị tương tác 2D thông qua `@xyflow/svelte`.
- **Lưu Tệp Nguyên Tử (Atomic Write):** Ghi file tạm rồi thực hiện hoán đổi tên nguyên tử (`os.Rename`), triệt tiêu hoàn toàn rủi ro hỏng tệp khi ứng dụng bị tắt đột ngột.

### Bảng So Sánh Hiệu Năng

| Tiêu chí | MDVault (Go + Svelte 5) | Obsidian | Notion Desktop | Logseq |
| :--- | :---: | :---: | :---: | :---: |
| **Tiêu hao RAM (chờ)** | **~35 - 50 MB** | 250 - 450 MB | 400 - 800 MB | 350 - 600 MB |
| **Thời gian khởi động** | **< 300 ms** | 1.5 - 3.5 s | 3.0 - 6.0 s | 2.5 - 5.0 s |
| **Kiến trúc Runtime** | **Wails v2 + Native WebView** | Electron | Electron | Electron |
| **Lưu trữ dữ liệu** | **100% Cục bộ (.md)** | Cục bộ (.md) | Đám mây (Định dạng đóng) | Cục bộ (.md/.org) |
| **Dung lượng file chạy** | **~15 - 20 MB** | ~120 MB | ~150 MB | ~180 MB |

### Hệ Thống Hồ Sơ Kỹ Thuật

Bộ tài liệu kỹ thuật chuẩn hóa nằm trong thư mục [`markdown/`](markdown/):

| Tài liệu | Tiêu đề | Tóm tắt nội dung |
| :--- | :--- | :--- |
| [BRD.md](markdown/BRD.md) | Business Requirements Document | Bài toán nghiệp vụ, phân tích đối thủ cạnh tranh, phạm vi dự án và mục tiêu đo lường |
| [PRB.md](markdown/PRB.md) | Product Requirements Brief | Chân dung người dùng, User Stories, tiêu chí nghiệm thu Gherkin, luồng giao diện |
| [SRS.md](markdown/SRS.md) | Software Requirements Specification | Đặc tả yêu cầu phần mềm IEEE 830: 7 yêu cầu chức năng và 4 nhóm phi chức năng |
| [HLD.md](markdown/HLD.md) | High-Level Design | Thiết kế kiến trúc tổng thể, mô hình Clean Architecture, Sequence Diagrams |
| [LLD.md](markdown/LLD.md) | Low-Level Design & API Specs | Cấu trúc dữ liệu chi tiết, Svelte stores, và toàn bộ 8 phương thức Wails IPC API |

### Hướng Dẫn Cài Đặt & Vận Hành

#### Yêu Cầu Môi Trường
- **Go 1.22+**: [Cài đặt Go](https://go.dev/dl/)
- **Node.js 18+ hoặc 20+ (LTS)**: [Cài đặt Node.js](https://nodejs.org/)
- **Wails CLI v2**:
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

#### Chế Độ Phát Triển (Development)
Khởi chạy ứng dụng với tính năng Hot Module Replacement (HMR) và tự động biên dịch Go:
```bash
wails dev
```

#### Đóng Gói Xuất Xưởng (Production)
Biên dịch thành file thực thi độc lập đã tối ưu kích thước:
```bash
wails build -clean -s -o mdvault.exe
```
File thực thi độc lập sẽ được tạo tại thư mục `build/bin/mdvault.exe`.

---

## Contributing

Please review [CONTRIBUTING.md](CONTRIBUTING.md) for details on code standards, commit formatting, and pull request procedures.

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
