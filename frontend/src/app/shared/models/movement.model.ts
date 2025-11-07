import { Product } from './product.model';
import { User } from './user.model';

export interface Movement {
  id: number;
  product_id: number;
  product?: Product;
  user_id: number;
  user?: User;
  type: 'entrada' | 'salida';
  quantity: number;
  description?: string;
  movement_date: string;
}

export interface MovementRequest {
  product_id: number;
  type: 'entrada' | 'salida';
  quantity: number;
  description?: string;
}

export interface MovementResponse {
  movements: Movement[];
  total: number;
}