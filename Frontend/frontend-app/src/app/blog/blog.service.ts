import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Blog } from './blog.model';
import { BlogDto } from './blog.dto';

@Injectable({
  providedIn: 'root',
})
export class BlogService {
  private apiUrl = 'http://localhost:8080/blogs';

  constructor(private http: HttpClient) {}

  getAllBlogs(): Observable<Blog[]> {
    return this.http.get<Blog[]>(`${this.apiUrl}`);
  }

  getById(id: string): Observable<Blog> {
    return this.http.get<Blog>(`${this.apiUrl}/${id}`);
  }

  create(blog: BlogDto, imageFiles?: File[]): Observable<any> {
    const formData = new FormData();
    formData.append('userId', blog.userId);
    formData.append('title', blog.title);
    formData.append('description', blog.description);

    console.log('Creating blog with:', {
      userId: blog.userId,
      title: blog.title,
      description: blog.description,
      imageCount: imageFiles?.length || 0,
    });

    if (imageFiles && imageFiles.length > 0) {
      imageFiles.forEach((file, index) => {
        console.log(
          `Appending image ${index}: ${file.name}, size: ${file.size}, type: ${file.type}`
        );
        formData.append('images', file);
      });
    }

    return this.http.post(`${this.apiUrl}`, formData);
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
