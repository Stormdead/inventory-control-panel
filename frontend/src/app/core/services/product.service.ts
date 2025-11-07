import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Product, ProductResponse, ProductRequest } from '../../shared/models/product.model';

@Injectable({
  providedIn: 'root'
})
export class ProductService {
  private readonly API_URL = `${environment.apiUrl}/products`;

  constructor(private http: HttpClient) {}

  getAll(): Observable<ProductResponse> {
    return this.http.get<ProductResponse>(this.API_URL);
  }

  getById(id: number): Observable<{ product: Product }> {
    return this.http.get<{ product: Product }>(`${this.API_URL}/${id}`);
  }

  getLowStock(): Observable<ProductResponse> {
    return this.http.get<ProductResponse>(`${this.API_URL}/low-stock`);
  }

  getByCategory(categoryId: number): Observable<ProductResponse> {
    return this.http.get<ProductResponse>(`${this.API_URL}/category/${categoryId}`);
  }

  create(data: ProductRequest): Observable<{ message: string; product: Product }> {
    return this.http.post<{ message: string; product: Product }>(this.API_URL, data);
  }

  update(id: number, data: ProductRequest): Observable<{ message: string; product: Product }> {
    return this.http.put<{ message: string; product: Product }>(`${this.API_URL}/${id}`, data);
  }

  delete(id: number): Observable<{ message: string }> {
    return this.http.delete<{ message: string }>(`${this.API_URL}/${id}`);
  }
}