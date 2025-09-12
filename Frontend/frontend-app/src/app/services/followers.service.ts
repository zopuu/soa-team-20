// src/app/services/followers.service.ts
import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class FollowersService {
  private apiUrl = 'http://localhost:7000/api/followers';

  constructor(private http: HttpClient) {}

  private getAuthHeaders(): HttpHeaders {
    const token = localStorage.getItem('jwt'); // store token after login
    return new HttpHeaders({
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    });
  }

  follow(userId: string, targetId: string): Observable<any> {
    const body = { user_id: userId, target_id: targetId };
    return this.http.post(`${this.apiUrl}/follow`, body, {
      headers: this.getAuthHeaders(),
    });
  }

  unfollow(userId: string, targetId: string): Observable<any> {
    const body = { user_id: userId, target_id: targetId };
    return this.http.post(`${this.apiUrl}/unfollow`, body, {
      headers: this.getAuthHeaders(),
    });
  }

  getRecommendations(userId: string): Observable<any> {
    return this.http.get(`${this.apiUrl}/recommendations/${userId}`, {
      headers: this.getAuthHeaders(),
    });
  }

  getFollowees(userId: string): Observable<any> {
    return this.http.get(`${this.apiUrl}/${userId}/following`, {
      headers: this.getAuthHeaders(),
    });
  }
}
