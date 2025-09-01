import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from './auth/auth.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent {
  constructor(private auth: AuthService, private router: Router) { }

  private getId(me: any): string | number | undefined {
    return me?.id ?? me?.Id ?? me?.userId ?? me?.uid;
  }

  goToMyProfile() {
    this.auth.whoAmI().subscribe({
      next: (me): void => {
        const id = this.getId(me);
        if (!id) {
          this.router.navigate(['/auth/login']);
          return; // <- bez vrednosti (void)
        }
        this.router.navigate(['/users', id, 'view']);
      },
      error: () => {
        this.router.navigate(['/auth/login']);
      }
    });
  }

  goToEditProfile() { // ispravljeno ime
    this.auth.whoAmI().subscribe({
      next: (me): void => {
        const id = this.getId(me);
        if (!id) { this.router.navigate(['/auth/login']); return; }
        this.router.navigate(['/users', id, 'edit']);
      },
      error: () => { this.router.navigate(['/auth/login']); }
    });
  }

  CreateBlog() {
    this.auth.whoAmI().subscribe({
      next: (me) => this.router.navigate(['/users', me.id, 'create-blog']),
      error: () => this.router.navigate(['/auth/login']),
    });
    console.log('Clicked');
  }

  CreateTour() {
    this.auth.whoAmI().subscribe({
      next: (me) => this.router.navigate(['/users', me.id, 'create-tour']),
      error: () => this.router.navigate(['/auth/login']),
    });
    console.log('Clicked');
  }

  goToAllTours() {
    this.router.navigate(['/tours']);
  }

  goToMyTours() {
    this.auth.whoAmI().subscribe({
      next: (me) => this.router.navigate(['/users', me.id, 'tours']),
      error: () => this.router.navigate(['/auth/login']),
    });
  }
}
