import { Component, OnInit, ViewChild, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, NavigationEnd } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatTableModule } from '@angular/material/table';
import { MatChipsModule } from '@angular/material/chips';
import { MatSidenavModule } from '@angular/material/sidenav';
import { BaseChartDirective } from 'ng2-charts';
import { ChartConfiguration, ChartData } from 'chart.js';
import { forkJoin, Subscription } from 'rxjs';
import { filter } from 'rxjs/operators';
import { DashboardService } from '../../core/services/dashboard.service';
import { ProductService } from '../../core/services/product.service';
import { NavbarComponent } from '../../shared/components/navbar/navbar.component';
import { SidebarComponent } from '../../shared/components/sidebar/sidebar.component';
import { LoadingComponent } from '../../shared/components/loading/loading.component';
import { DashboardStats, MovementSummary, TopProduct } from '../../shared/models/dashboard.model';
import { Movement } from '../../shared/models/movement.model';
import { Product } from '../../shared/models/product.model';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatIconModule,
    MatButtonModule,
    MatGridListModule,
    MatTableModule,
    MatChipsModule,
    MatSidenavModule,
    BaseChartDirective,
    NavbarComponent,
    SidebarComponent,
    LoadingComponent
  ],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss'
})
export class DashboardComponent implements OnInit, OnDestroy {
  loading = true;
  sidebarOpened = true;

  stats: DashboardStats | null = null;
  movementSummary: MovementSummary | null = null;
  recentMovements: Movement[] = [];
  lowStockProducts: Product[] = [];
  topProducts: TopProduct[] = [];

  // Charts
  @ViewChild('movementChart') movementChart?: BaseChartDirective;
  @ViewChild('topChart') topChart?: BaseChartDirective;

  public movementChartData: ChartData<'bar'> = {
    labels: ['Entradas', 'Salidas'],
    datasets: [{ data: [0, 0], label: 'Cantidad de Movimientos', backgroundColor: ['#4caf50', '#f44336'] }]
  };

  public movementChartOptions: ChartConfiguration<'bar'>['options'] = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: { legend: { display: true, position: 'top' }, title: { display: true, text: 'Movimientos del último mes' } }
  };

  public topProductsChartData: ChartData<'doughnut'> = {
    labels: [],
    datasets: [{ data: [], backgroundColor: ['#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF'] }]
  };

  public topProductsChartOptions: ChartConfiguration<'doughnut'>['options'] = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: { legend: { position: 'right' }, title: { display: true, text: 'Top 5 Productos más Movidos' } }
  };

  displayedColumns: string[] = ['product', 'type', 'quantity', 'date'];

  private routerSub?: Subscription;

  constructor(
    private dashboardService: DashboardService,
    private productService: ProductService,
    private router: Router,
    private cdr: ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    // Primera carga al init
    this.loadDashboardData();

    // Volver a cargar al navegar hacia la ruta /dashboard (evita condiciones en las que la instancia es reutilizada)
    this.routerSub = this.router.events
      .pipe(filter(event => event instanceof NavigationEnd))
      .subscribe((ev: NavigationEnd) => {
        // Si la URL actual es /dashboard (puede ajustar si tu ruta tiene prefijos)
        if (ev.urlAfterRedirects === '/dashboard' || ev.url === '/dashboard') {
          this.loadDashboardData();
        }
      });
  }

  ngOnDestroy(): void {
    this.routerSub?.unsubscribe();
  }

  loadDashboardData(): void {
    this.loading = true;

    forkJoin({
      stats: this.dashboardService.getStats(),
      summary: this.dashboardService.getMovementSummary(),
      recent: this.dashboardService.getRecentMovements(),
      lowStock: this.dashboardService.getLowStockAlerts(),
      top: this.dashboardService.getTopProducts()
    }).subscribe({
      next: (res) => {
        this.stats = res.stats.stats;
        this.movementSummary = res.summary.summary;
        this.recentMovements = res.recent.movements || [];
        this.lowStockProducts = res.lowStock.products || [];
        this.topProducts = res.top.products || [];

        this.updateMovementChart();
        this.updateTopProductsChart();

        // Forzar actualización visual inmediatamente y notificar a Angular
        setTimeout(() => {
          this.movementChart?.update();
          this.topChart?.update();
          this.cdr.detectChanges();
        }, 0);

        this.loading = false;
      },
      error: (err) => {
        console.error('Error cargando dashboard:', err);
        this.loading = false;
        this.cdr.detectChanges();
      }
    });
  }

  updateMovementChart(): void {
    if (this.movementSummary) {
      this.movementChartData.datasets[0].data = [
        this.movementSummary.cantidad_entradas ?? 0,
        this.movementSummary.cantidad_salidas ?? 0
      ];
    }
  }

  updateTopProductsChart(): void {
    this.topProductsChartData.labels = this.topProducts.map(p => p.product_name);
    this.topProductsChartData.datasets[0].data = this.topProducts.map(p => p.total_movements);
  }

  onToggleSidebar(): void { this.sidebarOpened = !this.sidebarOpened; }

  // Navegación desde template
  navigateTo(route: string): void {
    this.router.navigate([route]);
  }

  getMovementTypeClass(type: string): string {
    if (!type) return '';
    return type === 'entrada' ? 'chip-entrada' : 'chip-salida';
  }

  getStockClass(stock: number): string {
    if (stock === 0) return 'stock-critical';
    if (stock < 5) return 'stock-very-low';
    if (stock < 10) return 'stock-low';
    return 'stock-ok';
  }

  formatCurrency(value: number): string {
    return new Intl.NumberFormat('es-MX', { style: 'currency', currency: 'USD' }).format(value);
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleString('es-MX', {
      year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit'
    });
  }
}