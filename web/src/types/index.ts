export interface User {
  id: string;
  github_id: string;
  username: string;
  email: string;
  avatar_url: string;
  created_at: string;
  updated_at: string;
}

export interface Project {
  id: string;
  user_id: string;
  name: string;
  repo_url: string;
  default_branch: string;
  created_at: string;
  updated_at: string;
}

export type BuildStatus = "PENDING" | "RUNNING" | "SUCCESS" | "FAILED" | "CANCELLED";

export interface Build {
  id: string;
  project_id: string;
  commit_hash: string;
  commit_message: string;
  branch: string;
  status: BuildStatus;
  started_at?: string;
  finished_at?: string;
  created_at: string;
}

export interface Secret {
  id: string;
  project_id: string;
  key: string;
  created_at: string;
}
