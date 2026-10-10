# Contributing to MDVault

Thank you for your interest in contributing to MDVault. We welcome contributions that improve performance, maintain architectural cleanliness, or enhance user experience.

---

## 1. Code of Conduct

Please maintain a respectful, constructive, and inclusive tone in all interactions within issues, pull requests, and discussions.

---

## 2. Architectural Guidelines

Before making changes, review the architecture specifications in the `markdown/` directory:
- [BRD.md](markdown/BRD.md) - Business context and constraints
- [PRB.md](markdown/PRB.md) - Product requirements and user journeys
- [SRS.md](markdown/SRS.md) - Detailed functional and non-functional requirements
- [HLD.md](markdown/HLD.md) - Clean Architecture and system design
- [LLD.md](markdown/LLD.md) - Low-level design, data models, and IPC API specs

Key architectural principles:
- **Clean Architecture:** Keep domain entities and business rules isolated from transport and infrastructure details.
- **Anti-Corruption Layer:** Frontend UI components must interact with backend RPC through `WailsAdapter`, not direct bindings.
- **Concurrency Safety:** Protect shared state in the Go backend with `sync.Mutex`.
- **Zero-Bloat:** Minimize bundle size and avoid unnecessary heavy third-party dependencies.

---

## 3. Commit Message Conventions

We enforce [Conventional Commits (v1.0.0)](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short description>

[optional body]

[optional footer(s)]
```

### Allowed Types:
- `feat`: A new user-facing feature or API method
- `fix`: A bug fix
- `docs`: Documentation updates or additions
- `style`: Formatting, whitespace, or UI theme adjustments without functional change
- `refactor`: Code refactoring without changing functionality or fixing bugs
- `perf`: Code changes that improve performance
- `test`: Adding or updating tests
- `build`: Changes affecting build configurations, Go modules, or Vite configs
- `chore`: Maintenance tasks or repository housekeeping

### Examples:
- `feat(backend): implement full-text search indexing`
- `fix(watcher): prevent duplicate reload events with 50ms debounce`
- `docs(readme): add bilingual documentation and getting started guide`

---

## 4. Development Workflow

1. **Fork and clone** the repository.
2. **Create a feature branch**:
   ```bash
   git checkout -b feat/your-feature-name
   ```
3. **Install dependencies**:
   ```bash
   cd frontend && npm install && cd ..
   ```
4. **Run development mode**:
   ```bash
   wails dev
   ```
5. **Verify changes**: Ensure Go compiles without warnings, frontend builds cleanly (`npm run build`), and no race conditions occur.
6. **Submit a Pull Request**: Provide a clear description referencing related issues.
