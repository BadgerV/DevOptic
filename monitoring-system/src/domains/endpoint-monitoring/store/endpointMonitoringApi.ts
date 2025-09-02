import { createApi } from "@reduxjs/toolkit/query/react";
import { axiosBaseQuery } from "@/shared/services/api/baseQuery";
import {
  Endpoint,
  EndpointStatus,
  EndpointMetrics,
  EndpointsEssentials,
} from "../types";
import { ApiResponse, PaginatedResponse } from "@shared/types";
import { IsEndpointCheckRunning } from "../types";

export const endpointMonitoringApi = createApi({
  reducerPath: "endpointMonitoringApi",
  baseQuery: axiosBaseQuery({
    baseUrl: "/monitor",
    // prepareHeaders: (headers) => {
    //   const token = localStorage.getItem('token');
    //   if (token) {
    //     headers.set('authorization', `Bearer ${token}`);
    //   }
    //   return headers;
    // },
  }),
  tagTypes: [
    "Endpoint",
    "EndpointStatus",
    "EndpointMetrics",
    "EndpointEssentails",
  ],
  endpoints: (builder) => ({
    getEndpointRunningStatus: builder.query<
      ApiResponse<IsEndpointCheckRunning>,
      void
    >({
      query: () => "/check-scheduler-status",
      providesTags: ["EndpointStatus"],
    }),
    getAllEndpointEssentials: builder.query<
      ApiResponse<EndpointsEssentials[]>,
      void
    >({
      query: () => "/get-endpoint-essentials",
      providesTags: ["EndpointEssentails"],
    }),
    getEndpointByID: builder.query<
      ApiResponse<EndpointsEssentials[]>,
      string
    >({
      query: (id) => `/get-endpoint-by-id/${id}`,
      providesTags: ["EndpointEssentails"],
    }),
    updateEndpoint: builder.mutation<
      ApiResponse<EndpointsEssentials>,
      { id: string; data: Partial<EndpointsEssentials> }
    >({
      query: ({ id, data }) => ({
        url: `/update-endpoint/${id}`,
        method: "put",
        body: data,
      }),
      invalidatesTags: ["EndpointEssentails"],
    }),
  }),
});

export const {
  useGetEndpointRunningStatusQuery,
  useGetAllEndpointEssentialsQuery,
  useGetEndpointByIDQuery,
  useUpdateEndpointMutation,
} = endpointMonitoringApi;

