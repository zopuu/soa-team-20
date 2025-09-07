import { Component, OnInit } from '@angular/core';
import { UserService } from '../../services/user.service';
import { AuthService } from 'src/app/auth/auth.service';
import { User } from '../../services/user.service';

@Component({
  selector: 'app-list-users',
  templateUrl: './list-users.component.html',
  styleUrls: ['./list-users.component.css']
})
export class ListUsersComponent implements OnInit {
  users: User[] = [];
  currentUserId: number  = -1;
  defaultProfileImage = 'https://ui-avatars.com/api/?name=User';

  constructor(
    private userService: UserService,
    private authService: AuthService
  ) {}

  ngOnInit() {
    this.authService.whoAmI().subscribe(currentUser => {
      this.currentUserId = currentUser.id;
      this.loadUsers();
    });
  }

  loadUsers() {
    this.userService.getAllUsers().subscribe(users => {
      // Exclude current user
      console.log("Current user ID: " + this.currentUserId);
      this.users = (users || []).filter(u => u.Id !== this.currentUserId);
    });
  }

  followUser(user: User) {
    /*user.following = true;
    this.userService.followUser(user.id).subscribe({
      next: () => {
        user.isFollowed = true;
        user.following = false;
      },
      error: () => {
        user.following = false;
      }
    });*/
  }
}