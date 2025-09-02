import { BaseEntity } from "@shared/types";

export interface Endpoint extends BaseEntity {
  name: string;
  url: string;
  method: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
  expectedStatus: number;
  timeout: number;
  interval: number; // in seconds
  isActive: boolean;
  tags: string[];
}

export interface EndpointStatus {
  endpointId: string;
  status: "up" | "down" | "degraded";
  responseTime: number;
  statusCode: number;
  error?: string;
  timestamp: string;
  uptime: number; // percentage
}

export interface IsEndpointCheckRunning {
  scheduler_running: boolean;
  timestamp: string;
}

export interface EndpointMetrics {
  endpointId: string;
  averageResponseTime: number;
  uptime: number;
  totalChecks: number;
  failedChecks: number;
  lastCheck: string;
  trend: "improving" | "stable" | "degrading";
}

export interface EndpointsEssentials {
  endpointID: string;
  service_name: string;
  server_name: string;
  url : string;
  total_checks : string;
  uptime_percentage: number;
  downtime_count: string;
  successful_checks: string;
  avg_latency: string
  last_run: boolean
  failure_count: number
  id: string
  expected_code: any
  uptime: any
  api_method: any
  created_by: any
  created_at: any
  updated_at: any
  tags: any
  severity: any
}

export interface OverallStats {
  total_endpoints: number;
  total_checks: number;
  successful_checks: number;
  down_time_count: number;
  overall_uptime: number;
  average_latency: number;
}