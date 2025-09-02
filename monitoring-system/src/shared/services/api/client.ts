import axios, { AxiosInstance, AxiosResponse, AxiosError } from 'axios';
import { ApiError } from '@shared/types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors() {
    // Request interceptor
    this.client.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('token');
        console.log("the token is coming from client.ts", token)
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
          // config.headers.Authorization = `${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor
    // this.client.interceptors.response.use(
    //   (response: AxiosResponse) => response,
    //   (error: AxiosError) => {
    //     const apiError: ApiError = {
    //       message: error.message,
    //       code: error.code,
    //       details: error.response?.data,
    //     };

    //     if (error.response?.status === 401) {
    //       localStorage.removeItem('token');
    //       // window.location.href = '/login';

    //     }

    //     return Promise.reject(apiError);
    //   }
    // );
  }

  async get<T>(url: string, params?: Record<string, any>): Promise<T> {
    const response = await this.client.get<T>(url, { params });
    return response.data;
  }

  async post<T>(url: string, data?: any): Promise<T> {
    const response = await this.client.post<T>(url, data);
    return response.data;
  }

  async put<T>(url: string, data?: any): Promise<T> {
    const response = await this.client.put<T>(url, data);
    return response.data;
  }

  async patch<T>(url: string, data?: any): Promise<T> {
    const response = await this.client.patch<T>(url, data);
    return response.data;
  }

  async delete<T>(url: string): Promise<T> {
    const response = await this.client.delete<T>(url);
    return response.data;
  }
}

export const apiClient = new ApiClient();
