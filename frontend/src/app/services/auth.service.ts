// /frontend/src/app/services/auth.service.ts
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable, tap } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private tokenKey = 'token';
  private authSubject = new BehaviorSubject<boolean>(this.hasToken());

  authState = this.authSubject.asObservable();

  constructor(private http: HttpClient) {}

  register(username: string, password: string): Observable<any> {
    return this.http
      .post('http://localhost:8080/register', { username, password })
      .pipe(
        tap((res: any) => {
          if (res.token) {
            localStorage.setItem(this.tokenKey, res.token);
            this.authSubject.next(true);
          }
        })
      );
  }

  login(username: string, password: string): Observable<any> {
    return this.http
      .post('http://localhost:8080/login', { username, password })
      .pipe(
        tap((res: any) => {
          if (res.token) {
            localStorage.setItem(this.tokenKey, res.token);
            this.authSubject.next(true);
          }
        })
      );
  }

  logout() {
    localStorage.removeItem(this.tokenKey);
    this.authSubject.next(false);
  }

  getToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }

  private hasToken(): boolean {
    return !!localStorage.getItem(this.tokenKey);
  }
}
