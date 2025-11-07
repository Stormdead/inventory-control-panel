export interface Category {
  id: number;
  name: string;
  description: string;
  created_at?: string;
}

export interface CategoryResponse {
  categories: Category[];
}

export interface CategoryRequest {
  name: string;
  description: string;
}