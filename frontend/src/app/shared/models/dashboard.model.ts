export interface DashboardStats {
  total_products: number;
  total_categories: number;
  total_stock: number;
  low_stock_products: number;
  total_users: number;
  total_inventory_value: number;
}

export interface MovementSummary {
  total_entradas: number;
  total_salidas: number;
  cantidad_entradas: number;
  cantidad_salidas: number;
}

export interface TopProduct {
  product_id: number;
  product_name: string;
  total_movements: number;
  current_stock: number;
  category_name: string;
}

export interface DashboardResponse {
  stats: DashboardStats;
}