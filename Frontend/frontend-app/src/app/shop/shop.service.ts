import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { forkJoin, map, Observable, switchMap } from 'rxjs';
import { TourStatus } from '../tour/tour.model';
import { KeyPoint } from '../tour/keypoint.model';

export interface Tour {
  id: string;
  title: string;
  description: string;
  length: number;
  duration: number;
  status: number;
  imageUrl: string;
  price: number;
  firstKeypoint: string;
  keypoints: KeyPoint[];
  recensions: TourReview[];
}
export interface CartItem {
    tour_id: string;
    name: string;
    price: number;
}
export interface TourReview {
  id: string;
  touristName: string;
  touristEmail: string;
  rating: number;
  comment: string;
  tourId: string;
  commentedAt: string; // ISO date string
  visitedAt: string;   // ISO date string
  createdAt: string;   // ISO date string
}

export interface PurchaseToken {
  id: string;
  user_id: string;
  tour_id: string;
  tourName: string;
  price: number;
  purchased_at: string;
}

@Injectable({
  providedIn: 'root',
})
export class ShopService {
  private baseUrl = 'http://localhost:7000/api/shopping';

  constructor(private http: HttpClient) {}

  getAllTours(): Observable<Tour[]> {
    return this.http.get<Tour[]>(`http://localhost:7000/tours`).pipe(
      map((tours: Tour[]) => tours.filter(t => t.status === 1))
    );
  }
  /** Fetch all purchase tokens for a user */
  getTokens(userId: string): Observable<PurchaseToken[]> {
  return this.http
    .get<{ tokens: PurchaseToken[] }>(`${this.baseUrl}/tokens?userId=${userId}`)
    .pipe(map((res) => res.tokens || []));
}

  /** Fetch a tour by id */
  getTourById(tourId: string): Observable<Tour> {
    return this.http.get<Tour>(`http://localhost:7000/tours/${tourId}`);
  }

  /** Combine tokens + tour details */
  getMyTours(userId: string): Observable<Tour[]> {
    return this.getTokens(userId).pipe(
      switchMap((tokens) => {
        if (!tokens || tokens.length === 0) {
          return new Observable<Tour[]>((obs) => {
            obs.next([]);
            obs.complete();
          });
        }
        console.log('Tokens:', tokens);
        const tourRequests = tokens.map((token) =>
          this.getTourById(token.tour_id)
        );

        return forkJoin(tourRequests);
      })
    );
  }
  getReviews(tourId: string): Observable<any[]> {
    return this.http.get<any[]>(`http://localhost:7000/tours/${tourId}/reviews`);
  }

  addToCart(tourId: string, userId: string, name: string, price: number): Observable<any> {
    console.log('Adding to cart:', { tourId, userId, name, price });
    return this.http.post(`${this.baseUrl}/cart/add`, { userId,
        item: { tourId, name, price }
     });
  }

  removeFromCart(tourId: string, userId: string): Observable<any> {
    return this.http.post(`${this.baseUrl}/cart/remove`, { tourId, userId });
  }

  getCart(userId: string): Observable<any> {
    return this.http.get(`${this.baseUrl}/cart?userId=${userId}`);
  }

  checkout(userId: string): Observable<any> {
    return this.http.post(`${this.baseUrl}/checkout`, {userId});
  }
}
