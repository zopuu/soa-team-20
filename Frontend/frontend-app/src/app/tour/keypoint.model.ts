import { Image } from './image.model';

export interface Coordinates {
  latitude: number;
  longitude: number;
}

export interface KeyPoint {
  id: string;
  tourId: string;
  title: string;
  description: string;
  coordinates: Coordinates;
  image: Image;
  createdAt: Date;
}
