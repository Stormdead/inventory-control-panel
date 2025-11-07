import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { 
  DashboardResponse, 
  MovementSummary, 
  TopProduct 
} from '../../shared/models/dashboard.model';
import { MovementResponse } from '../../shared/models/movement.model';
import { ProductResponse } from '../../shared/models/product.model';

@Injectable({
  providedIn: 'root'
})
export class DashboardService {
  private readonly API_URL = `${environment.apiUrl}/dashboard`;

  constructor(private http: HttpClient) {}

  getStats(): Observable<DashboardResponse> {
    return this.http.get<DashboardResponse>(`${this.API_URL}/stats`);
  }

  getRecentMovements(): Observable<MovementResponse> {
    return this.http.get<MovementResponse>(`${this.API_URL}/recent-movements`);
  }

  getLowStockAlerts(): Observable<ProductResponse> {
    return this.http.get<ProductResponse>(`${this.API_URL}/low-stock-alerts`);
  }

  getMovementSummary(): Observable<{ summary: MovementSummary; period: string }> {
    return this.http.get<{ summary: MovementSummary; period: string }>(`${this.API_URL}/movement-summary`);
  }

  getTopProducts(): Observable<{ products: TopProduct[]; total: number }> {
    return this.http.get<{ products: TopProduct[]; total: number }>(`${this.API_URL}/top-products`);
  }
}