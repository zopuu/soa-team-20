export interface Coordinates {
  latitude: number;
  longitude: number;
}

export interface KeyPointDto {
  tourId: string;
  name: string;
  description: string;
  coordinates: Coordinates;
  image: string;
}
