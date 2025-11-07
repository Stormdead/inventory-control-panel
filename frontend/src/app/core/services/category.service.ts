import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Category, CategoryResponse, CategoryRequest } from '../../shared/models/category.model';

@Injectable({
  providedIn: 'root'
})
export class CategoryService {
  private readonly API_URL = `${environment.apiUrl}/categories`;

  constructor(private http: HttpClient) {}

  getAll(): Observable<CategoryResponse> {
    return this.http.get<CategoryResponse>(this.API_URL);
  }

  getById(id: number): Observable<{ category: Category }> {
    return this.http.get<{ category: Category }>(`${this.API_URL}/${id}`);
  }

  create(data: CategoryRequest): Observable<{ message: string; category: Category }> {
    return this.http.post<{ message: string; category: Category }>(this.API_URL, data);
  }

  update(id: number, data: CategoryRequest): Observable<{ message: string; category: Category }> {
    return this.http.put<{ message: string; category: Category }>(`${this.API_URL}/${id}`, data);
  }

  delete(id: number): Observable<{ message: string }> {
    return this.http.delete<{ message: string }>(`${this.API_URL}/${id}`);
  }
}