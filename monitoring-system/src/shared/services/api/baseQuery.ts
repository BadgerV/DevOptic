// src/services/axiosBaseQuery.ts
import { BaseQueryFn } from "@reduxjs/toolkit/query";
import { AxiosError } from "axios";
import { apiClient } from "./client"; // your ApiClient class

interface AxiosBaseQueryArgs {
  url: string;
  method?: "get" | "post" | "put" | "patch" | "delete";
  data?: any;
  params?: any;
}

export const axiosBaseQuery =
  (
    { baseUrl }: { baseUrl: string } = { baseUrl: "" }
  ): BaseQueryFn<AxiosBaseQueryArgs | string, unknown, unknown> =>
  async (args) => {
    try {
      let requestConfig: AxiosBaseQueryArgs;

      if (typeof args === "string") {
        // query: '/path'
        requestConfig = { url: baseUrl + args, method: "get" };
      } else {
        // query: { url, method, data, params }
        requestConfig = {
          ...args,
          url: baseUrl + args.url,
          method: args.method ?? "get",
        };
      }

      const method = requestConfig.method ?? "get";
      const result = await (apiClient as any)[method](
        requestConfig.url,
        requestConfig.data ?? requestConfig.params
      );

      //   console.log(result);

      return { data: result };
    } catch (axiosError) {
      const err = axiosError as AxiosError;
      return {
        error: {
          status: err.response?.status,
          data: err.response?.data || err.message,
        },
      };
    }
  };
