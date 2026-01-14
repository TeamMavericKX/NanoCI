# NanoCI 🥘

NanoCI is a single-binary CI/CD server written in Go. It's lightweight, container-native, and self-hosted.

## Features
- **Docker Isolation**: Each build step runs in its own container.
- **Real-time Logs**: Watch your builds flow via WebSockets.
- **Secrets Management**: AES-GCM encrypted secrets at rest.
- **GitHub Integration**: Trigger builds on push webhooks.
- **Redis Queue**: Scalable job distribution.

## Tech Stack
- **Backend**: Go 1.22 (Chi, Pgx, Zap)
- **Database**: PostgreSQL
- **Queue**: Redis
- **Isolation**: Docker SDK

## Getting Started

### 1. Prerequisites
- Docker & Docker Compose
- GitHub OAuth App (for authentication)

### 2. Configuration
Create a `.env` file:
```env
GITHUB_CLIENT_ID=your_id
GITHUB_CLIENT_SECRET=your_secret
ENCRYPTION_KEY=32_byte_key_here_...
```

### 3. Run
```bash
docker-compose up --build
```

### 4. Usage
Add a `.nanoci.yml` to your repository:
```yaml
image: alpine:latest
steps:
  - name: hello
    commands:
      - echo "Hello from NanoCI!"
```

## Architecture
See `docs/design/HLD.md` for details.
