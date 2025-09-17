export enum TourDifficulty {
  Beginner = 'Beginner',
  Intermediate = 'Intermediate',
  Advanced = 'Advanced',
  Pro = 'Pro',
}

export enum TourStatus {
  Draft = 'Draft',
  Published = 'Published',
  Archived = 'Archived',
}

export enum TransportType {
  Walking = 'Walking',
  Bicycle = 'Bicycle',
  Bus = 'Bus',
}

export interface Tour {
  id: string; // UUID
  authorId: string;
  title: string;
  description: string;
  difficulty: number;
  tags: string[];
  status: TourStatus;
  price: number;
  distance: number;
  duration: number;
  publishedAt: string;
  archivedAt: string;
  transportType: number;
}
