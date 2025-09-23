import { Component, OnInit } from '@angular/core';
import { ShopService, Tour } from '../shop.service';
import { AuthService } from 'src/app/auth/auth.service';
import { Router } from '@angular/router';
import { KeypointService } from 'src/app/tour/keypoint.service';


@Component({
  selector: 'app-my-tours',
  templateUrl: './my-tours.component.html',
  styleUrls: ['./my-tours.component.css'],
})
export class MyToursComponent implements OnInit {
  myTours: Tour[] = [];
  loading = true;
  currentUserId?: string;
  imageIndexes: { [tourId: string]: number } = {};

  constructor(private shopService: ShopService, private auth: AuthService, private router: Router,
    private keypointService: KeypointService
  ) {}

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
                this.loadKeypoints(t);
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

  loadKeypoints(tour: Tour): void {
    this.keypointService.getByTourSorted(tour.id).subscribe({
      next: (data) => {
        tour.keypoints = data;}
      })
}
  getFirstKeypointImage(tour: Tour): string {
  //  const keypoints = tour.keypoints || [];
  //const index = this.imageIndexes[tour.id] ?? 0;
  if(tour &&tour.keypoints && tour.keypoints.length > 0){
  const img = tour.keypoints[0].image;

  if (img?.data && img?.mimeType) {
    //return `data:${img.mimeType};base64,${img.data}`;
    
    return this.createImageDataUrl(img.data, img.mimeType);
  }
}

  // fallback to tour.imageUrl or placeholder
  return 'assets/placeholder.jpg';
}
  getCurrentImage(tour: Tour): string {
  const keypoints = tour.keypoints || [];
  const index = this.imageIndexes[tour.id] ?? 0;
  const img = keypoints[index]?.image;

  if (img?.data && img?.mimeType) {
    //return `data:${img.mimeType};base64,${img.data}`;
    return this.createImageDataUrl(img.data, img.mimeType);
  }

  // fallback to tour.imageUrl or placeholder
  return tour.imageUrl || 'assets/placeholder.jpg';
}
private createImageDataUrl(base64Data: string, mimeType: string): string {
    // If the base64 data doesn't have the data URL prefix, add it
    if (base64Data.startsWith('data:')) {
      return base64Data;
    }
    return `data:${mimeType};base64,${base64Data}`;
  }

  nextImage(tour: Tour): void {
  if (!tour.keypoints?.length) return;
  const current = this.imageIndexes[tour.id] ?? 0;
  this.imageIndexes[tour.id] = (current + 1) % tour.keypoints.length;
}
isInMyTours(tour: Tour): boolean {
    return this.myTours.some(t => t.id === tour.id);
  }
prevImage(tour: Tour): void {
  if (!tour.keypoints?.length) return;
  const current = this.imageIndexes[tour.id] ?? 0;
  this.imageIndexes[tour.id] = (current - 1 + tour.keypoints.length) % tour.keypoints.length;
}

  viewTour(tour: Tour): void {
   
        this.router.navigate(['/tours/view', tour.id, { ref: 'mytours' }]);
    
  
  }
}
