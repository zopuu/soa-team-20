import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { UserService, User } from '../../services/user.service';

@Component({
  selector: 'app-view-user',
  templateUrl: './view-user.component.html',
  styleUrls: ['./view-user.component.css'],
})
export class ViewUserComponent implements OnInit {
  user?: User;
  loading = true;
  error = '';

  constructor(private route: ActivatedRoute, private users: UserService) {}

  ngOnInit(): void {
    console.log(
      "this.route.snapshot.paramMap.get('id')",
      this.route.snapshot.paramMap.get('id')
    );
    const id = Number(this.route.snapshot.paramMap.get('id') ?? 1);

    this.users.getById(id).subscribe({
      next: (u) => {
        this.user = u;
        this.loading = false;
      },
      error: () => {
        this.error = 'Failed to load user.';
        this.loading = false;
      },
    });
  }
}
