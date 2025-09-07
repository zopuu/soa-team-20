import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Like } from './like.model';

@Injectable({
  providedIn: 'root'
})
export class LikeService {
  private apiUrl = 'http://localhost:5100/blogs/likes';

  constructor(private http: HttpClient) {}

  getLikesByBlog(blogId: string): Observable<Like[]> {
    return this.http.get<Like[]>(`${this.apiUrl}/${blogId}`);
  }

  createLike(like: { blogId: string; userId: string }): Observable<Like> {
    return this.http.post<Like>(this.apiUrl, like);
  }

  removeLike(blogId: string, userId: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${userId}/${blogId}`);
  }
}