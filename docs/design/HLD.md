# NanoCI - High-Level Design (HLD)

## 1. System Overview
NanoCI is a lightweight, self-hosted, container-native CI/CD platform. It provides a single-binary (or microservices capable) solution for executing build pipelines defined in `.nanoci.yml` files. It leverages Docker for build isolation, Redis for job queuing, and PostgreSQL for persistent storage.

## 2. Architecture Diagram

```mermaid
graph TD
    User[User / Developer] -->|UI / CLI| LoadBalancer
    GitHub[GitHub Webhook] -->|HTTP POST| LoadBalancer

    subgraph "NanoCI Cluster"
        LoadBalancer --> APIServer[API Server (Go)]
        
        APIServer -->|Read/Write| DB[(PostgreSQL)]
        APIServer -->|Enqueue Job| Redis[(Redis Queue)]
        
        subgraph "Worker Nodes"
            Worker[NanoCI Worker] -->|Poll| Redis
            Worker -->|Update Status / Logs| APIServer
            Worker -->|Execute| DockerEngine[Docker Engine]
        end
    end

    DockerEngine -->|Run| Container[Build Container]
```

## 3. Core Components

### 3.1. API Server (The Brain)
- **Responsibility**: 
  - Handles incoming HTTP requests (REST API).
  - Receives and verifies GitHub Webhooks.
  - Manages User Authentication (GitHub OAuth2).
  - Serves the Frontend assets (or proxies to Next.js).
  - Exposes WebSocket endpoints for real-time log streaming.
- **Tech**: Go (Chi/Echo), PostgreSQL, Redis Client.

### 3.2. Job Queue (The Nervous System)
- **Responsibility**: Decouples webhook reception from build execution.
- **Tech**: Redis (using lists or streams).

### 3.3. Worker / Runner (The Muscle)
- **Responsibility**:
  - Polls Redis for pending build jobs.
  - Clones the repository.
  - Parses `.nanoci.yml`.
  - Provision ephemeral Docker containers for each step.
  - Streams logs back to the API Server (or directly to storage) in real-time.
  - Updates build status (Pending -> Running -> Success/Failed).
- **Tech**: Go, Docker SDK.

### 3.4. Database (The Memory)
- **Responsibility**: Stores persistent data.
  - Users & Permissions.
  - Projects / Repositories.
  - Build History & Status.
  - Secrets (Encrypted).
- **Tech**: PostgreSQL.

### 3.5. Frontend (The Face)
- **Responsibility**: 
  - Dashboard for viewing builds.
  - Real-time log viewer.
  - Project configuration.
- **Tech**: React + TailwindCSS.

## 4. Key Workflows

### 4.1. Webhook to Build Trigger
1. GitHub sends a `push` event to `/api/v1/webhooks/github`.
2. API Server verifies the signature.
3. API Server looks up the repository in DB.
4. API Server creates a `Build` record in DB with status `PENDING`.
5. API Server pushes a job payload to Redis `nanoci:jobs` queue.

### 4.2. Build Execution
1. Worker pops job from Redis.
2. Worker updates Build status to `RUNNING` via API (or direct DB access if co-located).
3. Worker pulls code using git.
4. Worker reads `.nanoci.yml`.
5. For each step:
   - Create Docker container.
   - Execute command.
   - Stream stdout/stderr to Log Handler.
6. If all steps pass, update status to `SUCCESS`. Else `FAILED`.
7. Worker cleans up containers.

## 5. Security Considerations
- **Secrets**: Stored in DB encrypted with AES-GCM. Decrypted only by the worker at runtime and injected as env vars.
- **Isolation**: Every build runs in a fresh Docker container.
- **Authentication**: No local passwords. GitHub OAuth2 only for strict access control.

## 6. Scalability
- **Horizontal Scaling**: Multiple API Servers can sit behind a load balancer. Multiple Workers can run on different machines/nodes pointing to the same Redis/DB.
