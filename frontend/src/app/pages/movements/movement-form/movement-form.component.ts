import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatIconModule } from '@angular/material/icon'; // ← agregado
import { NavbarComponent } from '../../../shared/components/navbar/navbar.component';
import { SidebarComponent } from '../../../shared/components/sidebar/sidebar.component';
import { LoadingComponent } from '../../../shared/components/loading/loading.component';
import { MovementService } from '../../../core/services/movement.service';
import { ProductService } from '../../../core/services/product.service';
import { Product } from '../../../shared/models/product.model';

@Component({
  selector: 'app-movement-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatSidenavModule,
    MatSnackBarModule,
    MatIconModule, // ← agregado aquí
    NavbarComponent,
    SidebarComponent,
    LoadingComponent
  ],
  templateUrl: './movement-form.component.html',
  styleUrl: './movement-form.component.scss'
})
export class MovementFormComponent implements OnInit {
  movementForm: FormGroup;
  loading = true;
  submitting = false;
  sidebarOpened = true;
  products: Product[] = [];

  constructor(
    private fb: FormBuilder,
    private movementService: MovementService,
    private productService: ProductService,
    private router: Router,
    private snackBar: MatSnackBar
  ) {
    this.movementForm = this.fb.group({
      product_id: [null, [Validators.required]],
      type: ['entrada', [Validators.required]],
      quantity: [1, [Validators.required, Validators.min(1)]],
      description: ['']
    });
  }

  ngOnInit(): void {
    this.loadProducts();
  }

  loadProducts(): void {
    this.productService.getAll().subscribe({
      next: (res) => {
        this.products = res.products;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error cargando productos:', err);
        this.loading = false;
      }
    });
  }

  onSubmit(): void {
    if (this.movementForm.invalid) {
      this.movementForm.markAllAsTouched();
      return;
    }

    this.submitting = true;
    const payload = this.movementForm.value;
    payload.product_id = +payload.product_id;

    this.movementService.create(payload).subscribe({
      next: (res) => {
        this.showMessage(res.message || 'Movimiento creado');
        this.router.navigate(['/movements']);
      },
      error: (err) => {
        console.error('Error creando movimiento:', err);
        this.showMessage(err.error?.error || 'Error creando movimiento');
        this.submitting = false;
      }
    });
  }

  onToggleSidebar(): void {
    this.sidebarOpened = !this.sidebarOpened;
  }

  cancel(): void {
    this.router.navigate(['/movements']);
  }

  showMessage(message: string): void {
    this.snackBar.open(message, 'Cerrar', { duration: 3000 });
  }
}