// /frontend/src/app/components/auth/register/register.component.ts
import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AuthService } from '../../../services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-register',
  template: `
    <mat-card>
      <h2>Register</h2>
      <form [formGroup]="registerForm" (ngSubmit)="onSubmit()">
        <mat-form-field class="full-width">
          <mat-label>Username</mat-label>
          <input matInput formControlName="username" required />
          <mat-error *ngIf="registerForm.get('username')?.invalid">
            Username is required
          </mat-error>
        </mat-form-field>

        <mat-form-field class="full-width">
          <mat-label>Password</mat-label>
          <input matInput type="password" formControlName="password" required />
          <mat-error *ngIf="registerForm.get('password')?.invalid">
            Password is required
          </mat-error>
        </mat-form-field>

        <button
          mat-raised-button
          color="primary"
          type="submit"
          [disabled]="registerForm.invalid"
        >
          Register
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
export class RegisterComponent {
  registerForm: FormGroup;
  error: string = '';

  constructor(
    private fb: FormBuilder,
    private auth: AuthService,
    private router: Router
  ) {
    this.registerForm = this.fb.group({
      username: ['', Validators.required],
      password: ['', Validators.required],
    });
  }

  async onSubmit() {
    if (this.registerForm.invalid) return;

    const { username, password } = this.registerForm.value;
    this.auth.register(username, password).subscribe({
      next: (res) => {
        this.router.navigate(['/']);
      },
      error: (err) => {
        this.error = 'Registration failed';
        console.error('Registration error:', err);
      },
    });
  }
}
