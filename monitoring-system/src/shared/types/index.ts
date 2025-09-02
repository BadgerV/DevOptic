export interface BaseEntity {
  id: string;
  createdAt: string;
  updatedAt: string;
}

export interface ApiResponse<T> {
  data: T;
  message?: string;
  success: boolean;
}

export interface PaginatedResponse<T> extends ApiResponse<T[]> {
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export type LoadingState = 'idle' | 'loading' | 'succeeded' | 'failed';

export interface ApiError {
  message: string;
  code?: string;
  details?: any;
}
