import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatChipsModule } from '@angular/material/chips';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatTooltipModule } from '@angular/material/tooltip';
import { FormsModule } from '@angular/forms';
import { ProductService } from '../../../core/services/product.service';
import { CategoryService } from '../../../core/services/category.service';
import { AuthService } from '../../../core/services/auth.service';
import { NavbarComponent } from '../../../shared/components/navbar/navbar.component';
import { SidebarComponent } from '../../../shared/components/sidebar/sidebar.component';
import { LoadingComponent } from '../../../shared/components/loading/loading.component';
import { ConfirmDialogComponent } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { Product } from '../../../shared/models/product.model';
import { Category } from '../../../shared/models/category.model';

@Component({
  selector: 'app-product-list',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatTableModule,
    MatChipsModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatSidenavModule,
    MatDialogModule,
    NavbarComponent,
    MatTooltipModule,
    SidebarComponent,
    LoadingComponent
  ],
  templateUrl: './product-list.component.html',
  styleUrl: './product-list.component.scss'
})
export class ProductListComponent implements OnInit {
  loading = true;
  sidebarOpened = true;
  products: Product[] = [];
  filteredProducts: Product[] = [];
  categories: Category[] = [];
  
  // Filtros
  searchTerm = '';
  selectedCategory = 'all';
  stockFilter = 'all';

  displayedColumns: string[] = ['name', 'category', 'price', 'stock', 'actions'];

  constructor(
    private productService: ProductService,
    private categoryService: CategoryService,
    public authService: AuthService,
    private router: Router,
    private dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.loadProducts();
    this.loadCategories();
  }

  loadProducts(): void {
    this.loading = true;
    this.productService.getAll().subscribe({
      next: (response) => {
        this.products = response.products;
        this.filteredProducts = this.products;
        this.applyFilters();
        this.loading = false;
      },
      error: (error) => {
        console.error('Error cargando productos:', error);
        this.loading = false;
      }
    });
  }

  loadCategories(): void {
    this.categoryService.getAll().subscribe({
      next: (response) => {
        this.categories = response.categories;
      },
      error: (error) => console.error('Error cargando categorías:', error)
    });
  }

  applyFilters(): void {
    this.filteredProducts = this.products.filter(product => {
      // Filtro de búsqueda
      const matchesSearch = !this.searchTerm || 
        product.name.toLowerCase().includes(this.searchTerm.toLowerCase()) ||
        product.description.toLowerCase().includes(this.searchTerm.toLowerCase());

      // Filtro de categoría
      const matchesCategory = this.selectedCategory === 'all' || 
        product.category_id?.toString() === this.selectedCategory;

      // Filtro de stock
      let matchesStock = true;
      if (this.stockFilter === 'low') {
        matchesStock = product.stock < 10;
      } else if (this.stockFilter === 'out') {
        matchesStock = product.stock === 0;
      } else if (this.stockFilter === 'available') {
        matchesStock = product.stock >= 10;
      }

      return matchesSearch && matchesCategory && matchesStock;
    });
  }

  onSearchChange(): void {
    this.applyFilters();
  }

  onCategoryChange(): void {
    this.applyFilters();
  }

  onStockFilterChange(): void {
    this.applyFilters();
  }

  clearFilters(): void {
    this.searchTerm = '';
    this.selectedCategory = 'all';
    this.stockFilter = 'all';
    this.applyFilters();
  }

  getStockClass(stock: number): string {
    if (stock === 0) return 'stock-out';
    if (stock < 5) return 'stock-critical';
    if (stock < 10) return 'stock-low';
    return 'stock-good';
  }

  getStockLabel(stock: number): string {
    if (stock === 0) return 'Agotado';
    if (stock < 5) return 'Crítico';
    if (stock < 10) return 'Bajo';
    return 'Disponible';
  }

  viewProduct(id: number): void {
    this.router.navigate(['/products', id]);
  }

  editProduct(id: number): void {
    this.router.navigate(['/products/edit', id]);
  }

  deleteProduct(product: Product): void {
    const dialogRef = this.dialog.open(ConfirmDialogComponent, {
      width: '400px',
      data: {
        title: 'Eliminar Producto',
        message: `¿Estás seguro de que deseas eliminar el producto "${product.name}"?`,
        confirmText: 'Eliminar',
        cancelText: 'Cancelar'
      }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.productService.delete(product.id).subscribe({
          next: () => {
            console.log('Producto eliminado');
            this.loadProducts();
          },
          error: (error) => console.error('Error eliminando producto:', error)
        });
      }
    });
  }

  createProduct(): void {
    this.router.navigate(['/products/new']);
  }

  onToggleSidebar(): void {
    this.sidebarOpened = !this.sidebarOpened;
  }

  formatCurrency(value: number): string {
    return new Intl.NumberFormat('es-MX', {
      style: 'currency',
      currency: 'USD'
    }).format(value);
  }
}