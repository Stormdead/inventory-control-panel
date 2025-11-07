import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatDialogModule } from '@angular/material/dialog';
import { MatSnackBarModule, MatSnackBar } from '@angular/material/snack-bar';
import { MatTableModule } from '@angular/material/table'; // ← agregado
import { NavbarComponent } from '../../../shared/components/navbar/navbar.component';
import { SidebarComponent } from '../../../shared/components/sidebar/sidebar.component';
import { LoadingComponent } from '../../../shared/components/loading/loading.component';
import { ConfirmDialogComponent } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { ProductService } from '../../../core/services/product.service';
import { MovementService } from '../../../core/services/movement.service';
import { Product } from '../../../shared/models/product.model';
import { Movement } from '../../../shared/models/movement.model';
import { MatDialog } from '@angular/material/dialog';

@Component({
  selector: 'app-product-detail',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatIconModule,
    MatButtonModule,
    MatSidenavModule,
    MatDialogModule,
    MatSnackBarModule,
    MatTableModule, // ← agregado aquí
    NavbarComponent,
    SidebarComponent,
    LoadingComponent
  ],
  templateUrl: './product-detail.component.html',
  styleUrl: './product-detail.component.scss'
})
export class ProductDetailComponent implements OnInit {
  loading = true;
  sidebarOpened = true;
  product: Product | null = null;
  recentMovements: Movement[] = [];
  productId: number | null = null;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private productService: ProductService,
    private movementService: MovementService,
    private dialog: MatDialog,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      const id = params['id'];
      if (id) {
        this.productId = +id;
        this.loadProduct(this.productId);
        this.loadMovements(this.productId);
      } else {
        this.router.navigate(['/products']);
      }
    });
  }

  loadProduct(id: number): void {
    this.loading = true;
    this.productService.getById(id).subscribe({
      next: (res) => {
        this.product = res.product;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error cargando producto:', err);
        this.showMessage('Error al cargar el producto');
        this.router.navigate(['/products']);
      }
    });
  }

  loadMovements(id: number): void {
    this.movementService.getByProduct(id).subscribe({
      next: (res) => {
        this.recentMovements = res.movements.slice(0, 10);
      },
      error: (err) => console.error('Error cargando movimientos:', err)
    });
  }

  onToggleSidebar(): void {
    this.sidebarOpened = !this.sidebarOpened;
  }

  editProduct(): void {
    if (this.product) {
      this.router.navigate(['/products/edit', this.product.id]);
    }
  }

  deleteProduct(): void {
    if (!this.product) return;

    const dialogRef = this.dialog.open(ConfirmDialogComponent, {
      width: '400px',
      data: {
        title: 'Eliminar Producto',
        message: `¿Seguro que deseas eliminar "${this.product.name}"?`,
        confirmText: 'Eliminar',
        cancelText: 'Cancelar'
      }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.productService.delete(this.product!.id).subscribe({
          next: () => {
            this.showMessage('Producto eliminado');
            this.router.navigate(['/products']);
          },
          error: (err) => {
            console.error('Error eliminando producto:', err);
            this.showMessage('Error al eliminar producto');
          }
        });
      }
    });
  }

  // Método público para navegación desde el template (evita usar router directamente en template)
  navigateTo(route: string): void {
    this.router.navigate([route]);
  }

  formatCurrency(value: number): string {
    return new Intl.NumberFormat('es-MX', {
      style: 'currency',
      currency: 'USD'
    }).format(value);
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleString('es-MX', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  showMessage(message: string): void {
    this.snackBar.open(message, 'Cerrar', { duration: 3000 });
  }
}