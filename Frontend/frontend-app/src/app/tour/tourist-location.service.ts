import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Coordinates { latitude: number; longitude: number; }
export interface CurrentLocation { userId: string; coordinates: Coordinates; updatedAt: string; }

@Injectable({ providedIn: 'root' })
export class TouristLocationService {
  private api = 'http://localhost:7000/simulator/location';

  constructor(private http: HttpClient) {}

  get(userId: number): Observable<CurrentLocation | null> {
    return this.http.get<CurrentLocation | null>(`${this.api}/${userId}`);
  }

  set(userId: number, coordinates: Coordinates): Observable<void> {
    return this.http.put<void>(this.api, { userId, coordinates });
  }
}
