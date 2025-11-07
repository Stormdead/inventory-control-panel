import { Category } from './category.model';

export interface Product {
  id: number;
  name: string;
  description: string;
  category_id?: number;
  category?: Category;
  price: number;
  stock: number;
  image_url?: string;
  created_at?: string;
  updated_at?: string;
}

export interface ProductResponse {
  products: Product[];
  total: number;
}

export interface ProductRequest {
  name: string;
  description: string;
  category_id?: number;
  price: number;
  stock: number;
  image_url?: string;
}