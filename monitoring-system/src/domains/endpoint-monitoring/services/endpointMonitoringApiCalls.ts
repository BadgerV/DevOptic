// src/services/endpointMonitoringApiCalls.ts
import { apiClient } from "@/shared/services/api/client";

// Types for better type safety
interface SchedulerStatus {
  isRunning: boolean;
  lastCheck?: string;
  nextCheck?: string;
}

interface OverallStats {
  totalEndpoints: number;
  activeEndpoints: number;
  failedEndpoints: number;
  lastUpdateTime: string;
}
interface CreateEndpointPayload {
  service_name: string;
  url: string;
  server_name: string;
  api_method: string;
  expected_status_code: number;
  gitlab_url?: string | null;
  docker_container_name?: string | null;
  kubernetes_pod_name?: string | null;
  tags?: string[];
  description?: string | null;
  last_changed_by?: string | null;
}

interface CreateEndpointResponse {
  id: string;
  message: string;
  endpoint: CreateEndpointPayload;
}

// DRY helper for API calls with proper error handling and logging
const makeApiCall = async <T>(
  url: string,
  method: "get" | "post" | "put" | "patch" | "delete" = "get",
  data?: any,
  params?: Record<string, any>
): Promise<T> => {
  try {
    let response: T;

    switch (method) {
      case "get":
        response = await apiClient.get<T>(url, params);
        break;
      case "post":
        response = await apiClient.post<T>(url, data);
        break;
      case "put":
        response = await apiClient.put<T>(url, data);
        break;
      case "patch":
        response = await apiClient.patch<T>(url, data);
        break;
      case "delete":
        response = await apiClient.delete<T>(url);
        break;
      default:
        throw new Error(`Unsupported HTTP method: ${method}`);
    }

    console.log(`✅ API call successful: [${method.toUpperCase()}] ${url}`);
    return response;
  } catch (error) {
    console.error(
      `❌ API call failed: [${method.toUpperCase()}] ${url}`,
      error
    );
    throw error;
  }
};

// Helper for combined start/status operations
const performMonitorAction = async (
  actionUrl: string,
  actionName: string
): Promise<SchedulerStatus> => {
  try {
    await makeApiCall(actionUrl, "post");
    console.log(`✅ ${actionName} action completed successfully`);
  } catch (error: any) {
    console.error(`❌ ${actionName} action failed:`, error);
    throw error;
    // Continue to get status even if action fails
  }

  return await makeApiCall<SchedulerStatus>(
    "/monitor/check-scheduler-status",
    "get"
  );
};

// 🔹 Monitor API Calls with improved DRY implementation

export const startEndpointCheck = async (): Promise<SchedulerStatus> => {
  return performMonitorAction("/monitor/start-checks", "Start endpoint check");
};

export const stopEndpointCheck = async (): Promise<SchedulerStatus> => {
  return performMonitorAction("/monitor/stop-checks", "Stop endpoint check");
};

export const getSchedulerStatus = async (): Promise<SchedulerStatus> => {
  return makeApiCall<SchedulerStatus>("/monitor/check-scheduler-status", "get");
};

export const getOverallStats = async (): Promise<OverallStats> => {
  return makeApiCall<OverallStats>("/monitor/get-overall-stats", "get");
};

export const createEndpoint = async (endpointData: {
  service_name: string;
  url: string;
  server_name: string;
  api_method: string;
  expected_status_code: number;
  gitlab_url?: string | null;
  docker_container_name?: string | null;
  kubernetes_pod_name?: string | null;
  tags?: string[];
  description?: string | null;
  last_changed_by?: string | null;
}): Promise<CreateEndpointResponse> => {
  const payload: CreateEndpointPayload = {
    service_name: endpointData.service_name,
    url: endpointData.url,
    server_name: endpointData.server_name,
    api_method: endpointData.api_method,
    expected_status_code: endpointData.expected_status_code,
    gitlab_url: endpointData.gitlab_url || null,
    docker_container_name: endpointData.docker_container_name || null,
    kubernetes_pod_name: endpointData.kubernetes_pod_name || null,
    tags: endpointData.tags || [],
    description: endpointData.description || null,
    last_changed_by: endpointData.last_changed_by || null,
  };

  return makeApiCall<CreateEndpointResponse>(
    "/monitor/create-endpoint",
    "post",
    payload
  );
};

// Additional utility functions for common operations
export const refreshMonitoringData = async () => {
  const [schedulerStatus, overallStats] = await Promise.allSettled([
    getSchedulerStatus(),
    getOverallStats(),
  ]);

  return {
    schedulerStatus:
      schedulerStatus.status === "fulfilled" ? schedulerStatus.value : null,
    overallStats:
      overallStats.status === "fulfilled" ? overallStats.value : null,
    errors: [
      ...(schedulerStatus.status === "rejected"
        ? [schedulerStatus.reason]
        : []),
      ...(overallStats.status === "rejected" ? [overallStats.reason] : []),
    ],
  };
};

// Utility to validate endpoint creation parameters
export const validateEndpointParams = (
  service_name: string,
  url: string,
  server_name: string,
  api_method: string,
  expected_status_code: number
): { isValid: boolean; errors: string[] } => {
  const errors: string[] = [];

  if (!service_name?.trim()) errors.push("Service name is required");
  if (!url?.trim()) errors.push("URL is required");
  if (!server_name?.trim()) errors.push("Server name is required");
  if (!api_method?.trim()) errors.push("API method is required");
  if (
    !Number.isInteger(expected_status_code) ||
    expected_status_code < 100 ||
    expected_status_code > 599
  ) {
    errors.push(
      "Expected status code must be a valid HTTP status code (100-599)"
    );
  }

  try {
    new URL(url);
  } catch {
    errors.push("URL must be a valid URL format");
  }

  const validMethods = [
    "GET",
    "POST",
    "PUT",
    "PATCH",
    "DELETE",
    "HEAD",
    "OPTIONS",
  ];
  if (!validMethods.includes(api_method.toUpperCase())) {
    errors.push(`API method must be one of: ${validMethods.join(", ")}`);
  }

  return {
    isValid: errors.length === 0,
    errors,
  };
};

// src/services/endpointMonitoringApiCalls.ts

export interface LiveCheckResult {
  status: "success" | "failure";
  response_time?: number;
  status_code?: number;
  error?: string;
  timestamp: string;
}

/**
 * Checks if a specific endpoint is live by calling the backend.
 * Expects a backend endpoint like /monitor/check-endpoint/:id
 */
export const checkEndpointLive = async (
  endpointId: string | number
): Promise<LiveCheckResult> => {
  try {
    const response = await makeApiCall<{
      status: "up" | "unreachable";
      status_code?: number;
      latency_ms?: number;
      error?: string;
    }>(`/monitor/${endpointId}/check`, "get");

    return {
      status: response.status === "up" ? "success" : "failure",
      status_code: response.status_code,
      response_time: response.latency_ms,
      error: response.error,
      timestamp: new Date().toISOString(),
    };
  } catch (err: any) {
    return {
      status: "failure",
      error: err instanceof Error ? err.message : "Unknown error",
      timestamp: new Date().toISOString(),
    };
  }
};


// Safe create endpoint with validation
// export const createEndpointSafe = async (
//   service_name: string,
//   url: string,
//   server_name: string,
//   api_method: string,
//   expected_status_code: number
// ): Promise<CreateEndpointResponse> => {
//   const validation = validateEndpointParams(
//     service_name,
//     url,
//     server_name,
//     api_method,
//     expected_status_code
//   );

//   if (!validation.isValid) {
//     throw new Error(`Validation failed: ${validation.errors.join(", ")}`);
//   }

//   return createEndpoint(
//     service_name,
//     url,
//     server_name,
//     api_method,
//     expected_status_code
//   );
// };
