import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Movement, MovementRequest, MovementResponse } from '../../shared/models/movement.model';

@Injectable({
  providedIn: 'root'
})
export class MovementService {
  private readonly API_URL = `${environment.apiUrl}/movements`;

  constructor(private http: HttpClient) {}

  getAll(): Observable<MovementResponse> {
    return this.http.get<MovementResponse>(this.API_URL);
  }

  getById(id: number): Observable<{ movement: Movement }> {
    return this.http.get<{ movement: Movement }>(`${this.API_URL}/${id}`);
  }

  getByType(type: 'entrada' | 'salida'): Observable<MovementResponse> {
    return this.http.get<MovementResponse>(`${this.API_URL}/type/${type}`);
  }

  getByProduct(productId: number): Observable<MovementResponse> {
    return this.http.get<MovementResponse>(`${this.API_URL}/product/${productId}`);
  }

  create(data: MovementRequest): Observable<{ message: string; movement: Movement; nuevo_stock: number }> {
    return this.http.post<{ message: string; movement: Movement; nuevo_stock: number }>(this.API_URL, data);
  }

  delete(id: number): Observable<{ message: string }> {
    return this.http.delete<{ message: string }>(`${this.API_URL}/${id}`);
  }
}