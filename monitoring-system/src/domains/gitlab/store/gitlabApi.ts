import { createApi } from "@reduxjs/toolkit/query/react";
import { axiosBaseQuery } from "@/shared/services/api/baseQuery";
import {
  Service,
  PipelineUnit,
  ExecutionHistory,
  AuthorizationRequest,
  PipelineRunStatus,
  PipelineUnitWithServices,
  PipelineData,
} from "../types/index";

interface CreateServiceRequest {
  gitlab_repo_id: string;
  name: string;
  url: string;
  type: "micro" | "macro";
}

interface CreatePipelineUnitRequest {
  macro_service_id: string;
  micro_service_ids: string[];
}

interface TriggerPipelineUnitRequest {
  requester_id: string;
  selected_micro_service_ids: string[];
}

interface ApprovePipelineRunRequest {
  approver_id: string;
  comment: string;
}

interface RejectPipelineRunRequest {
  comment: string;
}

interface SuperAdminCheckResponse {
  isSuperAdmin: boolean;
}

export const gitlabApi = createApi({
  reducerPath: "gitlabApi",
  baseQuery: axiosBaseQuery({ baseUrl: "/gitlab" }),
  tagTypes: [
    "Service",
    "PipelineUnit",
    "Pipeline",
    "ExecutionHistory",
    "AuthorizationRequest",
    "PipelineStatus",
  ],
  endpoints: (builder) => ({
    checkSuperAdmin: builder.query<SuperAdminCheckResponse, string>({
      query: (userId) => `/${userId}/is-super-admin`,
    }),
    ggetServices: builder.query<Service[], void>({
      query: () => "/services",
      transformResponse: (response: Service[]) => {
        console.log("Raw response:", response);
        return response; // must return it so RTKQ still works
      },
      providesTags: ["Service"],
    }),
    getServices: builder.query<any, void>({
      query: () => "/services",
      transformResponse: (response: any) => {
        return response.data; // must return it so RTKQ still works
      },
      providesTags: ["Service"],
    }),

    createService: builder.mutation<Service, CreateServiceRequest>({
      query: (data) => ({
        url: "/services",
        method: "post",
        data,
      }),
      invalidatesTags: ["Service"],
    }),
    getPipelineUnits: builder.query<PipelineUnitWithServices[], void>({
      query: () => "/pipeline-units",
      transformResponse: (response: any) => {
        return response.data; // must return it so RTKQ still works
      },
      providesTags: ["PipelineUnit"],
    }),
    getPipelineUnit: builder.query<PipelineUnitWithServices, string>({
      query: (id) => `/pipeline-units/${id}`,
      providesTags: (result, error, id) => [{ type: "PipelineUnit", id }],
    }),
    getPipelineServices: builder.query<PipelineUnitWithServices, string>({
      query: (id) => `/pipeline-unit/get/${id}`,
      providesTags: (result, error, id) => [{ type: "PipelineUnit", id }],
      transformResponse: (response: any) => {
        return response.data; // must return it so RTKQ still works
      },
    }),
    createPipelineUnit: builder.mutation<
      PipelineUnit,
      CreatePipelineUnitRequest
    >({
      query: (data) => ({
        url: "/pipeline-units",
        method: "post",
        data,
      }),
      invalidatesTags: ["PipelineUnit"],
    }),
    triggerAuthorizationRequest: builder.mutation<
      AuthorizationRequest,
      { id: string; data: TriggerPipelineUnitRequest }
    >({
      query: ({ id, data }) => ({
        url: `/pipeline-units/${id}/trigger`,
        method: "post",
        data,
      }),
      invalidatesTags: ["AuthorizationRequest"],
    }),
    getAuthorizationRequests: builder.query<AuthorizationRequest[], void>({
      query: () => "/authorization-requests",
      transformResponse: (response: any) => {
        return response.data; // must return it so RTKQ still works
      },
      providesTags: ["AuthorizationRequest"],
    }),
    approveAuthorizationRequest: builder.mutation<
      void,
      { id: string; data: ApprovePipelineRunRequest }
    >({
      query: ({ id, data }) => ({
        url: `/authorization-requests/${id}/approve`,
        method: "post",
        data,
      }),
      invalidatesTags: ["AuthorizationRequest", "PipelineStatus"],
    }),
    rejectAuthorizationRequest: builder.mutation<
      void,
      { id: string; data: RejectPipelineRunRequest }
    >({
      query: ({ id, data }) => ({
        url: `/authorization-requests/${id}/reject`,
        method: "post",
        data,
      }),
      invalidatesTags: ["AuthorizationRequest", "PipelineStatus"],
    }),
    getExecutionHistory: builder.query<ExecutionHistory[], void>({
      query: () => `/pipeline-runs/history`,
      providesTags: (result) =>
        result
          ? [
              ...result.map((h) => ({
                type: "ExecutionHistory" as const,
                id: h.id,
              })),
              { type: "ExecutionHistory", id: "LIST" },
            ]
          : [{ type: "ExecutionHistory", id: "LIST" }],
      transformResponse: (response: any) => response.data, // ensure correct unwrap
    }),
    getPipelineRunStatus: builder.query<PipelineData, string>({
      query: (id) => `/pipeline-status/${id}`,
      providesTags: (result, error, id) => [{ type: "PipelineStatus", id }],
      transformResponse: (response: any) => response.data, // ensure correct unwrap
    }),
  }),
});

export const {
  useCheckSuperAdminQuery,
  useGetServicesQuery,
  useCreateServiceMutation,
  useGetPipelineUnitsQuery,
  useGetPipelineUnitQuery,
  useGetPipelineServicesQuery,
  useCreatePipelineUnitMutation,
  useTriggerAuthorizationRequestMutation,
  useGetAuthorizationRequestsQuery,
  useApproveAuthorizationRequestMutation,
  useRejectAuthorizationRequestMutation,
  useGetExecutionHistoryQuery,
  useGetPipelineRunStatusQuery,
} = gitlabApi;
