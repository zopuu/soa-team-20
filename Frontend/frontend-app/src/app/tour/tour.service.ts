import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Tour } from './tour.model';
import { TourDto } from './tour.dto';

@Injectable({
  providedIn: 'root',
})
export class TourService {
  private apiUrl = 'http://localhost:7000/tours';

  constructor(private http: HttpClient) {}

  getAllTours(): Observable<Tour[]> {
    return this.http.get<any[]>(`${this.apiUrl}`);
  }

  getById(id: string): Observable<Tour> {
    return this.http.get<any>(`${this.apiUrl}/${id}`);
  }

  create(tour: TourDto): Observable<any> {
    return this.http.post(`${this.apiUrl}`, tour);
  }

  update(id: string, tour: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/${id}`, tour);
  }

  delete(id: string): Observable<any> {
    return this.http.delete(`${this.apiUrl}/${id}`);
  }

  getAllByUser(userId: string): Observable<Tour[]> {
    return this.http.get<any[]>(`${this.apiUrl}/users/${userId}`);
  }
}
