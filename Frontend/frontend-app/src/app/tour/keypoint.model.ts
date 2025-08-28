export interface Coordinates {
  latitude: number;
  longitude: number;
}

export interface KeyPoint {
  id: string;
  tourId: string;
  name: string;
  description: string;
  coordinates: Coordinates;
  image: string;
  createdAt: Date;
}
