// /frontend/src/app/components/auth/login/login.component.ts
import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AuthService } from '../../../services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  template: `
    <mat-card>
      <h2>Login</h2>
      <form [formGroup]="loginForm" (ngSubmit)="onSubmit()">
        <mat-form-field class="full-width">
          <mat-label>Username</mat-label>
          <input matInput formControlName="username" required />
          <mat-error *ngIf="loginForm.get('username')?.invalid">
            Username is required
          </mat-error>
        </mat-form-field>

        <mat-form-field class="full-width">
          <mat-label>Password</mat-label>
          <input matInput type="password" formControlName="password" required />
          <mat-error *ngIf="loginForm.get('password')?.invalid">
            Password is required
          </mat-error>
        </mat-form-field>

        <button
          mat-raised-button
          color="primary"
          type="submit"
          [disabled]="loginForm.invalid"
        >
          Login
        </button>

        <div *ngIf="error" class="error">
          {{ error }}
        </div>
      </form>
    </mat-card>
  `,
  styles: [
    `
      mat-card {
        max-width: 400px;
        margin: 50px auto;
        padding: 20px;
      }

      .full-width {
        width: 100%;
      }

      .error {
        color: red;
        margin-top: 10px;
      }
    `,
  ],
})
export class LoginComponent {
  loginForm: FormGroup;
  error: string = '';

  constructor(
    private fb: FormBuilder,
    private auth: AuthService,
    private router: Router
  ) {
    this.loginForm = this.fb.group({
      username: ['', Validators.required],
      password: ['', Validators.required],
    });
  }

  async onSubmit() {
    if (this.loginForm.invalid) return;

    const { username, password } = this.loginForm.value;
    this.auth.login(username, password).subscribe({
      next: (res) => {
        this.router.navigate(['/']);
      },
      error: (err) => {
        this.error = 'Invalid credentials';
        console.error('Login error:', err);
      },
    });
  }
}
