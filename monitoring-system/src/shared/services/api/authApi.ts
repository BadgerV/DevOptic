// src/features/auth/services/authApi.ts
import { createApi } from "@reduxjs/toolkit/query/react";
import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { axiosBaseQuery } from "@/shared/services/api/baseQuery";
import { ApiResponse } from "@shared/types";

/* -------------------- Request/Response Types -------------------- */
interface LoginRequest {
  username: string;
  password: string;
}

interface RegisterRequest {
  username: string;
  password: string;
  email: string;
}

interface AuthResponse {
  token: string; // adjust if your backend returns something different
  user?: {
    id: string;
    username: string;
    email: string;
  };
}

interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

interface DeliveryEmailRequest {
  delivery_email: string;
}

interface DeliveryEmailResponse {
  delivery_email: string;
  message: string;
}

/* -------------------- RTK Query API -------------------- */
export const authApi = createApi({
  reducerPath: "authApi",
  baseQuery: axiosBaseQuery({
    baseUrl: "/auth",
  }),
  tagTypes: ["Auth"],
  endpoints: (builder) => ({
    login: builder.mutation<ApiResponse<AuthResponse>, LoginRequest>({
      query: (data) => ({
        url: "/login",
        method: "post",
        data,
      }),
      invalidatesTags: ["Auth"],
    }),
    validateToken: builder.query<any, void>({
      query: () => ({
        url: `/validate`,
        method: "get",
      }),
      providesTags: ["Auth"],
    }),
    register: builder.mutation<ApiResponse<AuthResponse>, RegisterRequest>({
      query: (data) => ({
        url: "/register",
        method: "post",
        data,
      }),
      invalidatesTags: ["Auth"],
    }),
    changePassword: builder.mutation<ApiResponse<{ message: string }>, ChangePasswordRequest>({
      query: (data) => ({
        url: "/change-password",
        method: "post",
        data,
      }),
    }),
    setDeliveryEmail: builder.mutation<ApiResponse<{ message: string }>, DeliveryEmailRequest>({
      query: (data) => ({
        url: "/set-delivery-email",
        method: "post",
        data,
      }),
      invalidatesTags: ["Auth"],
    }),
    getDeliveryEmail: builder.query<ApiResponse<DeliveryEmailResponse>, void>({
      query: () => ({
        url: "/get-delivery-email",
        method: "get",
      }),
      providesTags: ["Auth"],
    }),
  }),
});

/* Hooks exported for components */
export const {
  useLoginMutation,
  useRegisterMutation,
  useValidateTokenQuery,
  useChangePasswordMutation,
  useSetDeliveryEmailMutation,
  useGetDeliveryEmailQuery,
} = authApi;

/* -------------------- Regular Redux Slice -------------------- */
interface UserState {
  email: string;
  username: string;
  is_active: boolean;
  userId: string;
  isAdmin: boolean;
}

interface AuthState {
  user: UserState | null;
}

const initialState: AuthState = {
  user: null,
};

const authSlice = createSlice({
  name: "authSlice",
  initialState,
  reducers: {
    setUserState: (state, action: PayloadAction<UserState>) => {
      state.user = action.payload;
    },
    clearUserState: (state) => {
      state.user = null;
    },
  },
});

export const { setUserState, clearUserState } = authSlice.actions;
export const authReducer = authSlice.reducer;
