import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterLink, RouterLinkActive } from '@angular/router';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { AuthService } from '../../../core/services/auth.service';

interface MenuItem {
  title: string;
  icon: string;
  route: string;
  adminOnly?: boolean;
}

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [
    CommonModule,
    RouterLink,
    RouterLinkActive,
    MatSidenavModule,
    MatListModule,
    MatIconModule
  ],
  templateUrl: './sidebar.component.html',
  styleUrl: './sidebar.component.scss'
})
export class SidebarComponent {
  @Input() opened = true;

  menuItems: MenuItem[] = [
    { title: 'Dashboard', icon: 'dashboard', route: '/dashboard' },
    { title: 'Productos', icon: 'inventory', route: '/products' },
    { title: 'Categorías', icon: 'category', route: '/categories' },
    { title: 'Movimientos', icon: 'swap_horiz', route: '/movements' }
  ];

  constructor(public authService: AuthService) {}

  shouldShowItem(item: MenuItem): boolean {
    if (item.adminOnly) {
      return this.authService.isAdmin();
    }
    return true;
  }
}