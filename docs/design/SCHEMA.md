# NanoCI - Database Schema

## 1. Entity Relationship Diagram (Textual)

```mermaid
erDiagram
    USERS ||--o{ PROJECTS : owns
    PROJECTS ||--o{ BUILDS : has
    PROJECTS ||--o{ SECRETS : contains
    BUILDS ||--o{ STEPS : contains

    USERS {
        uuid id PK
        string github_id UK
        string username
        string email
        string avatar_url
        timestamp created_at
        timestamp updated_at
    }

    PROJECTS {
        uuid id PK
        uuid user_id FK
        string name
        string repo_url
        string github_repo_id UK
        string default_branch
        string webhook_secret
        timestamp created_at
        timestamp updated_at
    }

    SECRETS {
        uuid id PK
        uuid project_id FK
        string key
        string encrypted_value
        timestamp created_at
    }

    BUILDS {
        uuid id PK
        uuid project_id FK
        string commit_hash
        string commit_message
        string branch
        string status "PENDING, RUNNING, SUCCESS, FAILED"
        timestamp started_at
        timestamp finished_at
        timestamp created_at
    }

    STEPS {
        uuid id PK
        uuid build_id FK
        string name
        string status
        int exit_code
        text logs
        timestamp started_at
        timestamp finished_at
    }
```

## 2. Table Definitions (PostgreSQL)

### 2.1. Users
Stores user identity authenticated via GitHub.
- `id`: UUID, Primary Key.
- `github_id`: String, Unique, ID from GitHub API.
- `username`: String.
- `email`: String.
- `avatar_url`: String.
- `created_at`: Timestamp.
- `updated_at`: Timestamp.

### 2.2. Projects
Represents a GitHub repository that NanoCI is watching.
- `id`: UUID, Primary Key.
- `user_id`: UUID, Foreign Key -> Users.id (Owner).
- `name`: String (e.g., "princetheprogrammer/nanoci").
- `repo_url`: String (HTTPS clone URL).
- `github_repo_id`: String, Unique (GitHub's internal ID).
- `default_branch`: String (e.g., "main").
- `webhook_secret`: String (Used to verify signatures).
- `created_at`: Timestamp.
- `updated_at`: Timestamp.

### 2.3. Secrets
Environment variables encrypted at rest.
- `id`: UUID, Primary Key.
- `project_id`: UUID, Foreign Key -> Projects.id.
- `key`: String (e.g., "AWS_ACCESS_KEY").
- `encrypted_value`: String (Base64 encoded ciphertext).
- `created_at`: Timestamp.

### 2.4. Builds
A single execution of a pipeline.
- `id`: UUID, Primary Key.
- `project_id`: UUID, Foreign Key -> Projects.id.
- `commit_hash`: String.
- `commit_message`: String.
- `branch`: String.
- `status`: Enum (PENDING, RUNNING, SUCCESS, FAILED, CANCELLED).
- `started_at`: Timestamp (Nullable).
- `finished_at`: Timestamp (Nullable).
- `created_at`: Timestamp.

### 2.5. Steps (Optional/Advanced)
Granular tracking of each step in the pipeline.
- `id`: UUID, Primary Key.
- `build_id`: UUID, Foreign Key -> Builds.id.
- `name`: String (Step name from .nanoci.yml).
- `status`: Enum.
- `exit_code`: Integer.
- `started_at`: Timestamp.
- `finished_at`: Timestamp.
