import { createApi } from "@reduxjs/toolkit/query/react";
import { axiosBaseQuery } from "../../../shared/services/api/baseQuery";

// Types for the RBAC API
export interface Role {
  id: string;
  name: string;
  description: string;
  created_at?: string;
  updated_at?: string;
}

export interface Permission {
  id: string;
  resource: string;
  action: string;
  description: string;
  created_at?: string;
  updated_at?: string;
}

export interface CreateRoleRequest {
  name: string;
  description: string;
}

export interface CreatePermissionRequest {
  resource: string;
  action: string;
  description: string;
}

export interface AssignRoleRequest {
  user_id: string;
  role_id: string;
}

export interface RemoveRoleRequest {
  user_id: string;
  role_id: string;
}

export interface AssignPermissionRequest {
  role_id: string;
  permission_id: string;
}

export interface CheckPermissionRequest {
  user_id: string;
  resource: string;
  action: string;
}

export interface CheckPermissionResponse {
  allowed: boolean;
  message?: string;
}

export interface UserPermissionsResponse {
  user_id: string;
  permissions: Permission[];
}

export interface HealthCheckResponse {
  status: string;
  message: string;
}

export interface SuperAdminCheckResponse {
  is_super_admin: boolean;
  user_id: string;
}

// Using the existing axiosBaseQuery from services

// Create the RBAC API slice
export const rbacApi = createApi({
  reducerPath: "rbacApi",
  baseQuery: axiosBaseQuery({
    baseUrl: "/rbac",
  }),
  tagTypes: ["Role", "Permission", "UserPermissions"],
  endpoints: (builder) => ({
    // Health check endpoint
    healthCheck: builder.query<HealthCheckResponse, void>({
      query: () => "/health",
    }),

    // Role management endpoints
    createRole: builder.mutation<Role, CreateRoleRequest>({
      query: (body) => ({
        url: "/roles",
        method: "post",
        data: body,
      }),
      invalidatesTags: ["Role"],
    }),

    getAllRoles: builder.query<Role[], void>({
      query: () => "/roles",
      transformResponse: (response: { message: string; data: Role[] }) =>
        response.data,
      providesTags: ["Role"],
    }),

    // Permission management endpoints
    createPermission: builder.mutation<Permission, CreatePermissionRequest>({
      query: (body) => ({
        url: "/permissions",
        method: "post",
        data: body,
      }),
      invalidatesTags: ["Permission"],
    }),

    // Role assignment endpoints
    assignRole: builder.mutation<{ message: string }, AssignRoleRequest>({
      query: (body) => ({
        url: "/assign-role",
        method: "post",
        data: body,
      }),
      invalidatesTags: ["UserPermissions"],
    }),

    removeRole: builder.mutation<{ message: string }, RemoveRoleRequest>({
      query: (body) => ({
        url: "/remove-role",
        method: "post",
        data: body,
      }),
      invalidatesTags: ["UserPermissions"],
    }),

    // Permission assignment endpoint
    assignPermission: builder.mutation<
      { message: string },
      AssignPermissionRequest
    >({
      query: (body) => ({
        url: "/assign-permission",
        method: "post",
        data: body,
      }),
      invalidatesTags: ["UserPermissions"],
    }),

    // User permissions endpoints
    getUserPermissions: builder.query<UserPermissionsResponse, string>({
      query: (userId) => `/user/${userId}/permissions`,
      providesTags: ["UserPermissions"],
    }),

    // Permission checking endpoints
    checkPermission: builder.mutation<
      CheckPermissionResponse,
      CheckPermissionRequest
    >({
      query: (body) => ({
        url: "/check-permission",
        method: "post",
        data: body,
      }),
    }),

    checkUserPermission: builder.query<
      CheckPermissionResponse,
      {
        user_id: string;
        resource: string;
        action: string;
      }
    >({
      query: ({ user_id, resource, action }) => ({
        url: `/check-permission/${user_id}`,
        method: "get",
        params: { resource, action },
      }),
    }),

    // Super admin check endpoint (new endpoint to be added to backend)
    checkSuperAdmin: builder.query<SuperAdminCheckResponse, string>({
      query: (userId) => `/${userId}/is-super-admin`,
    }),
    // Super admin check endpoint (new endpoint to be added to backend)
    getUserUsernameAndId: builder.query<any, void>({
      query: () => `/get-all-users-usernames-id`,
    }),
  }),
});

// Export hooks for usage in components
export const {
  useHealthCheckQuery,
  useCreateRoleMutation,
  useGetAllRolesQuery,
  useCreatePermissionMutation,
  useAssignRoleMutation,
  useRemoveRoleMutation,
  useAssignPermissionMutation,
  useGetUserPermissionsQuery,
  useCheckPermissionMutation,
  useCheckUserPermissionQuery,
  useCheckSuperAdminQuery,
  useGetUserUsernameAndIdQuery
} = rbacApi;
