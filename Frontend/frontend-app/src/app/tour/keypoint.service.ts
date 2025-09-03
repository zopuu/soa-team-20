import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { KeyPoint } from './keypoint.model';
import { KeyPointDto } from './keypoint.dto';

@Injectable({ providedIn: 'root' })
export class KeypointService {
  private api = 'http://localhost:7000/keyPoints';

  constructor(private http: HttpClient) {}

  create(dto: KeyPointDto, imageFile?: File): Observable<KeyPoint> {
    const formData = new FormData();
    formData.append('tourId', dto.tourId);
    formData.append('title', dto.title);
    formData.append('description', dto.description);
    formData.append('latitude', dto.coordinates.latitude.toString());
    formData.append('longitude', dto.coordinates.longitude.toString());

    if (imageFile) {
      formData.append('image', imageFile);
    }

    return this.http.post<KeyPoint>(this.api, formData);
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

  update(id: string, dto: KeyPointDto, imageFile?: File): Observable<KeyPoint> {
    const formData = new FormData();
    formData.append('tourId', dto.tourId);
    formData.append('title', dto.title);
    formData.append('description', dto.description);
    formData.append('latitude', dto.coordinates.latitude.toString());
    formData.append('longitude', dto.coordinates.longitude.toString());

    if (imageFile) {
      formData.append('image', imageFile);
    }

    return this.http.put<KeyPoint>(`${this.api}/${id}`, formData);
  }

  delete(id: string): Observable<void> {
    return this.http.delete<void>(`${this.api}/${id}`);
  }

  getImageUrl(id: string): string {
    return `${this.api}/${id}/image`;
  }
}
