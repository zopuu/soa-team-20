import { Injectable } from '@angular/core';
import { CanActivate, Router,UrlTree } from '@angular/router';
import { AuthService } from './auth.service';
import { catchError, map, Observable, of } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class AdminGuard implements CanActivate {
    constructor(private authService: AuthService, private router: Router) {}

    canActivate(): Observable<boolean | UrlTree> {
    return this.authService.isAdmin().pipe(
      map(isAdmin => isAdmin ? true : this.router.parseUrl('/forbidden')),
      catchError(() => of(this.router.parseUrl('/login')))
    );
  }
}