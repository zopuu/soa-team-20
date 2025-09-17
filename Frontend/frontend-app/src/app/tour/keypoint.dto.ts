export interface Coordinates {
  latitude: number;
  longitude: number;
}

export interface KeyPointDto {
  tourId: string;
  title: string;
  description: string;
  coordinates: Coordinates;
  image?: File; // Optional file for image upload
}
