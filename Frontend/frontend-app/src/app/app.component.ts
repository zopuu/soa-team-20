import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from './auth/auth.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  constructor(private auth: AuthService, private router: Router) {}

  goToMyProfile() {
    this.auth.whoAmI().subscribe({
      next: me => this.router.navigate(['/users', me.id, 'view']),
      error: () => this.router.navigate(['/auth/login'])
    });
  }

  goToEditProgile() {
    this.auth.whoAmI().subscribe({
      next: me => this.router.navigate(['/users', me.id, 'edit']),
      error: () => this.router.navigate(['/auth/login'])
    });
  }
}
