import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Comment {
  id: string;
  userId: string;
  blogId: string;
  dateOfCreation: string;
  text: string;
  lastEdit: string;
}

@Injectable({
  providedIn: 'root'
})
export class CommentService {
  private apiUrl = 'http://localhost:8080/blogs/comments';

  constructor(private http: HttpClient) {}

  getCommentsByBlog(blogId: string): Observable<Comment[]> {
    return this.http.get<Comment[]>(`${this.apiUrl}/${blogId}`);
  }

  addComment(comment: { blogId: string; text: string; userId?: string }): Observable<Comment> {
    return this.http.post<Comment>(this.apiUrl, comment);
  }
}