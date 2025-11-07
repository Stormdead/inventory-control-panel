import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatCardModule } from '@angular/material/card';
import { MatTableModule } from '@angular/material/table';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatFormFieldModule } from '@angular/material/form-field';
import { NavbarComponent } from '../../../shared/components/navbar/navbar.component';
import { SidebarComponent } from '../../../shared/components/sidebar/sidebar.component';
import { LoadingComponent } from '../../../shared/components/loading/loading.component';
import { MovementService } from '../../../core/services/movement.service';
import { Movement } from '../../../shared/models/movement.model';

@Component({
  selector: 'app-movement-list',
  standalone: true,
  imports: [
    CommonModule,
    MatSidenavModule,
    MatCardModule,
    MatTableModule,
    MatIconModule,
    MatButtonModule,
    MatSelectModule,
    MatFormFieldModule,
    NavbarComponent,
    SidebarComponent,
    LoadingComponent
  ],
  templateUrl: './movement-list.component.html',
  styleUrl: './movement-list.component.scss'
})
export class MovementListComponent implements OnInit {
  loading = true;
  sidebarOpened = true;
  movements: Movement[] = [];
  filterType: 'all' | 'entrada' | 'salida' = 'all';
  displayedColumns: string[] = ['product', 'type', 'quantity', 'user', 'date'];

  constructor(private movementService: MovementService, private router: Router) {}

  ngOnInit(): void {
    this.loadMovements();
  }

  loadMovements(): void {
    this.loading = true;
    this.movementService.getAll().subscribe({
      next: (res) => {
        this.movements = res.movements;
        this.applyFilter();
        this.loading = false;
      },
      error: (err) => {
        console.error('Error cargando movimientos:', err);
        this.loading = false;
      }
    });
  }

  applyFilter(): void {
    if (this.filterType === 'all') {
      // already set
    } else {
      this.movements = this.movements.filter(m => m.type === this.filterType);
    }
  }

  onTypeChange(type: 'all' | 'entrada' | 'salida'): void {
    this.filterType = type;
    this.loadMovements();
  }

  createMovement(): void {
    this.router.navigate(['/movements/new']);
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleString('es-MX', {
      year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit'
    });
  }
}