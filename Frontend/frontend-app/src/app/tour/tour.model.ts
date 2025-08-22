export enum TourDifficulty {
  Easy = 'Easy',
  Medium = 'Medium',
  Hard = 'Hard',
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
