import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';
import { adminGuard } from './core/guards/admin.guard';

export const routes: Routes = [
  {
    path: '',
    redirectTo: '/dashboard',
    pathMatch: 'full'
  },
  {
    path: 'login',
    loadComponent: () => import('./pages/auth/login/login.component')
      .then(m => m.LoginComponent)
  },
  {
    path: 'register',
    loadComponent: () => import('./pages/auth/register/register.component')
      .then(m => m.RegisterComponent)
  },
  {
    path: 'dashboard',
    loadComponent: () => import('./pages/dashboard/dashboard.component')
      .then(m => m.DashboardComponent),
    canActivate: [authGuard]
  },
  {
    path: 'products',
    loadComponent: () => import('./pages/products/product-list/product-list.component')
      .then(m => m.ProductListComponent),
    canActivate: [authGuard]
  },
  {
    path: 'products/new',
    loadComponent: () => import('./pages/products/product-form/product-form.component')
      .then(m => m.ProductFormComponent),
    canActivate: [authGuard, adminGuard]
  },
  {
    path: 'products/edit/:id',
    loadComponent: () => import('./pages/products/product-form/product-form.component')
      .then(m => m.ProductFormComponent),
    canActivate: [authGuard]
  },
  {
    path: 'products/:id',
    loadComponent: () => import('./pages/products/product-detail/product-detail.component')
      .then(m => m.ProductDetailComponent),
    canActivate: [authGuard]
  },
  {
    path: 'categories',
    loadComponent: () => import('./pages/categories/category-list/category-list.component')
      .then(m => m.CategoryListComponent),
    canActivate: [authGuard]
  },
  {
    path: 'movements',
    loadComponent: () => import('./pages/movements/movement-list/movement-list.component')
      .then(m => m.MovementListComponent),
    canActivate: [authGuard]
  },
  {
    path: 'movements/new',
    loadComponent: () => import('./pages/movements/movement-form/movement-form.component')
      .then(m => m.MovementFormComponent),
    canActivate: [authGuard]
  },
  {
    path: '**',
    redirectTo: '/dashboard'
  }
];