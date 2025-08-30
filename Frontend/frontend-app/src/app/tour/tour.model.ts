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
  Bicyce = 'Bicycle',
  Bus = 'Bus',
}

export interface Tour {
  id: string; // UUID
  authorId: string;
  title: string;
  description: string;
  difficulty: TourDifficulty;
  tags: string[];
  status: TourStatus;
  price: number;
  distance: number;
  duration: number;
  publishedAt: Date;
  archivedAt: Date;
  transportType: TransportType;
}
