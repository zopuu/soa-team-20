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

export interface Tour {
  id: string; // UUID
  authorId: string;
  title: string;
  description: string;
  difficulty: TourDifficulty;
  tags: string[];
  status: TourStatus;
  price: number;
}
