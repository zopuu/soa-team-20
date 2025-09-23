import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Tour } from './tour.model';
import { TourDto } from './tour.dto';
import { Coordinates } from './keypoint.model';

export interface TourRating {
  id?: string;
  tourId?: string;
  rating: number;
  comment?: string;
  touristName?: string;
  touristEmail?: string;
  visitedAt?: string;     // yyyy-mm-dd (browser <input type="date">)
  commentedAt?: string;   // yyyy-mm-dd
  createdAt?: string;     // ISO from backend
  images?: string[];      // NEW
}

export type CreateTourRatingDto = Omit<TourRating, 'id' | 'tourId' | 'createdAt'>;

@Injectable({
  providedIn: 'root',
})
export class TourService {
  private apiUrl = 'http://localhost:7000/tours';
  private grpcUrl = 'http://localhost:7000/api/tours';

  constructor(private http: HttpClient) {}

  getAllTours(): Observable<Tour[]> {
    return this.http.get<any[]>(`${this.apiUrl}`);
  }

  getById(id: string): Observable<Tour> {
    return this.http.get<any>(`${this.apiUrl}/${id}`);
  }

  create(tour: TourDto): Observable<any> {
    return this.http.post(`${this.grpcUrl}`, tour);
  }

  update(id: string, tour: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/${id}`, tour);
  }

  delete(id: string): Observable<any> {
    return this.http.delete(`${this.grpcUrl}/${id}`);
  }

  getAllByUser(userId: string): Observable<Tour[]> {
    return this.http.get<any[]>(`${this.apiUrl}/users/${userId}`);
  }

  createReview(
    tourId: string,
    review: {
      rating: number;
      comment?: string;
      touristName?: string;
      touristEmail?: string;
      visitedAt?: string; // yyyy-mm-dd
      commentedAt?: string; // yyyy-mm-dd
    }
  ): Observable<any> {
    return this.http.post(`${this.apiUrl}/${tourId}/reviews`, review);
  }

  getReviews(tourId: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/${tourId}/reviews`);
  }

// getReviews(tourId: string): Observable<any[]> {
//   return this.http.get<any[]>(`${this.apiUrl}/${tourId}/reviews`);
// }
  getRatings(tourId: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/${tourId}/reviews`);
  }

  // tour-exec.service.ts
  startExecution(payload: { userId: string; tourId: string }) {
    return this.http.post(`${this.apiUrl}/tour-executions/start`, payload);
  }

  checkProximity(userId: number | string, coords: Coordinates) {
    return this.http.post<{
      reached: boolean;
      distanceMeters: number;
      nextKeyPoint?: { id: string; title: string };
      justCompletedPoint?: { id: string; title: string };
      remainingCount: number;
      completedSession: boolean;
    }>(`${this.apiUrl}/tour-executions/check`, {
      userId: String(userId),
      coords
    });
  }

  getActive(userId: number | string) {
    return this.http.post<any>(`${this.apiUrl}/tour-executions/active`, { userId: String(userId) });
  }

  abandonExecution(userId: number | string) {
    return this.http.post(`${this.apiUrl}/tour-executions/abandon`, { userId: String(userId) });
  }

  abandon(userId: number | string) {
    return this.http.post(`${this.apiUrl}/tour-executions/abandon`, { userId: String(userId) });
  }

  getActiveForTour(userId: string | number, tourId: string) {
    return this.http.post<any>(`${this.apiUrl}/tour-executions/active-for-tour`, {
      userId: String(userId),
      tourId
    });
  }

  purchaseAndStart(tourId: string, userId: string) {
    return this.http.post('/api/orchestrations/purchase-start', {
      userId, // however you expose the logged-in user
      tourId
    });
  }

}
