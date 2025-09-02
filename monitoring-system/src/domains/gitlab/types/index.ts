export interface Service {
  id: string;
  gitlab_repo_id: string;
  name: string;
  url: string;
  type: "micro" | "macro";
  created_at: string;
  updated_at: string;
}

export interface PipelineUnit {
  id: string;
  macro_service_id: string;
  micro_service_ids: string[];
  created_at: string;
  updated_at: string;
}

export interface PipelineUnitWithServices {
  pipeline_unit: PipelineUnit;
  MacroService: Service;
  MicroServices: Service[];
}

export interface ExecutionHistory {
  id: string;
  pipeline_run_id: string;
  requester_id: string;
  requester_name: string;
  approver_id: string;
  approver_name: string;
  status: "pending" | "running" | "success" | "failed";
  started_at: string;
  completed_at: string;
  execution_time: string;
  error_message?: string;
  macro_service_name: string;
}

export interface AuthorizationRequest {
  id: string;
  pipeline_run_id: string;
  requester_id: string;
  requester_name: string;
  approver_id: string;
  approver_name: string;
  status: "pending" | "accepted" | "rejected";
  created_at: string;
  updated_at: string;
  comment?: string;
}

export interface PipelineRunStatus {
  pipeline_run_id: string;
  current_service_id: string;
  status: "pending" | "running" | "success" | "failed";
  timestamp: string;
  message?: string;
  approver_id: any;
  comment: any;
}

export interface PipelineData {
  id: string;
  pipeline_unit_id: string;
  status: "completed" | "pending" | "failed" | string; // enum-like if you know all possible statuses
  macro_service_name: string;
  micro_service_names: string[];
  requester_name: string;
  approver_name: string;
  created_at: string; // ISO datetime string
  updated_at: string; // ISO datetime string
  gitlab_pipeline_id: number;
};
