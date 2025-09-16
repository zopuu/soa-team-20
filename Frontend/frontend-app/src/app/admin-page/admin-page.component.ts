// admin-page.component.ts
import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

interface User {
  id: number;
  username: string;
  email: string;
  role: string;
  status: string;
}

@Component({
  selector: 'app-admin-page',
  templateUrl: './admin-page.component.html',
  styleUrls: ['./admin-page.component.css']
})
export class AdminPageComponent implements OnInit {
  users: User[] = [];

  constructor(private http: HttpClient) {}

  ngOnInit(): void {
    this.http.get<User[]>('http://localhost:7000/api/admin')
      .subscribe(data => this.users = data);
  }

  blockUser(id: number) {
    this.http.put(`http://localhost:7000/api/admin/${id}/block`, {})
      .subscribe(() => this.refreshUsers());
  }

  unblockUser(id: number) {
    this.http.put(`http://localhost:7000/api/admin/${id}/unblock`, {})
      .subscribe(() => this.refreshUsers());
  }

  private refreshUsers() {
    this.http.get<User[]>('http://localhost:7000/api/admin')
      .subscribe(data => this.users = data);
  }
}
