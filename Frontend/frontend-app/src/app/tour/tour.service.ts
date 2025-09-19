import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Tour } from './tour.model';
import { TourDto } from './tour.dto';

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
}
