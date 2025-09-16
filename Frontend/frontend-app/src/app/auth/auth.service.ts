import { Injectable } from '@angular/core';
import { Observable, tap } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';


interface RegisterReq { username: string; email: string; password: string; role: string; }
interface LoginReq    { username: string; password: string; }
interface LoginRes    { token: string; }
interface WhoAmI  {
  username: string;
  role: string;
  email: string;
  id: number;
};


@Injectable({
  providedIn: 'root'
})
export class AuthService {

  private api = 'http://localhost:7000/api/auth'

  constructor(private http: HttpClient) { }

  login(dto: LoginReq): Observable<LoginRes> {
    return this.http.post<LoginRes>(`${this.api}/login`, dto)
      .pipe(tap(res => localStorage.setItem('token', res.token)));
  }

  register(dto: RegisterReq): Observable<any> {
    return this.http.post(`${this.api}/register`, dto);
  }

  logout() { localStorage.removeItem('token');}

  isLoggedIn() { return !!this.getToken(); }

  getToken() { return localStorage.getItem('token'); }

  whoAmI(): Observable<WhoAmI> {
    return this.http.get<WhoAmI>(`${this.api}/whoami`);
  }
  getRole(): Observable<string> {
    return this.whoAmI().pipe(
      map(user => user.role)
    );
  }
  isAdmin(): Observable<boolean> {
    return this.getRole().pipe(
      map(role => role === 'Admin')
    );
  }
}
