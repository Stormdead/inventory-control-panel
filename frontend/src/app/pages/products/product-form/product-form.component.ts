import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { ProductService } from '../../../core/services/product.service';
import { CategoryService } from '../../../core/services/category.service';
import { NavbarComponent } from '../../../shared/components/navbar/navbar.component';
import { SidebarComponent } from '../../../shared/components/sidebar/sidebar.component';
import { LoadingComponent } from '../../../shared/components/loading/loading.component';
import { Product } from '../../../shared/models/product.model';
import { Category } from '../../../shared/models/category.model';

@Component({
  selector: 'app-product-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatSelectModule,
    MatSidenavModule,
    MatSnackBarModule,
    NavbarComponent,
    SidebarComponent,
    LoadingComponent
  ],
  templateUrl: './product-form.component.html',
  styleUrl: './product-form.component.scss'
})
export class ProductFormComponent implements OnInit {
[x: string]: any;
  productForm: FormGroup;
  loading = true;
  submitting = false;
  sidebarOpened = true;
  isEditMode = false;
  productId: number | null = null;
  categories: Category[] = [];

  constructor(
    private fb: FormBuilder,
    private productService: ProductService,
    private categoryService: CategoryService,
    private router: Router,
    private route: ActivatedRoute,
    private snackBar: MatSnackBar
  ) {
    this.productForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3)]],
      description: ['', [Validators.required]],
      category_id: [null],
      price: [0, [Validators.required, Validators.min(0.01)]],
      stock: [0, [Validators.required, Validators.min(0)]],
      image_url: ['']
    });
  }

  ngOnInit(): void {
    this.loadCategories();
    
    // Verificar si es modo edición
    this.route.params.subscribe(params => {
      if (params['id']) {
        this.isEditMode = true;
        this.productId = +params['id'];
        this.loadProduct(this.productId);
      } else {
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

  loadProduct(id: number): void {
    this.productService.getById(id).subscribe({
      next: (response) => {
        const product = response.product;
        this.productForm.patchValue({
          name: product.name,
          description: product.description,
          category_id: product.category_id,
          price: product.price,
          stock: product.stock,
          image_url: product.image_url
        });
        this.loading = false;
      },
      error: (error) => {
        console.error('Error cargando producto:', error);
        this.showMessage('Error al cargar el producto');
        this.router.navigate(['/products']);
      }
    });
  }

  onSubmit(): void {
    if (this.productForm.invalid) {
      this.productForm.markAllAsTouched();
      return;
    }

    this.submitting = true;
    const formData = this.productForm.value;

    // Convertir category_id a número o null
    if (formData.category_id) {
      formData.category_id = +formData.category_id;
    }

    const request = this.isEditMode && this.productId
      ? this.productService.update(this.productId, formData)
      : this.productService.create(formData);

    request.subscribe({
      next: (response) => {
        const message = this.isEditMode 
          ? 'Producto actualizado exitosamente' 
          : 'Producto creado exitosamente';
        this.showMessage(message);
        this.router.navigate(['/products']);
      },
      error: (error) => {
        console.error('Error guardando producto:', error);
        this.showMessage(error.error?.error || 'Error al guardar el producto');
        this.submitting = false;
      }
    });
  }

  cancel(): void {
    this.router.navigate(['/products']);
  }

  onToggleSidebar(): void {
    this.sidebarOpened = !this.sidebarOpened;
  }

  showMessage(message: string): void {
    this.snackBar.open(message, 'Cerrar', {
      duration: 3000,
      horizontalPosition: 'right',
      verticalPosition: 'top'
    });
  }

  // Getters para validaciones
  get name() {
    return this.productForm.get('name');
  }

  get description() {
    return this.productForm.get('description');
  }

  get price() {
    return this.productForm.get('price');
  }

  get stock() {
    return this.productForm.get('stock');
  }

  // ← Agregar este método nuevo
  onImageError(event: Event): void {
    const target = event.target as HTMLImageElement;
    target.src = 'https://via.placeholder.com/300x300?text=Sin+Imagen';
  }
  
}