---
title: "Thiết Kế Kiến Trúc Tổng Thể (HLD) - MDVault"
tags: ["documentation", "architecture", "hld", "system-design", "clean-architecture"]
icon: "layers"
---

# MDVault - High-Level Design (HLD) & System Architecture Specification

> **Mã tài liệu:** HLD-MDVAULT-ARCH-2026  
> **Dự án:** MDVault - Professional Desktop Markdown Vault  
> **Phiên bản:** 2.0.0-PRO (Industrial Standard)  
> **Trạng thái:** Approved Architecture  
> **Kiến trúc cốt lõi:** Clean Architecture (Go) + Feature-Driven Design Modular SPA (Svelte 5)  

---

## 1. Tầm Nhìn Kiến Trúc & Triết Lý Thiết Kế (Architectural Vision & Principles)

**MDVault** được định hình dựa trên triết lý **Local-First, Lightweight & Zero-Bloat** (Tự chủ dữ liệu cục bộ, Siêu nhẹ và Không cồng kềnh). 

Để khắc phục hoàn toàn điểm yếu tiêu tốn tài nguyên khổng lồ của các phần mềm xây dựng trên nền tảng Chromium/Electron (thường ngốn 300MB - 800MB RAM chỉ để đọc file tĩnh), MDVault tận dụng sức mạnh biên dịch nhị phân trực tiếp của **Go 1.22+** kết hợp với WebView gốc của hệ điều hành (Microsoft Edge WebView2 trên Windows / WebKit trên macOS/Linux) thông qua framework **Wails v2**.

```mermaid
graph TD
    subgraph Traditional_Electron ["Kiến Trúc Electron Truyền Thống (Cồng Kềnh)"]
        NodeJS["Node.js Runtime (~50MB)"]
        Chromium["Embedded Chromium Browser (~250-400MB)"]
        AppCode1["JavaScript App Code"]
        TotalMem1["Tổng Tiêu Hao: ~350MB - 800MB RAM"]
    end

    subgraph MDVault_Wails ["Kiến Trúc MDVault (Tối Ưu & Siêu Nhẹ)"]
        GoCore["Go Native Compiled Binary (~10MB)"]
        NativeWebView["OS Native WebView2 (Shared DLL Memory ~25MB)"]
        Svelte5["Svelte 5 Runes (Zero-Runtime Virtual DOM)"]
        TotalMem2["Tổng Tiêu Hao: ~35MB - 50MB RAM"]
    end
```

### Các nguyên tắc kiến trúc chủ đạo:
1. **Dependency Inversion Principle (DIP):** Lớp bên ngoài phụ thuộc vào lớp bên trong. Logic nghiệp vụ (Domain & UseCase) hoàn toàn độc lập với chi tiết kỹ thuật (FileSystem, OS Dialog, Wails IPC).
2. **Anti-Corruption Layer (ACL):** Giao diện Frontend Svelte không phụ thuộc trực tiếp vào các hàm sinh tự động RPC của Wails. Lớp chuyển đổi `WailsAdapter` đứng ở giữa chịu trách nhiệm chuẩn hóa dữ liệu và xử lý bắt lỗi.
3. **Event-Driven Kernel Observation:** Giám sát tệp tin dựa trên sự kiện vi mô từ nhân hệ điều hành (OS Kernel Signals) qua `fsnotify`, loại bỏ hoàn toàn cơ chế Polling tốn CPU.
4. **Atomic Concurrency & Thread Safety:** Đảm bảo toàn vẹn dữ liệu khi đọc/ghi đồng thời bằng khóa độc quyền `sync.Mutex` và cơ chế ghi tệp tạm trước khi hoán đổi (Atomic Rename).

---

## 2. Bản Đồ Kiến Trúc Hệ Thống Tổng Thể (System Architecture Blueprint)

```mermaid
flowchart TB
    subgraph UI_Layer ["TẦNG GIAO DIỆN NGƯỜI DÙNG (Frontend - Svelte 5 + Vite)"]
        direction TB
        subgraph Components ["Feature-Driven Components"]
            Navbar["Navbar.svelte (Search / Actions)"]
            Sidebar["Sidebar.svelte (Folders & Tags)"]
            Editor["NoteEditor.svelte (Live Editor / Typography)"]
            Graph["GraphView.svelte (@xyflow/svelte)"]
        end
        subgraph StateMgmt ["Reactive State Management"]
            VaultStore["vaultStore (Notes, Folders, Active Note)"]
            ViewStore["viewStore (List Mode / Graph Mode)"]
            ThemeStore["themeStore (Dark / Light Tokens)"]
        end
        subgraph LibAdapters ["Lớp Cản Phá & Tiện Ích"]
            ACL["WailsAdapter (Anti-Corruption Layer)"]
            MDParser["markdown-it + typography plugin"]
            GraphEngine["graphBuilder.ts (Wiki-link Parser)"]
        end
    end

    subgraph IPC_Bridge ["CẦU NỐI TRUYỀN THÔNG (Wails v2 IPC Engine)"]
        IPC_Calls["Two-Way RPC Callbacks (Invoke / Promise)"]
        IPC_Events["Asynchronous Event Bus (runtime.EventsEmit)"]
    end

    subgraph Backend_Layer ["TẦNG LÕI HỆ THỐNG (Backend Go - Clean Architecture)"]
        direction TB
        subgraph Controller_Pkg ["Tầng Controller (Entrypoint & Transport)"]
            Ctrl["VaultController (Wails Binding Hub)"]
        end
        subgraph UseCase_Pkg ["Tầng UseCase (Application Orchestration)"]
            UC["VaultUseCase (Search, Cache, Traversal, Watcher Wire)"]
        end
        subgraph Domain_Pkg ["Tầng Domain (Core Models & Contracts)"]
            DomainModels["model.Note / model.Folder / model.WorkspaceData"]
            DomainRepo["repository.VaultRepository (Contract Interface)"]
        end
        subgraph Infra_Pkg ["Tầng Infrastructure (Physical Implementations)"]
            RepoImpl["VaultRepositoryImpl"]
            FSNotify["fsnotify.Watcher (OS Kernel Listener)"]
            GoYAML["goccy/go-yaml (Frontmatter Decoder)"]
        end
    end

    subgraph OS_Hardware ["TÀI NGUYÊN HỆ ĐIỀU HÀNH (Local Machine)"]
        LocalDisk[("Tệp Markdown Cục Bộ (.md)")]
        OS_Kernel["OS Kernel (ReadDirectoryChangesW / inotify)"]
    end

    %% Mối quan hệ liên tầng
    Components --> StateMgmt
    StateMgmt --> ACL
    Components --> MDParser
    Components --> GraphEngine
    ACL <--> IPC_Calls
    IPC_Events --> ACL

    IPC_Calls <--> Ctrl
    Ctrl --> IPC_Events
    Ctrl --> UC
    UC --> DomainRepo
    RepoImpl -.->|Implements| DomainRepo
    UC --> DomainModels
    RepoImpl --> DomainModels
    RepoImpl --> FSNotify
    RepoImpl --> GoYAML
    RepoImpl <--> LocalDisk
    OS_Kernel --> FSNotify
```

---

## 3. Kiến Trúc Phân Lớp Backend Go (Clean Architecture Decomposition)

Backend Go được tổ chức tuân thủ nghiêm ngặt nguyên lý 4 phân lớp của Clean Architecture:

```
d:\Tu-hoc\project\MDVault/
├── domain/                    # TẦNG 1: NHÂN DOANH NGHIỆP (INMOST LAYER)
│   ├── model/                 # Các thực thể cốt lõi (Note, Folder, WorkspaceData, Frontmatter)
│   └── repository/            # Giao diện hợp đồng (VaultRepository Contract Interface)
│
├── usecase/                   # TẦNG 2: NGHIỆP VỤ ỨNG DỤNG (APPLICATION ORCHESTRATION)
│   └── vault_usecase.go       # Triển khai thuật toán điều hướng, bộ nhớ đệm và tìm kiếm
│
├── infrastructure/            # TẦNG 3: TRIỂN KHAI NGOẠI VI (INFRASTRUCTURE ADAPTERS)
│   └── filesystem/            # Triển khai VaultRepository bằng `os`, `fsnotify`, `goccy/go-yaml`
│
├── controller/                # TẦNG 4: ĐIỂM TIẾP LỰC NGOÀI (TRANSPORT & WAILS CONTROLLER)
│   └── vault_controller.go    # Quản lý sự kiện Wails, chuyển hóa lỗi và gọi OS Dialog
│
└── main.go                    # ENTRY POINT & DEPENDENCY INJECTION ROOT
```

### 3.1. Tầng Domain (`domain/`)
- Chứa các kiểu dữ liệu thực thể thuần túy của hệ thống, không import bất kỳ package ngoại vi nào ngoài thư viện chuẩn.
- Định nghĩa interface `VaultRepository`:
  - `ScanWorkspace(dirPath string) (*model.WorkspaceData, error)`
  - `GetNote(filePath string) (*model.Note, error)`
  - `SaveNote(filePath string, content string, title string, tags []string) error`
  - `CreateFolder(parentPath string, name string) error`
  - `CreateNote(folderPath string, title string) (*model.Note, error)`
  - `StartWatching(dirPath string, onChange func(path string)) error`
  - `StopWatching() error`

### 3.2. Tầng UseCase (`usecase/`)
- Đảm nhận vai trò phối hợp logic nghiệp vụ:
  - Duy trì bộ nhớ đệm an toàn đồng thời cho workspace đang hoạt động (`currentWsPath`, `currentData`).
  - Cung cấp thuật toán tìm kiếm đa tiêu chí `SearchNotes(query string)` trên tập hợp ghi chú đã được nạp sẵn trong bộ nhớ.
  - Tiếp nhận callback từ Repository và phân phối tín hiệu tới tầng Controller.

### 3.3. Tầng Infrastructure (`infrastructure/filesystem/`)
- Triển khai interface `VaultRepository`:
  - Thuật toán duyệt đệ quy `filepath.WalkDir` tối ưu, chỉ thu thập file `.md`.
  - Bộ giải mã siêu dữ liệu `goccy/go-yaml` bóc tách Frontmatter trong thời gian micro-giây.
  - Module giám sát `fsnotify.Watcher` chạy trên Goroutine nền độc lập để bắt sự kiện từ OS Kernel.
  - Sử dụng cơ chế khóa độc quyền `sync.Mutex` bao bọc mọi thao tác I/O.

### 3.4. Tầng Controller (`controller/`)
- Đóng vai trò là Wails IPC Controller, tiếp nhận trực tiếp các yêu cầu từ giao diện frontend.
- Kích hoạt cửa sổ hộp thoại hệ thống `runtime.OpenDirectoryDialog`.
- Phát các tín hiệu thời gian thực `runtime.EventsEmit(ctx, "on-file-change", filePath)`.

---

## 4. Kiến Trúc Frontend Modular & Feature-Driven (Svelte 5)

Frontend áp dụng kiến trúc **Feature-Driven Design (FDD)**, cấu trúc theo từng nhóm tính năng độc lập nhằm tối đa hóa khả năng bảo trì và tái sử dụng:

```
frontend/src/
├── assets/                    # Styling Tokens & Css Architecture
│   ├── index.css              # Custom tokens
│   └── app.css                # @import "tailwindcss" (Tailwind v4) & Typography
├── lib/
│   ├── types/
│   │   └── domain.ts          # TypeScript Interfaces ánh xạ 1-1 với Go Structs
│   ├── stores/                # QUẢN TRỊ TRẠNG THÁI PHẢN ỨNG (REACTIVE STORES)
│   │   ├── vaultStore.ts      # Quản lý Cây Thư Mục, Danh Sách Notes, Active Note, Filter Tags
│   │   ├── viewStore.ts       # Chuyển đổi ngữ cảnh hiển thị (Chế độ Danh sách <=> Đồ thị Graph)
│   │   └── themeStore.ts      # Trình quản lý Theme Dark / Light theo Realtime Token
│   ├── wails/
│   │   └── wailsAdapter.ts    # LỚP CẢN PHÁ IPC (Anti-Corruption Layer)
│   └── utils/
│       ├── markdown.ts        # Pipeline render Markdown siêu tốc với markdown-it
│       └── graphBuilder.ts    # Bộ giải mã [[Wiki-link]] sinh cấu trúc node/edge cho @xyflow/svelte
├── components/                # FEATURE MODULES
│   ├── layout/Navbar.svelte   # Thanh điều khiển đa nhiệm, ô tìm kiếm và kích hoạt Workspace
│   ├── sidebar/               # Quản trị không gian lưu trữ
│   │   ├── Sidebar.svelte     # Bảng điều khiển phân nhóm & bộ lọc Tag pills
│   │   └── FolderTreeItem.svelte # Đệ quy hiển thị Thư mục phân cấp của Svelte 5
│   ├── note/                  # Không gian tương tác văn bản
│   │   ├── NotesList.svelte   # Khung cảm biến tóm tắt và danh sách ghi chú dạng thẻ/danh sách
│   │   └── NoteEditor.svelte  # Trình soạn thảo văn bản hỗ trợ Live Preview & Markdown Typography
│   └── graph/                 # Đồ thị Tri thức trực quan
│       └── GraphView.svelte   # Canvas đồ thị tương tác 2D (@xyflow/svelte)
├── App.svelte                 # Master Component, Khởi tạo Lifecycle & Lắng nghe sự kiện IPC
└── main.ts                    # Điểm nạp ứng dụng Svelte
```

---

## 5. Các Luồng Thực Thi Trọng Yếu (End-to-End Sequence Workflows)

### 5.1. Luồng Khởi Động & Nạp Thư Mục Workspace

```mermaid
sequenceDiagram
    autonumber
    actor User as Người Dùng
    participant UI as Svelte UI (Navbar)
    participant Adapter as WailsAdapter (ACL)
    participant IPC as Wails IPC Bridge
    participant Ctrl as VaultController (Go)
    participant UC as VaultUseCase (Go)
    participant Repo as VaultRepositoryImpl (Go)
    participant Disk as Local File System (.md)

    User->>UI: Nhấp vào nút "Chọn Workspace"
    UI->>Adapter: WailsAdapter.selectWorkspace()
    Adapter->>IPC: Invoke "SelectWorkspace"
    IPC->>Ctrl: (c *VaultController) SelectWorkspace()
    Ctrl->>User: runtime.OpenDirectoryDialog() (Cửa sổ OS)
    User-->>Ctrl: Chọn thư mục "D:\Notes"
    Ctrl->>UC: useCase.OpenWorkspace("D:\Notes")
    UC->>Repo: repo.ScanWorkspace("D:\Notes")
    Repo->>Disk: filepath.WalkDir đệ quy
    Disk-->>Repo: Đọc các file .md & parse YAML Frontmatter
    Repo->>Repo: fsnotify.Watcher.Add("D:\Notes") (Đăng ký Daemon)
    Repo-->>UC: *model.WorkspaceData
    UC-->>Ctrl: *model.WorkspaceData
    Ctrl-->>IPC: Trả về kết quả JSON
    IPC-->>Adapter: Giải mã đối tượng Promise<WorkspaceData>
    Adapter-->>UI: Cập nhật vaultStore
    UI-->>User: Render cây thư mục, danh sách Note và bộ lọc Tag
```

### 5.2. Luồng Phản Ứng Sự Kiện Tệp Ngoại Vi (Realtime Watcher Trigger)

```mermaid
sequenceDiagram
    autonumber
    actor ExternalUser as Trình Soạn Ngoài (VS Code / Vim)
    participant Disk as Local File System (.md)
    participant Kernel as OS Kernel Notification
    participant Watcher as fsnotify Daemon (Go Goroutine)
    participant Repo as VaultRepositoryImpl (Go)
    participant Ctrl as VaultController (Go)
    participant IPC as Wails Event Bus
    participant Adapter as WailsAdapter (Svelte)
    participant Store as vaultStore (Svelte Store)
    participant UI as Svelte Editor / NotesList

    ExternalUser->>Disk: Chỉnh sửa và lưu tệp "clean_arch.md"
    Disk->>Kernel: Sinh sự kiện Modify Event
    Kernel->>Watcher: Bắt tín hiệu ReadDirectoryChangesW
    Watcher->>Repo: Bắt sự kiện fsnotify.Event
    Repo->>Repo: Kiểm tra đuôi .md & debounce 50ms
    Repo->>Ctrl: Kích hoạt onChange callback(filePath)
    Ctrl->>IPC: runtime.EventsEmit("on-file-change", filePath)
    IPC-->>Adapter: EventsOn Callback được kích hoạt
    Adapter->>Store: Yêu cầu đồng bộ lại ghi chú
    Store->>Ctrl: WailsAdapter.getNoteContent(filePath)
    Ctrl-->>Store: Trả về dữ liệu mới nhất
    Store-->>UI: Reactive Update UI trong vòng dưới 15ms (Không giật màn hình)
```

### 5.3. Luồng Lưu Ghi Chú Nguyên Tử (Atomic Persistence Flow)

```mermaid
sequenceDiagram
    autonumber
    actor User as Người Dùng
    participant Editor as NoteEditor.svelte
    participant Adapter as WailsAdapter (ACL)
    participant Ctrl as VaultController (Go)
    participant UC as VaultUseCase (Go)
    participant Repo as VaultRepositoryImpl (Go)
    participant Disk as Local Disk (.md & .tmp)

    User->>Editor: Nhấp nút "Lưu" hoặc bấm phím Ctrl+S
    Editor->>Adapter: WailsAdapter.saveNoteContent(path, content, title, tags)
    Adapter->>Ctrl: Invoke "SaveNoteContent"
    Ctrl->>UC: useCase.SaveNote(...)
    UC->>Repo: repo.SaveNote(...)
    Repo->>Repo: Lock Mutex (sync.Mutex)
    Repo->>Disk: Tạo tệp tạm: path + ".tmp"
    Repo->>Disk: Ghi cấu trúc YAML Frontmatter + Nội dung mới
    Repo->>Disk: file.Sync() (Xả bộ nhớ đệm ổ đĩa)
    Repo->>Disk: os.Rename(path + ".tmp", path) (Thao tác nguyên tử)
    Repo->>Repo: Unlock Mutex
    Repo-->>UC: Trả về nil (Thành công)
    UC-->>Ctrl: Trả về nil
    Ctrl-->>Adapter: Promise resolved
    Adapter-->>Editor: Hiển thị Toast thông báo "Đã lưu thành công"
```

---

## 6. Xử Lý Các Mối Quan Tâm Xuyên Suốt (Cross-Cutting Concerns)

### 6.1. Kiểm Soát An Toàn Đồng Thời (Concurrency Guards)
- Mọi tài nguyên chia sẻ bao gồm: đường dẫn Workspace hiện tại (`currentWsPath`), bộ nhớ đệm tệp, và danh sách các thư mục đang được đăng ký với `fsnotify` đều được bảo vệ bằng con trỏ khóa `sync.Mutex`.
- Tránh deadlock bằng cách chỉ khóa Mutex trong phạm vi hẹp nhất của thao tác I/O hoặc truy xuất mảng, tuyệt đối không giữ khóa qua các tác vụ gọi network hoặc dialog hệ thống.

### 6.2. Kiểm Soát Lỗi Tập Trung (Centralized Error Strategy)
- Tầng Infrastructure bọc mọi lỗi I/O của hệ điều hành vào thông báo có ý nghĩa.
- Tầng Controller chuyển đổi `error` của Go thành chuỗi thông điệp rõ ràng trước khi trả về qua RPC Wails.
- Tầng Frontend `WailsAdapter` bắt lỗi trong khối `try-catch`, cung cấp cơ chế dự phòng an toàn (Fallback Gracefully) giúp ứng dụng không bao giờ bị đứng hình hay mất dữ liệu đang soạn thảo.

### 6.3. Tối Ưu Hóa Bộ Nhớ & Tốc Độ Render Đồ Thị
- Thuật toán trích xuất Wiki-link chỉ chạy trên tập ghi chú đã nạp vào bộ nhớ.
- Đồ thị `@xyflow/svelte` áp dụng kỹ thuật tính toán tọa độ Euclid phân bổ đồng tâm (Concentric Circle Distribution), đảm bảo khả năng hiển thị mượt mà 60 FPS ngay cả khi có hàng trăm nút liên kết.
