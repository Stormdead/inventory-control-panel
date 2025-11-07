export interface User {
  id: number;
  username: string;
  email: string;
  role: 'admin' | 'employee';
  created_at?: string;
  updated_at?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  role?: 'admin' | 'employee';
}

export interface AuthResponse {
  token: string;
  user: User;
  message: string;
}