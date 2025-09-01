import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Blog } from './blog.model';
import { BlogDto } from './blog.dto';

@Injectable({
  providedIn: 'root',
})
export class BlogService {
  private apiUrl = 'http://localhost:7000/blogs';

  constructor(private http: HttpClient) {}

  getAllBlogs(): Observable<Blog[]> {
    return this.http.get<Blog[]>(`${this.apiUrl}`);
  }

  getById(id: string): Observable<Blog> {
    return this.http.get<Blog>(`${this.apiUrl}/${id}`);
  }

  create(blog: BlogDto): Observable<any> {
    return this.http.post(`${this.apiUrl}`, blog);
  }

  update(id: string, blog: Blog): Observable<any> {
    return this.http.put(`${this.apiUrl}/${id}`, blog);
  }

  delete(id: string): Observable<any> {
    return this.http.delete(`${this.apiUrl}/${id}`);
  }

  getAllByUser(userId: string): Observable<Blog[]> {
    return this.http.get<Blog[]>(`${this.apiUrl}/users/${userId}`);
  }
}
