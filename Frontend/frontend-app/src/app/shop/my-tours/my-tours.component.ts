import { Component, OnInit } from '@angular/core';
import { ShopService, Tour } from '../shop.service';
import { AuthService } from 'src/app/auth/auth.service';
import { Router } from '@angular/router';


@Component({
  selector: 'app-my-tours',
  templateUrl: './my-tours.component.html',
  styleUrls: ['./my-tours.component.css'],
})
export class MyToursComponent implements OnInit {
  myTours: Tour[] = [];
  loading = true;
  currentUserId?: string;

  constructor(private shopService: ShopService, private auth: AuthService, private router: Router) {}

  ngOnInit(): void {

    this.auth.whoAmI().subscribe({
      next: (user) => {
        this.currentUserId = user.id?.toString();
        this.shopService.getMyTours(this.currentUserId).subscribe({
        next: (tours) => {
            this.myTours = tours;
            this.loading = false;
            this.myTours.forEach(t => {
                this.shopService.getReviews(t.id).subscribe({
                next: (reviews) => {
                    console.log('Reviews for tour', t.id, reviews);
                    t.recensions = reviews; }
                });
            });
        },
        error: (err) => {
            console.error('Failed to fetch my tours', err);
            this.loading = false;
        },
        });
      },
      error: () => {
        this.currentUserId = undefined;
        
      },
    });
    
  }
  viewTour(tourId: string): void {
    this.router.navigate(['/tours/view', tourId, { ref: 'mytours' }]);
  }
}
