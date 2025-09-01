import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { KeyPoint } from './keypoint.model';
import { KeyPointDto } from './keypoint.dto';

@Injectable({ providedIn: 'root' })
export class KeypointService {
  private api = 'http://localhost:7000/keyPoints';

  constructor(private http: HttpClient) {}

  create(dto: KeyPointDto): Observable<KeyPoint> {
    return this.http.post<KeyPoint>(this.api, dto as any);
  }

  getAll(): Observable<KeyPoint[]> {
    return this.http.get<KeyPoint[]>(this.api);
  }

  getById(id: string): Observable<KeyPoint> {
    return this.http.get<KeyPoint>(`${this.api}/${id}`);
  }

  getByTour(tourId: string): Observable<KeyPoint[]> {
    return this.http.get<KeyPoint[]>(`${this.api}/tours/${tourId}`);
  }

  getByTourSorted(tourId: string): Observable<KeyPoint[]> {
    return this.http.get<KeyPoint[]>(
      `${this.api}/tours/${tourId}/sortedByCreatedAt`
    );
  }

  update(id: string, dto: KeyPointDto): Observable<KeyPoint> {
    return this.http.put<KeyPoint>(`${this.api}/${id}`, dto as any);
  }

  delete(id: string): Observable<void> {
    return this.http.delete<void>(`${this.api}/${id}`);
  }
}
