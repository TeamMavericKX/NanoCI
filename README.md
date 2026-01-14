# 🥘 NanoCI: Your Own Private GitHub Actions

**NanoCI** is a high-performance, single-binary (or containerized) CI/CD server designed for developers who want the power of GitHub Actions without the cost or the bloat of Jenkins.

Built from the ground up with **Go**, **Docker**, and **WebSockets**, NanoCI provides a modern, container-native experience for automating your builds on a $5 VPS or a Raspberry Pi.

### 🧠 Why NanoCI?
*   **Zero Bloat:** Written in Go for raw speed and minimal memory footprint.
*   **Docker-Native:** Every build step runs in a fresh, isolated container.
*   **Real-Time Everything:** Watch your logs stream in real-time via WebSockets and Redis Pub/Sub.
*   **Security First:** AES-GCM 256-bit encryption for project secrets and GitHub OAuth2 authentication.
*   **Simple YAML:** Define pipelines in `.nanoci.yml` just like you’re used to.

### 🛠️ The "Modern" Stack
-   **Backend:** Go 1.22+ (Chi, pgx, Zap)
-   **Queue:** Redis (High-speed job distribution)
-   **Muscle:** Docker SDK (Ephemeral container orchestration)
-   **Frontend:** React + TypeScript + TailwindCSS (Vercel-style UI)
-   **Database:** PostgreSQL (Robust persistence)

### 🚀 Quick Start
```bash
# 1. Clone the power
git clone https://github.com/10xdev4u-alt/NanoCI.git && cd NanoCI

# 2. Set your secrets in .env
# GITHUB_CLIENT_ID=...
# GITHUB_CLIENT_SECRET=...
# ENCRYPTION_KEY=...

# 3. Launch the beast
docker-compose up --build
```

### 📋 Usage
Add a `.nanoci.yml` to your repository:
```yaml
image: alpine:latest
steps:
  - name: hello
    commands:
      - echo "Hello from NanoCI!"
```

## 🏗️ Architecture
See `docs/design/HLD.md` for details.

---
**Let's Cook.** 🥘 Built with ❤️ by PrinceTheProgrammer.