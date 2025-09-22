import { Component, OnInit } from '@angular/core';
import { CartItem, ShopService, Tour } from './shop.service';
import { AuthService } from '../auth/auth.service';
import { Router } from '@angular/router';
import { KeypointService } from '../tour/keypoint.service';
import { Image } from '../tour/image.model';



@Component({
  selector: 'app-shop',
  templateUrl: './shop.component.html',
  styleUrls: ['./shop.component.css'],
})
export class ShopComponent implements OnInit {
  tours: Tour[] = [];
  cart: CartItem[] = [];
  currentUserId?: string;
  cartOpen = false;
  myTours: Tour[] = [];
  imageIndexes: { [tourId: string]: number } = {};

  constructor(private shopService: ShopService, private auth: AuthService, private router: Router,
    private keypointService: KeypointService
  ) {}

  ngOnInit(): void {
    this.loadTours();
    
    this.auth.whoAmI().subscribe({
      next: (user) => {
        this.currentUserId = user.id?.toString();
        this.loadCart();
        this.shopService.getMyTours(this.currentUserId).subscribe({
        next: (tours) => {
            this.myTours = tours;
            console.log('myTours before loop:', this.myTours);
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
            
        },
        });
      },
      error: () => {
        this.currentUserId = undefined;
        
      },
    });
    
  }
  isInMyTours(tour: Tour): boolean {
    return this.myTours.some(t => t.id === tour.id);
  }
  loadTours(): void {
    this.shopService.getAllTours().subscribe({
      next: (data) => {
        console.log('Fetched tours:', data);
        this.tours = data
        this.tours.forEach(t => {
                this.shopService.getReviews(t.id).subscribe({
                next: (reviews) => {
                    console.log('Reviews for tour', t.id, reviews);
                    t.recensions = reviews; }
                });
                this.loadKeypoints(t);
            });
    },
      error: (err) => console.error('Error fetching tours', err),
    });
  }
  loadKeypoints(tour: Tour): void {
    this.keypointService.getByTourSorted(tour.id).subscribe({
      next: (data) => {
        tour.keypoints = data;}
      })
}
  viewTour(tour: Tour): void {
    if(this.isInMyTours(tour)){
        this.router.navigate(['/tours/view', tour.id, { ref: 'mytours' }]);
    }else{
        this.router.navigate(['/tours/view', tour.id, { ref: '' }]);
    }
    
  }

  loadCart(): void {
    if(!this.currentUserId) return;
    this.shopService.getCart(this.currentUserId).subscribe({
      
      next: (data) => {
        console.log('Cart data:', data);
        this.cart = data.cart.items || [];
      },
      error: (err) => console.error('Error fetching tours', err),
    });
  }

  toggleCart(): void {
    this.cartOpen = !this.cartOpen;
  }

  addToCart(tour: Tour): void {
    console.log('Adding to cart:', tour);
    this.cart.push({ tour_id: tour.id, name: tour.title, price: tour.price });
    if(!this.currentUserId) {
      alert('You must be logged in to add items to the cart.');
      return;
    }
    this.shopService.addToCart(tour.id, this.currentUserId, tour.title, tour.price).subscribe();
  }

  removeFromCart(tour: CartItem): void {
    console.log('Removing from cart:', tour);
    //this.cart = this.cart.filter((t) => t.id !== tour.id);
    if(!this.currentUserId) return;
    this.shopService.removeFromCart(tour.tour_id, this.currentUserId).subscribe();
  }
  remove(): void{
    console.log('Removing from cart:');
  }

  getTotalPrice(): number {
    return this.cart.reduce((acc, tour) => acc + tour.price, 0);
  }

  checkout(): void {
    if(!this.currentUserId) {
      alert('You must be logged in to checkout.');
      return;
    }
    this.shopService.checkout(this.currentUserId).subscribe({
      next: (res) => {
        alert('Checkout successful!');
        this.cart = [];
        this.cartOpen = false;
      },
      error: (err) => {
        console.error('Checkout failed', err);
        alert('Checkout failed');
      },
    });
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

prevImage(tour: Tour): void {
  if (!tour.keypoints?.length) return;
  const current = this.imageIndexes[tour.id] ?? 0;
  this.imageIndexes[tour.id] = (current - 1 + tour.keypoints.length) % tour.keypoints.length;
}
}
