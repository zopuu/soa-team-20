// src/app/services/user.service.ts
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface User {
  Id: number;
  Username: string;
  Email: string;
  Role: string;
  CreatedAt: string;
  Description: string;
  FirstName: string;
  LastName: string;
  Moto: string;
  ProfilePhoto: string;
}

export interface UpdateUserPayload {
  FirstName: string;
  LastName: string;
  ProfilePhoto: string;
  Description: string;
  Moto: string;
}

@Injectable({ providedIn: 'root' })
export class UserService {
  constructor(private http: HttpClient) {}

  getById(id: number): Observable<User> {
    return this.http.get<User>(`http://localhost:7000/api/users/userById/${id}`);
  }

  updateUser(id: number, payload: UpdateUserPayload): Observable<User> {
    return this.http.put<User>(`http://localhost:7000/api/users/updateUser/${id}`, payload);
  }

  getAllUsers(): Observable<User[]> {
    return this.http.get<User[]>(`http://localhost:7000/api/users/allUsers`);
  }
}
