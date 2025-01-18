// /frontend/src/app/components/common/navbar/navbar.component.ts
import { Component } from '@angular/core';
import { AuthService } from '../../../services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-navbar',
  template: `
    <mat-toolbar color="primary">
      <span>Maliaki Social App</span>
      <span class="spacer"></span>
      <button mat-button *ngIf="!authService.getToken()" routerLink="/login">
        Login
      </button>
      <button mat-button *ngIf="!authService.getToken()" routerLink="/register">
        Register
      </button>
      <button mat-button *ngIf="authService.getToken()" (click)="logout()">
        Logout
      </button>
    </mat-toolbar>
  `,
  styles: [
    `
      .spacer {
        flex: 1 1 auto;
      }
    `,
  ],
})
export class NavbarComponent {
  constructor(public authService: AuthService, private router: Router) {}

  logout() {
    this.authService.logout();
    this.router.navigate(['/login']);
  }
}
