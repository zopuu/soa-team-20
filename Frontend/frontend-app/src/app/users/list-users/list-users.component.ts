import { Component, OnInit } from '@angular/core';
import { UserService } from '../../services/user.service';
import { AuthService } from 'src/app/auth/auth.service';
import { FollowersService } from 'src/app/services/followers.service';
import { User } from '../../services/user.service';

@Component({
  selector: 'app-list-users',
  templateUrl: './list-users.component.html',
  styleUrls: ['./list-users.component.css']
})
export class ListUsersComponent implements OnInit {
  users: User[] = [];
  recommendedUsers: User[] = [];
  followees: Set<User> = new Set();
  currentUserId: number  = -1;
  defaultProfileImage = 'https://ui-avatars.com/api/?name=User';

  constructor(
    private userService: UserService,
    private authService: AuthService,
    private followersService: FollowersService
  ) {}

  ngOnInit() {
    this.authService.whoAmI().subscribe(currentUser => {
      this.currentUserId = currentUser.id;
      this.loadUsers();
      this.loadFollowees();
      this.loadRecommendations();
    });
  }

  loadRecommendations() {
    this.followersService.getRecommendations(this.currentUserId.toString()).subscribe(result => {
      if(result.user_ids && Array.isArray(result.user_ids)){
        result.user_ids.forEach((r: any) => {
          const curr = this.users.find(u => u.Id == r);
          console.log("Current recommendation: ", curr);
          if (curr) 
          this.recommendedUsers.push(curr);
        });
      }
    })
  }

  loadFollowees() {
    
    this.followersService.getFollowees(this.currentUserId.toString()).subscribe(result => {
      
      if(result.user_ids && Array.isArray(result.user_ids)){
        result.user_ids.forEach((f: any) => {
          const curr = this.users.find(u => u.Id == f);
          console.log("Current followee: ", curr);
          if (curr) 
          this.followees.add(curr);
          
        });
    }
      console.log("Followees: ", this.followees);
    });
  }

  loadUsers() {
    this.userService.getAllUsers().subscribe(users => {
      // Exclude current user
      console.log("Current user ID: " + this.currentUserId);
      this.users = (users || []).filter(u => u.Id != this.currentUserId);
    });
  }

  followUser(user: User) {
    this.followersService.follow(this.currentUserId.toString(), user.Id.toString()).subscribe(() => {
      console.log(`Followed user ${user.Username}`);
      this.followees.add(user);
    });
  }

  unfollowUser(user: User) {
    this.followersService.unfollow(this.currentUserId.toString(), user.Id.toString()).subscribe(() => {
      console.log(`Unfollowed user ${user.Username}`);
      this.followees.delete(user);
    });
  }
}